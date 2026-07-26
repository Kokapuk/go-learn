package auth

import (
	"errors"
	"go-learn/internal/httputil"
	"go-learn/internal/validation"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SignUp(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if message, hasError := validation.MakeValidationErrorMessage(err); hasError {
			c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: message})
			return
		}

		c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "Invalid JSON payload"})
		return
	}

	user, err := h.service.createUser(c, req)
	if err != nil {
		if errors.Is(err, errUsernameTaken) {
			c.JSON(http.StatusConflict, httputil.ErrorResponse{Message: errUsernameTaken.Error()})
			return
		}

		log.Println(err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, user)
}
