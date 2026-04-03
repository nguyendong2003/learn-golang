package handler

import (
	"net/http"

	"elearning-api/apperror"
	"elearning-api/dto"
	"elearning-api/service"
	"elearning-api/util"

	"github.com/gin-gonic/gin"
)

type FileUploadHandler interface {
	UploadImage() gin.HandlerFunc
	PresignUploadURL() gin.HandlerFunc
}

type fileUploadHandler struct {
	uploadService service.UploadService
}

func NewFileUploadHandler(
	uploadService service.UploadService,
) FileUploadHandler {
	return &fileUploadHandler{
		uploadService: uploadService,
	}
}

// UploadImage godoc
// @Summary Upload an image
// @Description Upload an image file to MinIO storage
// @Tags file-upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file"
// @Success 200 {object} dto.FileUploadResponse
// @Failure 400 {object} any
// @Failure 500 {object} any
// @Router /api/v1/upload/image [post]
// @Security BearerAuth
func (h *fileUploadHandler) UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("No file provided"))
			return
		}

		url, err := h.uploadService.UploadImage(c.Request.Context(), file)
		if err != nil {
			_ = c.Error(err)
			return
		}
		resp := dto.NewApiResponse(c)
		resp.Data = dto.FileUploadResponse{
			URL:      url,
			Filename: file.Filename,
			Size:     file.Size,
		}
		c.JSON(http.StatusOK, resp)
	}
}

// PresignUploadURL godoc
// @Summary Get a presigned URL for uploading a file
// @Description Get a presigned URL for uploading a file to MinIO storage
// @Tags file-upload
// @Accept json
// @Produce json
// @Param request body dto.PresignUploadRequest true "Presign upload URL request"
// @Success 200 {object} dto.PresignUploadURLResponse
// @Failure 400 {object} any
// @Failure 500 {object} any
// @Router /api/v1/upload/presign [post]
// @Security BearerAuth
func (h *fileUploadHandler) PresignUploadURL() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.PresignUploadRequest
		if err := util.BindAndValidateJSON(c, &req); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid request body"))
			return
		}

		url, err := h.uploadService.PresignUploadURL(c.Request.Context(), req.Filename, req.Filetype)
		if err != nil {
			_ = c.Error(err)
			return
		}

		resp := dto.NewApiResponse(c)
		resp.Data = dto.PresignUploadURLResponse{
			URL: url,
		}
		c.JSON(http.StatusOK, resp)
	}
}
