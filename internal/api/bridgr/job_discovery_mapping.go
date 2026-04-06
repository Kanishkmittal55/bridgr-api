package bridgr

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"

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
		Uuid:                  uid,
		Id:                    row.ID,
		UserId:                row.UserID,
		TargetRole:            row.TargetRole,
		Location:              row.Location,
		SourceBoard:           types.BridgrJobSearchProfileSourceBoard(row.SourceBoard),
		CareerSwitch:          row.CareerSwitch,
		CompanyStage:          types.BridgrJobSearchProfileCompanyStage(row.CompanyStage),
		SeniorityGoal:         types.BridgrJobSearchProfileSeniorityGoal(row.SeniorityGoal),
		CompensationGoal:      types.BridgrJobSearchProfileCompensationGoal(row.CompensationGoal),
		SoftwareStackMustHave: cloneStringSlice(row.SoftwareStackMustHave),
		CreatedAt:             pgTimestampToTime(row.CreatedAt),
		UpdatedAt:             pgTimestampToTime(row.UpdatedAt),
	}
	if row.CanonicalCvAnalysisUuid.Valid {
		o, err := uuid.ConvertPgUUIDToOapiUUID(row.CanonicalCvAnalysisUuid)
		if err != nil {
			return types.BridgrJobSearchProfile{}, err
		}
		out.CanonicalCvAnalysisUuid = &o
	}
	return out, nil
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
		UserID:                userID,
		TargetRole:            "",
		Location:              "",
		SourceBoard:           "indeed",
		CareerSwitch:          false,
		CompanyStage:          "any",
		SeniorityGoal:         "any",
		CompensationGoal:      "any",
		SoftwareStackMustHave: []string{},
	}
	if req.TargetRole != nil {
		p.TargetRole = *req.TargetRole
	}
	if req.Location != nil {
		p.Location = *req.Location
	}
	if req.SourceBoard != nil {
		p.SourceBoard = string(*req.SourceBoard)
	}
	if req.CareerSwitch != nil {
		p.CareerSwitch = *req.CareerSwitch
	}
	if req.CompanyStage != nil {
		p.CompanyStage = string(*req.CompanyStage)
	}
	if req.SeniorityGoal != nil {
		p.SeniorityGoal = string(*req.SeniorityGoal)
	}
	if req.CompensationGoal != nil {
		p.CompensationGoal = string(*req.CompensationGoal)
	}
	if req.SoftwareStackMustHave != nil {
		p.SoftwareStackMustHave = cloneStringSlice(*req.SoftwareStackMustHave)
	}
	if req.CanonicalCvAnalysisUuid != nil {
		pg, err := uuid.ConvertOapiUUIDToPgUUID(*req.CanonicalCvAnalysisUuid)
		if err != nil {
			return p, err
		}
		p.CanonicalCvAnalysisUuid = pg
	}
	return p, nil
}

func profileToUpdateByUserParams(existing *sqlc.BridgrJobSearchProfile, req types.UpsertBridgrJobSearchProfileRequest) (sqlc.UpdateJobSearchProfileByUserIDParams, error) {
	p := sqlc.UpdateJobSearchProfileByUserIDParams{
		UserID:                  existing.UserID,
		TargetRole:              existing.TargetRole,
		Location:                existing.Location,
		SourceBoard:             existing.SourceBoard,
		CareerSwitch:            existing.CareerSwitch,
		CompanyStage:            existing.CompanyStage,
		SeniorityGoal:           existing.SeniorityGoal,
		CompensationGoal:        existing.CompensationGoal,
		SoftwareStackMustHave:   cloneStringSlice(existing.SoftwareStackMustHave),
		CanonicalCvAnalysisUuid: existing.CanonicalCvAnalysisUuid,
	}
	if req.TargetRole != nil {
		p.TargetRole = *req.TargetRole
	}
	if req.Location != nil {
		p.Location = *req.Location
	}
	if req.SourceBoard != nil {
		p.SourceBoard = string(*req.SourceBoard)
	}
	if req.CareerSwitch != nil {
		p.CareerSwitch = *req.CareerSwitch
	}
	if req.CompanyStage != nil {
		p.CompanyStage = string(*req.CompanyStage)
	}
	if req.SeniorityGoal != nil {
		p.SeniorityGoal = string(*req.SeniorityGoal)
	}
	if req.CompensationGoal != nil {
		p.CompensationGoal = string(*req.CompensationGoal)
	}
	if req.SoftwareStackMustHave != nil {
		p.SoftwareStackMustHave = cloneStringSlice(*req.SoftwareStackMustHave)
	}
	if req.CanonicalCvAnalysisUuid != nil {
		pg, err := uuid.ConvertOapiUUIDToPgUUID(*req.CanonicalCvAnalysisUuid)
		if err != nil {
			return p, err
		}
		p.CanonicalCvAnalysisUuid = pg
	}
	return p, nil
}

func profileToUpdateByUserIDAndUUIDParams(existing *sqlc.BridgrJobSearchProfile, req types.UpsertBridgrJobSearchProfileRequest) (sqlc.UpdateJobSearchProfileByUserIDAndUUIDParams, error) {
	p := sqlc.UpdateJobSearchProfileByUserIDAndUUIDParams{
		UserID:                  existing.UserID,
		Uuid:                    existing.Uuid,
		TargetRole:              existing.TargetRole,
		Location:                existing.Location,
		SourceBoard:             existing.SourceBoard,
		CareerSwitch:            existing.CareerSwitch,
		CompanyStage:            existing.CompanyStage,
		SeniorityGoal:           existing.SeniorityGoal,
		CompensationGoal:        existing.CompensationGoal,
		SoftwareStackMustHave:   cloneStringSlice(existing.SoftwareStackMustHave),
		CanonicalCvAnalysisUuid: existing.CanonicalCvAnalysisUuid,
	}
	if req.TargetRole != nil {
		p.TargetRole = *req.TargetRole
	}
	if req.Location != nil {
		p.Location = *req.Location
	}
	if req.SourceBoard != nil {
		p.SourceBoard = string(*req.SourceBoard)
	}
	if req.CareerSwitch != nil {
		p.CareerSwitch = *req.CareerSwitch
	}
	if req.CompanyStage != nil {
		p.CompanyStage = string(*req.CompanyStage)
	}
	if req.SeniorityGoal != nil {
		p.SeniorityGoal = string(*req.SeniorityGoal)
	}
	if req.CompensationGoal != nil {
		p.CompensationGoal = string(*req.CompensationGoal)
	}
	if req.SoftwareStackMustHave != nil {
		p.SoftwareStackMustHave = cloneStringSlice(*req.SoftwareStackMustHave)
	}
	if req.CanonicalCvAnalysisUuid != nil {
		pg, err := uuid.ConvertOapiUUIDToPgUUID(*req.CanonicalCvAnalysisUuid)
		if err != nil {
			return p, err
		}
		p.CanonicalCvAnalysisUuid = pg
	}
	return p, nil
}

func cloneStringSlice(s []string) []string {
	if len(s) == 0 {
		return []string{}
	}
	out := make([]string, len(s))
	copy(out, s)
	return out
}

// jobSearchProfileNaturalKey is the per-user uniqueness key for job search profiles.
type jobSearchProfileNaturalKey struct {
	TargetRole              string
	Location                string
	SourceBoard             string
	CareerSwitch            bool
	CompanyStage            string
	SeniorityGoal           string
	CompensationGoal        string
	SoftwareStackMustHave   []string
	CanonicalCvAnalysisUuid pgtype.UUID
}

func naturalKeyFromRow(row *sqlc.BridgrJobSearchProfile) jobSearchProfileNaturalKey {
	return jobSearchProfileNaturalKey{
		TargetRole:              row.TargetRole,
		Location:                row.Location,
		SourceBoard:             row.SourceBoard,
		CareerSwitch:            row.CareerSwitch,
		CompanyStage:            row.CompanyStage,
		SeniorityGoal:           row.SeniorityGoal,
		CompensationGoal:        row.CompensationGoal,
		SoftwareStackMustHave:   cloneStringSlice(row.SoftwareStackMustHave),
		CanonicalCvAnalysisUuid: row.CanonicalCvAnalysisUuid,
	}
}

func naturalKeyFromCreate(p sqlc.CreateJobSearchProfileParams) jobSearchProfileNaturalKey {
	return jobSearchProfileNaturalKey{
		TargetRole:              p.TargetRole,
		Location:                p.Location,
		SourceBoard:             p.SourceBoard,
		CareerSwitch:            p.CareerSwitch,
		CompanyStage:            p.CompanyStage,
		SeniorityGoal:           p.SeniorityGoal,
		CompensationGoal:        p.CompensationGoal,
		SoftwareStackMustHave:   cloneStringSlice(p.SoftwareStackMustHave),
		CanonicalCvAnalysisUuid: p.CanonicalCvAnalysisUuid,
	}
}

func naturalKeyFromUpdate(p sqlc.UpdateJobSearchProfileByUserIDAndUUIDParams) jobSearchProfileNaturalKey {
	return jobSearchProfileNaturalKey{
		TargetRole:              p.TargetRole,
		Location:                p.Location,
		SourceBoard:             p.SourceBoard,
		CareerSwitch:            p.CareerSwitch,
		CompanyStage:            p.CompanyStage,
		SeniorityGoal:           p.SeniorityGoal,
		CompensationGoal:        p.CompensationGoal,
		SoftwareStackMustHave:   cloneStringSlice(p.SoftwareStackMustHave),
		CanonicalCvAnalysisUuid: p.CanonicalCvAnalysisUuid,
	}
}

func (a jobSearchProfileNaturalKey) equals(b jobSearchProfileNaturalKey) bool {
	if strings.TrimSpace(a.TargetRole) != strings.TrimSpace(b.TargetRole) {
		return false
	}
	if strings.TrimSpace(a.Location) != strings.TrimSpace(b.Location) {
		return false
	}
	if strings.TrimSpace(a.SourceBoard) != strings.TrimSpace(b.SourceBoard) {
		return false
	}
	if a.CareerSwitch != b.CareerSwitch {
		return false
	}
	if strings.TrimSpace(a.CompanyStage) != strings.TrimSpace(b.CompanyStage) {
		return false
	}
	if strings.TrimSpace(a.SeniorityGoal) != strings.TrimSpace(b.SeniorityGoal) {
		return false
	}
	if strings.TrimSpace(a.CompensationGoal) != strings.TrimSpace(b.CompensationGoal) {
		return false
	}
	if !stringSlicesEqualSorted(a.SoftwareStackMustHave, b.SoftwareStackMustHave) {
		return false
	}
	x, y := a.CanonicalCvAnalysisUuid, b.CanonicalCvAnalysisUuid
	if !x.Valid && !y.Valid {
		return true
	}
	if !x.Valid || !y.Valid {
		return false
	}
	return bytes.Equal(x.Bytes[:], y.Bytes[:])
}

func stringSlicesEqualSorted(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aa := append([]string(nil), a...)
	bb := append([]string(nil), b...)
	sort.Strings(aa)
	sort.Strings(bb)
	for i := range aa {
		if aa[i] != bb[i] {
			return false
		}
	}
	return true
}
