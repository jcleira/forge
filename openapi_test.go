package forge

import (
	"errors"
	"reflect"
	"testing"
)

func TestOpenAPI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path    string
		want    openAPI
		wantErr error
	}{
		{
			path: "openapi.yaml",
			want: validOpenAPI(),
		},
		{
			path:    "invalid.yaml",
			wantErr: errors.New("error opening file: open invalid.yaml: no such file or directory"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()

			got, err := loadOpenAPI(tt.path)
			if tt.wantErr != nil && err != tt.wantErr {
				t.Errorf("TestOpenAPI\n error = '%v'\n wantErr '%v'", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadOpenAPI()\n got = '%v'\n want= '%v'", got, tt.want)
			}
		})
	}
}

func validOpenAPI() openAPI {
	return openAPI{
		OpenAPI: "3.0.3",
		Info: info{
			Title:       "API - Clients",
			Description: "Clients",
			License: license{
				Name: "proprietary",
				URL:  "N/A",
			},
			Version: "1.0",
		},
		Servers: []server{
			{
				URL:         "https://api.clients.com",
				Description: "Production",
			},
			{
				URL:         "http://api-{branch}.{cluster}.clients.co",
				Description: "Staging",
				Variables: map[string]variable{
					"branch": {
						Default: "master",
					},
					"cluster": {
						Default: "staging",
						Enum:    []string{"staging", "staging2", "staging3"},
					},
				},
			},
		},
		Paths: map[string]pathItem{
			"/v1/clients": {
				Post: &operation{
					Tags:        []string{"Clients"},
					Summary:     "Create a client",
					Description: "Creates a client.\n\nThis client will be used to create quotes, credit notes.",
					OperationID: "create-client",
					Security:    []map[string][]string{},
					RequestBody: &requestBody{
						Content: map[string]mediaType{
							"application/json": {
								Schema: schema{
									Type: "object",
									Properties: map[string]property{
										"data": {
											Ref: "#/components/schemas/ClientPayloadCreate",
										},
									},
									Required: []string{"data"},
								},
							},
						},
						Required: true,
					},
					Responses: map[string]response{
						"200": {
							Description: "Client created",
							Content: map[string]mediaType{
								"application/json": {
									Schema: schema{
										Type: "object",
										Properties: map[string]property{
											"data": {
												Ref: "#/components/schemas/ClientResponse",
											},
										},
									},
								},
							},
						},
						"400": {Ref: "#/components/responses/BadRequest"},
						"401": {Ref: "#/components/responses/Unauthorized"},
						"403": {Ref: "#/components/responses/Forbidden"},
						"412": {Ref: "#/components/responses/FailedPrecondition"},
						"422": {Ref: "#/components/responses/UnprocessableContentClient"},
						"500": {Ref: "#/components/responses/InternalServerError"},
					},
				},
				Get: &operation{
					Tags:        []string{"Clients"},
					Summary:     "List clients",
					OperationID: "list-clients",
					Security:    []map[string][]string{},
					Parameters: []parameter{
						{Name: "filter", In: "query", Schema: schema{Ref: "#/components/parameters/ClientListFilter"}},
						{Name: "sort_by", In: "query", Schema: schema{Ref: "#/components/parameters/ClientListSort"}},
						{Name: "page", In: "query", Schema: schema{Ref: "#/components/parameters/ListPage"}},
						{Name: "per_page", In: "query", Schema: schema{Ref: "#/components/parameters/ListPerPage"}},
					},
					Responses: map[string]response{
						"200": {
							Description: "An array of clients",
							Content: map[string]mediaType{
								"application/json": {
									Schema: schema{
										Type: "object",
										Properties: map[string]property{
											"data": {
												Type: "array",
												Items: &schema{
													Ref: "#/components/schemas/ClientResponse",
												},
											},
											"meta": {
												Ref: "#/components/schemas/ListMeta",
											},
										},
									},
								},
							},
						},
						"400": {Ref: "#/components/responses/BadRequest"},
						"401": {Ref: "#/components/responses/Unauthorized"},
						"403": {Ref: "#/components/responses/Forbidden"},
						"404": {Ref: "#/components/responses/NotFound"},
						"412": {Ref: "#/components/responses/FailedPrecondition"},
						"422": {Ref: "#/components/responses/UnprocessableContentClient"},
						"500": {Ref: "#/components/responses/InternalServerError"},
					},
				},
			},
		},
		Components: components{
			SecuritySchemes: map[string]securityScheme{
				"clients.read": {
					Description: "Membership should have clients.read permission to see the list of Invoice",
					Type:        "http",
					Scheme:      "permissions",
				},
			},
			Parameters: map[string]parameter{
				"ClientListFilter": {
					Name:        "filter",
					Description: "Attributes to filter by.",
					In:          "query",
					Schema: schema{
						Type: "object",
						Properties: map[string]property{
							"organization_id":           {Type: "string", Format: "uuid"},
							"tax_identification_number": {Type: "string"},
							"vat_number":                {Type: "string"},
							"email":                     {Type: "string"},
							"name":                      {Type: "string"},
							"created_at_from":           {Type: "string"},
							"created_at_to":             {Type: "string"},
						},
						Example: map[string]interface{}{
							"created_at_from": "2022-01-21T12:01:02Z",
							"created_at_to":   "2022-01-21T12:01:02Z",
						},
					},
					Style:   "deepObject",
					Explode: true,
				},
				"ClientListSort": {
					Name:        "sort_by",
					Description: `Attributes to sort by. Format is "field:order". Available orders are "asc" (Ascending) and "desc" (Descending).`,
					In:          "query",
					Schema: schema{
						Type: "array",
						Items: &schema{
							Type: "string",
							Enum: []string{"name", "created_at"},
						},
					},
					Style:   "form",
					Explode: false,
				},
				"ListPage": {
					Name:        "page",
					Description: "The page to retrieve",
					In:          "query",
					Schema: schema{
						Type: "integer",
					},
				},
				"ListPerPage": {
					Name:        "per_page",
					Description: "Number of items per page",
					In:          "query",
					Schema: schema{
						Type: "integer",
					},
				},
				"ClientID": {
					Name:        "client_id",
					In:          "path",
					Description: "Client unique identifier",
					Required:    true,
					Schema: schema{
						Type:   "string",
						Format: "uuid",
					},
				},
			},
			Schemas: map[string]schema{
				"ClientResponse": {
					Type: "object",
					Properties: map[string]property{
						"name":          {Type: "string", Description: "The name of the client. It is a concatenation of first and last name.", Example: "John Doe"},
						"first_name":    {Type: "string", Example: "John"},
						"last_name":     {Type: "string", Example: "Doe"},
						"kind":          {Type: "string", Enum: []string{"individual", "freelancer"}, Example: "individual"},
						"email":         {Type: "string", Format: "email", Example: "john.doe@clients.eu"},
						"locale":        {Type: "string", Example: "fr", Description: "The locale of the client."},
						"address":       {Type: "string", Description: "The address of the client. (eg street, number, floor, door, etc)", Example: "123 Main Street"},
						"city":          {Type: "string", Example: "Paris"},
						"zip_code":      {Type: "string", Example: "75009"},
						"province_code": {Type: "string", Description: "Represents the province code of the client."},
						"country_code":  {Type: "string", Example: "fr"},
					},
				},
				"Error": {
					Type: "object",
					Properties: map[string]property{
						"errors": {
							Type: "array",
							Items: &schema{
								Type: "object",
								Properties: map[string]property{
									"status": {Type: "string"},
									"code":   {Type: "string"},
									"detail": {Type: "string"},
									"source": {
										Type: "object",
										Properties: map[string]property{
											"pointer": {Type: "string"},
										},
									},
								},
								Required: []string{"status", "code", "detail"},
							},
						},
					},
				},
				"ClientValidationError": {
					Type:     "object",
					Required: []string{"errors"},
					Properties: map[string]property{
						"errors": {
							Type: "array",
							Items: &schema{
								Type:     "object",
								Required: []string{"code"},
								Properties: map[string]property{
									"code":   {Type: "string", Enum: []string{"required_unless", "required", "max", "iso3166_1_alpha2", "bcp47_language_tag"}},
									"detail": {Type: "string"},
									"status": {Type: "string", Example: "422"},
								},
							},
						},
					},
				},
				"ClientPayloadCreate": {
					Type:        "object",
					Description: "A client that is an individual or a freelancer.",
					Properties: map[string]property{
						"type": {Type: "string", Enum: []string{"customers"}},
						"attributes": {
							Type: "object",
							Properties: map[string]property{
								"first_name":    {Type: "string", MaxLength: 60, Example: "John"},
								"last_name":     {Type: "string", MaxLength: 60, Example: "Doe"},
								"kind":          {Type: "string", Enum: []string{"individual", "freelancer"}, Example: "individual"},
								"email":         {Type: "string", Format: "email", Example: "john.doe@clients.eu", MaxLength: 250},
								"locale":        {Type: "string", MaxLength: 2, MinLength: 2, Example: "fr", Description: "The locale of the client. It is used to generate the invoice in this language."},
								"address":       {Type: "string", MaxLength: 250, Description: "The address of the client. (eg street, number, floor, door, etc)"},
								"city":          {Type: "string", MaxLength: 50, Example: "Paris"},
								"zip_code":      {Type: "string", MaxLength: 20, Example: "75009"},
								"province_code": {Type: "string", MaxLength: 2, Description: "Represents the province code of the client. It is required only for Italian organizations"},
								"country_code":  {Type: "string", MaxLength: 2, Example: "fr"},
							},
							Required: []string{"first_name", "last_name", "email"},
						},
						"relationships": {
							Type: "object",
							Properties: map[string]property{
								"organization": {
									Type: "object",
									Properties: map[string]property{
										"data": {
											Type: "object",
											Properties: map[string]property{
												"type": {Type: "string", Enum: []string{"organizations"}},
												"id":   {Type: "string", Format: "uuid"},
											},
											Required: []string{"type", "id"},
										},
									},
									Required: []string{"data"},
								},
							},
							Required: []string{"organization"},
						},
					},
					Required: []string{"type", "attributes", "relationships"},
				},
				"ListMeta": {
					Type: "object",
					Properties: map[string]property{
						"current_page": {Type: "integer", Example: 2},
						"next_page":    {Type: "integer", Example: nil},
						"prev_page":    {Type: "integer", Example: 1},
						"total_pages":  {Type: "integer", Example: 2},
						"total_count":  {Type: "integer", Example: 150},
						"per_page":     {Type: "integer", Example: 100},
					},
				},
			},
			Responses: map[string]response{
				"BadRequest": {
					Description: "Bad request",
					Content: map[string]mediaType{
						"application/json": {
							Schema: schema{
								Ref: "#/components/schemas/Error",
							},
						},
					},
				},
				"Unauthorized": {
					Description: "Unauthorized",
					Content: map[string]mediaType{
						"application/json": {
							Schema: schema{
								Ref: "#/components/schemas/Error",
							},
						},
					},
				},
				"Forbidden": {
					Description: "Forbidden",
					Content: map[string]mediaType{
						"application/json": {
							Schema: schema{
								Ref: "#/components/schemas/Error",
							},
						},
					},
				},
				"NotFound": {
					Description: "Not found",
					Content: map[string]mediaType{
						"application/json": {
							Schema: schema{
								Ref: "#/components/schemas/Error",
							},
						},
					},
				},
				"FailedPrecondition": {
					Description: "Failed precondition",
					Content: map[string]mediaType{
						"application/json": {
							Schema: schema{
								Ref: "#/components/schemas/Error",
							},
						},
					},
				},
				"UnprocessableContentClient": {
					Description: "Unprocessable content",
					Content: map[string]mediaType{
						"application/json": {
							Schema: schema{
								Ref: "#/components/schemas/ClientValidationError",
							},
						},
					},
				},
				"InternalServerError": {
					Description: "Internal server error",
					Content: map[string]mediaType{
						"application/json": {
							Schema: schema{
								Ref: "#/components/schemas/Error",
							},
						},
					},
				},
			},
		},
	}
}
