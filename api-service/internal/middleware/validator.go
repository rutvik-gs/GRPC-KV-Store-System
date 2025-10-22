package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

type ValidationMiddleware struct {
	router routers.Router
}

func NewValidationMiddleware(specPath string) (*ValidationMiddleware, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(specPath)
	if err != nil {
		return nil, err
	}

	if err := doc.Validate(loader.Context); err != nil {
		return nil, err
	}

	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		return nil, err
	}

	log.Println("OpenAPI spec loaded and validated successfully")

	return &ValidationMiddleware{
		router: router,
	}, nil
}

func (m *ValidationMiddleware) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, pathParams, err := m.router.FindRoute(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    r,
			PathParams: pathParams,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(context.Background(), requestValidationInput); err != nil {
			log.Printf("Validation error: %v", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Request validation failed: " + err.Error(),
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
