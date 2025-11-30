package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetNoteByID_Existing(t *testing.T) {
	// ID semeado em sql/notes.sql
	const seededID = "11111111-1111-1111-1111-111111111111"

	req := httptest.NewRequest(http.MethodGet, "/notes/"+seededID, nil)
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperado status 200, recebeu %d. Body: %s", w.Code, w.Body.String())
	}

	var res noteResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("falha ao parsear resposta: %v", err)
	}
	if res.ID != seededID {
		t.Fatalf("ID diferente do esperado. got=%q want=%q", res.ID, seededID)
	}
	if res.Name == "" {
		t.Fatalf("nome não deveria estar vazio")
	}
}

func Test_GetNoteByID_NotFound(t *testing.T) {
	const nonExistingID = "99999999-9999-9999-9999-999999999999"

	req := httptest.NewRequest(http.MethodGet, "/notes/"+nonExistingID, nil)
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("esperado status 404, recebeu %d. Body: %s", w.Code, w.Body.String())
	}
}
