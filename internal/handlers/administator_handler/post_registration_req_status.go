package administator_handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	"uir_draft/internal/handlers/authorization_handler/request_models"
)

func (h *AdministratorHandler) ReviewRegistrationRequest(ctx *gin.Context) {
	_, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithStatus(401)
		return
	}
	var body struct {
		RequestID string `json:"request_id" binding:"required,uuid"`
		Approve   bool   `json:"approve"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	rid, _ := uuid.Parse(body.RequestID)
	if err := h.registration.ReviewRequest(ctx.Request.Context(), rid, body.Approve); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "cannot update status"})
		return
	}

	if !body.Approve {
		ctx.Status(http.StatusOK)
		return
	}

	// Получаем данные самой заявки
	rr, err := h.registration.GetRequestByID(ctx.Request.Context(), rid)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "cannot load request"})
		return
	}

	newUser := model.Users{
		UserID:     uuid.New(),
		Email:      rr.Email,
		Password:   rr.PasswordHash,
		KasperID:   uuid.New(),
		UserType:   model.UserType(rr.UserType),
		Registered: true,
	}
	if err := h.authenticator.CreateUser(ctx.Request.Context(), &newUser); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "cannot create user"})
		return
	}

	// Привязываем токен (анонимный) к новому user TODO: пока токен не храним
	//if err := h.authenticator.AttachTokenToUser(ctx.Request.Context(), rr.TokenNumber, newUser.UserID); err != nil {
	//	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "cannot activate token"})
	//	return
	//}

	//от типа заявки инициализируем студента или руководителя
	switch rr.UserType {
	case model.UserType_Student.String():
		// распарсим payload в request_models.FirstStudentRegistry
		var req request_models.FirstStudentRegistry
		if err := json.Unmarshal(rr.Payload, &req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "bad payload"})
			return
		}
		if err := h.student.InitStudent(ctx, newUser, req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	case model.UserType_Supervisor.String():
		var req request_models.FirstSupervisorRegistry
		if err := json.Unmarshal(rr.Payload, &req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "bad payload"})
			return
		}
		if err := h.supervisor.InitSupervisor(ctx, newUser, req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.Status(http.StatusOK)
}
