package jobs

import (
	"errors"
	"go-learn/internal/httputil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) EnqueueJob(c *gin.Context) {
	jobID, err := h.service.enqueueJob()
	if err != nil {
		if errors.Is(err, errQueueFull) {
			c.JSON(http.StatusServiceUnavailable, httputil.ErrorResponse{Message: err.Error()})
			return
		}

		log.Panicln(err)
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Something went wrong"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"jobId": jobID})
}

func (h *Handler) GetJobStats(c *gin.Context) {
	jobID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, httputil.ErrorResponse{Message: "Invalid job id"})
		return
	}

	jobStatus, err := h.service.getJobStatus(jobID)
	if err != nil {
		if errors.Is(err, errJobDoesNotExist) {
			c.JSON(http.StatusNotFound, httputil.ErrorResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jobStatus": jobStatus})

}
