package types

// NotFoundErrorResponse matches the Bridgr OpenAPI 404 body (not_found_by_api).
type NotFoundErrorResponse struct {
	Error         string `json:"error"`
	NotFoundByApi bool   `json:"not_found_by_api"`
}

func NewNotFoundErrorResponse(msg string) *NotFoundErrorResponse {
	return &NotFoundErrorResponse{Error: msg, NotFoundByApi: true}
}
