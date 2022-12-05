package v1

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ibrat-muslim/booking-service/api/models"
	"github.com/ibrat-muslim/booking-service/pkg/utils"
	"github.com/ibrat-muslim/booking-service/storage/repo"
)

// @Security ApiKeyAuth
// @Router /users [post]
// @Summary Create a user
// @Description Create a user
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User"
// @Success 201 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateUser(ctx *gin.Context) {
	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if payload.UserType != repo.UserTypeSuperAdmin {
		ctx.JSON(http.StatusForbidden, errorResponse(ErrForbidden))
		return
	}

	var req models.CreateUserRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp, err := h.storage.User().Create(&repo.User{
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		DateOfBirth:     req.DateOfBirth,
		Email:           req.Email,
		PhoneNumber:     req.PhoneNumber,
		Gender:          req.Gender,
		Password:        hashedPassword,
		ProfileImageUrl: req.ProfileImageUrl,
		Address:         req.Address,
		Type:            req.Type,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, parseUserToModel(resp))
}

// @Router /users/{id} [get]
// @Summary Get a user by id
// @Description Get a user by id
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetUser(ctx *gin.Context) {

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.User().Get(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, parseUserToModel(resp))
}

// @Security ApiKeyAuth
// @Router /users/me [get]
// @Summary Get a user by token
// @Description Get a user by token
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetUserProfile(ctx *gin.Context) {
	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp, err := h.storage.User().Get(payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, parseUserToModel(resp))
}

// @Router /users [get]
// @Summary Get users
// @Description Get users
// @Tags user
// @Accept json
// @Produce json
// @Param filter query models.GetAllParamsRequest false "Filter"
// @Success 200 {object} models.GetUsersResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetUsers(ctx *gin.Context) {
	request, err := validateGetAllParamsRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.storage.User().GetAll(&repo.GetUsersParams{
		Limit:  request.Limit,
		Page:   request.Page,
		Search: request.Search,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, getUsersResponse(result))
}

func getUsersResponse(data *repo.GetUsersResult) *models.GetUsersResponse {
	response := models.GetUsersResponse{
		Users: make([]*models.User, 0),
		Count: data.Count,
	}

	for _, user := range data.Users {
		u := parseUserToModel(user)
		response.Users = append(response.Users, &u)
	}

	return &response
}

// @Security ApiKeyAuth
// @Router /users/{id} [put]
// @Summary Update a user
// @Description Update a user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param user body models.CreateUserRequest true "User"
// @Success 200 {object} models.OKResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) UpdateUser(ctx *gin.Context) {
	var req models.CreateUserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = h.storage.User().Update(&repo.User{
		ID:              id,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		DateOfBirth:     req.DateOfBirth,
		Email:           req.Email,
		PhoneNumber:     req.PhoneNumber,
		Gender:          req.Gender,
		Password:        hashedPassword,
		ProfileImageUrl: req.ProfileImageUrl,
		Address:         req.Address,
		Type:            req.Type,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.OKResponse{
		Message: "successfully updated",
	})
}

// @Security ApiKeyAuth
// @Router /users/{id} [delete]
// @Summary Delete a user
// @Description Delete a user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.OKResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = h.storage.User().Delete(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.OKResponse{
		Message: "successfully deleted",
	})
}

func parseUserToModel(user *repo.User) models.User {
	return models.User{
		ID:              user.ID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		DateOfBirth:     user.DateOfBirth,
		Email:           user.Email,
		PhoneNumber:     user.PhoneNumber,
		Gender:          user.Gender,
		ProfileImageUrl: user.ProfileImageUrl,
		Address:         user.Address,
		Type:            user.Type,
		CreatedAt:       user.CreatedAt,
	}
}
