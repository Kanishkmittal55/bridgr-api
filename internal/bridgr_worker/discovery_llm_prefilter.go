package bridgr_worker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Kanishkmittal55/bridgr-api/internal/llm/openaijson"
)

const (
	jobPrefilterSystemPrompt = `You gate job listings for a job seeker. You must compare the user's discovery preferences to each job (title, company, job description text).

Dimensions to evaluate (use "ok", "mismatch", or "unknown" for each key in by_dimension):
- career_switch: If the user is open to a career change, the role should plausibly fit a pivot; if not switching, the role should stay in the same professional lane implied by target_role.
- company_stage: e.g. startup / series / big tech / any — job or company signals should align with the user's stage preference when it is specific; "any" from user means ok unless the JD clearly contradicts.
- seniority_goal: user's target level (e.g. senior, staff) vs what the role implies.
- compensation_goal: user's pay expectation vs salary/comp signals in the JD if present; if JD missing pay info, use "unknown".
- stack: user's required technologies must appear in the JD (or clear equivalents). If user listed no stack requirements, use "ok".

Return ONLY a JSON object with exactly these keys:
{"pass": <boolean>, "confidence": <number between 0 and 1>, "by_dimension": { "career_switch": "...", "company_stage": "...", "seniority_goal": "...", "compensation_goal": "...", "stack": "..." }}

Set pass=true only if there is no strong mismatch on any dimension the user specified. If information is missing, prefer "unknown" rather than inventing. If the user marked a dimension as "any" or left it open, treat it as ok unless the JD clearly conflicts.`

	maxJdRunesForPrefilterLLM = 14000
)

// JobPrefilterLLMResult is the strict JSON shape returned by the prefilter model.
type JobPrefilterLLMResult struct {
	Pass        bool              `json:"pass"`
	Confidence  *float64          `json:"confidence"`
	ByDimension map[string]string `json:"by_dimension"`
}

func discoveryPrefsSnapshotForLLM(p DiscoveryRequestParams) map[string]interface{} {
	m := map[string]interface{}{
		"target_role":       strings.TrimSpace(p.TargetRole),
		"search_query":      strings.TrimSpace(p.SearchQuery),
		"location":          strings.TrimSpace(p.Location),
		"source_board":      strings.TrimSpace(p.SourceBoard),
		"company_stage":     strings.TrimSpace(p.CompanyStage),
		"seniority_goal":    strings.TrimSpace(p.SeniorityGoal),
		"compensation_goal": strings.TrimSpace(p.CompensationGoal),
		"stack_must_have":   p.SoftwareStackMustHave,
	}
	if p.CareerSwitch != nil {
		m["career_switch"] = *p.CareerSwitch
	} else {
		m["career_switch"] = nil
	}
	return m
}

func runJobPrefilterLLM(
	ctx context.Context,
	client *openaijson.Client,
	apiKey, model, baseURL, userContent string,
) (JobPrefilterLLMResult, error) {
	var zero JobPrefilterLLMResult
	if client == nil {
		client = &openaijson.Client{}
	}
	if baseURL != "" {
		client.BaseURL = baseURL
	}
	raw, err := client.ChatCompletionJSON(ctx, apiKey, model, jobPrefilterSystemPrompt, userContent)
	if err != nil {
		return zero, err
	}
	raw = stripJSONFences(raw)
	var out JobPrefilterLLMResult
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return zero, fmt.Errorf("prefilter json: %w", err)
	}
	if out.ByDimension == nil {
		out.ByDimension = map[string]string{}
	}
	return out, nil
}

func buildPrefilterUserPayload(prefs DiscoveryRequestParams, title, company, jdText string) (string, error) {
	job := map[string]interface{}{
		"title":   strings.TrimSpace(title),
		"company": strings.TrimSpace(company),
		"jd_text": truncateRunes(strings.TrimSpace(jdText), maxJdRunesForPrefilterLLM),
	}
	wrapper := map[string]interface{}{
		"user_preferences": discoveryPrefsSnapshotForLLM(prefs),
		"job":              job,
	}
	b, err := json.Marshal(wrapper)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func stripJSONFences(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "```") {
		return s
	}
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "json")
	s = strings.TrimSpace(s)
	if idx := strings.LastIndex(s, "```"); idx >= 0 {
		s = strings.TrimSpace(s[:idx])
	}
	return s
}

func prefilterConfidenceAccept(result JobPrefilterLLMResult, min float64) bool {
	if !result.Pass {
		return false
	}
	conf := 1.0
	if result.Confidence != nil {
		conf = *result.Confidence
	}
	return conf+1e-9 >= min
}
