package handlers

import (
	services"Zhooze/pkg/usecase"
	"Zhooze/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)
type ImageHandler struct {
	ImageUseCase services.ImageUseCase
}

func NewImageHandler(useCase services.ImageUseCase) *ImageHandler {
	return &ImageHandler{
		ImageUseCase: useCase,
	}
}
