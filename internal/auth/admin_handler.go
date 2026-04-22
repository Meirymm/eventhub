package auth

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	repo *UserRepository
}

func NewAdminHandler(repo *UserRepository) *AdminHandler {
	return &AdminHandler{repo: repo}
}

func (h *AdminHandler) GetUsers(c *gin.Context) {
	users, err := h.repo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
