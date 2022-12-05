package v1

import (
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type File struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

// @Security ApiKeyAuth
// @Router /file-upload [post]
// @Summary File upload
// @Description File upload
// @Tags file-upload
// @Accept json
// @Produce json
// @Param file formData file true "File"
// @Success 201 {object} models.OKResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) UploadFile(ctx *gin.Context) {
	var file File

	err := ctx.ShouldBind(&file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id := uuid.New()
	fileName := id.String() + filepath.Ext(file.File.Filename)
	dst, _ := os.Getwd()

	if _, err := os.Stat(dst + "/media"); os.IsNotExist(err) {
		os.Mkdir(dst+"/media", os.ModePerm)
	}

	filePath := "/media/" + fileName
	err = ctx.SaveUploadedFile(file.File, dst+filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"filename": filePath,
	})
}
