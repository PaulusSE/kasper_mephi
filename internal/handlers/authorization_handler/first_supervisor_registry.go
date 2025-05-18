package authorization_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"uir_draft/internal/handlers/authorization_handler/request_models"
	"uir_draft/internal/pkg/helpers"
)

// FirstSupervisorRegistry
//
//	@Summary		Первичная регистрация научного руководителя
//
//	@Description	Первичная регистрация научного руководителя
//
//	@Tags			NEW
//	@Accept			json
//
//	@Produce		json
//
//	@Success		200
//	@Param			token	path		string									true	"Токен пользователя"
//	@Param			input	body		request_models.FirstSupervisorRegistry	true	"Данные"
//	@Failure		400		{string}	string									"Неверный формат данных"
//	@Failure		401		{string}	string									"Токен протух"
//	@Failure		204		{string}	string									"Нет записей в БД"
//	@Failure		500		{string}	string									"Ошибка на стороне сервера"
//	@Router			/authorize/registration/supervisor/{token} [post]
func (h *AuthorizationHandler) FirstSupervisorRegistry(ctx *gin.Context) {
	if err := h.ValidateToken(ctx); err != nil {
		return
	}
	anonToken := helpers.GetToken(ctx)

	var req request_models.FirstSupervisorRegistry
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
	log.Printf("first_student_registry request body: %v", req)
	log.Printf("email value: %v", lo.FromPtr(&req.Email))
	if err := h.registration.CreateSupervisorRequest(ctx.Request.Context(), hashed, req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "cannot save request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Ваша заявка принята и ожидает подтверждения",
		"token":   anonToken,
	})
}
