package post

import (
	"base_project/domain/enums"
	"base_project/request"
	"base_project/service/post"
	sharedErrs "base_project/shared/errors"
	"base_project/shared/fhttp"
	"base_project/shared/fhttp/middleware"
	"base_project/util"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Controller manages the authentication operations, such as login, logout, etc.
type Controller struct {
	post post.Service
}

func (c *Controller) RegisterRoutes(r *mux.Router) {
	r.Use(middleware.Authentication())
	r.Handle("", fhttp.AppHandler(c.Index)).Methods(http.MethodGet)
	r.Handle("/{post_id:[0-9]+}", fhttp.AppHandler(c.Detail)).Methods(http.MethodGet)
	r.Handle("", fhttp.AppHandler(c.Store)).Methods(http.MethodPost)
	r.Handle("/{post_id:[0-9]+}", fhttp.AppHandler(c.Update)).Methods(http.MethodPut)
	r.Handle("/{post_id:[0-9]+}/publish", fhttp.AppHandler(c.Publish)).Methods(http.MethodPut)
	r.Handle("/{post_id:[0-9]+}/draft", fhttp.AppHandler(c.Draft)).Methods(http.MethodPut)
	r.Handle("/{post_id:[0-9]+}", fhttp.AppHandler(c.Delete)).Methods(http.MethodDelete)
	r.Handle("/user", fhttp.AppHandler(c.MyPosts)).Methods(http.MethodGet)
}

func (c *Controller) Index(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	var (
		input request.FilterPost
		err   error
	)

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

	input.Search = r.URL.Query().Get("search")

	if err = util.Validate(input); err != nil {
		return nil, err
	}

	result, err := c.post.FilterPosts(ctx, &input)
	if err != nil {
		return nil, err
	}

	return &fhttp.Response{Data: result, Status: http.StatusOK}, nil
}

func (c *Controller) MyPosts(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	var (
		input request.FilterPost
		err   error
	)

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
		input.Page, err = strconv.Atoi(data)
		if err != nil || input.Page <= 0 {
			return nil, fhttp.NewErrorResponse(
				http.StatusBadRequest,
				sharedErrs.ErrKindValidation.String(),
				"Please provide a valid page as integer")
		}
	}

	input.Search = r.URL.Query().Get("search")

	result, err := c.post.FilterMyPosts(ctx, &input)
	if err != nil {
		return nil, err
	}

	return &fhttp.Response{Data: result, Status: http.StatusOK}, nil
}

func (c *Controller) Detail(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	postID, err := strconv.ParseUint(mux.Vars(r)["post_id"], 10, 64)
	if err != nil {
		return nil, fhttp.NewErrorResponse(
			http.StatusBadRequest,
			sharedErrs.ErrKindValidation.String(),
			"Please provide the correct post_id as integer")
	}

	result, err := c.post.Detail(ctx, postID)
	if err != nil {
		return nil, err
	}

	return &fhttp.Response{Data: result, Status: http.StatusOK}, nil
}

func (c *Controller) Store(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	var input request.UpsertPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, fhttp.NewErrorResponse(
			http.StatusUnprocessableEntity,
			sharedErrs.ErrKindValidation.String(),
			fmt.Sprintf("Invalid request body: %v", err))
	}

	if err := util.Validate(input); err != nil {
		return nil, err
	}

	result, err := c.post.Store(ctx, &input)
	if err != nil {
		return nil, err
	}

	return &fhttp.Response{
		Data:    result,
		Status:  http.StatusOK,
		Message: fmt.Sprintf("Post %d is created successfully.", result.ID),
	}, nil
}

func (c *Controller) Update(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	var (
		input request.UpsertPost
		err   error
	)
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, fhttp.NewErrorResponse(
			http.StatusUnprocessableEntity,
			sharedErrs.ErrKindValidation.String(),
			fmt.Sprintf("Invalid request body: %v", err))
	}

	input.ID, err = strconv.ParseUint(mux.Vars(r)["post_id"], 10, 64)
	if err != nil {
		return nil, fhttp.NewErrorResponse(
			http.StatusBadRequest,
			sharedErrs.ErrKindValidation.String(),
			"Please provide the correct post_id as integer")
	}

	if err = util.Validate(input); err != nil {
		return nil, err
	}

	result, err := c.post.Update(ctx, &input)
	if err != nil {
		return nil, err
	}

	return &fhttp.Response{
		Data:    result,
		Message: "Post is updated successfully.",
		Status:  http.StatusOK}, nil
}

func (c *Controller) Delete(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	postID, err := strconv.ParseUint(mux.Vars(r)["post_id"], 10, 64)
	if err != nil {
		return nil, fhttp.NewErrorResponse(
			http.StatusBadRequest,
			sharedErrs.ErrKindValidation.String(),
			"Please provide the correct post_id as integer")
	}

	if err = c.post.Delete(ctx, postID); err != nil {
		return nil, err
	}

	return &fhttp.Response{
		Data:   "Post is deleted successfully.",
		Status: http.StatusOK,
	}, nil
}

func (c *Controller) Publish(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	postID, err := strconv.ParseUint(mux.Vars(r)["post_id"], 10, 64)
	if err != nil {
		return nil, fhttp.NewErrorResponse(
			http.StatusBadRequest,
			sharedErrs.ErrKindValidation.String(),
			"Please provide the correct post_id as integer")
	}

	if err = c.post.UpdateStatus(ctx, postID, enums.PUBLISHEDPostStatus); err != nil {
		return nil, err
	}

	return &fhttp.Response{
		Message: "Post is published successfully.",
		Status:  http.StatusOK,
	}, nil
}

func (c *Controller) Draft(r *http.Request) (*fhttp.Response, error) {
	ctx := r.Context()

	postID, err := strconv.ParseUint(mux.Vars(r)["post_id"], 10, 64)
	if err != nil {
		return nil, fhttp.NewErrorResponse(
			http.StatusBadRequest,
			sharedErrs.ErrKindValidation.String(),
			"Please provide the correct post_id as integer")
	}

	if err = c.post.UpdateStatus(ctx, postID, enums.DRAFTPostStatus); err != nil {
		return nil, err
	}

	return &fhttp.Response{
		Message: "Post is drafted successfully.",
		Status:  http.StatusOK,
	}, nil
}
