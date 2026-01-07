package auth

import (
	"base_project/request"
	"base_project/response"
	"base_project/service/auth"
	sharedErrs "base_project/shared/errors"
	"base_project/shared/fhttp"
	"base_project/shared/fhttp/middleware"
	"base_project/util"
	"base_project/util/constant"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Controller manages the authentication operations, such as login, logout, etc.
type Controller struct {
	auth auth.Service
}

func (c *Controller) RegisterRoutes(r *mux.Router) {
	r.Handle("/login", fhttp.AppHandler(c.Login)).Methods(http.MethodPost)
	r.Handle("/register", fhttp.AppHandler(c.Register)).Methods(http.MethodPost)

	currentUserAPI := r.PathPrefix("/me").Subrouter()
	currentUserAPI.Handle("", fhttp.AppHandler(c.CurrentUser)).Methods(http.MethodGet)
	currentUserAPI.Use(middleware.Authentication())

	logoutAPI := r.PathPrefix("/logout").Subrouter()
	logoutAPI.Handle("", fhttp.AppHandler(c.Logout)).Methods(http.MethodPost)
	logoutAPI.Use(middleware.Authentication())
}

func (c *Controller) Login(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	var input request.Login

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, fhttp.NewErrorResponse(
			http.StatusUnprocessableEntity,
			sharedErrs.ErrKindValidation.String(),
			fmt.Sprintf("Invalid request body: %v", err))
	}

	if err := util.Validate(input); err != nil {
		return nil, err
	}

	result, err := c.auth.Login(ctx, &input)
	if err != nil {
		return nil, err
	}

	return &fhttp.Response{Data: result, Status: http.StatusOK}, nil
}

func (c *Controller) Logout(r *http.Request) (*fhttp.Response, error) {
	if err := c.auth.Logout(r.Context()); err != nil {
		return nil, err
	}

	return &fhttp.Response{Message: "Logout Successful.", Status: http.StatusOK}, nil
}

func (c *Controller) CurrentUser(r *http.Request) (*fhttp.Response, error) {
	authUser := constant.AuthUserFromCtx(r.Context())
	if authUser == nil {
		return nil, sharedErrs.InvalidTokenErr
	}

	return &fhttp.Response{
		Data: response.User{
			ID:       authUser.ID,
			Username: authUser.Username,
			FullName: authUser.FullName,
			Role:     authUser.Role.String(),
		},
		Status: http.StatusOK,
	}, nil
}

func (c *Controller) Register(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	var input request.Register
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, fhttp.NewErrorResponse(
			http.StatusUnprocessableEntity,
			sharedErrs.ErrKindValidation.String(),
			fmt.Sprintf("Invalid request body: %v", err))
	}

	if err := util.Validate(input); err != nil {
		return nil, err
	}

	result, err := c.auth.Register(ctx, &input)
	if err != nil {
		return nil, err
	}

	return &fhttp.Response{Data: result, Status: http.StatusOK}, nil
}
