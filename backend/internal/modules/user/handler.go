package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"authapp/internal/platform/authn"
)

// Handler serves user HTTP endpoints.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Me godoc
// @Summary      Current user
// @Description  Returns the authenticated user profile
// @Tags         user
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  MeResponse
// @Failure      401  {object}  map[string]string
// @Router       /me [get]
func (h *Handler) Me(c *gin.Context) {
	v, ok := c.Get(authn.ContextUserIDKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	s, ok := v.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id, err := uuid.Parse(s)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	u, err := h.svc.GetByID(id)
	if err != nil || u == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, MeResponse{
		ID:    u.ID.String(),
		Email: u.Email,
	})
}

// MeResponse is returned by GET /me for Swagger.
type MeResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
