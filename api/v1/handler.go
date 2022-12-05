package v1

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ibrat-muslim/booking-service/api/models"
	"github.com/ibrat-muslim/booking-service/config"
	"github.com/ibrat-muslim/booking-service/storage"
)

var (
	ErrWrongEmailOrPass = errors.New("wrong email or password")
	ErrUserNotVerified  = errors.New("user not verified")
	ErrEmailExists      = errors.New("email already exists")
	ErrIncorrectCode    = errors.New("incorrect verification code")
	ErrCodeExpired      = errors.New("verification code has been expired")
	ErrForbidden        = errors.New("forbidden")
)

type handlerV1 struct {
	cfg      *config.Config
	storage  storage.StorageI
	inMemory storage.InMemoryStorageI
}

type HandlerV1Options struct {
	Cfg      *config.Config
	Storage  storage.StorageI
	InMemory storage.InMemoryStorageI
}

func New(options *HandlerV1Options) *handlerV1 {
	return &handlerV1{
		cfg:      options.Cfg,
		storage:  options.Storage,
		inMemory: options.InMemory,
	}
}

func errorResponse(err error) *models.ErrorResponse {
	return &models.ErrorResponse{
		Error: err.Error(),
	}
}

func validateGetAllParamsRequest(ctx *gin.Context) (*models.GetAllParamsRequest, error) {
	var (
		limit int64 = 10
		page  int64 = 1
		err   error
	)

	if ctx.Query("limit") != "" {
		limit, err = strconv.ParseInt(ctx.Query("limit"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Query("page") != "" {
		page, err = strconv.ParseInt(ctx.Query("page"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	return &models.GetAllParamsRequest{
		Limit:  int32(limit),
		Page:   int32(page),
		Search: ctx.Query("search"),
	}, nil
}
