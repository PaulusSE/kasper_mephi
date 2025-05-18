package supervisor_handler

import (
	"github.com/gin-gonic/gin"
)

// ListRegistrationRequests - список заявок на регистрацию
func (h *SupervisorHandler) ListRegistrationRequests(ctx *gin.Context) {
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
	ctx.JSON(200, requests)
}
