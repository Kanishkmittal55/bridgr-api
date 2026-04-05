package bridgr

import (
	"net/http"
)

func (s *server) V1GetBridgrPing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok","product":"bridgr"}`))
}
