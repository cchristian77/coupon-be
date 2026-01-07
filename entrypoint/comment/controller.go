package comment

import (
	"coupon_be/request"
	"coupon_be/service/comment"
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/shared/fhttp"
	"coupon_be/util"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Controller manages the authentication operations, such as login, logout, etc.
type Controller struct {
	comment comment.Service
}

func (c *Controller) RegisterRoutes(r *mux.Router) {
	r.Handle("/{post_id:[0-9]+}/posts", fhttp.AppHandler(c.FilterComments)).Methods(http.MethodGet)
	r.Handle("/{comment_id:[0-9]+}", fhttp.AppHandler(c.Detail)).Methods(http.MethodGet)
	r.Handle("", fhttp.AppHandler(c.Store)).Methods(http.MethodPost)
	r.Handle("/{comment_id:[0-9]+}", fhttp.AppHandler(c.Update)).Methods(http.MethodPut)
	r.Handle("/{comment_id:[0-9]+}", fhttp.AppHandler(c.Delete)).Methods(http.MethodDelete)
}

func (c *Controller) FilterComments(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	var (
		input request.FilterComment
		err   error
	)

	input.PostID, err = strconv.ParseUint(mux.Vars(r)["post_id"], 10, 64)
	if err != nil {
		return nil, fhttp.NewErrorResponse(
			http.StatusBadRequest,
			sharedErrs.ErrKindValidation.String(),
			"Please provide the correct post_id as integer")
	}

	if data := r.URL.Query().Get("page"); data != "" {
		input.Page, err = strconv.Atoi(data)
		if err != nil || input.Page <= 0 {
			return nil, fhttp.NewErrorResponse(
				http.StatusBadRequest,
				sharedErrs.ErrKindValidation.String(),
				"Please provide a valid page as integer")
		}
	}

	if data := r.URL.Query().Get("per_page"); data != "" {
		input.PerPage, err = strconv.Atoi(data)
		if err != nil || input.Page <= 0 {
			return nil, fhttp.NewErrorResponse(
				http.StatusBadRequest,
				sharedErrs.ErrKindValidation.String(),
				"Please provide a valid page as integer")
		}
	}

	result, err := c.comment.FilterComments(ctx, &input)
	if err != nil {
		return nil, err
	}

	return &fhttp.Response{Data: result, Status: http.StatusOK}, nil
}

func (c *Controller) Detail(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	commentID, err := strconv.ParseUint(mux.Vars(r)["comment_id"], 10, 64)
	if err != nil {
		return nil, fhttp.NewErrorResponse(
			http.StatusBadRequest,
			sharedErrs.ErrKindValidation.String(),
			"Please provide the correct comment_id as integer")
	}

	result, err := c.comment.Detail(ctx, commentID)
	if err != nil {
		return nil, err
	}

	return &fhttp.Response{Data: result, Status: http.StatusOK}, nil
}

func (c *Controller) Store(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	var input request.CreateComment
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, fhttp.NewErrorResponse(
			http.StatusUnprocessableEntity,
			sharedErrs.ErrKindValidation.String(),
			fmt.Sprintf("Invalid request body: %v", err))
	}

	if err := util.Validate(input); err != nil {
		return nil, err
	}

	result, err := c.comment.Store(ctx, &input)
	if err != nil {
		return nil, err
	}

	return &fhttp.Response{
		Data:    result,
		Status:  http.StatusOK,
		Message: fmt.Sprintf("Comment %d is created successfully.", result.ID),
	}, nil
}

func (c *Controller) Update(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	var (
		input request.UpdateComment
		err   error
	)
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, fhttp.NewErrorResponse(
			http.StatusUnprocessableEntity,
			sharedErrs.ErrKindValidation.String(),
			fmt.Sprintf("Invalid request body: %v", err))
	}

	input.ID, err = strconv.ParseUint(mux.Vars(r)["comment_id"], 10, 64)
	if err != nil {
		return nil, fhttp.NewErrorResponse(
			http.StatusBadRequest,
			sharedErrs.ErrKindValidation.String(),
			"Please provide the correct comment_id as integer")
	}

	if err = util.Validate(input); err != nil {
		return nil, err
	}

	result, err := c.comment.Update(ctx, &input)
	if err != nil {
		return nil, err
	}

	return &fhttp.Response{Data: result, Status: http.StatusOK}, nil
}

func (c *Controller) Delete(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	commentID, err := strconv.ParseUint(mux.Vars(r)["comment_id"], 10, 64)
	if err != nil {
		return nil, fhttp.NewErrorResponse(
			http.StatusBadRequest,
			sharedErrs.ErrKindValidation.String(),
			"Please provide the correct comment_id as integer")
	}

	if err = c.comment.Delete(ctx, commentID); err != nil {
		return nil, err
	}

	return &fhttp.Response{
		Data:   "Comment is deleted successfully.",
		Status: http.StatusOK,
	}, nil
}
