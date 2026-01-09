package coupon

import (
	"coupon_be/request"
	"coupon_be/service/coupon"
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/shared/external/redis"
	"coupon_be/shared/fhttp"
	"coupon_be/util"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Controller manages the authentication operations, such as login, logout, etc.
type Controller struct {
	coupon    coupon.Service
	redisLock redis.ILock
}

func (c *Controller) RegisterRoutes(r *mux.Router) {
	r.Handle("", fhttp.AppHandler(c.Index)).Methods(http.MethodGet)
	r.Handle("/{coupon_name}", fhttp.AppHandler(c.Detail)).Methods(http.MethodGet)
	r.Handle("", fhttp.AppHandler(c.Store)).Methods(http.MethodPost)
	r.Handle("/claim", fhttp.AppHandler(c.Claim)).Methods(http.MethodPost)
}

func (c *Controller) Index(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	var input request.FilterCoupon

	input.Search = r.URL.Query().Get("search")

	result, err := c.coupon.Filter(ctx, &input)
	if err != nil {
		return nil, err
	}

	return &fhttp.Response{Data: result, Status: http.StatusOK}, nil
}

func (c *Controller) Detail(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	code := mux.Vars(r)["coupon_name"]
	if code == "" {
		return nil, fhttp.NewErrorResponse(
			http.StatusBadRequest,
			sharedErrs.ErrKindValidation.String(),
			"Please provide the correct coupon_name as string")
	}

	result, err := c.coupon.Detail(ctx, code)
	if err != nil {
		return nil, err
	}

	return &fhttp.Response{Data: result, Status: http.StatusOK}, nil
}

func (c *Controller) Store(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	var input request.UpsertCoupon
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, fhttp.NewErrorResponse(
			http.StatusUnprocessableEntity,
			sharedErrs.ErrKindValidation.String(),
			fmt.Sprintf("Invalid request body: %v", err))
	}

	if err := util.Validate(input); err != nil {
		return nil, err
	}

	result, err := c.coupon.Store(ctx, &input)
	if err != nil {
		return nil, err
	}

	return &fhttp.Response{
		Data:    result,
		Status:  http.StatusCreated,
		Message: fmt.Sprintf("Coupon %s is created successfully.", result.Name),
	}, nil
}

func (c *Controller) Claim(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	var input request.ClaimCoupon
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, fhttp.NewErrorResponse(
			http.StatusUnprocessableEntity,
			sharedErrs.ErrKindValidation.String(),
			fmt.Sprintf("Invalid request body: %v", err))
	}

	if err := util.Validate(input); err != nil {
		return nil, err
	}

	couponClaimKey := "claim:coupon:%s"
	err := c.redisLock.WithLock(ctx, fmt.Sprintf(couponClaimKey, input.CouponName), func() error {
		if err := c.coupon.Claim(ctx, &input); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &fhttp.Response{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("Coupon %s is successfully claimed by user %s.", input.CouponName, input.UserName),
	}, nil
}
