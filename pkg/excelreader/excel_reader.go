package excelreader

import (
	"excel-scraper/config"
	"excel-scraper/internal/domain"
	"excel-scraper/pkg/dir"
	"github.com/tealeg/xlsx"
	"go.uber.org/zap"
	"path/filepath"
)

type ExcelReader struct {
	Dir    *dir.Files
	logger *zap.Logger
	path   string
}

func NewExcelReader(cfg *config.Config, dir *dir.Files, logger *zap.Logger) *ExcelReader {
	return &ExcelReader{path: cfg.ScanningOpts.Path, Dir: dir, logger: logger}
}

func (er *ExcelReader) ReadAllExcelFiles() ([]*domain.Employees, error) {
	files, err := er.Dir.List(er.path)
	if err != nil {
		er.logger.Error("Error listing files", zap.Error(err))
		return nil, err
	}

	var employees []*domain.Employees

	for _, file := range files {
		if filepath.Ext(file) == ".xlsx" {
			err := er.readSpecificExcelFile(file, &employees)
			if err != nil {
				er.logger.Error("Error reading Excel file", zap.String("file", file), zap.Error(err))
				continue
			}
		}
	}

	return employees, nil
}

func (er *ExcelReader) readSpecificExcelFile(filePath string, employees *[]*domain.Employees) error {
	file, err := xlsx.OpenFile(filePath)
	if err != nil {
		er.logger.Error("Error opening Excel file", zap.String("file", filePath), zap.Error(err))
		return err
	}

	for _, sheet := range file.Sheets {
		isFirstRow := true
		for _, row := range sheet.Rows {
			if isFirstRow {
				isFirstRow = false
				continue
			}
			employee := mapToEmployee(row.Cells)
			*employees = append(*employees, employee)
		}
	}

	return nil
}

func mapToEmployee(cells []*xlsx.Cell) *domain.Employees {
	employee := &domain.Employees{}
	employee.FullName = cells[0].String()
	employee.Preferred = cells[1].String()
	employee.Email = cells[2].String()
	employee.UniqueIdentifier = cells[3].String()
	employee.ManagersEmail = cells[4].String()
	employee.StartDate = cells[5].String()
	employee.Tenure = cells[6].String()
	employee.Language = cells[7].String()
	return employee
}
