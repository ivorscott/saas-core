package handlers

import (
	"context"
	"fmt"
	"github.com/devpies/devpie-client-core/users/domain/users"
	"github.com/devpies/devpie-client-core/users/platform/auth0"
	mockAuth "github.com/devpies/devpie-client-core/users/platform/auth0/mocks"
	"github.com/devpies/devpie-client-core/users/platform/database"
	th "github.com/devpies/devpie-client-core/users/platform/testhelpers"
	"github.com/devpies/devpie-client-core/users/platform/web"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type queryMock struct {
	mock.Mock
}

func (m *queryMock) Create(ctx context.Context, repo database.Storer, nu users.NewUser, now time.Time) (users.User, error) {
	args := m.Called(ctx, repo, nu, now)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *queryMock) RetrieveByEmail(repo database.Storer, email string) (users.User, error) {
	args := m.Called(repo, email)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *queryMock) RetrieveMe(ctx context.Context, repo database.Storer, uid string) (users.User, error) {
	args := m.Called(ctx, repo, uid)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *queryMock) RetrieveMeByAuthID(ctx context.Context, repo database.Storer, aid string) (users.User, error) {
	args := m.Called(ctx, repo, aid)
	return args.Get(0).(users.User), args.Error(1)
}

type deps struct {
	service *Users
	repo *database.Repository
	auth0  *mockAuth.Auther
	query *queryMock
}

func setupMocks() *deps {
	mockRepo := th.Repo()
	mockAuth0 := &mockAuth.Auther{}
	mockQueries := &queryMock{}
	return &deps{
		repo:  mockRepo,
		auth0: mockAuth0,
		query: mockQueries,
		service: &Users{
			repo: mockRepo,
			auth0: mockAuth0,
			query: mockQueries,
		},
	}
}

func user() users.User {
	return users.User{
		ID: "a4b54ec1-57f9-4c39-ab53-d936dbb6c177",
		Auth0ID:   "auth0|60a666916089a00069b2a773",
		Email:     "testuser@devpie.io",
		FirstName: th.StringPointer("testuser"),
		Picture:   th.StringPointer("https://s.gravatar.com/avatar/xxxxxxxxxxxx.png"),
		EmailVerified: false,
		LastName:      th.StringPointer(""),
		Locale:        th.StringPointer(""),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func TestUsers_RetrieveMe_200(t *testing.T) {
	// setup mocks
	u := user()
	fake := setupMocks()
	fake.auth0.On("UserByID", context.Background()).Return(u.ID)
	fake.query.On("RetrieveMe", context.Background(), fake.repo, u.ID).Return(u, nil)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		u := fake.service
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
		fake.auth0.AssertExpectations(t)
		fake.query.AssertExpectations(t)
	})
}

func TestUsers_RetrieveMe_404_Missing_ID(t *testing.T) {
	// setup mocks
	fake := setupMocks()
	fake.auth0.On("UserByID", context.Background()).Return("")

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		u := fake.service

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
		fake.auth0.AssertExpectations(t)
	})
}

func TestUsers_RetrieveMe_404_Missing_Record(t *testing.T) {
	// setup mocks
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	fake := setupMocks()
	fake.auth0.On("UserByID", context.Background()).Return(uid)
	fake.query.On("RetrieveMe", context.Background(), fake.repo, uid).Return(users.User{}, users.ErrNotFound)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		u := fake.service

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
		fake.auth0.AssertExpectations(t)
		fake.query.AssertExpectations(t)
	})
}

func TestUsers_RetrieveMe_404_Invalid_ID(t *testing.T) {
	// setup mocks
	uid := "123"
	fake := setupMocks()
	fake.auth0.On("UserByID", context.Background()).Return(uid)
	fake.query.On("RetrieveMe", context.Background(), fake.repo, uid).Return(users.User{}, users.ErrInvalidID)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		u := fake.service

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
		fake.auth0.AssertExpectations(t)
		fake.query.AssertExpectations(t)
	})
}

func TestUsers_RetrieveMe_500_Uncaught_Error(t *testing.T) {
	cause := errors.New("Something went wrong")

	// setup mocks
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	fake := setupMocks()
	fake.auth0.On("UserByID", context.Background()).Return(uid)
	fake.query.On("RetrieveMe", context.Background(), fake.repo, uid).Return(users.User{}, cause)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		u := fake.service
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
		fake.auth0.AssertExpectations(t)
		fake.query.AssertExpectations(t)
	})
}

func TestUsers_Create_201(t *testing.T) {
	// setup mocks
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	nu := users.NewUser{
		Auth0ID:   "auth0|60a666916089a00069b2a773",
		Email:     "testuser@devpie.io",
		FirstName: th.StringPointer("testuser"),
		Picture:   th.StringPointer("https://s.gravatar.com/avatar/xxxxxxxxxxxx.png"),
	}

	fake := setupMocks()
	mockToken := auth0.Token{}

	fake.auth0.
		On("GenerateToken").Return(mockToken, nil).
		On("UpdateUserAppMetaData", mockToken, nu.Auth0ID, uid).Return(nil)

	fake.query.
		On("RetrieveMeByAuthID", context.Background(), fake.repo, nu.Auth0ID).
		Return(users.User{}, users.ErrNotFound).
		On("Create", context.Background(), fake.repo, nu, mock.AnythingOfType("time.Time")).
		Return(user(), nil)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		u := fake.service
		_ = u.Create(w, r)
	})

	// make request
	json := fmt.Sprintf(`{ "auth0Id": "%s", "email": "%s", "firstName": "%s", "picture": "%s" }`,
		nu.Auth0ID, nu.Email, *nu.FirstName, *nu.Picture)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(json))
	mux.ServeHTTP(writer, request)

	t.Run("Assert Handler Response", func(t *testing.T) {
		assert.Equal(t, http.StatusCreated, writer.Code)
	})

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.AssertExpectations(t)
		fake.query.AssertExpectations(t)
	})
}

func TestUsers_Create_400(t *testing.T) {
	// setup mocks
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	nu := users.NewUser{
		Auth0ID:   "auth0|60a666916089a00069b2a773",
		Email:     "testuser@devpie.io",
		FirstName: th.StringPointer("testuser"),
		Picture:   th.StringPointer("https://s.gravatar.com/avatar/xxxxxxxxxxxx.png"),
	}
	fake := setupMocks()
	mockToken := auth0.Token{}

	fake.auth0.
		On("GenerateToken").Return(mockToken, nil).
		On("UpdateUserAppMetaData", mockToken, nu.Auth0ID, uid).Return(auth0.ErrInvalidID)

	fake.query.
		On("RetrieveMeByAuthID", context.Background(), fake.repo, nu.Auth0ID).
		Return(users.User{}, users.ErrNotFound).
		On("Create", context.Background(), fake.repo, nu, mock.AnythingOfType("time.Time")).
		Return(user(), nil)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		u := fake.service

		var webErr *web.Error
		err := u.Create(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, auth0.ErrInvalidID))
			assert.Equal(t, http.StatusBadRequest, err.(*web.Error).Status)
		})
	})

	// make request
	json := fmt.Sprintf(`{ "auth0Id": "%s", "email": "%s", "firstName": "%s", "picture": "%s" }`,
		nu.Auth0ID, nu.Email, *nu.FirstName, *nu.Picture)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(json))
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.AssertExpectations(t)
		fake.query.AssertExpectations(t)
	})
}
