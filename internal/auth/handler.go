package auth

import (
	"errors"
	"go-learn/internal/httputil"
	"go-learn/internal/validation"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	service *Service
}

type authResponse struct {
	Token string `json:"token"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SignUp(c *gin.Context) {
	var req signUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if message, hasError := validation.MakeValidationErrorMessage(err); hasError {
			c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: message})
			return
		}

		c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "Invalid JSON payload"})
		return
	}

	user, err := h.service.signUp(c, req)
	if err != nil {
		if errors.Is(err, errUsernameTaken) {
			c.JSON(http.StatusConflict, httputil.ErrorResponse{Message: errUsernameTaken.Error()})
			return
		}

		log.Println(err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Something went wrong"})
		return
	}

	token, err := createToken(user.ID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Something went wrong"})
		return
	}
	response := authResponse{Token: token}
	c.JSON(http.StatusCreated, response)
}

func (h *Handler) SignIn(c *gin.Context) {
	var req signInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if message, hasError := validation.MakeValidationErrorMessage(err); hasError {
			c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: message})
			return
		}

		c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "Invalid JSON payload"})
		return
	}

	user, err := h.service.signIn(c, req)
	if err != nil {
		if errors.Is(err, errInvalidUsernameOrPassword) {
			c.JSON(http.StatusUnauthorized, httputil.ErrorResponse{Message: errInvalidUsernameOrPassword.Error()})
			return
		}

		log.Println(err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Something went wrong"})
		return
	}

	token, err := createToken(user.ID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Something went wrong"})
		return
	}
	response := authResponse{Token: token}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetSelf(c *gin.Context) {
	userID := c.GetString("userID")

	user, err := h.service.getSelf(c, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, httputil.ErrorResponse{Message: "Unauthorized"})
			return
		}

		log.Panicln(err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) RequireAuth(c *gin.Context) {
	header := c.GetHeader("Authorization")
	tokenString, found := strings.CutPrefix(header, "Bearer ")
	if !found || tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httputil.ErrorResponse{Message: "Unauthorized"})
		return
	}

	claims, err := verifyToken(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httputil.ErrorResponse{Message: "Unauthorized"})
		return
	}

	sub, err := claims.GetSubject()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httputil.ErrorResponse{Message: "Unauthorized"})
		return
	}
	c.Set("userID", sub)
	c.Next()
}
