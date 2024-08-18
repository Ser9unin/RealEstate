package render

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// responseJSON отправляет ответ в формате json.
func ResponseJSON(w http.ResponseWriter, _ *http.Request, status int, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(status)

	_, err := w.Write(buf.Bytes())
	if err != nil {
		log.Println(err)
	}
}

// JSONMap is a map alias.
type JSONMap map[string]interface{}

// ErrorJSON отправляет ошибку в формате json.
func ErrorJSON(w http.ResponseWriter, r *http.Request, httpStatusCode int, err error, details string) {
	ResponseJSON(w, r, httpStatusCode, JSONMap{"error": err.Error(), "details": details})
}

// NoContent отправляет ответ что контента нет.
func NoContent(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

var (
	ErrNotFound            = errors.New("your requested item is not found")
	ErrInternalServerError = errors.New("internal server error")
)

// StatusCode получает http статус из ошибки.
func StatusCode(err error) int {
	if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
