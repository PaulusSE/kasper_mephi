package authorization_handler

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"uir_draft/internal/handlers/authorization_handler/request_models"
	"uir_draft/internal/pkg/helpers"
)

// FirstStudentRegistry
//
//	@Summary		Первичная регистрация студента
//
//	@Description	Первичная регистрация студента
//
//	@Tags			Authorization
//	@Accept			json
//
//	@Produce		json
//
//	@Success		200
//	@Param			token	path		string								true	"Токен пользователя"
//	@Param			input	body		request_models.FirstStudentRegistry	true	"Данные"
//	@Failure		400		{string}	string								"Неверный формат данных"
//	@Failure		401		{string}	string								"Токен протух"
//	@Failure		204		{string}	string								"Нет записей в БД"
//	@Failure		500		{string}	string								"Ошибка на стороне сервера"
//	@Router			/authorize/registration/student/{token} [post]
func (h *AuthorizationHandler) FirstStudentRegistry(ctx *gin.Context) {
	if err := h.ValidateToken(ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	anonToken := helpers.GetToken(ctx)

	req := request_models.FirstStudentRegistry{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "hash error"})
		return
	}
	if err := h.registration.CreateStudentRequest(ctx.Request.Context(), hashed, req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "cannot save request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Ваша заявка принята и ожидает подтверждения",
		"token":   anonToken,
	})
}
