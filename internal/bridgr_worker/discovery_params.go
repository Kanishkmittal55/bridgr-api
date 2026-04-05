package bridgr_worker

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	jobsearchv1 "github.com/Kanishkmittal55/bridgr-api/internal/gen/radar/services/job_search/v1"
)

// urlHashHex stable dedupe key for (user_id, url_hash).
func urlHashHex(jobURL string) string {
	s := strings.TrimSpace(strings.ToLower(jobURL))
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// findJobsRequestFromDiscoveryParams maps bridgr.job_search_discovery_runs.request_params JSON
// (from buildDiscoveryRequestParams) to Radar FindJobsRequest. Returns preferred source_board label.
func findJobsRequestFromDiscoveryParams(requestParams []byte) (*jobsearchv1.FindJobsRequest, string, error) {
	if len(requestParams) == 0 {
		return nil, "radar", fmt.Errorf("empty request_params")
	}
	var m map[string]interface{}
	if err := json.Unmarshal(requestParams, &m); err != nil {
		return nil, "radar", err
	}

	roles := stringSliceFromMixed("target_roles", m)
	loc := firstLocationString(m)
	maxN := maxResultsFromParams(m)
	board := firstBoard(m)

	searchQuery := strings.TrimSpace(strings.Join(roles, " "))
	if searchQuery == "" {
		searchQuery = "software engineer"
	}

	req := &jobsearchv1.FindJobsRequest{
		TargetRoles: roles,
		SearchQuery: searchQuery,
		Location:    loc,
		MaxResults:  maxN,
	}
	return req, board, nil
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

func maxResultsFromParams(m map[string]interface{}) int32 {
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
	}
	return 10
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
