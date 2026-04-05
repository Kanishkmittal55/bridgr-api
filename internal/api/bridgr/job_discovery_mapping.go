package bridgr

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	"github.com/Kanishkmittal55/bridgr-api/internal/uuid"
	types "github.com/Kanishkmittal55/bridgr-api/pkg/types"
	"github.com/jackc/pgx/v5/pgtype"
)

func jobSearchProfileFromRow(row *sqlc.BridgrJobSearchProfile) (types.BridgrJobSearchProfile, error) {
	uid, err := oapiUUIDFromPg(row.Uuid)
	if err != nil {
		return types.BridgrJobSearchProfile{}, err
	}
	out := types.BridgrJobSearchProfile{
		Uuid:            uid,
		Id:              row.ID,
		UserId:          row.UserID,
		MaxSurfacedJobs: row.MaxSurfacedJobs,
		CreatedAt:       pgTimestampToTime(row.CreatedAt),
		UpdatedAt:       pgTimestampToTime(row.UpdatedAt),
	}
	if row.CanonicalCvAnalysisUuid.Valid {
		o, err := uuid.ConvertPgUUIDToOapiUUID(row.CanonicalCvAnalysisUuid)
		if err != nil {
			return types.BridgrJobSearchProfile{}, err
		}
		out.CanonicalCvAnalysisUuid = &o
	}
	if tr, ok := decodeStringSlice(row.TargetRoles); ok {
		out.TargetRoles = &tr
	}
	if loc, ok := decodeLocationSlice(row.Locations); ok {
		out.Locations = &loc
	}
	if be, ok := decodeStringSlice(row.BoardsEnabled); ok {
		out.BoardsEnabled = &be
	}
	if m := jsonBytesToMap(row.Matching); m != nil {
		out.Matching = m
	}
	return out, nil
}

func decodeStringSlice(b []byte) ([]string, bool) {
	if len(b) == 0 {
		return nil, false
	}
	var s []string
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, false
	}
	return s, true
}

func decodeLocationSlice(b []byte) ([]map[string]interface{}, bool) {
	if len(b) == 0 {
		return nil, false
	}
	var raw []json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, false
	}
	out := make([]map[string]interface{}, 0, len(raw))
	for _, r := range raw {
		var m map[string]interface{}
		if err := json.Unmarshal(r, &m); err != nil {
			continue
		}
		out = append(out, m)
	}
	return out, true
}

func mustJSONBytes(v interface{}) []byte {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return []byte("null")
	}
	return b
}

func discoveryRunFromRow(row *sqlc.BridgrJobSearchDiscoveryRun) (types.BridgrJobSearchDiscoveryRun, error) {
	uid, err := oapiUUIDFromPg(row.Uuid)
	if err != nil {
		return types.BridgrJobSearchDiscoveryRun{}, err
	}
	out := types.BridgrJobSearchDiscoveryRun{
		Uuid:              uid,
		Id:                row.ID,
		UserId:            row.UserID,
		Status:            row.Status,
		RawCandidateCount: row.RawCandidateCount,
		NewCandidateCount: row.NewCandidateCount,
		CreatedAt:         pgTimestampToTime(row.CreatedAt),
		UpdatedAt:         pgTimestampToTime(row.UpdatedAt),
	}
	if len(row.RequestParams) > 0 {
		out.RequestParams = jsonBytesToMap(row.RequestParams)
	}
	if len(row.RadarMeta) > 0 {
		out.RadarMeta = jsonBytesToMap(row.RadarMeta)
	}
	if row.StartedAt.Valid {
		t := row.StartedAt.Time.UTC()
		out.StartedAt = &t
	}
	if row.CompletedAt.Valid {
		t := row.CompletedAt.Time.UTC()
		out.CompletedAt = &t
	}
	out.ErrorCode = pgTextStringPtr(row.ErrorCode)
	out.ErrorDetail = pgTextStringPtr(row.ErrorDetail)
	out.SqsMessageId = pgTextStringPtr(row.SqsMessageID)
	return out, nil
}

func jobCandidateFromRow(row *sqlc.BridgrJobCandidate) (types.BridgrJobCandidate, error) {
	uid, err := oapiUUIDFromPg(row.Uuid)
	if err != nil {
		return types.BridgrJobCandidate{}, err
	}
	out := types.BridgrJobCandidate{
		Uuid:            uid,
		Id:              row.ID,
		UserId:          row.UserID,
		JobUrl:          row.JobUrl,
		UrlHash:         row.UrlHash,
		IngestionStatus: row.IngestionStatus,
		CreatedAt:       pgTimestampToTime(row.CreatedAt),
		UpdatedAt:       pgTimestampToTime(row.UpdatedAt),
	}
	if row.DiscoveryRunUuid.Valid {
		dr, err := uuid.ConvertPgUUIDToOapiUUID(row.DiscoveryRunUuid)
		if err != nil {
			return types.BridgrJobCandidate{}, err
		}
		out.DiscoveryRunUuid = &dr
	}
	sb := row.SourceBoard
	out.SourceBoard = &sb
	out.SourceJobId = pgTextStringPtr(row.SourceJobID)
	out.ContentHash = pgTextStringPtr(row.ContentHash)
	out.Title = pgTextStringPtr(row.Title)
	out.Company = pgTextStringPtr(row.Company)
	out.Location = pgTextStringPtr(row.Location)
	out.JdText = pgTextStringPtr(row.JdText)
	out.JdS3Uri = pgTextStringPtr(row.JdS3Uri)
	if row.FetchedAt.Valid {
		t := row.FetchedAt.Time.UTC()
		out.FetchedAt = &t
	}
	if len(row.RadarPayload) > 0 {
		out.RadarPayload = jsonBytesToMap(row.RadarPayload)
	}
	out.ApplicationUrl = pgTextStringPtr(row.ApplicationUrl)
	return out, nil
}

func jobNotificationFromRow(row *sqlc.BridgrJobNotification) (types.BridgrJobNotification, error) {
	uid, err := oapiUUIDFromPg(row.Uuid)
	if err != nil {
		return types.BridgrJobNotification{}, err
	}
	jc, err := uuid.ConvertPgUUIDToOapiUUID(row.JobCandidateUuid)
	if err != nil {
		return types.BridgrJobNotification{}, err
	}
	out := types.BridgrJobNotification{
		Uuid:             uid,
		Id:               row.ID,
		UserId:           row.UserID,
		JobCandidateUuid: jc,
		Channel:          row.Channel,
		Status:           row.Status,
		CreatedAt:        pgTimestampToTime(row.CreatedAt),
		UpdatedAt:        pgTimestampToTime(row.UpdatedAt),
	}
	if len(row.Payload) > 0 {
		out.Payload = jsonBytesToMap(row.Payload)
	}
	if row.SentAt.Valid {
		t := row.SentAt.Time.UTC()
		out.SentAt = &t
	}
	if row.SeenAt.Valid {
		t := row.SeenAt.Time.UTC()
		out.SeenAt = &t
	}
	out.ErrorDetail = pgTextStringPtr(row.ErrorDetail)
	return out, nil
}

func profileToCreateParams(userID int32, req types.UpsertBridgrJobSearchProfileRequest) (sqlc.CreateJobSearchProfileParams, error) {
	p := sqlc.CreateJobSearchProfileParams{
		UserID:          userID,
		TargetRoles:     []byte("[]"),
		Locations:       []byte("[]"),
		BoardsEnabled:   []byte("[]"),
		Matching:        []byte("{}"),
		MaxSurfacedJobs: 3,
	}
	if req.TargetRoles != nil {
		p.TargetRoles = mustJSONBytes(*req.TargetRoles)
	}
	if req.Locations != nil {
		p.Locations = mustJSONBytes(*req.Locations)
	}
	if req.BoardsEnabled != nil {
		p.BoardsEnabled = mustJSONBytes(*req.BoardsEnabled)
	}
	if req.Matching != nil {
		b, err := json.Marshal(*req.Matching)
		if err != nil {
			return p, fmt.Errorf("matching: %w", err)
		}
		p.Matching = b
	}
	if req.CanonicalCvAnalysisUuid != nil {
		pg, err := uuid.ConvertOapiUUIDToPgUUID(*req.CanonicalCvAnalysisUuid)
		if err != nil {
			return p, err
		}
		p.CanonicalCvAnalysisUuid = pg
	}
	if req.MaxSurfacedJobs != nil {
		p.MaxSurfacedJobs = *req.MaxSurfacedJobs
	}
	return p, nil
}

func profileToUpdateByUserParams(existing *sqlc.BridgrJobSearchProfile, req types.UpsertBridgrJobSearchProfileRequest) (sqlc.UpdateJobSearchProfileByUserIDParams, error) {
	p := sqlc.UpdateJobSearchProfileByUserIDParams{
		UserID:                  existing.UserID,
		TargetRoles:             existing.TargetRoles,
		Locations:               existing.Locations,
		BoardsEnabled:           existing.BoardsEnabled,
		Matching:                existing.Matching,
		CanonicalCvAnalysisUuid: existing.CanonicalCvAnalysisUuid,
		MaxSurfacedJobs:         existing.MaxSurfacedJobs,
	}
	if req.TargetRoles != nil {
		p.TargetRoles = mustJSONBytes(*req.TargetRoles)
	}
	if req.Locations != nil {
		p.Locations = mustJSONBytes(*req.Locations)
	}
	if req.BoardsEnabled != nil {
		p.BoardsEnabled = mustJSONBytes(*req.BoardsEnabled)
	}
	if req.Matching != nil {
		b, err := json.Marshal(*req.Matching)
		if err != nil {
			return p, fmt.Errorf("matching: %w", err)
		}
		p.Matching = b
	}
	if req.CanonicalCvAnalysisUuid != nil {
		pg, err := uuid.ConvertOapiUUIDToPgUUID(*req.CanonicalCvAnalysisUuid)
		if err != nil {
			return p, err
		}
		p.CanonicalCvAnalysisUuid = pg
	} else {
		// explicit null in JSON omitted — keep existing
	}
	if req.MaxSurfacedJobs != nil {
		p.MaxSurfacedJobs = *req.MaxSurfacedJobs
	}
	return p, nil
}

func discoveryRateWindowStart() pgtype.Timestamp {
	return pgtype.Timestamp{Time: time.Now().UTC().Add(-1 * time.Hour), Valid: true}
}
