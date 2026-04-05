package response

import (
	"net/http"

	"github.com/Kanishkmittal55/bridgr-api/internal/apierrors"
	"github.com/Kanishkmittal55/bridgr-api/internal/httpx"
	"github.com/Kanishkmittal55/bridgr-api/internal/logger"
	"github.com/Kanishkmittal55/bridgr-api/pkg/types"
)

const unauthorizedMsg = "Unauthorized"

// DefaultErrorHandler validates OpenAPI / chi wrapper failures.
func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	httpx.WriteResponse(w, r, http.StatusBadRequest, types.NewErrorResponse(err.Error()))
}

// DefaultUnauthorizedHandler is used when API key auth fails.
func DefaultUnauthorizedHandler(w http.ResponseWriter, r *http.Request) {
	httpx.WriteResponse(w, r, http.StatusUnauthorized, types.NewErrorResponse(unauthorizedMsg))
}

// DefaultPanicHandler logs and returns 500.
func DefaultPanicHandler(w http.ResponseWriter, r *http.Request, rvr any) {
	logger.Get().Errorw("panic occurred", "err", rvr)
	httpx.WriteResponse(w, r, http.StatusInternalServerError, types.NewErrorResponse(apierrors.InternalErrorString))
}
