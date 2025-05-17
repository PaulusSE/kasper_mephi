package student_handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	"uir_draft/internal/handlers/supervisor_handler/request_models"
	"uir_draft/internal/pkg/models"
)

// Recommendation для ответа
type Recommendation struct {
	Title string `json:"title"`
	Link  string `json:"url"`
}

type RecommendationProfile struct {
	StudentID       string               `json:"student_id"`
	Dissertation    string               `json:"dissertation"`
	Publications    []models.Publication `json:"publications"`
	Conferences     []models.Conference  `json:"conferences"`
	Patents         []models.Patent      `json:"patents"`
	Comments        []string             `json:"comments"`
	Progressiveness int                  `json:"progressiveness"`
	Semester        int32                `json:"semester"`
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

// GetRecommendedArticles — новая версия, собирающая профиль студента и передающая его в Python
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
			ctx.AbortWithError(400, err)
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
			ctx.AbortWithError(400, err)
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
		ctx.Status(204)
		return
	}
	dissertationPath := fmt.Sprintf(
		"./dissertations/%s/semester%d/%s",
		dis.StudentID.String(),
		dis.Semester,
		lo.FromPtr(dis.FileName),
	)

	// Проверяем существование файла диссертации
	if _, err := os.Stat(dissertationPath); os.IsNotExist(err) {
		log.Printf("Dissertation file does not exist: %s", dissertationPath)
		ctx.JSON(404, gin.H{"error": "Dissertation file not found"})
		dissertationPath = ""
	}

	// Научные работы
	scientificWorks, err := h.scientific.GetScientificWorks(ctx, req.StudentID)
	if err != nil {
		log.Printf("GetScientificWorks error: %v\n", err)
	}
	var pubs []models.Publication
	var confs []models.Conference
	var patents []models.Patent
	for _, sw := range scientificWorks {
		if sw.Semester == int(req.Semester) {
			pubs = sw.Publications
			confs = sw.Conferences
			patents = sw.Patents
			break
		}
	}

	// Прогресс
	progressModel, err := h.dissertation.GetDissertationPage(ctx, req.StudentID)
	if err != nil {
		log.Printf("GetDissertationPage error: %v\n", err)
	}
	progressArr := progressModel.Progresses
	progress := 0
	for _, p := range progressArr {
		if p.Semester == req.Semester {
			progress = int(p.Progressiveness)
		}
	}

	// Комментарии
	reportComments, err := h.report.GetReportComments(ctx, req.StudentID)
	if err != nil {
		log.Printf("GetReportComments error: %v\n", err)
	}
	commentStrs := []string{}
	for _, c := range reportComments.DissertationComments {
		if c.Commentary != nil && *c.Commentary != "" {
			commentStrs = append(commentStrs, *c.Commentary)
		}
	}

	recoProfile := RecommendationProfile{
		StudentID:       req.StudentID.String(),
		Dissertation:    dissertationPath,
		Publications:    pubs,
		Conferences:     confs,
		Patents:         patents,
		Comments:        commentStrs,
		Progressiveness: progress,
		Semester:        req.Semester,
	}

	// Маршалинг профиля
	profileJson, err := json.Marshal(recoProfile)
	if err != nil {
		log.Printf("Failed to marshal RecommendationProfile: %v\n", err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	log.Printf("Передаю в Python:\n%s\n", string(profileJson))

	// Запуск Python-скрипта, читаем отдельно stdout и stderr
	scriptPath := "./reco_model_2.py"
	picklePath := "./parsed_articles.pkl"
	topN := ctx.DefaultQuery("top_n", "5")
	args := []string{"-i", picklePath, "-n", topN}
	pythonPath := "/usr/local/bin/python3"

	topNInt, _ := strconv.Atoi(topN)
	cached, updatedAt, err := h.recommendationCache.GetCachedRecommendations(ctx, req.StudentID, req.Semester, topNInt)
	cacheValid := false
	if err == nil && cached != nil && len(cached) > 0 && updatedAt != nil {
		if time.Since(*updatedAt) < 7*24*time.Hour {
			cacheValid = true
		} else {
			log.Printf("Кэш устарел: recommendations старше недели (updatedAt: %v)", *updatedAt)
		}
	}
	if cacheValid {
		log.Printf("Отдаем рекомендации из кэша (student=%s, semester=%d, top_n=%d)", req.StudentID, req.Semester, topNInt)
		ctx.JSON(200, cached)
		return
	}

	cmd := exec.Command(pythonPath, append([]string{scriptPath}, args...)...)
	cmd.Env = append(cmd.Env, "PYTHONIOENCODING=UTF-8")
	cmd.Stdin = bytes.NewReader(profileJson)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	// Всегда логируем stderr (ворнинги/ошибки Python)
	if stderr.Len() > 0 {
		log.Printf("Script stderr:\n%s\n", stderr.String())
	}

	if err != nil {
		log.Printf("Script error: %v\n", err)
		ctx.JSON(500, gin.H{
			"error": fmt.Sprintf("Script error: %v\nStderr:\n%s", err, stderr.String()),
		})
		return
	}

	// Теперь парсим только stdout
	var recs []Recommendation
	if err := json.Unmarshal(stdout.Bytes(), &recs); err != nil {
		log.Printf("Failed to unmarshal JSON: %v\n", err)
		ctx.JSON(500, gin.H{
			"error": fmt.Sprintf("Failed to parse JSON: %v. Output: %s", err, stdout.String()),
		})
		return
	}

	if err := h.recommendationCache.SaveCachedRecommendations(ctx, req.StudentID, req.Semester, topNInt, recs); err != nil {
		log.Printf("Не удалось сохранить рекомендации в кэш: %v", err)
	}
	log.Printf("Parsed %d recommendations\n", len(recs))
	ctx.JSON(200, recs)
	log.Println("=== GetRecommendedArticles end ===")
}
