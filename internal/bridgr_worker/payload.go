package bridgr_worker

// AnalysisJobPayload is the JSON body posted to the Bridgr skill-gap queue.
type AnalysisJobPayload struct {
	AnalysisUUID string `json:"analysis_uuid"`
}
