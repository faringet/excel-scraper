package http

import (
	"excel-scraper/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouterImpl struct {
	employeeController  *EmployeeController
	thirdPartController *ThirdPartController
	server              *gin.Engine
	logger              *zap.Logger
	url                 string
}

func NewRouter(cfg *config.Config, logger *zap.Logger, employeeController *EmployeeController, thirdPartController *ThirdPartController) *RouterImpl {
	return &RouterImpl{url: cfg.LocalURL, employeeController: employeeController, thirdPartController: thirdPartController, logger: logger}
}

func (r *RouterImpl) RegisterRoutes() {
	router := gin.Default()

	router.GET("/employee", func(c *gin.Context) {
		r.employeeController.GetScanned(c)
	})

	router.GET("/forwardEmployee", func(c *gin.Context) {
		r.thirdPartController.SendData(c)
	})

	r.server = router
}

func (r *RouterImpl) Start() error {
	return r.server.Run(r.url)
}
