package register

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Ser9unin/RealEstate/internal/render"
	repository "github.com/Ser9unin/RealEstate/internal/storage/repo"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	errEmailRequired       = errors.New("email is required")
	errRoleRequired        = errors.New("user role is required")
	errPasswordRequired    = errors.New("password is required")
	errEmailOrUUIDRequired = errors.New("email or id is required")
)

type UserService struct {
	storage Storage
	logger  Logger
}

type Storage interface {
	NewUser(ctx context.Context, arg repository.User) (repository.User, error)
	UserByID(ctx context.Context, userID string) (repository.User, error)
	UserByEmail(ctx context.Context, email string) (repository.User, error)
	UserByIDAndRole(ctx context.Context, userID, userRole string) (bool, error)
}
type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

var empty struct{}
var roleMap = map[string]struct{}{
	"client":    empty,
	"moderator": empty,
}

func NewUserService(s Storage, logger Logger) *UserService {
	return &UserService{storage: s, logger: logger}
}

func (s *UserService) Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var user repository.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		s.logger.Error(err.Error())
		render.ErrorJSON(w, r, http.StatusBadRequest, err, "Invalid request payload")
		return
	}

	if err := validateUserPayload(user); err != nil {
		s.logger.Error(err.Error())
		render.ErrorJSON(w, r, http.StatusBadRequest, err, "not enought validation data")
		return
	}

	user.HashPass, err = render.HashPassword(user.HashPass)
	if err != nil {
		s.logger.Error(err.Error())
		render.ErrorJSON(w, r, http.StatusInternalServerError, err, "Error creating user")
		return
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		s.logger.Error(err.Error())
		render.ErrorJSON(w, r, http.StatusInternalServerError, err, "Error creating user")
		return
	}

	user.UserID = uuid.String()

	_, err = s.storage.NewUser(r.Context(), user)
	if err != nil {
		s.logger.Error(err.Error())
		render.ErrorJSON(w, r, http.StatusInternalServerError, err, "Error creating user")
		return
	}

	render.ResponseJSON(w, r, http.StatusOK, uuid)
}

func (s *UserService) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var user repository.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		s.logger.Error(err.Error())
		render.ErrorJSON(w, r, http.StatusBadRequest, err, "invalid request payload")
		return
	}

	if err := validateLoginPayload(user); err != nil {
		s.logger.Error(err.Error())
		render.ErrorJSON(w, r, http.StatusBadRequest, err, "not enought validation data")
		return
	}

	password := user.HashPass

	if user.UserID != "" {
		user, err = s.storage.UserByID(r.Context(), user.UserID)
	} else {
		user, err = s.storage.UserByEmail(r.Context(), user.Email)
	}
	if err != nil {
		s.logger.Error(err.Error())
		render.ErrorJSON(w, r, http.StatusInternalServerError, err, "error login user")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPass), []byte(password)); err != nil {
		s.logger.Error(err.Error())
		render.ErrorJSON(w, r, http.StatusInternalServerError, err, "error login user")
		return
	}

	token, err := createJWT(user)
	if err != nil {
		s.logger.Error(err.Error())
		render.ErrorJSON(w, r, http.StatusInternalServerError, err, "error login user")
		return
	}

	render.ResponseJSON(w, r, http.StatusOK, token)
}

func createJWT(user repository.User) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    user.UserID,
		"role":      user.Role,
		"expiresAt": time.Now().Add(time.Hour * 24 * 120).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func validateUserPayload(user repository.User) error {
	if user.Email == "" {
		return errEmailRequired
	}

	if user.HashPass == "" {
		return errPasswordRequired
	}

	if user.Role == "" {
		return errRoleRequired
	}

	if _, ok := roleMap[user.Role]; !ok {
		return errRoleRequired
	}

	return nil
}

func validateLoginPayload(user repository.User) error {
	if user.Email == "" && user.UserID == "" {
		return errEmailOrUUIDRequired
	}

	if user.HashPass == "" {
		return errPasswordRequired
	}

	return nil
}
