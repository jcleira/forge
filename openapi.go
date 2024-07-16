package forge

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type openAPI struct {
	OpenAPI           string              `yaml:"openapi"`
	Info              info                `yaml:"info"`
	JsonSchemaDialect string              `yaml:"jsonSchemaDialect"`
	Servers           []server            `yaml:"servers"`
	Paths             map[string]pathItem `yaml:"paths"`
	Components        components          `yaml:"components"`
}

type info struct {
	Title         string  `yaml:"title"`
	Summary       string  `yaml:"summary"`
	Description   string  `yaml:"description"`
	TermOfService string  `yaml:"termOfService"`
	Contact       contact `yaml:"contact"`
	License       license `yaml:"license"`
	Version       string  `yaml:"version"`
}

type contact struct {
	Name  string `yaml:"name"`
	URL   string `yaml:"url"`
	Email string `yaml:"email"`
}

type license struct {
	Name       string `yaml:"name"`
	Identifier string `yaml:"identifier"`
	URL        string `yaml:"url"`
}

type server struct {
	URL         string              `yaml:"url"`
	Description string              `yaml:"description"`
	Variables   map[string]variable `yaml:"variables,omitempty"`
}

type variable struct {
	Default     string   `yaml:"default"`
	Enum        []string `yaml:"enum,omitempty"`
	Description string   `yaml:"description"`
}

type pathItem struct {
	Ref         string      `yaml:"$ref,omitempty"`
	Summary     string      `yaml:"summary,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Get         *operation  `yaml:"get,omitempty"`
	Put         *operation  `yaml:"put,omitempty"`
	Post        *operation  `yaml:"post,omitempty"`
	Delete      *operation  `yaml:"delete,omitempty"`
	Options     *operation  `yaml:"options,omitempty"`
	Head        *operation  `yaml:"head,omitempty"`
	Patch       *operation  `yaml:"patch,omitempty"`
	Trace       *operation  `yaml:"trace,omitempty"`
	Servers     []server    `yaml:"servers,omitempty"`
	Parameters  []parameter `yaml:"parameters,omitempty"`
}

type operation struct {
	Tags        []string     `yaml:"tags"`
	Summary     string       `yaml:"summary"`
	Description string       `yaml:"description,omitempty"`
	ExternalDoc externalDocs `yaml:"externalDocs,omitempty"`
	Parameters  []parameter  `yaml:"parameters,omitempty"`

	OperationID string                `yaml:"operationId"`
	Security    []map[string][]string `yaml:"security"`
	RequestBody *requestBody          `yaml:"requestBody,omitempty"`
	Responses   map[string]response   `yaml:"responses"`
}

type externalDocs struct {
	Description string `yaml:"description"`
	URL         string `yaml:"url"`
}

type parameter struct {
	Ref string `yaml:"$ref,omitempty"`

	Name            string `yaml:"name"`
	In              string `yaml:"in"`
	Description     string `yaml:"description"`
	Required        bool   `yaml:"required,omitempty"`
	Deprecated      bool   `yaml:"deprecated"`
	AllowEmptyValue bool   `yaml:"allowEmptyValue"`
	Style           string `yaml:"style"`
	Explode         bool   `yaml:"explode"`
	allowReserved   bool   `yaml:"allowReserved"`
	// TODO we are here
	Schema schema `yaml:"schema"`
}

type schema struct {
	Ref         string              `yaml:"$ref,omitempty"`
	Type        string              `yaml:"type"`
	Properties  map[string]property `yaml:"properties,omitempty"`
	Items       *schema             `yaml:"items,omitempty"`
	Required    []string            `yaml:"required,omitempty"`
	Description string              `yaml:"description,omitempty"`
	Enum        []string            `yaml:"enum,omitempty"`
	Format      string              `yaml:"format,omitempty"`
	Example     interface{}         `yaml:"example,omitempty"`
}

type requestBody struct {
	Content  map[string]mediaType `yaml:"content"`
	Required bool                 `yaml:"required"`
}

type mediaType struct {
	Schema schema `yaml:"schema"`
}

type property struct {
	Type        string              `yaml:"type"`
	Description string              `yaml:"description,omitempty"`
	Example     interface{}         `yaml:"example,omitempty"`
	Items       *schema             `yaml:"items,omitempty"`
	Properties  map[string]property `yaml:"properties,omitempty"`
	Required    []string            `yaml:"required,omitempty"`
	Format      string              `yaml:"format,omitempty"`
	MaxLength   int                 `yaml:"maxLength,omitempty"`
	MinLength   int                 `yaml:"minLength,omitempty"`
	Ref         string              `yaml:"$ref,omitempty"`
	Enum        []string            `yaml:"enum,omitempty"`
}

type response struct {
	Description string               `yaml:"description"`
	Content     map[string]mediaType `yaml:"content,omitempty"`
	Ref         string               `yaml:"$ref,omitempty"`
}

type components struct {
	SecuritySchemes map[string]securityScheme `yaml:"securitySchemes"`
	Parameters      map[string]parameter      `yaml:"parameters"`
	Schemas         map[string]schema         `yaml:"schemas"`
	Responses       map[string]response       `yaml:"responses"`
}

type securityScheme struct {
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
	Scheme      string `yaml:"scheme"`
}

// loadOpenAPI loads the openAPI specification from a file.
func loadOpenAPI(file string) (openAPI, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return openAPI{}, fmt.Errorf("error opening file: %v", err)
	}

	var api openAPI
	err = yaml.Unmarshal(data, &api)
	if err != nil {
		return openAPI{}, fmt.Errorf("error unmarshalling yaml: %v", err)
	}

	return api, nil
}
