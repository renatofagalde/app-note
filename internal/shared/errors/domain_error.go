package domainerror

import "fmt"

type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func New(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

var (
	ErrInvalidInput = New("INVALID_INPUT", "Input inválido")
	ErrNotFound     = New("NOT_FOUND", "Registro não encontrado")
	ErrConflict     = New("CONFLICT", "Registro já existente")
)
