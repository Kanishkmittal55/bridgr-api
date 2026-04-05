package response

import (
	"net/http"

	"github.com/hassleskip/bridgr-api/internal/logger"
	"github.com/hassleskip/bridgr-api/pkg/types"
	hserr "github.com/hassleskip/hassle-go/pkg/errors"
	hshttp "github.com/hassleskip/hassle-go/pkg/http"
)

const unauthorizedMsg = "Unauthorized"

// DefaultErrorHandler validates OpenAPI / chi wrapper failures.
func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	hshttp.WriteResponse(w, r, http.StatusBadRequest, types.NewErrorResponse(err.Error()))
}

// DefaultUnauthorizedHandler is used when API key auth fails.
func DefaultUnauthorizedHandler(w http.ResponseWriter, r *http.Request) {
	hshttp.WriteResponse(w, r, http.StatusUnauthorized, types.NewErrorResponse(unauthorizedMsg))
}

// DefaultPanicHandler logs and returns 500.
func DefaultPanicHandler(w http.ResponseWriter, r *http.Request, rvr any) {
	logger.Get().Errorw("panic occurred", "err", rvr)
	hshttp.WriteResponse(w, r, http.StatusInternalServerError, types.NewErrorResponse(hserr.InternalErrorString))
}
