package authorization_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *AuthorizationHandler) CreateAnonymousToken(ctx *gin.Context) {
	// генерим новый UUID
	token := uuid.New().String()

	// запись в БД: authorization_token { token_number: token, is_active: false, user_id: NULL }
	if err := h.authenticator.CreateAnonymousToken(ctx, token); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
