package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"excel-scraper/config"
	"excel-scraper/internal/domain"
	"excel-scraper/pkg/excelreader"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"excel-scraper/transport/http/models"
)

type ThirdPartController struct {
	ExcelReader *excelreader.ExcelReader
	logger      *zap.Logger
	receiverURL string
}

func NewThirdPartController(cfg *config.Config, logger *zap.Logger, excelReader *excelreader.ExcelReader) *ThirdPartController {
	return &ThirdPartController{logger: logger, ExcelReader: excelReader, receiverURL: cfg.PlatformAPI.URL}
}

func (tpc *ThirdPartController) SendData(c *gin.Context) {
	data, err := tpc.ExcelReader.ReadAllExcelFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading Excel files"})
		tpc.logger.Error("Error reading Excel files")
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshalling JSON data"})
		tpc.logger.Error("Error marshalling JSON data")
		return
	}

	employees := toEmployee(data)

	if err := tpc.sendDataToReceiver(jsonData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending data to receiver"})
		tpc.logger.Error("Error sending data to receiver", zap.Error(err))
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data sent successfully", "employees": employees})
	tpc.logger.Info("Data sent successfully")
}

func toEmployee(data []*domain.Employees) []models.Employee {
	var employees []models.Employee
	for _, d := range data {
		employee := models.Employee{
			FullName:         d.FullName,
			Preferred:        d.Preferred,
			Email:            d.Email,
			UniqueIdentifier: d.UniqueIdentifier,
			ManagersEmail:    d.ManagersEmail,
			StartDate:        d.StartDate,
			Tenure:           d.Tenure,
			Language:         d.Language,
		}
		employees = append(employees, employee)
	}
	return employees
}

func (tpc *ThirdPartController) sendDataToReceiver(data []byte) error {
	resp, err := http.Post(tpc.receiverURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected status code from receiver")
	}

	return nil
}
