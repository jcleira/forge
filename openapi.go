package forge

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type openAPI struct {
	OpenAPI    string              `yaml:"openapi"`
	Info       info                `yaml:"info"`
	Servers    []server            `yaml:"servers"`
	Paths      map[string]pathItem `yaml:"paths"`
	Components components          `yaml:"components"`
}

type info struct {
	Title       string  `yaml:"title"`
	Description string  `yaml:"description"`
	License     license `yaml:"license"`
	Version     string  `yaml:"version"`
}

type license struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type server struct {
	URL         string              `yaml:"url"`
	Description string              `yaml:"description"`
	Variables   map[string]variable `yaml:"variables,omitempty"`
}

type variable struct {
	Default string   `yaml:"default"`
	Enum    []string `yaml:"enum,omitempty"`
}

type pathItem struct {
	Post  *operation `yaml:"post,omitempty"`
	Get   *operation `yaml:"get,omitempty"`
	Patch *operation `yaml:"patch,omitempty"`
}

type operation struct {
	Tags        []string              `yaml:"tags"`
	Summary     string                `yaml:"summary"`
	Description string                `yaml:"description,omitempty"`
	OperationID string                `yaml:"operationId"`
	Security    []map[string][]string `yaml:"security"`
	RequestBody *requestBody          `yaml:"requestBody,omitempty"`
	Responses   map[string]response   `yaml:"responses"`
	Parameters  []parameter           `yaml:"parameters,omitempty"`
}

type requestBody struct {
	Content  map[string]mediaType `yaml:"content"`
	Required bool                 `yaml:"required"`
}

type mediaType struct {
	Schema schema `yaml:"schema"`
}

type schema struct {
	Type        string              `yaml:"type"`
	Properties  map[string]property `yaml:"properties,omitempty"`
	Ref         string              `yaml:"$ref,omitempty"`
	Items       *schema             `yaml:"items,omitempty"`
	Required    []string            `yaml:"required,omitempty"`
	Description string              `yaml:"description,omitempty"`
	Enum        []string            `yaml:"enum,omitempty"`
	Format      string              `yaml:"format,omitempty"`
	Example     interface{}         `yaml:"example,omitempty"`
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

type parameter struct {
	Name        string `yaml:"name"`
	In          string `yaml:"in"`
	Description string `yaml:"description"`
	Schema      schema `yaml:"schema"`
	Style       string `yaml:"style"`
	Explode     bool   `yaml:"explode"`
	Required    bool   `yaml:"required,omitempty"`
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
