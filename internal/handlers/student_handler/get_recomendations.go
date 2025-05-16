package student_handler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"log"
	"os/exec"
	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	"uir_draft/internal/handlers/supervisor_handler/request_models"
	"uir_draft/internal/pkg/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

// Recommendation представляет структуру рекомендации
type Recommendation struct {
	Title string `json:"title"`
	Link  string `json:"url"`
}

type RecoDTO struct {
	StudentID uuid.UUID `json:"student_id,omitempty"`
	Semester  int32     `json:"semester,omitempty"`
}

// GetRecommendedArticles
//
//	@Summary		Получение рекомендованных статей
//	@Description	Получение списка рекомендованных статей (пока случайный пример)
//	@Tags			Supervisor.RecommendedArticles
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	[]RecommendedArticle	"Данные"
//	@Param			token	path		string					true	"Токен пользователя"
//	@Failure		400		{string}	string					"Неверный формат данных"
//	@Failure		401		{string}	string					"Токен протух"
//	@Failure		500		{string}	string					"Ошибка на стороне сервера"
//	@Router			/supervisors/recommended-articles/{token} [get]
func (h *StudentHandler) GetRecommendedArticles(ctx *gin.Context) {
	log.Println("=== GetRecommendedArticles start ===")

	user, err := h.authenticate(ctx)
	if err != nil {
		log.Printf("Authentication failed: %v\n", err)
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}
	log.Printf("Authenticated user: %v\n", user)

	req := RecoDTO{}
	if user.UserType == model.UserType_Student {
		reqBody := request_models.DownloadDissertationRequestSupOnlySemester{}
		if err = ctx.ShouldBindJSON(&reqBody); err != nil {
			log.Printf("Failed to bind JSON: %v\n", err)
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		req = RecoDTO{
			StudentID: user.KasperID,
			Semester:  reqBody.Semester,
		}
	} else {
		reqBody := request_models.DownloadDissertationRequestSup{}
		if err = ctx.ShouldBindJSON(&reqBody); err != nil {
			log.Printf("Failed to bind JSON: %v\n", err)
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		req = RecoDTO{
			StudentID: reqBody.StudentID,
			Semester:  reqBody.Semester,
		}
	}

	log.Printf("Request body: studentID=%s, semester=%d\n", req.StudentID, req.Semester)

	dis, err := h.dissertation.GetDissertationData(ctx, req.StudentID, req.Semester)
	if err != nil {
		log.Printf("GetDissertationData error: %v\n", err)
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}
	if dis.FileName == nil {
		log.Println("No dissertation file: returning 204")
		ctx.Status(http.StatusNoContent)
		return
	}
	log.Printf("Found dissertation record: %+v\n", dis)

	filePath := fmt.Sprintf(
		"./dissertations/%s/semester%d/%s",
		dis.StudentID.String(),
		dis.Semester,
		lo.FromPtr(dis.FileName),
	)
	log.Printf("Dissertation file path: %s\n", filePath)

	topN := ctx.DefaultQuery("top_n", "5")
	scriptPath := "./recomendations.py"
	picklePath := "./parsed_articles.pkl"
	args := []string{"-i", picklePath, "-f", filePath, "-n", topN}
	log.Printf("Executing: python3 %s %v\n", scriptPath, args)

	cmd := exec.Command("python3", append([]string{scriptPath}, args...)...)
	cmd.Env = append(cmd.Env, "PYTHONIOENCODING=UTF-8")
	out, err := cmd.CombinedOutput()
	log.Printf("Script stdout+stderr:\n%s\n", string(out))
	if err != nil {
		log.Printf("Script error: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Script error: %v\nOutput:\n%s", err, string(out)),
		})
		return
	}

	var recs []Recommendation
	if err := json.Unmarshal(out, &recs); err != nil {
		log.Printf("Failed to unmarshal JSON: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to parse JSON: %v. Output: %s", err, string(out)),
		})
		return
	}
	log.Printf("Parsed %d recommendations\n", len(recs))

	ctx.JSON(http.StatusOK, recs)
	log.Println("=== GetRecommendedArticles end ===")
}
