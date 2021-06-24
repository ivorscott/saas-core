package handlers

import (
	"context"
	"fmt"
	"github.com/devpies/devpie-client-core/users/domain/users"
	"github.com/devpies/devpie-client-core/users/platform/auth0"
	mockAuth "github.com/devpies/devpie-client-core/users/platform/auth0/mocks"
	"github.com/devpies/devpie-client-core/users/platform/database"
	mockDB "github.com/devpies/devpie-client-core/users/platform/database/mocks"
	th "github.com/devpies/devpie-client-core/users/platform/testhelpers"
	"github.com/devpies/devpie-client-core/users/platform/web"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type QueryMock struct {
	mock.Mock
}

func (m *QueryMock) Create(ctx context.Context, repo database.Storer, nu users.NewUser, now time.Time) (users.User, error) {
	args := m.Called(ctx, repo, nu, now)
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

	t.Run("Assert Handler Response", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		mockAuth0.AssertExpectations(t)
		mockQueries.AssertExpectations(t)
	})
}

func TestUsers_RetrieveMe_404_Missing_ID(t *testing.T) {
	// setup mocks
	mockAuth0 := &mockAuth.Auther{}
	mockAuth0.On("GetUserByID", context.Background()).Return("")

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		u := &Users{
			auth0: mockAuth0,
		}
		var webErr *web.Error
		err := u.RetrieveMe(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, users.ErrNotFound))
			assert.Equal(t, http.StatusNotFound, err.(*web.Error).Status)
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users/me", nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		mockAuth0.AssertExpectations(t)
	})
}

func TestUsers_RetrieveMe_404_Missing_Record(t *testing.T) {
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
	mockQueries.On("RetrieveMe", context.Background(), mockRepo, uid).Return(users.User{}, users.ErrNotFound)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		u := &Users{
			repo:  mockRepo,
			auth0: mockAuth0,
			query: mockQueries,
		}
		var webErr *web.Error
		err := u.RetrieveMe(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, users.ErrNotFound))
			assert.Equal(t, http.StatusNotFound, err.(*web.Error).Status)
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users/me", nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		mockAuth0.AssertExpectations(t)
		mockQueries.AssertExpectations(t)
	})
}

func TestUsers_RetrieveMe_404_Invalid_ID(t *testing.T) {
	// setup mocks
	uid := "123"
	mockAuth0 := &mockAuth.Auther{}
	mockAuth0.On("GetUserByID", context.Background()).Return(uid)
	mockRepo := &database.Repository{
		SqlxStorer: &mockDB.SqlxStorer{},
		Squirreler: &mockDB.Squirreler{},
		URL:        url.URL{},
	}
	mockQueries := &QueryMock{}
	mockQueries.On("RetrieveMe", context.Background(), mockRepo, uid).Return(users.User{}, users.ErrInvalidID)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		u := &Users{
			repo:  mockRepo,
			auth0: mockAuth0,
			query: mockQueries,
		}
		var webErr *web.Error
		err := u.RetrieveMe(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, users.ErrInvalidID))
			assert.Equal(t, http.StatusBadRequest, err.(*web.Error).Status)
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users/me", nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		mockAuth0.AssertExpectations(t)
		mockQueries.AssertExpectations(t)
	})
}

func TestUsers_RetrieveMe_500_Uncaught_Error(t *testing.T) {
	cause := errors.New("Something went wrong")

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
	mockQueries.On("RetrieveMe", context.Background(), mockRepo, uid).Return(users.User{}, cause)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		u := &Users{
			repo:  mockRepo,
			auth0: mockAuth0,
			query: mockQueries,
		}
		err := u.RetrieveMe(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
			assert.Equal(t, err.Error(), fmt.Sprintf(`looking for user "%s": %s`, uid, cause))
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users/me", nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		mockAuth0.AssertExpectations(t)
		mockQueries.AssertExpectations(t)
	})
}

func TestUsers_Create_201(t *testing.T) {
	// setup mocks
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	newUser := users.NewUser{
		Auth0ID:   "auth0|60a666916089a00069b2a773",
		Email:     "testuser@devpie.io",
		FirstName: th.StringPointer("testuser"),
		Picture:   th.StringPointer("https://s.gravatar.com/avatar/xxxxxxxxxxxx.png"),
	}
	mockToken := auth0.Token{}
	mockAuth0 := &mockAuth.Auther{}
	mockAuth0.
		On("GetOrCreateToken").Return(mockToken, nil).
		On("UpdateUserAppMetaData", mockToken, newUser.Auth0ID, uid).Return(nil).Once()

	mockRepo := &database.Repository{
		SqlxStorer: &mockDB.SqlxStorer{},
		Squirreler: &mockDB.Squirreler{},
		URL:        url.URL{},
	}
	mockQueries := &QueryMock{}
	mockQueries.
		On("RetrieveMeByAuthID", context.Background(), mockRepo, newUser.Auth0ID).
		Return(users.User{}, users.ErrNotFound).
		On("Create", context.Background(), mockRepo, newUser, mock.AnythingOfType("time.Time")).
		Return(users.User{
			ID:            uid,
			Auth0ID:       newUser.Auth0ID,
			Email:         newUser.Email,
			FirstName:     newUser.FirstName,
			Picture:       newUser.Picture,
			EmailVerified: false,
			LastName:      th.StringPointer(""),
			Locale:        th.StringPointer(""),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}, nil)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		u := &Users{
			repo:  mockRepo,
			auth0: mockAuth0,
			query: mockQueries,
		}
		_ = u.Create(w, r)
	})

	// make request
	json := fmt.Sprintf(`{ "auth0Id": "%s", "email": "%s", "firstName": "%s", "picture": "%s" }`,
		newUser.Auth0ID, newUser.Email, *newUser.FirstName, *newUser.Picture)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(json))
	mux.ServeHTTP(writer, request)

	t.Run("Assert Handler Response", func(t *testing.T) {
		assert.Equal(t, http.StatusCreated, writer.Code)
	})

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		mockAuth0.AssertExpectations(t)
		mockQueries.AssertExpectations(t)
	})
}
