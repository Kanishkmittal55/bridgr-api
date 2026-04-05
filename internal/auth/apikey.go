package auth

import (
	"net/http"

	"github.com/Kanishkmittal55/bridgr-api/internal/httpx"
	"github.com/Kanishkmittal55/bridgr-api/pkg/types"
)

const DefaultXAPIKeyHeader = "X-API-KEY"

type ErrorHandler func(w http.ResponseWriter, r *http.Request)

type Options struct {
	ApiKeyHeader *string
	ErrorHandler ErrorHandler
}

func ApiKeyAuthenticator(apiKeys []string) func(http.Handler) http.Handler {
	return ApiKeyAuthenticatorWithOptions(apiKeys, nil)
}

func ApiKeyAuthenticatorWithOptions(apiKeys []string, opts *Options) func(http.Handler) http.Handler {
	header := DefaultXAPIKeyHeader
	errorHandler := func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteResponse(w, r, http.StatusUnauthorized, types.NewErrorResponse("Insufficient permissions"))
	}

	if opts != nil {
		if opts.ApiKeyHeader != nil {
			header = *opts.ApiKeyHeader
		}
		if opts.ErrorHandler != nil {
			errorHandler = opts.ErrorHandler
		}
	}

	keyMap := make(map[string]struct{}, len(apiKeys))
	for _, k := range apiKeys {
		keyMap[k] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := keyMap[r.Header.Get(header)]; ok {
				next.ServeHTTP(w, r)
			} else {
				errorHandler(w, r)
			}
		})
	}
}
