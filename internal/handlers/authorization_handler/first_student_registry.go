package authorization_handler

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	"uir_draft/internal/handlers/authorization_handler/request_models"
	"uir_draft/internal/pkg/helpers"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
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
	//user, err := h.authenticateStudent(ctx)
	//if err != nil {
	//	ctx.AbortWithError(models.MapErrorToCode(err), err)
	//	return
	//}
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
	newUser := model.Users{
		UserID:     uuid.New(),
		Email:      req.Email,
		Password:   string(hashed),
		KasperID:   uuid.New(),
		UserType:   model.UserType_Student,
		Registered: true,
	}
	if err := h.authenticator.CreateUser(ctx.Request.Context(), &newUser); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "cannot create user"})
		return
	}
	if err := h.authenticator.AttachTokenToUser(ctx.Request.Context(), anonToken, newUser.UserID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "cannot activate token"})
		return
	}
	if err = h.student.InitStudent(ctx, newUser, req); err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": anonToken})
}
