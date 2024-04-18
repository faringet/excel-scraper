package scheduler

import (
	"excel-scraper/config"
	"excel-scraper/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"net/http/httptest"
	"time"
)

type Scheduler struct {
	thirdPartController *http.ThirdPartController
	updateTime          int
	checkTime           int
	logger              *zap.Logger
}

func NewScheduler(cfg *config.Config, thirdPartController *http.ThirdPartController, logger *zap.Logger) *Scheduler {
	return &Scheduler{thirdPartController: thirdPartController, updateTime: cfg.Scheduler.Update, checkTime: cfg.Scheduler.Update, logger: logger}
}

func (r *Scheduler) Run() {
	s := gocron.NewScheduler(time.UTC)

	interval := time.Duration(r.updateTime) * time.Second
	_, err := s.Every(interval).WaitForSchedule().Do(r.callSendData)
	if err != nil {
		r.logger.Error("Error scheduling InsertUpdatedForecastFromWeatherData", zap.Error(err))
		return
	}

	s.StartAsync()

	r.logger.Info("Scheduler started successfully")
}

func (r *Scheduler) callSendData() {
	// тут фейковый http запрос
	req := httptest.NewRequest("GET", "/dummy", nil)
	w := httptest.NewRecorder()

	// для контекста Gin
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	// чтобы sendData отработал
	r.thirdPartController.SendData(ctx)
}
