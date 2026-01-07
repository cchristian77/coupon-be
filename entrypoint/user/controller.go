package user

import (
	"coupon_be/request"
	"coupon_be/service/user"
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/shared/fhttp"
	"coupon_be/util"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Controller manages the authentication operations, such as login, logout, etc.
type Controller struct {
	auth user.Service
}

func (c *Controller) RegisterRoutes(r *mux.Router) {
	r.Handle("/register", fhttp.AppHandler(c.Register)).Methods(http.MethodPost)
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
