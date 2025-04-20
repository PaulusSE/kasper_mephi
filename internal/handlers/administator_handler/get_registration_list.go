package administator_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type requestDTO struct {
	RequestID uuid.UUID `json:"request_id"`
	Email     string    `json:"email"`
	UserType  string    `json:"user_type"`
	CreatedAt time.Time `json:"created_at"`
}

// ListRegistrationRequests - список заявок на регистрацию
func (h *AdministratorHandler) ListRegistrationRequests(ctx *gin.Context) {
	user, err := h.authenticate(ctx) // или AuthenticateWithUserType Admin
	if err != nil {
		ctx.AbortWithStatus(401)
		return
	}
	requests, err := h.registration.ListRequests(ctx.Request.Context(), user.UserType.String())
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"error": "db error"})
		return
	}
	out := make([]requestDTO, len(requests))
	for i, r := range requests {
		out[i] = requestDTO{
			RequestID: r.RequestID,
			Email:     r.Email,
			UserType:  r.UserType,
			CreatedAt: r.CreatedAt,
		}
	}
	ctx.JSON(http.StatusOK, out)
}
