package auth

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	repository "github.com/Ser9unin/RealEstate/internal/storage/repo"
// )

// type ErrorResponse struct {
// 	Error string `json:"error"`
// }

// func TestUserByToken(t *testing.T) {

// 	t.Run("should validate if the email is not empty", func(t *testing.T) {
// 		payload := &repository.User{
// 			Email:    "",
// 			HashPass: "secretpass",
// 			Role:     "client",
// 		}

// 		b, err := json.Marshal(payload)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(b))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()
// 		router := http.NewServeMux()

// 		router.HandleFunc("/register", service.Register)

// 		router.ServeHTTP(rr, req)

// 		if rr.Code != http.StatusBadRequest {
// 			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
// 		}

// 		var response ErrorResponse
// 		err = json.NewDecoder(rr.Body).Decode(&response)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		if response.Error != errEmailRequired.Error() {
// 			t.Errorf("expected error message %s, got %s", response.Error, errEmailRequired.Error())
// 		}
// 	})

// 	t.Run("should create a user", func(t *testing.T) {
// 		payload := &repository.User{
// 			Email:    "realestate@avito.ru",
// 			HashPass: "secretpass",
// 			Role:     "client",
// 		}

// 		b, err := json.Marshal(payload)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(b))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()
// 		router := http.NewServeMux()

// 		router.HandleFunc("/register", service.Register)

// 		router.ServeHTTP(rr, req)

// 		if rr.Code != http.StatusCreated {
// 			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
// 		}
// 	})
// }

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	repository "github.com/Ser9unin/RealEstate/internal/storage/repo"
// )

// type ErrorResponse struct {
// 	Error string `json:"error"`
// }

// func TestUserByToken(t *testing.T) {

// 	t.Run("should validate if the email is not empty", func(t *testing.T) {
// 		payload := &repository.User{
// 			Email:    "",
// 			HashPass: "secretpass",
// 			Role:     "client",
// 		}

// 		b, err := json.Marshal(payload)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(b))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()
// 		router := http.NewServeMux()

// 		router.HandleFunc("/register", service.Register)

// 		router.ServeHTTP(rr, req)

// 		if rr.Code != http.StatusBadRequest {
// 			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
// 		}

// 		var response ErrorResponse
// 		err = json.NewDecoder(rr.Body).Decode(&response)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		if response.Error != errEmailRequired.Error() {
// 			t.Errorf("expected error message %s, got %s", response.Error, errEmailRequired.Error())
// 		}
// 	})

// 	t.Run("should create a user", func(t *testing.T) {
// 		payload := &repository.User{
// 			Email:    "realestate@avito.ru",
// 			HashPass: "secretpass",
// 			Role:     "client",
// 		}

// 		b, err := json.Marshal(payload)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(b))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()
// 		router := http.NewServeMux()

// 		router.HandleFunc("/register", service.Register)

// 		router.ServeHTTP(rr, req)

// 		if rr.Code != http.StatusCreated {
// 			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
// 		}
// 	})
// }
