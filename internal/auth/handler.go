package auth
import (
    "net/http"
    "strconv"
    "eventhub/internal/models"
    "github.com/gin-gonic/gin"
)
type AuthHandler struct {
    service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
    return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := h.service.Register(&req)
    if err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := h.service.Login(&req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
    userID := c.GetInt("user_id")

    user, err := h.service.GetMe(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UpdateMe(c *gin.Context) {
    userID := c.GetInt("user_id")

    var req models.UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.service.UpdateMe(userID, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
        return
    }

    c.JSON(http.StatusOK, user)
}
func (h *AuthHandler) GetAllUsers(c *gin.Context) {
    users, err := h.service.GetAllUsers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, users)
}

func (h *AuthHandler) UpdateUserRole(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    var req struct {
        Role string `json:"role" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if req.Role != "user" && req.Role != "organizer" && req.Role != "admin" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "role must be: user, organizer, admin"})
        return
    }

    if err := h.service.UpdateUserRole(id, req.Role); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "role updated"})
}