package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	authmiddleware "github.com/hassleskip/hassle-go/pkg/middleware/auth"

	"github.com/hassleskip/bridgr-api/internal/api/bridgr"
	"github.com/hassleskip/bridgr-api/internal/api/deps"
	"github.com/hassleskip/bridgr-api/internal/api/response"
)

// Routes is the route tree for bridgr-api (health, /v1); global middleware is applied in ApplyDefaultMiddleware.
func Routes(d *deps.Deps) http.Handler {
	r := chi.NewRouter()
	ApplyDefaultMiddleware(r)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	h := bridgr.NewServer(d)
	w := bridgr.ServerInterfaceWrapper{
		Handler:          h,
		ErrorHandlerFunc: response.DefaultErrorHandler,
	}

	r.Route("/v1", func(r chi.Router) {
		r.Route("/autograph", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(authmiddleware.ApiKeyAuthenticatorWithOptions(d.AccessToApiKeys[authmiddleware.AccessWrite], &authmiddleware.Options{ErrorHandler: response.DefaultUnauthorizedHandler}))
				r.Post("/upload-url", w.V1PostGenerateUploadUrl)
			})
		})
		r.Route("/bridgr", mountBridgr(&w, d))
	})

	return r
}

func mountBridgr(w *bridgr.ServerInterfaceWrapper, d *deps.Deps) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/ping", w.V1GetBridgrPing)
		r.Group(func(r chi.Router) {
			r.Use(authmiddleware.ApiKeyAuthenticatorWithOptions(d.AccessToApiKeys[authmiddleware.AccessRead], &authmiddleware.Options{ErrorHandler: response.DefaultUnauthorizedHandler}))
			r.Get("/analyses/{analysisUUID}", w.V1GetBridgrAnalysis)
			r.Get("/analyses/{analysisUUID}/coverage", w.V1GetBridgrAnalysisCoverage)
			r.Get("/analyses/{analysisUUID}/graphs", w.V1GetBridgrAnalysisGraphs)
			r.Get("/analyses/{analysisUUID}/graphs/{kind}", w.V1GetBridgrAnalysisGraphByKind)
			r.Get("/analyses/{analysisUUID}/learning-path", w.V1GetBridgrAnalysisLearningPath)
			r.Get("/analyses/{analysisUUID}/nodes/matched", w.V1GetBridgrAnalysisNodesMatched)
			r.Get("/analyses/{analysisUUID}/nodes/unmatched", w.V1GetBridgrAnalysisNodesUnmatched)
			r.Get("/graphs/{graphUUID}/nodes", w.V1GetBridgrGraphNodes)
			r.Get("/graphs/{graphUUID}/nodes/by-key/{nodeKey}", w.V1GetBridgrGraphNodeByKey)
			r.Get("/graphs/{graphUUID}/edges", w.V1GetBridgrGraphEdges)
			r.Get("/paths/{pathUUID}/steps", w.V1GetBridgrPathSteps)
			r.Get("/users/{userID}/analyses", w.V1GetBridgrUserAnalyses)
			r.Get("/users/{userID}/coverage", w.V1GetBridgrUserCoverage)
		})
		r.Group(func(r chi.Router) {
			r.Use(authmiddleware.ApiKeyAuthenticatorWithOptions(d.AccessToApiKeys[authmiddleware.AccessWrite], &authmiddleware.Options{ErrorHandler: response.DefaultUnauthorizedHandler}))
			r.Post("/analyses", w.V1PostBridgrAnalyses)
			r.Delete("/analyses/{analysisUUID}", w.V1DeleteBridgrAnalysis)
			r.Patch("/analyses/{analysisUUID}/status", w.V1PatchBridgrAnalysisStatus)
			r.Post("/analyses/{analysisUUID}/graphs", w.V1PostBridgrAnalysisGraphs)
			r.Delete("/analyses/{analysisUUID}/graphs", w.V1DeleteBridgrAnalysisGraphs)
			r.Post("/analyses/{analysisUUID}/learning-path", w.V1PostBridgrAnalysisLearningPath)
			r.Delete("/analyses/{analysisUUID}/learning-path", w.V1DeleteBridgrAnalysisLearningPath)
			r.Post("/graphs/{graphUUID}/nodes", w.V1PostBridgrGraphNodes)
			r.Post("/graphs/{graphUUID}/edges", w.V1PostBridgrGraphEdges)
		})
	}
}
