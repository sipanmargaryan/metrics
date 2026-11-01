package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"metrics/internal/helpers"
	"metrics/internal/middleware"
	"metrics/internal/order"
	"metrics/internal/product"
	"metrics/internal/warehouse"
)

func TestFieldValidationResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	if err := helpers.RegisterValidators(); err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	// Invalid requests must return before accessing these nil services.
	product.RegisterRoutes(router.Group("/products"), product.InitProductHandler(nil))
	warehouse.RegisterRoutes(router.Group("/warehouses"), warehouse.InitWarehouseHandler(nil))
	order.RegisterRoutes(router.Group("/orders"), order.InitOrderHandler(nil))
	for _, tc := range []struct {
		name   string
		path   string
		body   string
		error  string
		fields map[string]string
	}{
		{"missing product fields", "/products", `{}`, "validation failed", map[string]string{"name": "is required", "warehouse_id": "is required"}},
		{"blank name and missing warehouse", "/products", `{"name":" \t\n "}`, "validation failed", map[string]string{"name": "must not be blank", "warehouse_id": "is required"}},
		{"blank warehouse name", "/warehouses", `{"name":" "}`, "validation failed", map[string]string{"name": "must not be blank"}},
		{"negative product quantity", "/products", `{"name":"Supplies","warehouse_id":1,"quantity":-1}`, "validation failed", map[string]string{"quantity": "must be greater than or equal to 0"}},
		{"invalid order fields", "/orders", `{"quantity":-1}`, "validation failed", map[string]string{"product_id": "is required", "quantity": "must be greater than 0"}},
		{"wrong field type", "/products", `{"name":"Supplies","warehouse_id":"abc"}`, "validation failed", map[string]string{"warehouse_id": "has an invalid type"}},
		{"malformed JSON", "/products", `{"name":`, "invalid JSON body", nil},
		{"empty body", "/products", ``, "invalid JSON body", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tc.path, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("status=%d, body=%s", recorder.Code, recorder.Body.String())
			}
			var response helpers.ErrorResponse
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			if response.Error != tc.error || !reflect.DeepEqual(response.Fields, tc.fields) {
				t.Fatalf("response=%+v, want error=%q fields=%v", response, tc.error, tc.fields)
			}
		})
	}
}
