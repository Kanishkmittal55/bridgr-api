package bridgr

import (
	"net/http"

	hsHttp "github.com/hassleskip/hassle-go/pkg/http"
)

func (s *server) writeCreated(w http.ResponseWriter, r *http.Request, body interface{}) {
	hsHttp.WriteResponse(w, r, http.StatusCreated, body)
}
