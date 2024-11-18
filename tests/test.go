package tests

import (
	"lab5/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllNodes(t *testing.T) {
	req, err := http.NewRequest("GET", "/nodes", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer your-secret-token")

	rr := httptest.NewRecorder()
	graphController := &handlers.GraphController{Neo4j: nil} // TODO: fix
	handler := http.HandlerFunc(graphController.GetAllNodes)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Дополнительная проверка JSON-ответа
}
