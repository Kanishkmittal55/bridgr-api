package bridgr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hassleskip/bridgr-api/internal/repository/sqlc"
	"github.com/hassleskip/bridgr-api/internal/uuid"
	hserr "github.com/hassleskip/hassle-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s *server) loadAnalysis(ctx context.Context, analysisUUID openapi_types.UUID) (pgtype.UUID, *sqlc.HskipUsersBridgrSkillGapAnalysis, error) {
	if err := s.requireStore(); err != nil {
		return pgtype.UUID{}, nil, err
	}
	pgid, err := uuid.ConvertOapiUUIDToPgUUID(analysisUUID)
	if err != nil {
		return pgtype.UUID{}, nil, fmt.Errorf("%w: invalid analysis uuid: %w", hserr.ErrBadRequest, err)
	}
	row, err := s.deps.Repo.GetSkillGapAnalysisByUUID(ctx, s.querier(), pgid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgtype.UUID{}, nil, fmt.Errorf("%w: analysis not found", hserr.ErrNotFound)
		}
		return pgtype.UUID{}, nil, fmt.Errorf("%w: %w", hserr.ErrInternal, err)
	}
	return pgid, row, nil
}

func (s *server) requireStore() error {
	if s.deps.Repo == nil {
		return fmt.Errorf("%w: database not configured", hserr.ErrServiceUnavailable)
	}
	return nil
}

func (s *server) loadGraph(ctx context.Context, graphUUID openapi_types.UUID) (pgtype.UUID, *sqlc.HskipUsersBridgrSkillGapGraph, error) {
	if err := s.requireStore(); err != nil {
		return pgtype.UUID{}, nil, err
	}
	pgid, err := uuid.ConvertOapiUUIDToPgUUID(graphUUID)
	if err != nil {
		return pgtype.UUID{}, nil, fmt.Errorf("%w: invalid graph uuid: %w", hserr.ErrBadRequest, err)
	}
	row, err := s.deps.Repo.GetSkillGapGraphByUUID(ctx, s.querier(), pgid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgtype.UUID{}, nil, fmt.Errorf("%w: graph not found", hserr.ErrNotFound)
		}
		return pgtype.UUID{}, nil, fmt.Errorf("%w: %w", hserr.ErrInternal, err)
	}
	return pgid, row, nil
}

func (s *server) loadPath(ctx context.Context, pathUUID openapi_types.UUID) (pgtype.UUID, *sqlc.HskipUsersBridgrSkillGapLearningPath, error) {
	if err := s.requireStore(); err != nil {
		return pgtype.UUID{}, nil, err
	}
	pgid, err := uuid.ConvertOapiUUIDToPgUUID(pathUUID)
	if err != nil {
		return pgtype.UUID{}, nil, fmt.Errorf("%w: invalid path uuid: %w", hserr.ErrBadRequest, err)
	}
	row, err := s.deps.Repo.GetSkillGapLearningPathByUUID(ctx, s.querier(), pgid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgtype.UUID{}, nil, fmt.Errorf("%w: learning path not found", hserr.ErrNotFound)
		}
		return pgtype.UUID{}, nil, fmt.Errorf("%w: %w", hserr.ErrInternal, err)
	}
	return pgid, row, nil
}

func (s *server) querier() sqlc.Querier {
	if s.deps.HsQuerier == nil {
		return nil
	}
	return s.deps.HsQuerier
}

func bridgrUserLimitOffset(limit, offset *int32) (int32, int32) {
	l := int32(20)
	if limit != nil && *limit > 0 {
		l = *limit
	}
	o := int32(0)
	if offset != nil && *offset > 0 {
		o = *offset
	}
	return l, o
}
