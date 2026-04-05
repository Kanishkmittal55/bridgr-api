package openapi

import (
	_ "embed"

	"github.com/getkin/kin-openapi/openapi3"
)

var (
	//go:embed OpenAPISpec.gen.yaml
	spec []byte
)

// LoadSpec parses the bundled Bridgr API OpenAPI document.
func LoadSpec() *openapi3.T {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	doc, err := loader.LoadFromData(spec)
	if err != nil {
		panic(err)
	}
	return doc
}
