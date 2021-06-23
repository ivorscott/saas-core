package handlers

import (
	"github.com/devpies/devpie-client-core/users/domain/users"
	mockAuth "github.com/devpies/devpie-client-core/users/platform/auth0/mocks"
	"github.com/devpies/devpie-client-core/users/platform/database"
	mockDB "github.com/devpies/devpie-client-core/users/platform/database/mocks"
	th "github.com/devpies/devpie-client-core/users/platform/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type QueryMock struct {
	mock.Mock
}

func (m *QueryMock) Create(ctx context.Context, repo database.Storer, nu users.NewUser, aid string, now time.Time) (users.User, error) {
	args := m.Called(ctx, repo, nu, aid, now)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *QueryMock) RetrieveByEmail(repo database.Storer, email string) (users.User, error) {
	args := m.Called(repo, email)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *QueryMock) RetrieveMe(ctx context.Context, repo database.Storer, uid string) (users.User, error) {
	args := m.Called(ctx, repo, uid)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *QueryMock) RetrieveMeByAuthID(ctx context.Context, repo database.Storer, aid string) (users.User, error) {
	args := m.Called(ctx, repo, aid)
	return args.Get(0).(users.User), args.Error(1)
}

func TestUsers_RetrieveMe_200(t *testing.T) {
	// setup mocks
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	mockAuth0 := &mockAuth.Auther{}
	mockAuth0.On("GetUserByID", context.Background()).Return(uid)
	mockRepo := &database.Repository{
		SqlxStorer: &mockDB.SqlxStorer{},
		Squirreler: &mockDB.Squirreler{},
		URL:        url.URL{},
	}
	mockQueries := &QueryMock{}
	mockQueries.On("RetrieveMe", context.Background(), mockRepo, uid).Return(users.User{
		ID:            uid,
		Auth0ID:       "auth0|60a666916089a00069b2a773",
		Email:         "testuser@devpie.io",
		EmailVerified: false,
		FirstName:     th.StringPointer("testuser"),
		LastName:      th.StringPointer(""),
		Picture:       th.StringPointer("https://s.gravatar.com/avatar/xxxxxxxxxxxx.png"),
		Locale:        th.StringPointer(""),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		u := &Users{
			repo:  mockRepo,
			auth0: mockAuth0,
			query: mockQueries,
		}
		_ = u.RetrieveMe(w, r)
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users/me", nil)
	mux.ServeHTTP(writer, request)

	mockAuth0.AssertExpectations(t)
	mockQueries.AssertExpectations(t)

	assert.Equal(t, writer.Code, http.StatusOK)
}
