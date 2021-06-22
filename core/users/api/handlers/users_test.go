package handlers

import (
	"context"
	"github.com/devpies/devpie-client-core/users/domain/users"
	mockQueries "github.com/devpies/devpie-client-core/users/domain/users/mocks"
	"github.com/devpies/devpie-client-core/users/platform/auth0/mocks"
	mockRepo "github.com/devpies/devpie-client-core/users/platform/database/mocks"
	th "github.com/devpies/devpie-client-core/users/platform/testhelpers"

	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUsers_RetrieveMe_200(t *testing.T) {
	// mocks
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	mockReturnUser := users.User{
		ID: uid,
		Auth0ID: "auth0|60a666916089a00069b2a773",
		Email: "testuser@devpie.io",
		EmailVerified: false,
		FirstName: th.StringPointer("testuser"),
		LastName: th.StringPointer(""),
		Picture: th.StringPointer("https://s.gravatar.com/avatar/xxxxxxxxxxxx.png"),
		Locale: th.StringPointer(""),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockAuth0 := &mocks.Auther{}
	mockAuth0.On("GetUserByID", context.Background()).Return(uid)
	mockRepository := &mockRepo.DataStorer{}
	mockUserQueries := &mockQueries.UserQuerier{}
	mockUserQueries.On("RetrieveMe", context.Background(), mockRepository, uid).Return(mockReturnUser, nil)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request ) {
		u := &Users{auth0: mockAuth0}
		_ = u.RetrieveMe(w,r)
	})

	// execute request
	writer := httptest.NewRecorder()
	request , _ := http.NewRequest(http.MethodGet,"/users/me", nil)
	mux.ServeHTTP(writer, request)

	// make assertions
	//assert := assertFunc.New(t)

	//assert.Equal(http.StatusOK, writer.Code)
}
