package register

import (
	"context"

	repository "github.com/Ser9unin/Apartments/internal/storage/repo"
)

type MockStore struct {
}

func (s *MockStore) NewUser(ctx context.Context, arg repository.User) (repository.User, error) {
	return arg, nil
}

func (s *MockStore) UserByID(ctx context.Context, userID string) (repository.User, error) {
	return repository.User{}, nil
}

func (s *MockStore) UserByEmail(ctx context.Context, email string) (repository.User, error) {
	return repository.User{}, nil
}

func (s *MockStore) UserByIDAndRole(ctx context.Context, userID, userRole string) (bool, error) {
	return true, nil
}

type MockLogger struct{}

func (l *MockLogger) Info(msg string)  {}
func (l *MockLogger) Error(msg string) {}
func (l *MockLogger) Debug(msg string) {}
func (l *MockLogger) Warn(msg string)  {}
