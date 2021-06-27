package handlers

import (
	"context"
	"fmt"
	mockQuery "github.com/devpies/devpie-client-core/users/domain/mocks"
	"github.com/devpies/devpie-client-core/users/domain/users"
	"github.com/devpies/devpie-client-core/users/platform/auth0"
	mockAuth "github.com/devpies/devpie-client-core/users/platform/auth0/mocks"
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

func setupUserMocks() *User {
	return &User{
		repo:  th.Repo(),
		auth0: &mockAuth.Auther{},
		query: UserQueries{&mockQuery.UserQuerier{}},
	}
}

func newUser() users.NewUser {
	return users.NewUser{
		Auth0ID:   "auth0|60a666916089a00069b2a773",
		Email:     "testuser@devpie.io",
		FirstName: th.StringPointer("testuser"),
		Picture:   th.StringPointer("https://s.gravatar.com/avatar/xxxxxxxxxxxx.png"),
	}
}

func user() users.User {
	return users.User{
		ID:            "a4b54ec1-57f9-4c39-ab53-d936dbb6c177",
		Auth0ID:       "auth0|60a666916089a00069b2a773",
		Email:         "testuser@devpie.io",
		FirstName:     th.StringPointer("testuser"),
		Picture:       th.StringPointer("https://s.gravatar.com/avatar/xxxxxxxxxxxx.png"),
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
	fake := setupUserMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", context.Background()).Return(u.ID)
	fake.query.user.(*mockQuery.UserQuerier).On("RetrieveMe", context.Background(), fake.repo, u.ID).Return(u, nil)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		_ = fake.RetrieveMe(w, r)
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users/me", nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Handler Response", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.user.(*mockQuery.UserQuerier).AssertExpectations(t)
	})
}

func TestUsers_RetrieveMe_404_Missing_ID(t *testing.T) {
	// setup mocks
	fake := setupUserMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", context.Background()).Return("")

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.RetrieveMe(w, r)

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
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
	})
}

func TestUsers_RetrieveMe_404_Missing_Record(t *testing.T) {
	// setup mocks
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	fake := setupUserMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", context.Background()).Return(uid)
	fake.query.user.(*mockQuery.UserQuerier).On("RetrieveMe", context.Background(), fake.repo, uid).Return(users.User{}, users.ErrNotFound)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.RetrieveMe(w, r)

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
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.user.(*mockQuery.UserQuerier).AssertExpectations(t)
	})
}

func TestUsers_RetrieveMe_404_Invalid_ID(t *testing.T) {
	// setup mocks
	uid := "123"
	fake := setupUserMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", context.Background()).Return(uid)
	fake.query.user.(*mockQuery.UserQuerier).On("RetrieveMe", context.Background(), fake.repo, uid).Return(users.User{}, users.ErrInvalidID)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.RetrieveMe(w, r)

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
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.user.(*mockQuery.UserQuerier).AssertExpectations(t)
	})
}

func TestUsers_RetrieveMe_500_Uncaught_Error(t *testing.T) {
	cause := errors.New("Something went wrong")

	// setup mocks
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	fake := setupUserMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", context.Background()).Return(uid)
	fake.query.user.(*mockQuery.UserQuerier).On("RetrieveMe", context.Background(), fake.repo, uid).Return(users.User{}, cause)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		err := fake.RetrieveMe(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
			assert.Equal(t, fmt.Sprintf(`looking for user "%s": %s`, uid, cause), err.Error())
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users/me", nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.user.(*mockQuery.UserQuerier).AssertExpectations(t)
	})
}

func payload(nu users.NewUser) string {
	return fmt.Sprintf(`{ "auth0Id": "%s", "email": "%s", "firstName": "%s", "picture": "%s" }`,
		nu.Auth0ID, nu.Email, *nu.FirstName, *nu.Picture)
}

func TestUsers_Create_201(t *testing.T) {
	// setup mocks
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	nu := newUser()
	fake := setupUserMocks()

	fake.auth0.(*mockAuth.Auther).
		On("GenerateToken").Return(auth0.Token{}, nil).
		On("UpdateUserAppMetaData", auth0.Token{}, nu.Auth0ID, uid).Return(nil)

	fake.query.user.(*mockQuery.UserQuerier).
		On("RetrieveMeByAuthID", context.Background(), fake.repo, nu.Auth0ID).
		Return(users.User{}, users.ErrNotFound).
		On("Create", context.Background(), fake.repo, nu, mock.AnythingOfType("time.Time")).
		Return(user(), nil)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		_ = fake.Create(w, r)
	})

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(payload(nu)))
	mux.ServeHTTP(writer, request)

	t.Run("Assert Handler Response", func(t *testing.T) {
		assert.Equal(t, http.StatusCreated, writer.Code)
	})

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.user.(*mockQuery.UserQuerier).AssertExpectations(t)
	})
}

func TestUsers_Create_400_Missing_Payload(t *testing.T) {
	// setup mocks
	fake := setupUserMocks()

	testcases := []struct {
		name string
		arg  string
	}{
		{"empty payload", ""},
		{"empty object", "{}"},
	}

	for _, v := range testcases {
		// setup server
		mux := http.NewServeMux()
		mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
			err := fake.Create(w, r)

			t.Run(fmt.Sprintf("Assert Handler Response/%s", v.name), func(t *testing.T) {
				assert.NotNil(t, err)
			})
		})

		writer := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(v.arg))
		mux.ServeHTTP(writer, request)

		t.Run(fmt.Sprintf("Assert Server Response/%s", v.name), func(t *testing.T) {
			assert.Equal(t, http.StatusBadRequest, writer.Code)
		})
	}
}

func TestUsers_Create_400_Invalid_ID_For_UpdateUserAppMetadata(t *testing.T) {
	// setup mocks
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	nu := newUser()
	fake := setupUserMocks()

	fake.auth0.(*mockAuth.Auther).
		On("GenerateToken").Return(auth0.Token{}, nil).
		On("UpdateUserAppMetaData", auth0.Token{}, nu.Auth0ID, uid).Return(auth0.ErrInvalidID)

	fake.query.user.(*mockQuery.UserQuerier).
		On("RetrieveMeByAuthID", context.Background(), fake.repo, nu.Auth0ID).
		Return(users.User{}, users.ErrNotFound).
		On("Create", context.Background(), fake.repo, nu, mock.AnythingOfType("time.Time")).
		Return(user(), nil)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.Create(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, auth0.ErrInvalidID))
			assert.Equal(t, http.StatusBadRequest, err.(*web.Error).Status)
		})
	})

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(payload(nu)))
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.user.(*mockQuery.UserQuerier).AssertExpectations(t)
	})
}

func TestUsers_Create_500_Uncaught_Error_For_GenerateToken(t *testing.T) {
	cause := errors.New("Something went wrong")

	// setup mocks
	nu := newUser()
	fake := setupUserMocks()
	fake.auth0.(*mockAuth.Auther).On("GenerateToken").Return(auth0.Token{}, cause)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		err := fake.Create(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
			assert.Equal(t, cause.Error(), err.Error())
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(payload(nu)))
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.user.(*mockQuery.UserQuerier).AssertExpectations(t)
	})
}

func TestUsers_Create_500_Uncaught_Error_For_Create(t *testing.T) {
	cause := errors.New("Something went wrong")

	// setup mocks
	nu := newUser()
	fake := setupUserMocks()
	fake.query.user.(*mockQuery.UserQuerier).
		On("RetrieveMeByAuthID", context.Background(), fake.repo, nu.Auth0ID).
		Return(users.User{}, users.ErrNotFound).
		On("Create", context.Background(), fake.repo, nu, mock.AnythingOfType("time.Time")).
		Return(users.User{}, cause)

	fake.auth0.(*mockAuth.Auther).On("GenerateToken").Return(auth0.Token{}, nil)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		err := fake.Create(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
			assert.Equal(t, fmt.Sprintf("failed to create user: %s", cause), err.Error())
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(payload(nu)))
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.user.(*mockQuery.UserQuerier).AssertExpectations(t)
	})
}

func TestUsers_Create_500_Uncaught_Error_For_UpdateUserAppMetadata(t *testing.T) {
	cause := errors.New("Something went wrong")

	// setup mocks
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	nu := newUser()
	fake := setupUserMocks()
	fake.query.user.(*mockQuery.UserQuerier).
		On("RetrieveMeByAuthID", context.Background(), fake.repo, nu.Auth0ID).
		Return(users.User{}, users.ErrNotFound).
		On("Create", context.Background(), fake.repo, nu, mock.AnythingOfType("time.Time")).
		Return(user(), nil)

	fake.auth0.(*mockAuth.Auther).On("GenerateToken").Return(auth0.Token{}, nil).
		On("UpdateUserAppMetaData", auth0.Token{}, nu.Auth0ID, uid).Return(cause)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		err := fake.Create(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
			assert.Equal(t, fmt.Sprintf("failed to update user app metadata: %s", cause), err.Error())
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(payload(nu)))
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.user.(*mockQuery.UserQuerier).AssertExpectations(t)
	})
}
