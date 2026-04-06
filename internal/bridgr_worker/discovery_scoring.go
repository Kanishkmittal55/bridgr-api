package bridgr_worker

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/Kanishkmittal55/bridgr-api/internal/llm/openaijson"
)

const (
	jobScoreSystemPrompt = `You refine a job relevance score for one listing vs the candidate's CV and discovery preferences.

You receive JSON with: user_preferences, job (title, company, jd_text), deterministic_scores (skill_match, experience_match, location_match, recency, board_quality, composite), deterministic_matched_skills, deterministic_gap_skills, deterministic_gap_severity.

Return ONLY a JSON object with exactly these keys:
{"composite_score": <number 0-1>, "matched_skills": <array of short skill strings>, "gap_skills": <array of short skill strings>, "gap_severity": <one of: none, minor, moderate, major>}

Guidelines:
- composite_score should be consistent with how well the CV fits the JD; stay close to the deterministic composite unless you see clear evidence to adjust.
- matched_skills: skills the candidate clearly has that the role needs (subset or refinement of the deterministic list is ok).
- gap_skills: important required skills the candidate appears to lack.
- gap_severity reflects how blocking the gaps are for this role.`

	maxJdRunesForScoreLLM = 14000
	maxCvRunesForScoreLLM = 12000
)

const discoveryScoringVersion = "discovery-score-p3-v1"

// jobDiscoveryScoreOutput is the full scoring payload for bridgr.job_scores.
type jobDiscoveryScoreOutput struct {
	SkillMatchScore      float32
	ExperienceMatchScore float32
	LocationMatchScore   float32
	RecencyScore         float32
	BoardQualityScore    float32
	CompositeScore       float32
	MatchedSkillsJSON    []byte
	GapSkillsJSON        []byte
	GapSeverity          string
	ScoringModel         string
	ScoringVersion       string
}

// jobScoreLLMResult is the strict JSON from the optional score-refinement model.
type jobScoreLLMResult struct {
	CompositeScore *float64 `json:"composite_score"`
	MatchedSkills  []string `json:"matched_skills"`
	GapSkills      []string `json:"gap_skills"`
	GapSeverity    string   `json:"gap_severity"`
}

var wordSplitRe = regexp.MustCompile(`[^a-z0-9+#.]+`)

var scoreStopwords = map[string]struct{}{
	"the": {}, "and": {}, "for": {}, "with": {}, "that": {}, "this": {}, "from": {}, "have": {},
	"has": {}, "are": {}, "was": {}, "were": {}, "been": {}, "being": {}, "will": {}, "your": {},
	"you": {}, "our": {}, "all": {}, "any": {}, "not": {}, "but": {}, "can": {}, "may": {},
	"one": {}, "two": {}, "per": {}, "via": {}, "into": {}, "also": {}, "more": {}, "most": {},
	"some": {}, "such": {}, "than": {}, "then": {}, "them": {}, "they": {}, "their": {},
	"work": {}, "team": {}, "role": {}, "job": {}, "years": {}, "year": {}, "experience": {},
	"including": {}, "leading": {}, "using": {}, "used": {}, "use": {},
}

func computeDiscoveryScore(prefs DiscoveryRequestParams, resumeText, jdText, jobTitle, jobCompany string) jobDiscoveryScoreOutput {
	jd := strings.ToLower(jdText) + " " + strings.ToLower(jobTitle) + " " + strings.ToLower(jobCompany)
	cv := strings.ToLower(resumeText)

	skill := skillMatchScore(prefs, cv, jd)
	exp := experienceMatchScore(prefs, jobTitle, jdText)
	loc := locationMatchScore(prefs, jd)
	rec := float32(0.75)
	board := float32(0.82)

	composite := 0.42*skill + 0.28*exp + 0.12*loc + 0.10*rec + 0.08*board
	if composite > 1 {
		composite = 1
	}
	if composite < 0 {
		composite = 0
	}

	matched, gaps := matchedAndGapSkills(prefs, cv, jd)
	gapSev := gapSeverityFromCounts(len(matched), len(gaps))

	ms, _ := json.Marshal(matched)
	gs, _ := json.Marshal(gaps)

	return jobDiscoveryScoreOutput{
		SkillMatchScore:      skill,
		ExperienceMatchScore: exp,
		LocationMatchScore:   loc,
		RecencyScore:         rec,
		BoardQualityScore:    board,
		CompositeScore:       composite,
		MatchedSkillsJSON:    ms,
		GapSkillsJSON:        gs,
		GapSeverity:          gapSev,
		ScoringModel:         "deterministic_v1",
		ScoringVersion:       discoveryScoringVersion,
	}
}

func skillMatchScore(prefs DiscoveryRequestParams, cvLower, jdBlobLower string) float32 {
	stack := normalizeSkillTerms(prefs.SoftwareStackMustHave)
	if len(stack) == 0 {
		return wordJaccard(cvLower, jdBlobLower)
	}
	matched := 0
	for _, term := range stack {
		if term == "" {
			continue
		}
		if strings.Contains(jdBlobLower, term) {
			matched++
		}
	}
	stackCover := float32(matched) / float32(max(1, len(stack)))
	jacc := wordJaccard(cvLower, jdBlobLower)
	return clamp01(stackCover*0.72 + jacc*0.28)
}

func experienceMatchScore(prefs DiscoveryRequestParams, title, jd string) float32 {
	goal := seniorityRank(strings.ToLower(strings.TrimSpace(prefs.SeniorityGoal)))
	if goal < 0 {
		return 0.78
	}
	jobLv := inferJobSeniority(strings.ToLower(strings.TrimSpace(title)), strings.ToLower(jd))
	diff := jobLv - goal
	switch {
	case diff >= 0 && diff <= 1:
		return 1
	case diff > 1:
		pen := 0.18 * float32(diff-1)
		return clamp01(1 - pen)
	default:
		pen := 0.22 * float32(-diff)
		return clamp01(1 - pen)
	}
}

func locationMatchScore(prefs DiscoveryRequestParams, jdBlobLower string) float32 {
	loc := strings.TrimSpace(prefs.Location)
	if loc == "" {
		return 0.88
	}
	ln := strings.ToLower(loc)
	if strings.Contains(jdBlobLower, ln) {
		return 1
	}
	for _, tok := range strings.Fields(ln) {
		if len(tok) >= 3 && strings.Contains(jdBlobLower, tok) {
			return 0.92
		}
	}
	return 0.4
}

func seniorityRank(s string) int {
	if s == "" || s == "any" {
		return -1
	}
	if strings.Contains(s, "intern") {
		return 0
	}
	if strings.Contains(s, "junior") || strings.Contains(s, "jr") || strings.Contains(s, "entry") {
		return 1
	}
	if strings.Contains(s, "mid") || strings.Contains(s, "intermediate") || strings.Contains(s, "ii") {
		return 2
	}
	if strings.Contains(s, "staff") || strings.Contains(s, "principal") {
		return 4
	}
	if strings.Contains(s, "director") || strings.Contains(s, "vp") || strings.Contains(s, "vice president") {
		return 5
	}
	if strings.Contains(s, "senior") || strings.Contains(s, "sr") {
		return 3
	}
	if strings.Contains(s, "lead") || strings.Contains(s, "head") {
		return 3
	}
	return 2
}

func inferJobSeniority(titleLower, jdLower string) int {
	s := titleLower + " " + jdLower
	switch {
	case strings.Contains(s, "intern"):
		return 0
	case strings.Contains(s, "junior") || strings.Contains(s, " jr") || strings.Contains(titleLower, "entry"):
		return 1
	case strings.Contains(s, "staff engineer") || strings.Contains(s, "principal"):
		return 4
	case strings.Contains(s, "director") || strings.Contains(s, "vp ") || strings.Contains(s, "vice president"):
		return 5
	case strings.Contains(s, "senior") || strings.Contains(s, " sr"):
		return 3
	case strings.Contains(s, "lead "):
		return 3
	default:
		return 2
	}
}

func normalizeSkillTerms(in []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, raw := range in {
		t := strings.ToLower(strings.TrimSpace(raw))
		if t == "" {
			continue
		}
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
	}
	return out
}

func matchedAndGapSkills(prefs DiscoveryRequestParams, cvLower, jdBlobLower string) (matched, gaps []string) {
	stack := normalizeSkillTerms(prefs.SoftwareStackMustHave)
	if len(stack) == 0 {
		return nil, nil
	}
	for _, term := range stack {
		inJD := strings.Contains(jdBlobLower, term)
		if !inJD {
			continue
		}
		inCV := strings.Contains(cvLower, term)
		if inCV {
			matched = append(matched, term)
		} else {
			gaps = append(gaps, term)
		}
	}
	return matched, gaps
}

func gapSeverityFromCounts(matchedN, gapN int) string {
	total := matchedN + gapN
	if total == 0 || gapN == 0 {
		return "none"
	}
	r := float64(gapN) / float64(total)
	switch {
	case r < 0.25:
		return "minor"
	case r < 0.5:
		return "moderate"
	default:
		return "major"
	}
}

func wordJaccard(a, b string) float32 {
	sa := tokenizeForOverlap(a)
	sb := tokenizeForOverlap(b)
	if len(sa) == 0 || len(sb) == 0 {
		return 0.45
	}
	inter := 0
	for t := range sa {
		if _, ok := sb[t]; ok {
			inter++
		}
	}
	union := len(sa) + len(sb) - inter
	if union == 0 {
		return 0.45
	}
	return clamp01(float32(inter) / float32(union))
}

func tokenizeForOverlap(s string) map[string]struct{} {
	parts := wordSplitRe.Split(strings.ToLower(s), -1)
	m := make(map[string]struct{})
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if len(p) < 3 {
			continue
		}
		if _, stop := scoreStopwords[p]; stop {
			continue
		}
		if len(p) == 3 && !containsLetter(p) {
			continue
		}
		m[p] = struct{}{}
	}
	return m
}

func containsLetter(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func clamp01(x float32) float32 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

func runJobScoreLLM(ctx context.Context, client *openaijson.Client, apiKey, model, baseURL, userJSON string) (jobScoreLLMResult, error) {
	var zero jobScoreLLMResult
	if client == nil {
		client = &openaijson.Client{}
	}
	if baseURL != "" {
		client.BaseURL = baseURL
	}
	raw, err := client.ChatCompletionJSON(ctx, apiKey, model, jobScoreSystemPrompt, userJSON)
	if err != nil {
		return zero, err
	}
	raw = stripJSONFences(raw)
	var out jobScoreLLMResult
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return zero, fmt.Errorf("score llm json: %w", err)
	}
	return out, nil
}

func buildScoreLLMUserPayloadWithCV(prefs DiscoveryRequestParams, title, company, jdText, cvText string, det jobDiscoveryScoreOutput) (string, error) {
	var matched, gaps []string
	_ = json.Unmarshal(det.MatchedSkillsJSON, &matched)
	_ = json.Unmarshal(det.GapSkillsJSON, &gaps)
	wrapper := map[string]interface{}{
		"user_preferences": discoveryPrefsSnapshotForLLM(prefs),
		"job": map[string]interface{}{
			"title":   strings.TrimSpace(title),
			"company": strings.TrimSpace(company),
			"jd_text": truncateRunes(strings.TrimSpace(jdText), maxJdRunesForScoreLLM),
		},
		"cv_excerpt":                   truncateRunes(strings.TrimSpace(cvText), maxCvRunesForScoreLLM),
		"deterministic_scores":         map[string]float32{"skill_match": det.SkillMatchScore, "experience_match": det.ExperienceMatchScore, "location_match": det.LocationMatchScore, "recency": det.RecencyScore, "board_quality": det.BoardQualityScore, "composite": det.CompositeScore},
		"deterministic_matched_skills": matched,
		"deterministic_gap_skills":     gaps,
		"deterministic_gap_severity":   det.GapSeverity,
	}
	b, err := json.Marshal(wrapper)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func mergeDeterministicWithLLM(det jobDiscoveryScoreOutput, llm jobScoreLLMResult, modelName string) jobDiscoveryScoreOutput {
	out := det
	out.ScoringModel = "deterministic_v1+llm"
	if modelName != "" {
		out.ScoringModel = "deterministic_v1+" + modelName
	}
	if llm.CompositeScore != nil {
		c := float32(*llm.CompositeScore)
		if c < 0 {
			c = 0
		}
		if c > 1 {
			c = 1
		}
		out.CompositeScore = clamp01(det.CompositeScore*0.62 + c*0.38)
	}
	if len(llm.MatchedSkills) > 0 {
		ms, _ := json.Marshal(llm.MatchedSkills)
		out.MatchedSkillsJSON = ms
	}
	if len(llm.GapSkills) > 0 {
		gs, _ := json.Marshal(llm.GapSkills)
		out.GapSkillsJSON = gs
	}
	sev := strings.ToLower(strings.TrimSpace(llm.GapSeverity))
	switch sev {
	case "none", "minor", "moderate", "major":
		out.GapSeverity = sev
	}
	return out
}
