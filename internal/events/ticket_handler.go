package events

import (
	"net/http"
	"strconv"

	"eventhub/internal/models"
	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	service *TicketService
}

func NewTicketHandler(service *TicketService) *TicketHandler {
	return &TicketHandler{service: service}
}

func (h *TicketHandler) BookTicket(c *gin.Context) {
	var req models.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ticket, err := h.service.BookTicket(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ticket)
}

func (h *TicketHandler) GetMyTickets(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}
	tickets, err := h.service.GetTicketsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)
}