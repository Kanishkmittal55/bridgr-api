package bridgr_worker

import (
	"encoding/json"
	"testing"
)

func TestComputeDiscoveryScore_stackMatchAndGap(t *testing.T) {
	prefs := DiscoveryRequestParams{
		SeniorityGoal:         "senior",
		Location:              "San Francisco",
		SoftwareStackMustHave: []string{"Go", "PostgreSQL"},
	}
	cv := "Senior engineer with Go, PostgreSQL, and Kubernetes experience."
	jd := "We need strong Go and PostgreSQL for our backend in San Francisco."
	out := computeDiscoveryScore(prefs, cv, jd, "Senior Backend Engineer", "Acme")

	if out.CompositeScore < 0.5 || out.CompositeScore > 1 {
		t.Fatalf("composite out of range: %v", out.CompositeScore)
	}
	var matched, gaps []string
	if err := json.Unmarshal(out.MatchedSkillsJSON, &matched); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(out.GapSkillsJSON, &gaps); err != nil {
		t.Fatal(err)
	}
	if len(matched) != 2 || len(gaps) != 0 {
		t.Fatalf("matched=%#v gaps=%#v", matched, gaps)
	}
	if out.GapSeverity != "none" {
		t.Fatalf("gap severity: %s", out.GapSeverity)
	}
}

func TestComputeDiscoveryScore_gapSeverity(t *testing.T) {
	prefs := DiscoveryRequestParams{
		SeniorityGoal:         "senior",
		SoftwareStackMustHave: []string{"Go", "Rust", "Python"},
	}
	cv := "Go developer."
	jd := "Looking for Go, Rust, and Python experts."
	out := computeDiscoveryScore(prefs, cv, jd, "Engineer", "Co")
	var gaps []string
	_ = json.Unmarshal(out.GapSkillsJSON, &gaps)
	if len(gaps) < 2 {
		t.Fatalf("expected gaps, got %#v", gaps)
	}
	if out.GapSeverity == "none" {
		t.Fatal("expected non-none gap_severity")
	}
}

func TestMergeDeterministicWithLLM_compositeBlend(t *testing.T) {
	det := computeDiscoveryScore(DiscoveryRequestParams{}, "", "kubernetes docker", "DevOps", "X")
	c := 0.2
	merged := mergeDeterministicWithLLM(det, jobScoreLLMResult{CompositeScore: &c}, "test-model")
	if merged.CompositeScore >= det.CompositeScore {
		t.Fatalf("expected blended composite below deterministic when LLM is low: det=%v merged=%v", det.CompositeScore, merged.CompositeScore)
	}
	if merged.ScoringModel != "deterministic_v1+test-model" {
		t.Fatalf("model: %s", merged.ScoringModel)
	}
}
