package events

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	service *EventService
}

func NewEventHandler(service *EventService) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) GetEvents(c *gin.Context) {
	events, err := h.service.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}
