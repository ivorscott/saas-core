package handlers

import (
	"errors"
	"fmt"
	"github.com/devpies/devpie-client-core/users/platform/web"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/devpies/devpie-client-core/users/domain/memberships"
	mockQuery "github.com/devpies/devpie-client-core/users/domain/mocks"
	mockAuth "github.com/devpies/devpie-client-core/users/platform/auth0/mocks"
	th "github.com/devpies/devpie-client-core/users/platform/testhelpers"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newMembership() memberships.NewMembership {
	return memberships.NewMembership{
		TeamID: "39541c75-ca3e-4e2b-9728-54327772d001",
		UserID: "a4b54ec1-57f9-4c39-ab53-d936dbb6c177",
		Role:   "administrator",
	}
}

func membership() memberships.Membership {
	return memberships.Membership{
		ID:        "085cb8a0-b221-4a6d-95be-592eb68d5670",
		TeamID:    "39541c75-ca3e-4e2b-9728-54327772d001",
		UserID:    "a4b54ec1-57f9-4c39-ab53-d936dbb6c177",
		Role:      "administrator",
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
}

func membershipEnhanced() memberships.MembershipEnhanced {
	return memberships.MembershipEnhanced{
		ID:        "085cb8a0-b221-4a6d-95be-592eb68d5670",
		TeamID:    "39541c75-ca3e-4e2b-9728-54327772d001",
		UserID:    "a4b54ec1-57f9-4c39-ab53-d936dbb6c177",
		FirstName: th.StringPointer("testuser"),
		LastName:  th.StringPointer("smith"),
		Picture:   th.StringPointer("https://s.gravatar.com/avatar/xxxxxxxxxxxx.png"),
		Email:     "example@devpie.io",
		Role:      "administrator",
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
}

func setupMembershipMocks() *Membership {
	return &Membership{
		repo:  th.Repo(),
		auth0: &mockAuth.Auther{},
		query: MembershipQueries{&mockQuery.MembershipQuerier{}},
	}
}

func TestMembership_RetrieveMemberships_200(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tid := "39541c75-ca3e-4e2b-9728-54327772d001"
	m := membershipEnhanced()
	ms := []memberships.MembershipEnhanced{m}

	// setup mocks
	fake := setupMembershipMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("RetrieveMemberships", mock.AnythingOfType("*context.valueCtx"), fake.repo, uid, tid).Return(ms, nil)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		_ = fake.RetrieveMemberships(w, r)
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tid), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Handler Response", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.membership.(*mockQuery.MembershipQuerier).AssertExpectations(t)
	})
}

func TestMembership_RetrieveMemberships_400(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tid := "39541c75-ca3e-4e2b-9728-54327772d001"

	// setup mocks
	fake := setupMembershipMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("RetrieveMemberships", mock.AnythingOfType("*context.valueCtx"), fake.repo, uid, tid).Return([]memberships.MembershipEnhanced{}, memberships.ErrInvalidID)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.RetrieveMemberships(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, memberships.ErrInvalidID))
			assert.Equal(t, http.StatusBadRequest, err.(*web.Error).Status)
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tid), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.membership.(*mockQuery.MembershipQuerier).AssertExpectations(t)
	})
}

func TestMembership_RetrieveMemberships_404(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tid := "39541c75-ca3e-4e2b-9728-54327772d001"

	// setup mocks
	fake := setupMembershipMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("RetrieveMemberships", mock.AnythingOfType("*context.valueCtx"), fake.repo, uid, tid).Return([]memberships.MembershipEnhanced{}, memberships.ErrNotFound)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.RetrieveMemberships(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, memberships.ErrNotFound))
			assert.Equal(t, http.StatusNotFound, err.(*web.Error).Status)
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tid), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.membership.(*mockQuery.MembershipQuerier).AssertExpectations(t)
	})
}

func TestMembership_RetrieveMemberships_500_Uncaught_Error(t *testing.T) {
	cause := errors.New("something went wrong")

	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tid := "39541c75-ca3e-4e2b-9728-54327772d001"

	// setup mocks
	fake := setupMembershipMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("RetrieveMemberships", mock.AnythingOfType("*context.valueCtx"), fake.repo, uid, tid).Return([]memberships.MembershipEnhanced{}, cause)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		err := fake.RetrieveMemberships(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
			assert.Equal(t, fmt.Sprintf(`failed to retrieve memberships: %s`, cause), err.Error())
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tid), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.membership.(*mockQuery.MembershipQuerier).AssertExpectations(t)
	})
}
