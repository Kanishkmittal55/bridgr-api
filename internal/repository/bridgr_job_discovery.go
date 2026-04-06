package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// Job discovery tables: wrappers around sqlc.Querier (see sqlc/querier.go).

// =============================================================================
// Job search profiles
// =============================================================================

// CreateJobSearchProfile inserts a per-user search preference row.
func (r *Repo) CreateJobSearchProfile(ctx context.Context, querier sqlc.Querier, params sqlc.CreateJobSearchProfileParams) (*sqlc.BridgrJobSearchProfile, error) {
	row, err := querier.CreateJobSearchProfile(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("CreateJobSearchProfile: %w", err)
	}
	return &row, nil
}

// GetJobSearchProfileByUUID loads a profile by primary UUID.
func (r *Repo) GetJobSearchProfileByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) (*sqlc.BridgrJobSearchProfile, error) {
	row, err := querier.GetJobSearchProfileByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("GetJobSearchProfileByUUID: %w", err)
	}
	return &row, nil
}

// GetJobSearchProfileByID loads a profile by BIGSERIAL id.
func (r *Repo) GetJobSearchProfileByID(ctx context.Context, querier sqlc.Querier, id int64) (*sqlc.BridgrJobSearchProfile, error) {
	row, err := querier.GetJobSearchProfileByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetJobSearchProfileByID: %w", err)
	}
	return &row, nil
}

// UpdateJobSearchProfile updates mutable columns by profile id.
func (r *Repo) UpdateJobSearchProfile(ctx context.Context, querier sqlc.Querier, params sqlc.UpdateJobSearchProfileParams) (*sqlc.BridgrJobSearchProfile, error) {
	row, err := querier.UpdateJobSearchProfile(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateJobSearchProfile: %w", err)
	}
	return &row, nil
}

// UpdateJobSearchProfileByUserID updates mutable columns by user id.
func (r *Repo) UpdateJobSearchProfileByUserID(ctx context.Context, querier sqlc.Querier, params sqlc.UpdateJobSearchProfileByUserIDParams) (*sqlc.BridgrJobSearchProfile, error) {
	row, err := querier.UpdateJobSearchProfileByUserID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateJobSearchProfileByUserID: %w", err)
	}
	return &row, nil
}

// DeleteJobSearchProfileByUserID deletes the profile for a user.
func (r *Repo) DeleteJobSearchProfileByUserID(ctx context.Context, querier sqlc.Querier, userID int32) error {
	if err := querier.DeleteJobSearchProfileByUserID(ctx, userID); err != nil {
		return fmt.Errorf("DeleteJobSearchProfileByUserID: %w", err)
	}
	return nil
}

// DeleteJobSearchProfileByID deletes a profile by id.
func (r *Repo) DeleteJobSearchProfileByID(ctx context.Context, querier sqlc.Querier, id int64) error {
	if err := querier.DeleteJobSearchProfileByID(ctx, id); err != nil {
		return fmt.Errorf("DeleteJobSearchProfileByID: %w", err)
	}
	return nil
}

// ListJobSearchProfilesByUserID returns profiles for a user (ordered by id).
func (r *Repo) ListJobSearchProfilesByUserID(ctx context.Context, querier sqlc.Querier, userID int32) ([]sqlc.BridgrJobSearchProfile, error) {
	rows, err := querier.ListJobSearchProfilesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ListJobSearchProfilesByUserID: %w", err)
	}
	return rows, nil
}

// GetJobSearchProfileByUserIDAndUUID loads a profile scoped to user_id + uuid.
func (r *Repo) GetJobSearchProfileByUserIDAndUUID(ctx context.Context, querier sqlc.Querier, params sqlc.GetJobSearchProfileByUserIDAndUUIDParams) (*sqlc.BridgrJobSearchProfile, error) {
	row, err := querier.GetJobSearchProfileByUserIDAndUUID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("GetJobSearchProfileByUserIDAndUUID: %w", err)
	}
	return &row, nil
}

// UpdateJobSearchProfileByUserIDAndUUID updates mutable columns for a user's profile row.
func (r *Repo) UpdateJobSearchProfileByUserIDAndUUID(ctx context.Context, querier sqlc.Querier, params sqlc.UpdateJobSearchProfileByUserIDAndUUIDParams) (*sqlc.BridgrJobSearchProfile, error) {
	row, err := querier.UpdateJobSearchProfileByUserIDAndUUID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateJobSearchProfileByUserIDAndUUID: %w", err)
	}
	return &row, nil
}

// DeleteJobSearchProfileByUserIDAndUUID deletes a profile scoped to user_id + uuid.
func (r *Repo) DeleteJobSearchProfileByUserIDAndUUID(ctx context.Context, querier sqlc.Querier, params sqlc.DeleteJobSearchProfileByUserIDAndUUIDParams) error {
	if err := querier.DeleteJobSearchProfileByUserIDAndUUID(ctx, params); err != nil {
		return fmt.Errorf("DeleteJobSearchProfileByUserIDAndUUID: %w", err)
	}
	return nil
}

// =============================================================================
// Job search discovery runs
// =============================================================================

// CreateJobSearchDiscoveryRun inserts a discovery run audit row.
func (r *Repo) CreateJobSearchDiscoveryRun(ctx context.Context, querier sqlc.Querier, params sqlc.CreateJobSearchDiscoveryRunParams) (*sqlc.BridgrJobSearchDiscoveryRun, error) {
	row, err := querier.CreateJobSearchDiscoveryRun(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("CreateJobSearchDiscoveryRun: %w", err)
	}
	return &row, nil
}

// GetJobSearchDiscoveryRunByUUID loads a run by UUID.
func (r *Repo) GetJobSearchDiscoveryRunByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) (*sqlc.BridgrJobSearchDiscoveryRun, error) {
	row, err := querier.GetJobSearchDiscoveryRunByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("GetJobSearchDiscoveryRunByUUID: %w", err)
	}
	return &row, nil
}

// GetJobSearchDiscoveryRunByID loads a run by BIGSERIAL id.
func (r *Repo) GetJobSearchDiscoveryRunByID(ctx context.Context, querier sqlc.Querier, id int64) (*sqlc.BridgrJobSearchDiscoveryRun, error) {
	row, err := querier.GetJobSearchDiscoveryRunByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetJobSearchDiscoveryRunByID: %w", err)
	}
	return &row, nil
}

// ListJobSearchDiscoveryRunsByUser paginates runs for a user (newest first).
func (r *Repo) ListJobSearchDiscoveryRunsByUser(ctx context.Context, querier sqlc.Querier, params sqlc.ListJobSearchDiscoveryRunsByUserParams) ([]sqlc.BridgrJobSearchDiscoveryRun, error) {
	rows, err := querier.ListJobSearchDiscoveryRunsByUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("ListJobSearchDiscoveryRunsByUser: %w", err)
	}
	return rows, nil
}

// CountJobSearchDiscoveryRunsByUserSince counts runs created at or after the given time (inclusive).
func (r *Repo) CountJobSearchDiscoveryRunsByUserSince(ctx context.Context, querier sqlc.Querier, params sqlc.CountJobSearchDiscoveryRunsByUserSinceParams) (int64, error) {
	n, err := querier.CountJobSearchDiscoveryRunsByUserSince(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("CountJobSearchDiscoveryRunsByUserSince: %w", err)
	}
	return n, nil
}

// SetJobSearchDiscoveryRunStatus updates status only (does not touch other columns).
func (r *Repo) SetJobSearchDiscoveryRunStatus(ctx context.Context, querier sqlc.Querier, params sqlc.SetJobSearchDiscoveryRunStatusParams) (*sqlc.BridgrJobSearchDiscoveryRun, error) {
	row, err := querier.SetJobSearchDiscoveryRunStatus(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("SetJobSearchDiscoveryRunStatus: %w", err)
	}
	return &row, nil
}

// PatchJobSearchDiscoveryRun sets all mutable execution fields; merge from an existing row for fields you are not changing.
func (r *Repo) PatchJobSearchDiscoveryRun(ctx context.Context, querier sqlc.Querier, params sqlc.PatchJobSearchDiscoveryRunParams) (*sqlc.BridgrJobSearchDiscoveryRun, error) {
	row, err := querier.PatchJobSearchDiscoveryRun(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("PatchJobSearchDiscoveryRun: %w", err)
	}
	return &row, nil
}

// DeleteJobSearchDiscoveryRunByID deletes a run row.
func (r *Repo) DeleteJobSearchDiscoveryRunByID(ctx context.Context, querier sqlc.Querier, id int64) error {
	if err := querier.DeleteJobSearchDiscoveryRunByID(ctx, id); err != nil {
		return fmt.Errorf("DeleteJobSearchDiscoveryRunByID: %w", err)
	}
	return nil
}

// =============================================================================
// Job candidates
// =============================================================================

// InsertJobCandidate inserts a new candidate (fails on user_id+url_hash conflict).
func (r *Repo) InsertJobCandidate(ctx context.Context, querier sqlc.Querier, params sqlc.InsertJobCandidateParams) (*sqlc.BridgrJobCandidate, error) {
	row, err := querier.InsertJobCandidate(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("InsertJobCandidate: %w", err)
	}
	return &row, nil
}

// UpsertJobCandidate inserts or updates on (user_id, url_hash).
func (r *Repo) UpsertJobCandidate(ctx context.Context, querier sqlc.Querier, params sqlc.UpsertJobCandidateParams) (*sqlc.BridgrJobCandidate, error) {
	row, err := querier.UpsertJobCandidate(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpsertJobCandidate: %w", err)
	}
	return &row, nil
}

// GetJobCandidateByUUID loads a candidate by UUID.
func (r *Repo) GetJobCandidateByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) (*sqlc.BridgrJobCandidate, error) {
	row, err := querier.GetJobCandidateByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("GetJobCandidateByUUID: %w", err)
	}
	return &row, nil
}

// GetJobCandidateByID loads a candidate by BIGSERIAL id.
func (r *Repo) GetJobCandidateByID(ctx context.Context, querier sqlc.Querier, id int64) (*sqlc.BridgrJobCandidate, error) {
	row, err := querier.GetJobCandidateByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetJobCandidateByID: %w", err)
	}
	return &row, nil
}

// GetJobCandidateByUserAndURLHash dedupes lookup.
func (r *Repo) GetJobCandidateByUserAndURLHash(ctx context.Context, querier sqlc.Querier, params sqlc.GetJobCandidateByUserAndURLHashParams) (*sqlc.BridgrJobCandidate, error) {
	row, err := querier.GetJobCandidateByUserAndURLHash(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("GetJobCandidateByUserAndURLHash: %w", err)
	}
	return &row, nil
}

// ListJobCandidatesByUser paginates candidates for a user.
func (r *Repo) ListJobCandidatesByUser(ctx context.Context, querier sqlc.Querier, params sqlc.ListJobCandidatesByUserParams) ([]sqlc.BridgrJobCandidate, error) {
	rows, err := querier.ListJobCandidatesByUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("ListJobCandidatesByUser: %w", err)
	}
	return rows, nil
}

// ListJobCandidatesByUserFiltered paginates with optional discovery_run / status / board filters.
func (r *Repo) ListJobCandidatesByUserFiltered(ctx context.Context, querier sqlc.Querier, params sqlc.ListJobCandidatesByUserFilteredParams) ([]sqlc.BridgrJobCandidate, error) {
	rows, err := querier.ListJobCandidatesByUserFiltered(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("ListJobCandidatesByUserFiltered: %w", err)
	}
	return rows, nil
}

// ListJobCandidatesByDiscoveryRunUUID lists candidates for one discovery run.
func (r *Repo) ListJobCandidatesByDiscoveryRunUUID(ctx context.Context, querier sqlc.Querier, discoveryRunUuid pgtype.UUID) ([]sqlc.BridgrJobCandidate, error) {
	rows, err := querier.ListJobCandidatesByDiscoveryRunUUID(ctx, discoveryRunUuid)
	if err != nil {
		return nil, fmt.Errorf("ListJobCandidatesByDiscoveryRunUUID: %w", err)
	}
	return rows, nil
}

// UpdateJobCandidateIngestion updates JD body and ingestion fields.
func (r *Repo) UpdateJobCandidateIngestion(ctx context.Context, querier sqlc.Querier, params sqlc.UpdateJobCandidateIngestionParams) (*sqlc.BridgrJobCandidate, error) {
	row, err := querier.UpdateJobCandidateIngestion(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateJobCandidateIngestion: %w", err)
	}
	return &row, nil
}

// UpdateJobCandidateByID replaces mutable columns on a candidate.
func (r *Repo) UpdateJobCandidateByID(ctx context.Context, querier sqlc.Querier, params sqlc.UpdateJobCandidateByIDParams) (*sqlc.BridgrJobCandidate, error) {
	row, err := querier.UpdateJobCandidateByID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateJobCandidateByID: %w", err)
	}
	return &row, nil
}

// DeleteJobCandidateByUUID deletes a candidate by UUID.
func (r *Repo) DeleteJobCandidateByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) error {
	if err := querier.DeleteJobCandidateByUUID(ctx, uuid); err != nil {
		return fmt.Errorf("DeleteJobCandidateByUUID: %w", err)
	}
	return nil
}

// DeleteJobCandidateByID deletes a candidate by id.
func (r *Repo) DeleteJobCandidateByID(ctx context.Context, querier sqlc.Querier, id int64) error {
	if err := querier.DeleteJobCandidateByID(ctx, id); err != nil {
		return fmt.Errorf("DeleteJobCandidateByID: %w", err)
	}
	return nil
}

// =============================================================================
// Job notifications
// =============================================================================

// CreateJobNotification inserts an outbox notification row.
func (r *Repo) CreateJobNotification(ctx context.Context, querier sqlc.Querier, params sqlc.CreateJobNotificationParams) (*sqlc.BridgrJobNotification, error) {
	row, err := querier.CreateJobNotification(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("CreateJobNotification: %w", err)
	}
	return &row, nil
}

// GetJobNotificationByUUID loads a notification by UUID.
func (r *Repo) GetJobNotificationByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) (*sqlc.BridgrJobNotification, error) {
	row, err := querier.GetJobNotificationByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("GetJobNotificationByUUID: %w", err)
	}
	return &row, nil
}

// GetJobNotificationByID loads a notification by BIGSERIAL id.
func (r *Repo) GetJobNotificationByID(ctx context.Context, querier sqlc.Querier, id int64) (*sqlc.BridgrJobNotification, error) {
	row, err := querier.GetJobNotificationByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetJobNotificationByID: %w", err)
	}
	return &row, nil
}

// ListJobNotificationsByUser paginates all notifications for a user.
func (r *Repo) ListJobNotificationsByUser(ctx context.Context, querier sqlc.Querier, params sqlc.ListJobNotificationsByUserParams) ([]sqlc.BridgrJobNotification, error) {
	rows, err := querier.ListJobNotificationsByUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("ListJobNotificationsByUser: %w", err)
	}
	return rows, nil
}

// ListJobNotificationsByUserFiltered paginates with optional since / status filters.
func (r *Repo) ListJobNotificationsByUserFiltered(ctx context.Context, querier sqlc.Querier, params sqlc.ListJobNotificationsByUserFilteredParams) ([]sqlc.BridgrJobNotification, error) {
	rows, err := querier.ListJobNotificationsByUserFiltered(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("ListJobNotificationsByUserFiltered: %w", err)
	}
	return rows, nil
}

// ListJobNotificationsByUserAndStatus filters by status.
func (r *Repo) ListJobNotificationsByUserAndStatus(ctx context.Context, querier sqlc.Querier, params sqlc.ListJobNotificationsByUserAndStatusParams) ([]sqlc.BridgrJobNotification, error) {
	rows, err := querier.ListJobNotificationsByUserAndStatus(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("ListJobNotificationsByUserAndStatus: %w", err)
	}
	return rows, nil
}

// ListPendingJobNotificationsByUser returns pending rows (oldest first).
func (r *Repo) ListPendingJobNotificationsByUser(ctx context.Context, querier sqlc.Querier, params sqlc.ListPendingJobNotificationsByUserParams) ([]sqlc.BridgrJobNotification, error) {
	rows, err := querier.ListPendingJobNotificationsByUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("ListPendingJobNotificationsByUser: %w", err)
	}
	return rows, nil
}

// UpdateJobNotificationStatus updates delivery status fields.
func (r *Repo) UpdateJobNotificationStatus(ctx context.Context, querier sqlc.Querier, params sqlc.UpdateJobNotificationStatusParams) (*sqlc.BridgrJobNotification, error) {
	row, err := querier.UpdateJobNotificationStatus(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateJobNotificationStatus: %w", err)
	}
	return &row, nil
}

// UpdateJobNotificationSeen marks a notification seen.
func (r *Repo) UpdateJobNotificationSeen(ctx context.Context, querier sqlc.Querier, params sqlc.UpdateJobNotificationSeenParams) (*sqlc.BridgrJobNotification, error) {
	row, err := querier.UpdateJobNotificationSeen(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateJobNotificationSeen: %w", err)
	}
	return &row, nil
}

// DeleteJobNotificationByID deletes a notification row.
func (r *Repo) DeleteJobNotificationByID(ctx context.Context, querier sqlc.Querier, id int64) error {
	if err := querier.DeleteJobNotificationByID(ctx, id); err != nil {
		return fmt.Errorf("DeleteJobNotificationByID: %w", err)
	}
	return nil
}

// DeleteJobNotificationByUUID deletes by UUID.
func (r *Repo) DeleteJobNotificationByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) error {
	if err := querier.DeleteJobNotificationByUUID(ctx, uuid); err != nil {
		return fmt.Errorf("DeleteJobNotificationByUUID: %w", err)
	}
	return nil
}

// =============================================================================
// Analysis ↔ job candidate link
// =============================================================================

// CreateAnalysisJobLink ties an analysis UUID to a job candidate UUID.
func (r *Repo) CreateAnalysisJobLink(ctx context.Context, querier sqlc.Querier, params sqlc.CreateAnalysisJobLinkParams) (*sqlc.BridgrAnalysisJobLink, error) {
	row, err := querier.CreateAnalysisJobLink(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("CreateAnalysisJobLink: %w", err)
	}
	return &row, nil
}

// GetAnalysisJobLinkByUUID loads a link row by UUID.
func (r *Repo) GetAnalysisJobLinkByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) (*sqlc.BridgrAnalysisJobLink, error) {
	row, err := querier.GetAnalysisJobLinkByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("GetAnalysisJobLinkByUUID: %w", err)
	}
	return &row, nil
}

// GetAnalysisJobLinkByID loads a link by BIGSERIAL id.
func (r *Repo) GetAnalysisJobLinkByID(ctx context.Context, querier sqlc.Querier, id int64) (*sqlc.BridgrAnalysisJobLink, error) {
	row, err := querier.GetAnalysisJobLinkByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetAnalysisJobLinkByID: %w", err)
	}
	return &row, nil
}

// GetAnalysisJobLinkByPair returns the link for one analysis + candidate pair.
func (r *Repo) GetAnalysisJobLinkByPair(ctx context.Context, querier sqlc.Querier, params sqlc.GetAnalysisJobLinkByPairParams) (*sqlc.BridgrAnalysisJobLink, error) {
	row, err := querier.GetAnalysisJobLinkByPair(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("GetAnalysisJobLinkByPair: %w", err)
	}
	return &row, nil
}

// ListAnalysisJobLinksByUser paginates links for a user.
func (r *Repo) ListAnalysisJobLinksByUser(ctx context.Context, querier sqlc.Querier, params sqlc.ListAnalysisJobLinksByUserParams) ([]sqlc.BridgrAnalysisJobLink, error) {
	rows, err := querier.ListAnalysisJobLinksByUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("ListAnalysisJobLinksByUser: %w", err)
	}
	return rows, nil
}

// ListAnalysisJobLinksByAnalysisUUID lists links for a skill-gap analysis.
func (r *Repo) ListAnalysisJobLinksByAnalysisUUID(ctx context.Context, querier sqlc.Querier, analysisUuid pgtype.UUID) ([]sqlc.BridgrAnalysisJobLink, error) {
	rows, err := querier.ListAnalysisJobLinksByAnalysisUUID(ctx, analysisUuid)
	if err != nil {
		return nil, fmt.Errorf("ListAnalysisJobLinksByAnalysisUUID: %w", err)
	}
	return rows, nil
}

// ListAnalysisJobLinksByCandidateUUID lists links for a job candidate.
func (r *Repo) ListAnalysisJobLinksByCandidateUUID(ctx context.Context, querier sqlc.Querier, jobCandidateUuid pgtype.UUID) ([]sqlc.BridgrAnalysisJobLink, error) {
	rows, err := querier.ListAnalysisJobLinksByCandidateUUID(ctx, jobCandidateUuid)
	if err != nil {
		return nil, fmt.Errorf("ListAnalysisJobLinksByCandidateUUID: %w", err)
	}
	return rows, nil
}

// UpdateAnalysisJobLinkKind updates link_kind.
func (r *Repo) UpdateAnalysisJobLinkKind(ctx context.Context, querier sqlc.Querier, params sqlc.UpdateAnalysisJobLinkKindParams) (*sqlc.BridgrAnalysisJobLink, error) {
	row, err := querier.UpdateAnalysisJobLinkKind(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateAnalysisJobLinkKind: %w", err)
	}
	return &row, nil
}

// DeleteAnalysisJobLinkByUUID deletes a link by UUID.
func (r *Repo) DeleteAnalysisJobLinkByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) error {
	if err := querier.DeleteAnalysisJobLinkByUUID(ctx, uuid); err != nil {
		return fmt.Errorf("DeleteAnalysisJobLinkByUUID: %w", err)
	}
	return nil
}

// DeleteAnalysisJobLinkByID deletes a link by id.
func (r *Repo) DeleteAnalysisJobLinkByID(ctx context.Context, querier sqlc.Querier, id int64) error {
	if err := querier.DeleteAnalysisJobLinkByID(ctx, id); err != nil {
		return fmt.Errorf("DeleteAnalysisJobLinkByID: %w", err)
	}
	return nil
}

// DeleteAnalysisJobLinkByPair deletes by analysis + candidate pair.
func (r *Repo) DeleteAnalysisJobLinkByPair(ctx context.Context, querier sqlc.Querier, params sqlc.DeleteAnalysisJobLinkByPairParams) error {
	if err := querier.DeleteAnalysisJobLinkByPair(ctx, params); err != nil {
		return fmt.Errorf("DeleteAnalysisJobLinkByPair: %w", err)
	}
	return nil
}

// =============================================================================
// User job discovery exclusions (URL hash batch lookup)
// =============================================================================

// ListUserJobDiscoveryExclusionsByURLHashes returns each url_hash in urlHashes that already has an exclusion row for this user (indexed batch lookup).
func (r *Repo) ListUserJobDiscoveryExclusionsByURLHashes(ctx context.Context, querier sqlc.Querier, userID int32, urlHashes []string) ([]string, error) {
	if len(urlHashes) == 0 {
		return nil, nil
	}
	rows, err := querier.ListUserJobDiscoveryExclusionsByURLHashes(ctx, sqlc.ListUserJobDiscoveryExclusionsByURLHashesParams{
		UserID:    userID,
		UrlHashes: urlHashes,
	})
	if err != nil {
		return nil, fmt.Errorf("ListUserJobDiscoveryExclusionsByURLHashes: %w", err)
	}
	return rows, nil
}

// UserJobDiscoveryExcludedURLHashSet is a convenience for filtering FindJobs results: keys are excluded url_hash values present in urlHashes.
func (r *Repo) UserJobDiscoveryExcludedURLHashSet(ctx context.Context, querier sqlc.Querier, userID int32, urlHashes []string) (map[string]struct{}, error) {
	existing, err := r.ListUserJobDiscoveryExclusionsByURLHashes(ctx, querier, userID, urlHashes)
	if err != nil {
		return nil, err
	}
	out := make(map[string]struct{}, len(existing))
	for _, h := range existing {
		out[h] = struct{}{}
	}
	return out, nil
}

// BatchExcludedURLHashes is an alias for UserJobDiscoveryExcludedURLHashSet (one indexed batch lookup per call).
func (r *Repo) BatchExcludedURLHashes(ctx context.Context, querier sqlc.Querier, userID int32, urlHashes []string) (map[string]struct{}, error) {
	return r.UserJobDiscoveryExcludedURLHashSet(ctx, querier, userID, urlHashes)
}

// AddJobDiscoveryExclusion inserts user_job_discovery_exclusions (ON CONFLICT DO NOTHING). reason and discoveryRunUuid are optional.
func (r *Repo) AddJobDiscoveryExclusion(ctx context.Context, querier sqlc.Querier, userID int32, urlHash, reason string, discoveryRunUuid pgtype.UUID) error {
	reason = strings.TrimSpace(reason)
	reas := pgtype.Text{}
	if reason != "" {
		reas = pgtype.Text{String: reason, Valid: true}
	}
	return r.AddUserJobDiscoveryExclusion(ctx, querier, sqlc.AddUserJobDiscoveryExclusionParams{
		UserID:           userID,
		UrlHash:          urlHash,
		SourceBoard:      pgtype.Text{},
		Reason:           reas,
		DiscoveryRunUuid: discoveryRunUuid,
	})
}

// AddUserJobDiscoveryExclusion records an excluded URL hash for a user. Idempotent on (user_id, url_hash).
func (r *Repo) AddUserJobDiscoveryExclusion(ctx context.Context, querier sqlc.Querier, params sqlc.AddUserJobDiscoveryExclusionParams) error {
	if err := querier.AddUserJobDiscoveryExclusion(ctx, params); err != nil {
		return fmt.Errorf("AddUserJobDiscoveryExclusion: %w", err)
	}
	return nil
}

// GetUserJobDiscoveryExclusionByUserAndURLHash loads one exclusion row when present.
func (r *Repo) GetUserJobDiscoveryExclusionByUserAndURLHash(ctx context.Context, querier sqlc.Querier, params sqlc.GetUserJobDiscoveryExclusionByUserAndURLHashParams) (*sqlc.BridgrUserJobDiscoveryExclusion, error) {
	row, err := querier.GetUserJobDiscoveryExclusionByUserAndURLHash(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("GetUserJobDiscoveryExclusionByUserAndURLHash: %w", err)
	}
	return &row, nil
}

// DeleteUserJobDiscoveryExclusionByUserAndURLHash removes an exclusion row (e.g. tooling, tests).
func (r *Repo) DeleteUserJobDiscoveryExclusionByUserAndURLHash(ctx context.Context, querier sqlc.Querier, params sqlc.DeleteUserJobDiscoveryExclusionByUserAndURLHashParams) error {
	if err := querier.DeleteUserJobDiscoveryExclusionByUserAndURLHash(ctx, params); err != nil {
		return fmt.Errorf("DeleteUserJobDiscoveryExclusionByUserAndURLHash: %w", err)
	}
	return nil
}
