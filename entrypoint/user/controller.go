package user

import (
	"coupon_be/service/user"
	"coupon_be/shared/fhttp"
	"net/http"

	"github.com/gorilla/mux"
)

// Controller manages the authentication operations, such as login, logout, etc.
type Controller struct {
	user user.Service
}

func (c *Controller) RegisterRoutes(r *mux.Router) {
	r.Handle("/register", fhttp.AppHandler(c.Register)).Methods(http.MethodPost)
}

func (c *Controller) Register(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	//var input request.Register
	//if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
	//	return nil, fhttp.NewErrorResponse(
	//		http.StatusUnprocessableEntity,
	//		sharedErrs.ErrKindValidation.String(),
	//		fmt.Sprintf("Invalid request body: %v", err))
	//}
	//
	//if err := util.Validate(input); err != nil {
	//	return nil, err
	//}

	if err := c.user.Register(ctx); err != nil {
		return nil, err
	}

	return &fhttp.Response{Status: http.StatusOK}, nil
}
