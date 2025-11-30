package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_GetAllNotes_ReturnsSeededData(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/notes", nil)
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperado status 200, recebeu %d. Body: %s", w.Code, w.Body.String())
	}

	var res []noteResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("falha ao parsear lista: %v", err)
	}

	if len(res) == 0 {
		t.Fatalf("esperava pelo menos 1 note dos seeds")
	}
}

func Test_GetAllNotes_AfterCreate_IncludesNewNote(t *testing.T) {
	// 1) Cria uma note nova
	body := `{
		"name": "Note criada no teste GetAll",
		"content": { "section": "test", "value": 123 }
	}`

	reqCreate := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()

	testRouter.ServeHTTP(wCreate, reqCreate)

	if wCreate.Code != http.StatusCreated {
		t.Fatalf("falha ao criar note. status=%d body=%s", wCreate.Code, wCreate.Body.String())
	}

	var created noteResponse
	if err := json.Unmarshal(wCreate.Body.Bytes(), &created); err != nil {
		t.Fatalf("erro ao parsear resposta de criação: %v", err)
	}

	// 2) Chama GET /notes e verifica se novo ID está na lista
	reqList := httptest.NewRequest(http.MethodGet, "/notes", nil)
	wList := httptest.NewRecorder()

	testRouter.ServeHTTP(wList, reqList)

	if wList.Code != http.StatusOK {
		t.Fatalf("esperado status 200, recebeu %d. Body: %s", wList.Code, wList.Body.String())
	}

	var list []noteResponse
	if err := json.Unmarshal(wList.Body.Bytes(), &list); err != nil {
		t.Fatalf("erro ao parsear lista: %v", err)
	}

	if len(list) == 0 {
		t.Fatalf("lista não deveria estar vazia")
	}

	found := false
	for _, n := range list {
		if n.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("note criada (id=%s) não encontrada na lista", created.ID)
	}
}
