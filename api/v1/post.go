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
// @Router /posts [post]
// @Summary Create a post
// @Description Create a post
// @Tags post
// @Accept json
// @Produce json
// @Param post body models.CreatePostRequest true "Post"
// @Success 201 {object} models.Post
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreatePost(ctx *gin.Context) {

	var req models.CreatePostRequest

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

	resp, err := h.storage.Post().Create(&repo.Post{
		Title:       req.Title,
		Description: req.Description,
		ImageUrl:    req.ImageUrl,
		UserID:      payload.UserID,
		CategoryID:  req.CategoryID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, parsePostToModel(resp))
}

// @Router /posts/{id} [get]
// @Summary Get a post by id
// @Description Get a post by id
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.Post
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetPost(ctx *gin.Context) {

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.Post().Get(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	post := parsePostToModel(resp)

	likeInfo, err := h.storage.Like().GetLikesDislikesCount(post.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	post.LikeInfo = &models.PostLikeInfo{
		LikesCount:    likeInfo.LikesCount,
		DislikesCount: likeInfo.DislikesCount,
	}

	ctx.JSON(http.StatusOK, post)
}

func validateGetPostsParams(ctx *gin.Context) (*models.GetPostsParams, error) {
	var (
		limit      int64 = 10
		page       int64 = 1
		err        error
		userID     int64
		categoryID int64
		sortByDate string = "desc"
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

	if ctx.Query("user_id") != "" {
		userID, err = strconv.ParseInt(ctx.Query("user_id"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Query("category_id") != "" {
		categoryID, err = strconv.ParseInt(ctx.Query("category_id"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Query("sort_by_date") != "" &&
		(ctx.Query("sort_by_date") == "asc" || ctx.Query("sort_by_date") == "desc") {
		sortByDate = ctx.Query("sort_by_date")
	}

	return &models.GetPostsParams{
		Limit:      int32(limit),
		Page:       int32(page),
		Search:     ctx.Query("search"),
		UserID:     userID,
		CategoryID: categoryID,
		SortByDate: sortByDate,
	}, nil
}

// @Router /posts [get]
// @Summary Get posts
// @Description Get posts
// @Tags post
// @Accept json
// @Produce json
// @Param filter query models.GetPostsParams false "Filter"
// @Success 200 {object} models.GetPostsResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetPosts(ctx *gin.Context) {
	request, err := validateGetPostsParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.storage.Post().GetAll(&repo.GetPostsParams{
		Limit:      request.Limit,
		Page:       request.Page,
		Search:     request.Search,
		UserID:     request.UserID,
		CategoryID: request.CategoryID,
		SortByDate: request.SortByDate,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response, err := getPostsResponse(h, result)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func getPostsResponse(h *handlerV1, data *repo.GetPostsResult) (*models.GetPostsResponse, error) {
	response := models.GetPostsResponse{
		Posts: make([]*models.Post, 0),
		Count: data.Count,
	}

	for _, post := range data.Posts {
		p := parsePostToModel(post)

		likeInfo, err := h.storage.Like().GetLikesDislikesCount(p.ID)
		if err != nil {
			return nil, err
		}

		p.LikeInfo = &models.PostLikeInfo{
			LikesCount:    likeInfo.LikesCount,
			DislikesCount: likeInfo.DislikesCount,
		}

		response.Posts = append(response.Posts, &p)
	}

	return &response, nil
}

// @Security ApiKeyAuth
// @Router /posts/{id} [put]
// @Summary Update a post
// @Description Update a post
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param post body models.CreatePostRequest true "Post"
// @Success 200 {object} models.OKResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) UpdatePost(ctx *gin.Context) {
	var req models.CreatePostRequest

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

	err = h.storage.Post().Update(&repo.Post{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		ImageUrl:    req.ImageUrl,
		CategoryID:  req.CategoryID,
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
// @Router /posts/{id} [delete]
// @Summary Delete a post
// @Description Delete a post
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.OKResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) DeletePost(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = h.storage.Post().Delete(id)
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

func parsePostToModel(post *repo.Post) models.Post {
	return models.Post{
		ID:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		ImageUrl:    post.ImageUrl,
		UserID:      post.UserID,
		CategoryID:  post.CategoryID,
		CreatedAt:   post.CreatedAt,
		ViewsCount:  post.ViewsCount,
	}
}
