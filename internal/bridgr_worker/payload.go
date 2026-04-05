package bridgr_worker

// Queue message kinds for the shared Bridgr SQS queue (discriminated JSON).
const (
	KindSkillGapAnalysis = "skill_gap_analysis"
	KindJobDiscovery     = "job_discovery"
)

// QueuePayload is the JSON body for SQS messages (kind discriminates the worker branch).
type QueuePayload struct {
	Kind         string `json:"kind"`
	AnalysisUUID string `json:"analysis_uuid,omitempty"`
	RunUUID      string `json:"run_uuid,omitempty"`
	UserID       int32  `json:"user_id,omitempty"`
}
