package register

import (
	"context"

	repository "github.com/Ser9unin/RealEstate/internal/storage/repo"
)

type MockStore struct{}

func (s *MockStore) NewUser(_ context.Context, arg repository.User) (repository.User, error) {
	return arg, nil
}

func (s *MockStore) UserByID(_ context.Context, _ string) (repository.User, error) {
	return repository.User{}, nil
}

func (s *MockStore) UserByEmail(_ context.Context, _ string) (repository.User, error) {
	return repository.User{}, nil
}

func (s *MockStore) UserByIDAndRole(_ context.Context, _, _ string) (bool, error) {
	return true, nil
}

type MockLogger struct{}

func (l *MockLogger) Info(_ string)  {}
func (l *MockLogger) Error(_ string) {}
func (l *MockLogger) Debug(_ string) {}
func (l *MockLogger) Warn(_ string)  {}
