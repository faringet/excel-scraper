package http

import (
	"encoding/json"
	"excel-scraper/pkg/excelreader"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EmployeeController struct {
	ExcelReader *excelreader.ExcelReader
	logger      *zap.Logger
}

func NewConfigController(logger *zap.Logger, excelReader *excelreader.ExcelReader) *EmployeeController {
	return &EmployeeController{logger: logger, ExcelReader: excelReader}
}

func (ec *EmployeeController) GetScanned(c *gin.Context) {
	data, err := ec.ExcelReader.ReadAllExcelFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading Excel files"})
		ec.logger.Error("Error reading Excel files")
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshalling JSON data"})
		ec.logger.Error("Error marshalling JSON data")
		return
	}

	c.Data(http.StatusOK, "application/json", jsonData)
	ec.logger.Info("config sent successfully")
}
