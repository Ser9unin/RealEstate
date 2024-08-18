package register

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	repository "github.com/Ser9unin/Apartments/internal/storage/repo"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func TestValidateUserPayload(t *testing.T) {
	type args struct {
		user *repository.User
	}
	tests := []struct {
		name string
		args args
		want error
	}{

		{
			name: "should return error if email is empty",
			args: args{
				user: &repository.User{
					HashPass: "secretpass",
					Role:     "client",
				},
			},
			want: errEmailRequired,
		},
		{
			name: "should return error if role is empty",
			args: args{
				user: &repository.User{
					Email:    "realestate@avito.ru",
					HashPass: "secretpass",
				},
			},
			want: errRoleRequired,
		},
		{
			name: "should return error if the password is empty",
			args: args{
				user: &repository.User{
					Email: "realestate@avito.ru",
					Role:  "client",
				},
			},
			want: errPasswordRequired,
		},
		{
			name: "should return nil if all fields are present",
			args: args{
				user: &repository.User{
					Email:    "realestate@avito.ru",
					HashPass: "secretpass",
					Role:     "client",
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateUserPayload(*tt.args.user); got != tt.want {
				t.Errorf("validateUserPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	ms := &MockStore{}
	logger := &MockLogger{}
	service := NewUserService(ms, logger)

	t.Run("should validate if the email is not empty", func(t *testing.T) {
		payload := &repository.User{
			Email:    "",
			HashPass: "secretpass",
			Role:     "client",
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		router.HandleFunc("/register", service.Register)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		var response ErrorResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		if response.Error != errEmailRequired.Error() {
			t.Errorf("expected error message %s, got %s", response.Error, errEmailRequired.Error())
		}
	})

	t.Run("should create a user", func(t *testing.T) {
		payload := &repository.User{
			Email:    "realestate@avito.ru",
			HashPass: "secretpass",
			Role:     "client",
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		router.HandleFunc("/register", service.Register)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}
