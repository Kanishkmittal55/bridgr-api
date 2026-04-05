package bridgr

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hassleskip/bridgr-api/internal/repository/sqlc"
	"github.com/hassleskip/bridgr-api/internal/uuid"
	types "github.com/hassleskip/bridgr-api/pkg/types"
	"github.com/jackc/pgx/v5/pgtype"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func pgTextStringPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	s := t.String
	return &s
}

func oapiUUIDPtrFromPg(u pgtype.UUID) (*openapi_types.UUID, error) {
	if !u.Valid {
		return nil, nil
	}
	o, err := uuid.ConvertPgUUIDToOapiUUID(u)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func oapiUUIDFromPg(u pgtype.UUID) (openapi_types.UUID, error) {
	return uuid.ConvertPgUUIDToOapiUUID(u)
}

func pgTimestampToTime(ts pgtype.Timestamp) time.Time {
	if ts.Valid {
		return ts.Time
	}
	return time.Time{}
}

func jsonBytesToMap(b []byte) *map[string]interface{} {
	if len(b) == 0 {
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil || len(m) == 0 {
		return nil
	}
	return &m
}

func mapToJSONBytes(m *map[string]interface{}) ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(*m)
}

func numericToFloat32Ptr(n pgtype.Numeric) *float32 {
	fv, err := n.Float64Value()
	if err != nil || !fv.Valid {
		return nil
	}
	v := float32(fv.Float64)
	return &v
}

func float32ToNumericPtr(f *float32) pgtype.Numeric {
	if f == nil {
		return pgtype.Numeric{Valid: false}
	}
	var n pgtype.Numeric
	_ = n.Scan(fmt.Sprintf("%g", *f))
	return n
}

func countToInterfacePtr(v interface{}) *interface{} {
	if v == nil {
		return nil
	}
	x := interface{}(v)
	return &x
}

func int64PtrFromCountInterface(v interface{}) *int64 {
	if v == nil {
		return nil
	}
	switch x := v.(type) {
	case int64:
		return &x
	case int32:
		i := int64(x)
		return &i
	case uint64:
		i := int64(x)
		return &i
	default:
		return nil
	}
}

func analysisFromRow(row *sqlc.HskipUsersBridgrSkillGapAnalysis) (types.BridgrSkillGapAnalysis, error) {
	out := types.BridgrSkillGapAnalysis{
		Id:        row.ID,
		UserId:    row.UserID,
		Status:    types.BridgrSkillGapAnalysisStatus(row.Status),
		CreatedAt: pgTimestampToTime(row.CreatedAt),
		UpdatedAt: pgTimestampToTime(row.UpdatedAt),
	}
	var err error
	out.Uuid, err = oapiUUIDFromPg(row.Uuid)
	if err != nil {
		return types.BridgrSkillGapAnalysis{}, err
	}
	out.FounderPersonaUuid, err = oapiUUIDPtrFromPg(row.FounderPersonaUuid)
	if err != nil {
		return types.BridgrSkillGapAnalysis{}, err
	}
	out.PursuitUuid, err = oapiUUIDPtrFromPg(row.PursuitUuid)
	if err != nil {
		return types.BridgrSkillGapAnalysis{}, err
	}
	out.Title = pgTextStringPtr(row.Title)
	out.CvAssetUri = pgTextStringPtr(row.CvAssetUri)
	out.JdAssetUri = pgTextStringPtr(row.JdAssetUri)
	out.CvFingerprint = pgTextStringPtr(row.CvFingerprint)
	out.JdFingerprint = pgTextStringPtr(row.JdFingerprint)
	out.LlmModel = pgTextStringPtr(row.LlmModel)
	out.PromptVersion = pgTextStringPtr(row.PromptVersion)
	out.ExtractionPayload = jsonBytesToMap(row.ExtractionPayload)
	out.GapSummary = jsonBytesToMap(row.GapSummary)
	out.MermaidDiagram = pgTextStringPtr(row.MermaidDiagram)
	out.ErrorCode = pgTextStringPtr(row.ErrorCode)
	out.ErrorDetail = pgTextStringPtr(row.ErrorDetail)
	out.SqsMessageId = pgTextStringPtr(row.SqsMessageID)
	return out, nil
}

func graphFromRow(row *sqlc.HskipUsersBridgrSkillGapGraph) (types.BridgrSkillGapGraph, error) {
	out := types.BridgrSkillGapGraph{
		Kind:      types.BridgrSkillGapGraphKind(row.Kind),
		CreatedAt: pgTimestampToTime(row.CreatedAt),
		UpdatedAt: pgTimestampToTime(row.UpdatedAt),
	}
	var err error
	out.Uuid, err = oapiUUIDFromPg(row.Uuid)
	if err != nil {
		return types.BridgrSkillGapGraph{}, err
	}
	out.AnalysisUuid, err = oapiUUIDFromPg(row.AnalysisUuid)
	if err != nil {
		return types.BridgrSkillGapGraph{}, err
	}
	out.Metadata = jsonBytesToMap(row.Metadata)
	out.Id = row.ID
	return out, nil
}

func nodeFromRow(row *sqlc.HskipUsersBridgrSkillGapNode) (types.BridgrSkillGapNode, error) {
	out := types.BridgrSkillGapNode{
		NodeKey:     row.NodeKey,
		DisplayName: row.DisplayName,
		CreatedAt:   pgTimestampToTime(row.CreatedAt),
		UpdatedAt:   pgTimestampToTime(row.UpdatedAt),
		Id:          row.ID,
	}
	var err error
	out.Uuid, err = oapiUUIDFromPg(row.Uuid)
	if err != nil {
		return types.BridgrSkillGapNode{}, err
	}
	out.GraphUuid, err = oapiUUIDFromPg(row.GraphUuid)
	if err != nil {
		return types.BridgrSkillGapNode{}, err
	}
	out.Description = pgTextStringPtr(row.Description)
	out.ProficiencyHint = pgTextStringPtr(row.ProficiencyHint)
	out.Source = pgTextStringPtr(row.Source)
	out.Evidence = jsonBytesToMap(row.Evidence)
	out.Metadata = jsonBytesToMap(row.Metadata)
	if row.PositionX.Valid {
		out.PositionX = &row.PositionX.Int32
	}
	if row.PositionY.Valid {
		out.PositionY = &row.PositionY.Int32
	}
	return out, nil
}

func edgeFromRow(row *sqlc.HskipUsersBridgrSkillGapEdge) (types.BridgrSkillGapEdge, error) {
	out := types.BridgrSkillGapEdge{
		Relation:  row.Relation,
		CreatedAt: pgTimestampToTime(row.CreatedAt),
		UpdatedAt: pgTimestampToTime(row.UpdatedAt),
		Id:        row.ID,
	}
	var err error
	out.Uuid, err = oapiUUIDFromPg(row.Uuid)
	if err != nil {
		return types.BridgrSkillGapEdge{}, err
	}
	out.GraphUuid, err = oapiUUIDFromPg(row.GraphUuid)
	if err != nil {
		return types.BridgrSkillGapEdge{}, err
	}
	out.FromNodeUuid, err = oapiUUIDFromPg(row.FromNodeUuid)
	if err != nil {
		return types.BridgrSkillGapEdge{}, err
	}
	out.ToNodeUuid, err = oapiUUIDFromPg(row.ToNodeUuid)
	if err != nil {
		return types.BridgrSkillGapEdge{}, err
	}
	out.Metadata = jsonBytesToMap(row.Metadata)
	out.Weight = numericToFloat32Ptr(row.Weight)
	return out, nil
}

func learningPathFromRow(row *sqlc.HskipUsersBridgrSkillGapLearningPath) (types.BridgrSkillGapLearningPath, error) {
	out := types.BridgrSkillGapLearningPath{
		PathVersion: row.PathVersion,
		CreatedAt:   pgTimestampToTime(row.CreatedAt),
		UpdatedAt:   pgTimestampToTime(row.UpdatedAt),
		Id:          row.ID,
	}
	var err error
	out.Uuid, err = oapiUUIDFromPg(row.Uuid)
	if err != nil {
		return types.BridgrSkillGapLearningPath{}, err
	}
	out.AnalysisUuid, err = oapiUUIDFromPg(row.AnalysisUuid)
	if err != nil {
		return types.BridgrSkillGapLearningPath{}, err
	}
	out.Algorithm = pgTextStringPtr(row.Algorithm)
	out.Title = pgTextStringPtr(row.Title)
	out.PathMetadata = jsonBytesToMap(row.PathMetadata)
	return out, nil
}

func pathStepFromRow(row *sqlc.HskipUsersBridgrSkillGapPathStep) (types.BridgrSkillGapPathStep, error) {
	out := types.BridgrSkillGapPathStep{
		StepIndex: row.StepIndex,
		Title:     row.Title,
		CreatedAt: pgTimestampToTime(row.CreatedAt),
		UpdatedAt: pgTimestampToTime(row.UpdatedAt),
		Id:        row.ID,
	}
	var err error
	out.Uuid, err = oapiUUIDFromPg(row.Uuid)
	if err != nil {
		return types.BridgrSkillGapPathStep{}, err
	}
	out.PathUuid, err = oapiUUIDFromPg(row.PathUuid)
	if err != nil {
		return types.BridgrSkillGapPathStep{}, err
	}
	out.Rationale = pgTextStringPtr(row.Rationale)
	out.EstimatedHours = numericToFloat32Ptr(row.EstimatedHours)
	out.ResourceUri = pgTextStringPtr(row.ResourceUri)
	out.ResourceKind = pgTextStringPtr(row.ResourceKind)
	fli, err := oapiUUIDPtrFromPg(row.FounderLearningItemUuid)
	if err != nil {
		return types.BridgrSkillGapPathStep{}, err
	}
	out.FounderLearningItemUuid = fli
	cl, err := oapiUUIDPtrFromPg(row.CourseLessonUuid)
	if err != nil {
		return types.BridgrSkillGapPathStep{}, err
	}
	out.CourseLessonUuid = cl
	out.Metadata = jsonBytesToMap(row.Metadata)
	if len(row.LinkedNodeKeys) > 0 {
		var keys []string
		if err := json.Unmarshal(row.LinkedNodeKeys, &keys); err == nil && len(keys) > 0 {
			out.LinkedNodeKeys = &keys
		}
	}
	return out, nil
}

func coverageUserRowFromSQL(row sqlc.GetSkillGapCoverageByUserRow) (types.BridgrSkillGapUserCoverageRow, error) {
	s := types.BridgrSkillGapAnalysisStatus(row.Status)
	out := types.BridgrSkillGapUserCoverageRow{
		AnalysisId:      &row.AnalysisID,
		UserId:          &row.UserID,
		Status:          &s,
		Title:           pgTextStringPtr(row.Title),
		CandidateSkills: int64PtrFromCountInterface(row.CandidateSkills),
		RequiredSkills:  int64PtrFromCountInterface(row.RequiredSkills),
		MatchedSkills:   int64PtrFromCountInterface(row.MatchedSkills),
	}
	ca := pgTimestampToTime(row.CreatedAt)
	out.CreatedAt = &ca
	au, err := oapiUUIDPtrFromPg(row.AnalysisUuid)
	if err != nil {
		return types.BridgrSkillGapUserCoverageRow{}, err
	}
	out.AnalysisUuid = au
	out.CoveragePct = numericToFloat32Ptr(row.CoveragePct)
	return out, nil
}
