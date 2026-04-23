package events

import (
	"eventhub/internal/models"
	"net/http"
    "strconv"
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
	req.UserID = c.GetInt("user_id")

	ticket, err := h.service.BookTicket(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ticket)
}

func (h *TicketHandler) GetMyTickets(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
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
func (h *TicketHandler) ValidateQR(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
        return
    }

    ticket, err := h.service.ValidateTicket(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "valid":   false,
            "message": "ticket not found",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "valid":    true,
        "ticket":   ticket.ID,
        "event_id": ticket.EventID,
        "user_id":  ticket.UserID,
        "message":  "ticket is valid",
    })
}
