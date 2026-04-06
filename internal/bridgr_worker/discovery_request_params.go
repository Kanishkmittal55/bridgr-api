package bridgr_worker

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	jobsearchv1 "github.com/Kanishkmittal55/bridgr-api/internal/gen/radar/services/job_search/v1"
)

// DiscoveryRequestParams is the typed shape of job_search_discovery_runs.request_params
// (JSON from BuildDiscoveryRequestParams and API overrides).
type DiscoveryRequestParams struct {
	UserID                  int32    `json:"user_id"`
	JobSearchProfileUUID    string   `json:"job_search_profile_uuid"`
	TargetRole              string   `json:"target_role"`
	SearchQuery             string   `json:"search_query"`
	TargetRoles             []string `json:"target_roles"`
	Location                string   `json:"location"`
	SourceBoard             string   `json:"source_board"`
	BoardsEnabled           []string `json:"boards_enabled"`
	CareerSwitch            *bool    `json:"career_switch,omitempty"`
	CompanyStage            string   `json:"company_stage"`
	SeniorityGoal           string   `json:"seniority_goal"`
	CompensationGoal        string   `json:"compensation_goal"`
	SoftwareStackMustHave   []string `json:"software_stack_must_have"`
	CanonicalCvAnalysisUUID string   `json:"canonical_cv_analysis_uuid"`
	MaxSurfacedJobs         int32    `json:"max_surfaced_jobs"`
}

const maxResumeTextRunesForFindJobs = 48000

// ParseDiscoveryRequestParams unmarshals request_params JSON into a struct plus the raw map
// for legacy keys (e.g. locations).
func ParseDiscoveryRequestParams(requestParams []byte) (DiscoveryRequestParams, map[string]interface{}, error) {
	if len(requestParams) == 0 {
		return DiscoveryRequestParams{}, nil, fmt.Errorf("empty request_params")
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(requestParams, &raw); err != nil {
		return DiscoveryRequestParams{}, nil, err
	}
	var p DiscoveryRequestParams
	if err := json.Unmarshal(requestParams, &p); err != nil {
		return DiscoveryRequestParams{}, raw, err
	}
	if p.UserID == 0 {
		p.UserID = intFromMap(raw, "user_id")
	}
	if len(p.SoftwareStackMustHave) == 0 {
		p.SoftwareStackMustHave = stringSliceFromMixed("software_stack_must_have", raw)
	}
	if p.MaxSurfacedJobs == 0 {
		p.MaxSurfacedJobs = maxResultsFromParamsMap(raw)
	}
	return p, raw, nil
}

// BuildFindJobsRequest maps discovery params to Radar FindJobsRequest.
// resumeText may be empty when CV resolution is skipped or unavailable.
// runUUID should be the discovery run id for tracing / future CancelCrawl.
func BuildFindJobsRequest(p DiscoveryRequestParams, raw map[string]interface{}, runUUID string, resumeText string) (*jobsearchv1.FindJobsRequest, string, error) {
	if raw == nil {
		raw = map[string]interface{}{}
	}

	searchQuery := strings.TrimSpace(p.TargetRole)
	if searchQuery == "" {
		searchQuery = strings.TrimSpace(p.SearchQuery)
	}
	targetRoles := p.TargetRoles
	if searchQuery != "" && len(targetRoles) == 0 {
		targetRoles = []string{searchQuery}
	} else if len(targetRoles) > 0 && searchQuery == "" {
		searchQuery = strings.TrimSpace(strings.Join(targetRoles, " "))
	}
	if searchQuery == "" {
		targetRoles = stringSliceFromMixed("target_roles", raw)
		searchQuery = strings.TrimSpace(strings.Join(targetRoles, " "))
	}
	if searchQuery == "" {
		searchQuery = "software engineer"
		targetRoles = []string{searchQuery}
	} else if len(targetRoles) == 0 {
		targetRoles = []string{searchQuery}
	}

	loc := strings.TrimSpace(p.Location)
	if loc == "" {
		loc = firstLocationString(raw)
	}

	maxN := p.MaxSurfacedJobs
	if maxN <= 0 {
		maxN = maxResultsFromParamsMap(raw)
	}

	board := strings.TrimSpace(p.SourceBoard)
	if board == "" && len(p.BoardsEnabled) > 0 {
		board = strings.TrimSpace(p.BoardsEnabled[0])
	}
	if board == "" {
		board = firstBoard(raw)
	}

	req := &jobsearchv1.FindJobsRequest{
		ResumeText:  truncateRunes(resumeText, maxResumeTextRunesForFindJobs),
		TargetRoles: targetRoles,
		SearchQuery: searchQuery,
		Location:    loc,
		MaxResults:  maxN,
		RunUuid:     strings.TrimSpace(runUUID),
	}
	return req, board, nil
}

func truncateRunes(s string, maxRunes int) string {
	if maxRunes <= 0 || s == "" {
		return s
	}
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	return string(runes[:maxRunes]) + "\n…"
}

func intFromMap(m map[string]interface{}, key string) int32 {
	v, ok := m[key]
	if !ok || v == nil {
		return 0
	}
	switch x := v.(type) {
	case float64:
		return int32(x)
	case int:
		return int32(x)
	case int32:
		return x
	case int64:
		return int32(x)
	case json.Number:
		n, _ := x.Int64()
		return int32(n)
	case string:
		n, _ := strconv.ParseInt(strings.TrimSpace(x), 10, 32)
		return int32(n)
	default:
		return 0
	}
}

func maxResultsFromParamsMap(m map[string]interface{}) int32 {
	v, ok := m["max_surfaced_jobs"]
	if !ok {
		return 10
	}
	switch x := v.(type) {
	case float64:
		if x > 0 && x <= 100 {
			return int32(x)
		}
	case int:
		if x > 0 && x <= 100 {
			return int32(x)
		}
	case int32:
		if x > 0 && x <= 100 {
			return x
		}
	case json.Number:
		n, err := x.Int64()
		if err == nil && n > 0 && n <= 100 {
			return int32(n)
		}
	}
	return 10
}

func stringSliceFromMixed(key string, m map[string]interface{}) []string {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	raw, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, it := range raw {
		switch x := it.(type) {
		case string:
			if s := strings.TrimSpace(x); s != "" {
				out = append(out, s)
			}
		case float64:
			out = append(out, strconv.FormatInt(int64(x), 10))
		}
	}
	return out
}

func firstLocationString(m map[string]interface{}) string {
	v, ok := m["locations"]
	if !ok || v == nil {
		return ""
	}
	raw, ok := v.([]interface{})
	if !ok {
		return ""
	}
	for _, it := range raw {
		switch x := it.(type) {
		case string:
			if s := strings.TrimSpace(x); s != "" {
				return s
			}
		case map[string]interface{}:
			if s, _ := x["location"].(string); strings.TrimSpace(s) != "" {
				return strings.TrimSpace(s)
			}
		}
	}
	return ""
}

func firstBoard(m map[string]interface{}) string {
	v, ok := m["boards_enabled"]
	if !ok || v == nil {
		return "indeed"
	}
	raw, ok := v.([]interface{})
	if !ok {
		return "indeed"
	}
	for _, it := range raw {
		if s, ok := it.(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	return "indeed"
}
