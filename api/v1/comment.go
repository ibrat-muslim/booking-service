package v1

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ibrat-muslim/booking-service/api/models"
	"github.com/ibrat-muslim/booking-service/storage/repo"
)

// @Security ApiKeyAuth
// @Router /comments [post]
// @Summary Create a comment
// @Description Create a comment
// @Tags comment
// @Accept json
// @Produce json
// @Param comment body models.CreateCommentRequest true "Comment"
// @Success 201 {object} models.Comment
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateComment(ctx *gin.Context) {

	var req models.CreateCommentRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp, err := h.storage.Comment().Create(&repo.Comment{
		PostID:      req.PostID,
		UserID:      payload.UserID,
		Description: req.Description,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, parseCommentToModel(resp))
}

func validateGetCommentsParams(ctx *gin.Context) (*models.GetCommentsParams, error) {
	var (
		limit  int64 = 10
		page   int64 = 1
		postID int64
		userID int64
		err    error
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

	if ctx.Query("post_id") != "" {
		postID, err = strconv.ParseInt(ctx.Query("post_id"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Query("user_id") != "" {
		userID, err = strconv.ParseInt(ctx.Query("user_id"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	return &models.GetCommentsParams{
		Limit:  int32(limit),
		Page:   int32(page),
		PostID: postID,
		UserID: userID,
	}, nil
}

// @Router /comments [get]
// @Summary Get comments
// @Description Get comments
// @Tags comment
// @Accept json
// @Produce json
// @Param filter query models.GetCommentsParams false "Filter"
// @Success 200 {object} models.GetCommentsResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetComments(ctx *gin.Context) {
	request, err := validateGetCommentsParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.storage.Comment().GetAll(&repo.GetCommentsParams{
		Limit:  request.Limit,
		Page:   request.Page,
		PostID: request.PostID,
		UserID: request.UserID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, getCommentsResponse(result))
}

func getCommentsResponse(data *repo.GetCommentsResult) *models.GetCommentsResponse {
	response := models.GetCommentsResponse{
		Comments: make([]*models.Comment, 0),
		Count:    data.Count,
	}

	for _, comment := range data.Comments {
		c := parseCommentToModel(comment)

		c.User = &models.CommentUser{
			ID:              comment.UserID,
			FirstName:       comment.User.FirstName,
			LastName:        comment.User.LastName,
			Email:           comment.User.Email,
			ProfileImageUrl: comment.User.ProfileImageUrl,
		}

		response.Comments = append(response.Comments, &c)
	}

	return &response
}

// @Security ApiKeyAuth
// @Router /comments/{id} [put]
// @Summary Update a comment
// @Description Update a comment
// @Tags comment
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param comment body models.CreateCommentRequest true "Comment"
// @Success 200 {object} models.OKResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) UpdateComment(ctx *gin.Context) {
	var req models.CreateCommentRequest

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

	updatedAt := time.Now()

	err = h.storage.Comment().Update(&repo.Comment{
		ID:          id,
		Description: req.Description,
		UpdatedAt:   &updatedAt,
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
// @Router /comments/{id} [delete]
// @Summary Delete a comment
// @Description Delete a comment
// @Tags comment
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.OKResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) DeleteComment(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = h.storage.Comment().Delete(id)
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

func parseCommentToModel(comment *repo.Comment) models.Comment {
	return models.Comment{
		ID:          comment.ID,
		PostID:      comment.PostID,
		UserID:      comment.UserID,
		Description: comment.Description,
		CreatedAt:   comment.CreatedAt,
		UpdatedAt:   comment.UpdatedAt,
	}
}
