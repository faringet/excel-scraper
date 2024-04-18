package main

import (
	"excel-scraper/config"
	"excel-scraper/internal/scheduler"
	"excel-scraper/pkg/dir"
	"excel-scraper/pkg/excelreader"
	"excel-scraper/pkg/zaplogger"
	"excel-scraper/transport/http"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"log"
)

func main() {
	// Viper
	_, cfg, errViper := config.NewViper("conf_local")
	if errViper != nil {
		log.Fatal(errors.WithMessage(errViper, "Viper startup error"))
	}

	// Zap logger
	logger, loggerCleanup, errZapLogger := zaplogger.New(zaplogger.Mode(cfg.Logger.Development))
	if errZapLogger != nil {
		log.Fatal(errors.WithMessage(errZapLogger, "Zap logger startup error"))
	}
	defer loggerCleanup()

	// Working with the file system
	dir := dir.NewFilesFS(cfg.ScanningOpts.Path, cfg, logger)

	// ExcelReader
	excelReader := excelreader.NewExcelReader(cfg, dir, logger)

	// HTTP Controllers
	employeeController := http.NewConfigController(logger, excelReader)
	thirdPartController := http.NewThirdPartController(cfg, logger, excelReader)

	// HTTP router
	router := http.NewRouter(cfg, logger, employeeController, thirdPartController)
	router.RegisterRoutes()

	// Channel for error transmission
	errCh := make(chan error, 1)

	// Router in goroutine
	go func() {
		err := router.Start()
		if err != nil {
			logger.Error("Error starting router", zap.Error(err))
		}
		errCh <- err
	}()

	// Scheduler
	sch := scheduler.NewScheduler(cfg, thirdPartController, logger)
	sch.Run()

	// Router shut down or error
	if err := <-errCh; err != nil {
		logger.Error("Router exited with error", zap.Error(err))
	}
}
