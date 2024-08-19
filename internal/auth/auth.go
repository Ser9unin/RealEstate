package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Ser9unin/RealEstate/internal/render"

	"github.com/golang-jwt/jwt"
)

func AuthModerator(httpHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, err := userByToken(r)
		if err != nil {
			permissionDenied(w, r, err)
			return
		}

		if role != "moderator" {
			permissionDenied(w, r, fmt.Errorf("no rights for such action"))
			return
		}

		// Call the function if the token is valid
		w.Header().Set("role", role)

		httpHandler(w, r)
	}
}

func AuthAny(httpHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, err := userByToken(r)
		if err != nil {
			permissionDenied(w, r, err)
			return
		}

		// Call the function if the token is valid
		w.Header().Set("role", role)
		httpHandler(w, r)
	}
}

func userByToken(r *http.Request) (string, error) {
	tokenString := getTokenFromRequest(r)
	if tokenString == "" {
		return "", fmt.Errorf("there is no token")
	}

	token, err := validateJWT(tokenString)
	if err != nil {
		log.Printf("failed to validate token: %v", err)
		return "", err
	}

	if !token.Valid {
		log.Println("invalid token")
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	userRole := claims["role"].(string)

	// эта часть под вопросом, теоретически роль может не совпадать в БД и реальности и это надо проверять
	// var storage repository.Queries
	// userID := claims["userID"].(string)
	// _, err = storage.UserByIDAndRole(r.Context(), userID, userRole)
	// if err != nil {
	// 	log.Printf("id or role is wrong: %v", err)
	// 	return "", err
	// }

	return userRole, nil
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		tokenAuth = strings.TrimPrefix(tokenAuth, "bearer ")
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func permissionDenied(w http.ResponseWriter, r *http.Request, err error) {
	render.ErrorJSON(w, r, http.StatusUnauthorized, err, "permission denied")
}
