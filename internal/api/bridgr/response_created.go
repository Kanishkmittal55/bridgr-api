package bridgr

import (
	"net/http"

	"github.com/Kanishkmittal55/bridgr-api/internal/httpx"
)

func (s *server) writeCreated(w http.ResponseWriter, r *http.Request, body interface{}) {
	httpx.WriteResponse(w, r, http.StatusCreated, body)
}
