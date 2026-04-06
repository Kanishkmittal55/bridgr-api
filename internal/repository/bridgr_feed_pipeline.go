package repository

import (
	"context"
	"fmt"

	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// Feed pipeline + supported boards: wrappers around sqlc.Querier (harvest schedules, enrichments, scores, feed items, job boards).

// =============================================================================
// Job harvest schedules
// =============================================================================

// CreateJobHarvestSchedule inserts a per-user schedule for a job_search_profile.
func (r *Repo) CreateJobHarvestSchedule(ctx context.Context, querier sqlc.Querier, params sqlc.CreateJobHarvestScheduleParams) (*sqlc.BridgrJobHarvestSchedule, error) {
	row, err := querier.CreateJobHarvestSchedule(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("CreateJobHarvestSchedule: %w", err)
	}
	return &row, nil
}

// GetJobHarvestScheduleByUUID loads a harvest schedule by UUID.
func (r *Repo) GetJobHarvestScheduleByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) (*sqlc.BridgrJobHarvestSchedule, error) {
	row, err := querier.GetJobHarvestScheduleByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("GetJobHarvestScheduleByUUID: %w", err)
	}
	return &row, nil
}

// GetJobHarvestScheduleByID loads a harvest schedule by BIGSERIAL id.
func (r *Repo) GetJobHarvestScheduleByID(ctx context.Context, querier sqlc.Querier, id int64) (*sqlc.BridgrJobHarvestSchedule, error) {
	row, err := querier.GetJobHarvestScheduleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetJobHarvestScheduleByID: %w", err)
	}
	return &row, nil
}

// GetJobHarvestScheduleByUserAndProfileUUID returns the schedule for one user + profile pair.
func (r *Repo) GetJobHarvestScheduleByUserAndProfileUUID(ctx context.Context, querier sqlc.Querier, params sqlc.GetJobHarvestScheduleByUserAndProfileUUIDParams) (*sqlc.BridgrJobHarvestSchedule, error) {
	row, err := querier.GetJobHarvestScheduleByUserAndProfileUUID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("GetJobHarvestScheduleByUserAndProfileUUID: %w", err)
	}
	return &row, nil
}

// ListJobHarvestSchedulesByUser lists all harvest schedules for a user.
func (r *Repo) ListJobHarvestSchedulesByUser(ctx context.Context, querier sqlc.Querier, userID int32) ([]sqlc.BridgrJobHarvestSchedule, error) {
	rows, err := querier.ListJobHarvestSchedulesByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ListJobHarvestSchedulesByUser: %w", err)
	}
	return rows, nil
}

// ListJobHarvestSchedulesDue returns enabled schedules with next_run_at <= now() (oldest due first). Used by the harvest scheduler.
func (r *Repo) ListJobHarvestSchedulesDue(ctx context.Context, querier sqlc.Querier) ([]sqlc.BridgrJobHarvestSchedule, error) {
	rows, err := querier.ListJobHarvestSchedulesDue(ctx)
	if err != nil {
		return nil, fmt.Errorf("ListJobHarvestSchedulesDue: %w", err)
	}
	return rows, nil
}

// UpdateJobHarvestScheduleByID updates cadence, rotation, and run timestamps.
func (r *Repo) UpdateJobHarvestScheduleByID(ctx context.Context, querier sqlc.Querier, params sqlc.UpdateJobHarvestScheduleByIDParams) (*sqlc.BridgrJobHarvestSchedule, error) {
	row, err := querier.UpdateJobHarvestScheduleByID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateJobHarvestScheduleByID: %w", err)
	}
	return &row, nil
}

// DeleteJobHarvestScheduleByUUID deletes a harvest schedule by UUID.
func (r *Repo) DeleteJobHarvestScheduleByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) error {
	if err := querier.DeleteJobHarvestScheduleByUUID(ctx, uuid); err != nil {
		return fmt.Errorf("DeleteJobHarvestScheduleByUUID: %w", err)
	}
	return nil
}

// DeleteJobHarvestScheduleByID deletes a harvest schedule by id.
func (r *Repo) DeleteJobHarvestScheduleByID(ctx context.Context, querier sqlc.Querier, id int64) error {
	if err := querier.DeleteJobHarvestScheduleByID(ctx, id); err != nil {
		return fmt.Errorf("DeleteJobHarvestScheduleByID: %w", err)
	}
	return nil
}

// =============================================================================
// Job enrichments
// =============================================================================

// CreateJobEnrichment inserts structured JD extraction for a (user, candidate) pair.
func (r *Repo) CreateJobEnrichment(ctx context.Context, querier sqlc.Querier, params sqlc.CreateJobEnrichmentParams) (*sqlc.BridgrJobEnrichment, error) {
	row, err := querier.CreateJobEnrichment(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("CreateJobEnrichment: %w", err)
	}
	return &row, nil
}

// UpsertJobEnrichment inserts or updates on (job_candidate_uuid, user_id).
func (r *Repo) UpsertJobEnrichment(ctx context.Context, querier sqlc.Querier, params sqlc.UpsertJobEnrichmentParams) (*sqlc.BridgrJobEnrichment, error) {
	row, err := querier.UpsertJobEnrichment(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpsertJobEnrichment: %w", err)
	}
	return &row, nil
}

// GetJobEnrichmentByUUID loads an enrichment by UUID.
func (r *Repo) GetJobEnrichmentByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) (*sqlc.BridgrJobEnrichment, error) {
	row, err := querier.GetJobEnrichmentByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("GetJobEnrichmentByUUID: %w", err)
	}
	return &row, nil
}

// GetJobEnrichmentByID loads an enrichment by BIGSERIAL id.
func (r *Repo) GetJobEnrichmentByID(ctx context.Context, querier sqlc.Querier, id int64) (*sqlc.BridgrJobEnrichment, error) {
	row, err := querier.GetJobEnrichmentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetJobEnrichmentByID: %w", err)
	}
	return &row, nil
}

// GetJobEnrichmentByUserAndCandidateUUID loads enrichment for one user + candidate.
func (r *Repo) GetJobEnrichmentByUserAndCandidateUUID(ctx context.Context, querier sqlc.Querier, params sqlc.GetJobEnrichmentByUserAndCandidateUUIDParams) (*sqlc.BridgrJobEnrichment, error) {
	row, err := querier.GetJobEnrichmentByUserAndCandidateUUID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("GetJobEnrichmentByUserAndCandidateUUID: %w", err)
	}
	return &row, nil
}

// ListJobEnrichmentsByUser paginates enrichments for a user.
func (r *Repo) ListJobEnrichmentsByUser(ctx context.Context, querier sqlc.Querier, params sqlc.ListJobEnrichmentsByUserParams) ([]sqlc.BridgrJobEnrichment, error) {
	rows, err := querier.ListJobEnrichmentsByUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("ListJobEnrichmentsByUser: %w", err)
	}
	return rows, nil
}

// UpdateJobEnrichmentByID updates enrichment fields by id.
func (r *Repo) UpdateJobEnrichmentByID(ctx context.Context, querier sqlc.Querier, params sqlc.UpdateJobEnrichmentByIDParams) (*sqlc.BridgrJobEnrichment, error) {
	row, err := querier.UpdateJobEnrichmentByID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateJobEnrichmentByID: %w", err)
	}
	return &row, nil
}

// DeleteJobEnrichmentByUUID deletes an enrichment by UUID.
func (r *Repo) DeleteJobEnrichmentByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) error {
	if err := querier.DeleteJobEnrichmentByUUID(ctx, uuid); err != nil {
		return fmt.Errorf("DeleteJobEnrichmentByUUID: %w", err)
	}
	return nil
}

// DeleteJobEnrichmentByID deletes an enrichment by id.
func (r *Repo) DeleteJobEnrichmentByID(ctx context.Context, querier sqlc.Querier, id int64) error {
	if err := querier.DeleteJobEnrichmentByID(ctx, id); err != nil {
		return fmt.Errorf("DeleteJobEnrichmentByID: %w", err)
	}
	return nil
}

// =============================================================================
// Job scores
// =============================================================================

// CreateJobScore inserts a relevance score row for a (user, candidate) pair.
func (r *Repo) CreateJobScore(ctx context.Context, querier sqlc.Querier, params sqlc.CreateJobScoreParams) (*sqlc.BridgrJobScore, error) {
	row, err := querier.CreateJobScore(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("CreateJobScore: %w", err)
	}
	return &row, nil
}

// UpsertJobScore inserts or updates on (job_candidate_uuid, user_id).
func (r *Repo) UpsertJobScore(ctx context.Context, querier sqlc.Querier, params sqlc.UpsertJobScoreParams) (*sqlc.BridgrJobScore, error) {
	row, err := querier.UpsertJobScore(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpsertJobScore: %w", err)
	}
	return &row, nil
}

// GetJobScoreByUUID loads a score by UUID.
func (r *Repo) GetJobScoreByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) (*sqlc.BridgrJobScore, error) {
	row, err := querier.GetJobScoreByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("GetJobScoreByUUID: %w", err)
	}
	return &row, nil
}

// GetJobScoreByID loads a score by BIGSERIAL id.
func (r *Repo) GetJobScoreByID(ctx context.Context, querier sqlc.Querier, id int64) (*sqlc.BridgrJobScore, error) {
	row, err := querier.GetJobScoreByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetJobScoreByID: %w", err)
	}
	return &row, nil
}

// GetJobScoreByUserAndCandidateUUID loads score for one user + candidate.
func (r *Repo) GetJobScoreByUserAndCandidateUUID(ctx context.Context, querier sqlc.Querier, params sqlc.GetJobScoreByUserAndCandidateUUIDParams) (*sqlc.BridgrJobScore, error) {
	row, err := querier.GetJobScoreByUserAndCandidateUUID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("GetJobScoreByUserAndCandidateUUID: %w", err)
	}
	return &row, nil
}

// ListJobScoresByUser paginates scores for a user (ordered by composite_score in SQL).
func (r *Repo) ListJobScoresByUser(ctx context.Context, querier sqlc.Querier, params sqlc.ListJobScoresByUserParams) ([]sqlc.BridgrJobScore, error) {
	rows, err := querier.ListJobScoresByUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("ListJobScoresByUser: %w", err)
	}
	return rows, nil
}

// UpdateJobScoreByID updates score fields by id.
func (r *Repo) UpdateJobScoreByID(ctx context.Context, querier sqlc.Querier, params sqlc.UpdateJobScoreByIDParams) (*sqlc.BridgrJobScore, error) {
	row, err := querier.UpdateJobScoreByID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateJobScoreByID: %w", err)
	}
	return &row, nil
}

// DeleteJobScoreByUUID deletes a score by UUID.
func (r *Repo) DeleteJobScoreByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) error {
	if err := querier.DeleteJobScoreByUUID(ctx, uuid); err != nil {
		return fmt.Errorf("DeleteJobScoreByUUID: %w", err)
	}
	return nil
}

// DeleteJobScoreByID deletes a score by id.
func (r *Repo) DeleteJobScoreByID(ctx context.Context, querier sqlc.Querier, id int64) error {
	if err := querier.DeleteJobScoreByID(ctx, id); err != nil {
		return fmt.Errorf("DeleteJobScoreByID: %w", err)
	}
	return nil
}

// =============================================================================
// Feed items
// =============================================================================

// CreateFeedItem inserts a denormalized feed row for a (user, candidate) pair.
func (r *Repo) CreateFeedItem(ctx context.Context, querier sqlc.Querier, params sqlc.CreateFeedItemParams) (*sqlc.BridgrFeedItem, error) {
	row, err := querier.CreateFeedItem(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("CreateFeedItem: %w", err)
	}
	return &row, nil
}

// UpsertFeedItem inserts or updates on (job_candidate_uuid, user_id).
func (r *Repo) UpsertFeedItem(ctx context.Context, querier sqlc.Querier, params sqlc.UpsertFeedItemParams) (*sqlc.BridgrFeedItem, error) {
	row, err := querier.UpsertFeedItem(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpsertFeedItem: %w", err)
	}
	return &row, nil
}

// UpsertFeedItemFromDiscovery upserts a feed row for the discovery worker; on conflict it refreshes score and copy but keeps feed_status, seen_at, and first surfaced_at.
func (r *Repo) UpsertFeedItemFromDiscovery(ctx context.Context, querier sqlc.Querier, params sqlc.UpsertFeedItemFromDiscoveryParams) (*sqlc.BridgrFeedItem, error) {
	row, err := querier.UpsertFeedItemFromDiscovery(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpsertFeedItemFromDiscovery: %w", err)
	}
	return &row, nil
}

// GetFeedItemByUUID loads a feed item by UUID.
func (r *Repo) GetFeedItemByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) (*sqlc.BridgrFeedItem, error) {
	row, err := querier.GetFeedItemByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("GetFeedItemByUUID: %w", err)
	}
	return &row, nil
}

// GetFeedItemByID loads a feed item by BIGSERIAL id.
func (r *Repo) GetFeedItemByID(ctx context.Context, querier sqlc.Querier, id int64) (*sqlc.BridgrFeedItem, error) {
	row, err := querier.GetFeedItemByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetFeedItemByID: %w", err)
	}
	return &row, nil
}

// GetFeedItemByUserAndCandidateUUID loads feed item for one user + candidate.
func (r *Repo) GetFeedItemByUserAndCandidateUUID(ctx context.Context, querier sqlc.Querier, params sqlc.GetFeedItemByUserAndCandidateUUIDParams) (*sqlc.BridgrFeedItem, error) {
	row, err := querier.GetFeedItemByUserAndCandidateUUID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("GetFeedItemByUserAndCandidateUUID: %w", err)
	}
	return &row, nil
}

// ListFeedItemsByUser paginates feed items for a user (newest surfaced first).
func (r *Repo) ListFeedItemsByUser(ctx context.Context, querier sqlc.Querier, params sqlc.ListFeedItemsByUserParams) ([]sqlc.BridgrFeedItem, error) {
	rows, err := querier.ListFeedItemsByUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("ListFeedItemsByUser: %w", err)
	}
	return rows, nil
}

// UpdateFeedItemByID updates feed fields by id (e.g. status, seen_at).
func (r *Repo) UpdateFeedItemByID(ctx context.Context, querier sqlc.Querier, params sqlc.UpdateFeedItemByIDParams) (*sqlc.BridgrFeedItem, error) {
	row, err := querier.UpdateFeedItemByID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateFeedItemByID: %w", err)
	}
	return &row, nil
}

// DeleteFeedItemByUUID deletes a feed item by UUID.
func (r *Repo) DeleteFeedItemByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) error {
	if err := querier.DeleteFeedItemByUUID(ctx, uuid); err != nil {
		return fmt.Errorf("DeleteFeedItemByUUID: %w", err)
	}
	return nil
}

// DeleteFeedItemByID deletes a feed item by id.
func (r *Repo) DeleteFeedItemByID(ctx context.Context, querier sqlc.Querier, id int64) error {
	if err := querier.DeleteFeedItemByID(ctx, id); err != nil {
		return fmt.Errorf("DeleteFeedItemByID: %w", err)
	}
	return nil
}

// =============================================================================
// Supported job boards (global catalog)
// =============================================================================

// CreateSupportedJobBoard inserts a board definition row.
func (r *Repo) CreateSupportedJobBoard(ctx context.Context, querier sqlc.Querier, params sqlc.CreateSupportedJobBoardParams) (*sqlc.BridgrSupportedJobBoard, error) {
	row, err := querier.CreateSupportedJobBoard(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("CreateSupportedJobBoard: %w", err)
	}
	return &row, nil
}

// GetSupportedJobBoardByUUID loads a board by UUID.
func (r *Repo) GetSupportedJobBoardByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) (*sqlc.BridgrSupportedJobBoard, error) {
	row, err := querier.GetSupportedJobBoardByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("GetSupportedJobBoardByUUID: %w", err)
	}
	return &row, nil
}

// GetSupportedJobBoardByID loads a board by BIGSERIAL id.
func (r *Repo) GetSupportedJobBoardByID(ctx context.Context, querier sqlc.Querier, id int64) (*sqlc.BridgrSupportedJobBoard, error) {
	row, err := querier.GetSupportedJobBoardByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetSupportedJobBoardByID: %w", err)
	}
	return &row, nil
}

// GetSupportedJobBoardByBoardID loads a board by stable board_id slug.
func (r *Repo) GetSupportedJobBoardByBoardID(ctx context.Context, querier sqlc.Querier, boardID string) (*sqlc.BridgrSupportedJobBoard, error) {
	row, err := querier.GetSupportedJobBoardByBoardID(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("GetSupportedJobBoardByBoardID: %w", err)
	}
	return &row, nil
}

// ListSupportedJobBoards returns all boards (admin / diagnostics).
func (r *Repo) ListSupportedJobBoards(ctx context.Context, querier sqlc.Querier) ([]sqlc.BridgrSupportedJobBoard, error) {
	rows, err := querier.ListSupportedJobBoards(ctx)
	if err != nil {
		return nil, fmt.Errorf("ListSupportedJobBoards: %w", err)
	}
	return rows, nil
}

// ListActiveSupportedJobBoards returns boards where is_active is true (user-facing pickers).
func (r *Repo) ListActiveSupportedJobBoards(ctx context.Context, querier sqlc.Querier) ([]sqlc.BridgrSupportedJobBoard, error) {
	rows, err := querier.ListActiveSupportedJobBoards(ctx)
	if err != nil {
		return nil, fmt.Errorf("ListActiveSupportedJobBoards: %w", err)
	}
	return rows, nil
}

// UpdateSupportedJobBoardByID updates a board row by id.
func (r *Repo) UpdateSupportedJobBoardByID(ctx context.Context, querier sqlc.Querier, params sqlc.UpdateSupportedJobBoardByIDParams) (*sqlc.BridgrSupportedJobBoard, error) {
	row, err := querier.UpdateSupportedJobBoardByID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateSupportedJobBoardByID: %w", err)
	}
	return &row, nil
}

// DeleteSupportedJobBoardByUUID deletes a board by UUID.
func (r *Repo) DeleteSupportedJobBoardByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) error {
	if err := querier.DeleteSupportedJobBoardByUUID(ctx, uuid); err != nil {
		return fmt.Errorf("DeleteSupportedJobBoardByUUID: %w", err)
	}
	return nil
}

// DeleteSupportedJobBoardByID deletes a board by id.
func (r *Repo) DeleteSupportedJobBoardByID(ctx context.Context, querier sqlc.Querier, id int64) error {
	if err := querier.DeleteSupportedJobBoardByID(ctx, id); err != nil {
		return fmt.Errorf("DeleteSupportedJobBoardByID: %w", err)
	}
	return nil
}
