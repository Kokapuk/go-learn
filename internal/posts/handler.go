package posts

import (
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

func (h *Handler) CreatePost(c *gin.Context) {
	var req createPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if message, hasError := validation.MakeValidationErrorMessage(err); hasError {
			c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: message})
			return
		}

		c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "Invalid JSON payload"})
		return
	}

	params := createPostParams{Title: req.Title, Content: req.Content, AuthorID: c.GetString("userID")}
	post, err := h.service.createPost(c, params)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *Handler) GetPosts(c *gin.Context) {
	var params listParams
	if err := c.ShouldBindQuery(&params); err != nil {
		if message, hasError := validation.MakeValidationErrorMessage(err); hasError {
			c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: message})
			return
		}

		c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "Invalid Params"})
		return
	}

	posts, err := h.service.getPosts(c, params)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *Handler) GetOwningPosts(c *gin.Context) {
	var params listParams
	if err := c.ShouldBindQuery(&params); err != nil {
		if message, hasError := validation.MakeValidationErrorMessage(err); hasError {
			c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: message})
			return
		}

		c.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "Invalid Params"})
		return
	}

	posts, err := h.service.getPostsByAuthorID(c, params, c.GetString("userID"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, posts)
}
