package v1

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ibrat-muslim/booking-service/api/models"
	"github.com/ibrat-muslim/booking-service/storage/repo"
)

// @Security ApiKeyAuth
// @Router /categories [post]
// @Summary Create a category
// @Description Create a category
// @Tags category
// @Accept json
// @Produce json
// @Param category body models.CreateCategoryRequest true "Category"
// @Success 201 {object} models.Category
// @Failure 400 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateCategory(ctx *gin.Context) {

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if payload.UserType != repo.UserTypeSuperAdmin {
		ctx.JSON(http.StatusForbidden, errorResponse(ErrForbidden))
		return
	}

	var req models.CreateCategoryRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.Category().Create(&repo.Category{Title: req.Title})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, models.Category{
		ID:        resp.ID,
		Title:     resp.Title,
		CreatedAt: resp.CreatedAt,
	})
}

// @Router /categories/{id} [get]
// @Summary Get a category by id
// @Description Get a category by id
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.Category
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetCategory(ctx *gin.Context) {

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.Category().Get(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.Category{
		ID:        resp.ID,
		Title:     resp.Title,
		CreatedAt: resp.CreatedAt,
	})
}

// @Router /categories [get]
// @Summary Get categories
// @Description Get categories
// @Tags category
// @Accept json
// @Produce json
// @Param filter query models.GetAllParamsRequest false "Filter"
// @Success 200 {object} models.GetCategoriesResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetCategories(ctx *gin.Context) {
	request, err := validateGetAllParamsRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.storage.Category().GetAll(&repo.GetCategoriesParams{
		Limit:  request.Limit,
		Page:   request.Page,
		Search: request.Search,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, getCategoriesResponse(result))
}

func getCategoriesResponse(data *repo.GetCategoriesResult) *models.GetCategoriesResponse {
	response := models.GetCategoriesResponse{
		Categories: make([]*models.Category, 0),
		Count:      data.Count,
	}

	for _, c := range data.Categories {
		response.Categories = append(response.Categories, &models.Category{
			ID:        c.ID,
			Title:     c.Title,
			CreatedAt: c.CreatedAt,
		})
	}

	return &response
}

// @Security ApiKeyAuth
// @Router /categories/{id} [put]
// @Summary Update a category
// @Description Update a category
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param category body models.CreateCategoryRequest true "Category"
// @Success 200 {object} models.OKResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) UpdateCategory(ctx *gin.Context) {

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if payload.UserType != repo.UserTypeSuperAdmin {
		ctx.JSON(http.StatusForbidden, errorResponse(ErrForbidden))
		return
	}

	var req models.CreateCategoryRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = h.storage.Category().Update(&repo.Category{
		ID:    id,
		Title: req.Title,
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
// @Router /categories/{id} [delete]
// @Summary Delete a category
// @Description Delete a category
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.OKResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) DeleteCategory(ctx *gin.Context) {

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if payload.UserType != repo.UserTypeSuperAdmin {
		ctx.JSON(http.StatusForbidden, errorResponse(ErrForbidden))
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = h.storage.Category().Delete(id)
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
