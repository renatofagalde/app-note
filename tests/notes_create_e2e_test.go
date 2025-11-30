package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_CreateNote_Success(t *testing.T) {
	body := `{
		"name": "Teste E2E - Sucesso",
		"content": { "html": "<h1>Hello Test</h1>", "lang": "pt-BR" }
	}`

	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("esperado status 201, recebeu %d. Body: %s", w.Code, w.Body.String())
	}

	var res noteResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("falha ao parsear resposta JSON: %v", err)
	}

	if res.ID == "" {
		t.Fatalf("esperava ID preenchido")
	}
	if res.Name != "Teste E2E - Sucesso" {
		t.Fatalf("nome diferente do esperado. got=%q", res.Name)
	}
	if len(res.Content) == 0 {
		t.Fatalf("content não pode ser vazio")
	}
}

func Test_CreateNote_InvalidJSON(t *testing.T) {
	body := `{
		"name": "Teste JSON inválido",
		"content": { "html": "<h1>Oops</h1>" }` // JSON quebrado

	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperado status 400 para JSON inválido, recebeu %d. Body: %s", w.Code, w.Body.String())
	}
}

func Test_CreateNote_MissingFields(t *testing.T) {
	// name em branco e content vazio → deve bater no ErrInvalidInput
	body := `{
		"name": "   ",
		"content": {}
	}`

	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperado status 400 para input inválido, recebeu %d. Body: %s", w.Code, w.Body.String())
	}
}
