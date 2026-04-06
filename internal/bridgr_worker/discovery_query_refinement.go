package bridgr_worker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	jobsearchv1 "github.com/Kanishkmittal55/bridgr-api/internal/gen/radar/services/job_search/v1"
	"github.com/Kanishkmittal55/bridgr-api/internal/llm/openaijson"
)

// discoveryBatchStats counts rejections for one FindJobs batch (drives query refinement).
type discoveryBatchStats struct {
	ExclusionSkipped  int32
	PrefilterRejected int32
	LowScoreRejected  int32
	AcceptedDelta     int32
}

func copyFindJobsRequest(src *jobsearchv1.FindJobsRequest) *jobsearchv1.FindJobsRequest {
	if src == nil {
		return nil
	}
	tr := append([]string(nil), src.TargetRoles...)
	return &jobsearchv1.FindJobsRequest{
		ResumeText:  src.ResumeText,
		TargetRoles: tr,
		SearchQuery: src.SearchQuery,
		Location:    src.Location,
		MaxResults:  src.MaxResults,
		RunUuid:     src.RunUuid,
	}
}

// discoveryPageMaxResults chooses per-iteration MaxResults: enough headroom for filtering, capped by profile and pageCap.
func discoveryPageMaxResults(need int32, profileCeiling int32, pageCap int32) int32 {
	if need <= 0 {
		need = 1
	}
	n := need + 4
	if n < 8 {
		n = 8
	}
	if profileCeiling > 0 && n > profileCeiling {
		n = profileCeiling
	}
	capN := pageCap
	if capN <= 0 {
		capN = 20
	}
	if n > capN {
		n = capN
	}
	return n
}

func ruleBasedRefine(prefs DiscoveryRequestParams, searchQuery string, targetRoles []string, st discoveryBatchStats) (newQuery string, newRoles []string, note string) {
	q := strings.TrimSpace(searchQuery)
	r := append([]string(nil), targetRoles...)
	if len(r) == 0 && q != "" {
		r = []string{q}
	}

	low := st.LowScoreRejected
	pre := st.PrefilterRejected
	ex := st.ExclusionSkipped

	if low >= pre && low >= 2 && len(prefs.SoftwareStackMustHave) > 0 {
		qLower := strings.ToLower(q)
		for _, raw := range prefs.SoftwareStackMustHave {
			term := strings.ToLower(strings.TrimSpace(raw))
			if term == "" || strings.Contains(qLower, term) {
				continue
			}
			q = strings.TrimSpace(q + " " + raw)
			r = []string{q}
			return q, r, "rule:narrow_stack_term:" + term
		}
	}

	if pre >= low && pre >= 2 {
		fields := strings.Fields(q)
		if len(fields) > 2 {
			q = strings.Join(fields[:len(fields)-1], " ")
			r = []string{q}
			return q, r, "rule:broaden_trim_last_token"
		}
	}

	if ex >= 3 && low == 0 && pre == 0 {
		sg := strings.TrimSpace(prefs.SeniorityGoal)
		if sg != "" && !strings.Contains(strings.ToLower(q), strings.ToLower(sg)) {
			q = strings.TrimSpace(q + " " + sg)
			r = []string{q}
			return q, r, "rule:nudge_seniority"
		}
	}

	return q, r, "rule:none"
}

const queryRefineSystemPrompt = `You refine a job search query for the next FindJobs call.

Given JSON with: current_search_query, current_target_roles, user_preferences (target role, location, seniority, stack must-haves), batch_stats (counts of jobs excluded_previously, prefilter_rejected, low_score_rejected, accepted_this_batch), and rule_based_suggestion (optional query/roles from deterministic rules).

Return ONLY a JSON object:
{"search_query": "<non-empty string>", "target_roles": ["<role>", ...], "rationale": "<one short sentence>"}

The new query must stay in the same career lane as the user; prefer adding/removing concrete tools or seniority tokens over unrelated pivots. If the rule_based_suggestion is already good, you may return it verbatim with a brief rationale.`

type queryRefineLLMResult struct {
	SearchQuery string   `json:"search_query"`
	TargetRoles []string `json:"target_roles"`
	Rationale   string   `json:"rationale"`
}

func runQueryRefineLLM(
	ctx context.Context,
	client *openaijson.Client,
	apiKey, model, baseURL string,
	prefs DiscoveryRequestParams,
	currentQuery string,
	currentRoles []string,
	stats discoveryBatchStats,
	ruleQuery string,
	ruleRoles []string,
	ruleNote string,
) (rationale string, outQuery string, outRoles []string, err error) {
	if client == nil {
		client = &openaijson.Client{}
	}
	if baseURL != "" {
		client.BaseURL = baseURL
	}
	wrapper := map[string]interface{}{
		"current_search_query": currentQuery,
		"current_target_roles": currentRoles,
		"user_preferences":     discoveryPrefsSnapshotForLLM(prefs),
		"batch_stats": map[string]interface{}{
			"exclusion_skipped":  stats.ExclusionSkipped,
			"prefilter_rejected": stats.PrefilterRejected,
			"low_score_rejected": stats.LowScoreRejected,
			"accepted":           stats.AcceptedDelta,
		},
		"rule_based_suggestion": map[string]interface{}{
			"search_query": ruleQuery,
			"target_roles": ruleRoles,
			"note":         ruleNote,
		},
	}
	userJSON, jerr := json.Marshal(wrapper)
	if jerr != nil {
		return "", "", nil, jerr
	}
	raw, cerr := client.ChatCompletionJSON(ctx, apiKey, model, queryRefineSystemPrompt, string(userJSON))
	if cerr != nil {
		return "", "", nil, cerr
	}
	raw = stripJSONFences(raw)
	var out queryRefineLLMResult
	if uerr := json.Unmarshal([]byte(raw), &out); uerr != nil {
		return "", "", nil, fmt.Errorf("query refine json: %w", uerr)
	}
	outQuery = strings.TrimSpace(out.SearchQuery)
	out.TargetRoles = trimStringSlice(out.TargetRoles)
	if outQuery == "" {
		return "", "", nil, fmt.Errorf("empty search_query from model")
	}
	if len(out.TargetRoles) == 0 {
		outRoles = []string{outQuery}
	} else {
		outRoles = out.TargetRoles
	}
	return strings.TrimSpace(out.Rationale), outQuery, outRoles, nil
}

func trimStringSlice(in []string) []string {
	var out []string
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

func refineForNextIteration(
	ctx context.Context,
	p *Processor,
	prefs DiscoveryRequestParams,
	currentQuery string,
	currentRoles []string,
	stats discoveryBatchStats,
	llmCalls *int32,
	accum *discoveryRunAccum,
) (nextQuery string, nextRoles []string, ruleNote, llmRationale string) {
	ruleQ, ruleR, ruleNote := ruleBasedRefine(prefs, currentQuery, currentRoles, stats)
	nextQuery, nextRoles = ruleQ, ruleR
	llmRationale = ""

	if !p.workerOpts.DiscoveryLLMQueryRefineEnabled || p.openAIKey == "" {
		return nextQuery, nextRoles, ruleNote, ""
	}
	budget := p.workerOpts.DiscoveryMaxLLMCallsPerRun
	if budget > 0 && *llmCalls >= budget {
		return nextQuery, nextRoles, ruleNote, ""
	}
	rat, q, roles, err := runQueryRefineLLM(ctx, p.prefilterLLM, p.openAIKey, p.workerOpts.DiscoveryOpenAIModel, p.workerOpts.DiscoveryOpenAIBaseURL,
		prefs, currentQuery, currentRoles, stats, ruleQ, ruleR, ruleNote)
	if err != nil {
		return nextQuery, nextRoles, ruleNote, ""
	}
	*llmCalls++
	if accum != nil {
		accum.llmQueryRefineCalls++
	}
	return q, roles, ruleNote, rat
}
