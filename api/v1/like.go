package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ibrat-muslim/booking-service/api/models"
	"github.com/ibrat-muslim/booking-service/storage/repo"
)

// @Security ApiKeyAuth
// @Router /likes [post]
// @Summary Create or update like
// @Description Create or update like
// @Tags like
// @Accept json
// @Produce json
// @Param like body models.CreateOrUpdateLikeRequest true "Like"
// @Success 200 {object} models.Like
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateOrUpdateLike(ctx *gin.Context) {

	var req models.CreateOrUpdateLikeRequest

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

	err = h.storage.Like().CreateOrUpdate(&repo.Like{
		PostID: req.PostID,
		UserID: payload.UserID,
		Status: req.Status,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.OKResponse{
		Message: "Successfully finished",
	})
}

// @Security ApiKeyAuth
// @Router /likes/user-post [get]
// @Summary Get like by user and post
// @Description Get like by user and post
// @Tags like
// @Accept json
// @Produce json
// @Param post_id query int true "Post ID"
// @Success 200 {object} models.Like
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetLike(ctx *gin.Context) {

	postID, err := strconv.ParseInt(ctx.Query("post_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp, err := h.storage.Like().Get(postID, payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.Like{
		ID:     resp.ID,
		PostID: resp.PostID,
		UserID: resp.UserID,
		Status: resp.Status,
	})
}
