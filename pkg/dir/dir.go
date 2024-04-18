package dir

import (
	"excel-scraper/config"
	"go.uber.org/zap"
	"io/fs"
	"os"
	"strings"
)

type Files struct {
	FS     fs.FS
	logger *zap.Logger
	format string
}

func NewFilesFS(dirPath string, cfg *config.Config, logger *zap.Logger) *Files {
	fileSystem := os.DirFS(dirPath)
	return &Files{format: cfg.ScanningOpts.Format, FS: fileSystem, logger: logger}
}

func (f *Files) List(dir string) ([]string, error) {
	f.logger.Info("Listing "+f.format+" files in directory", zap.String("directory", dir))

	var fileNames []string
	err := fs.WalkDir(f.FS, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			f.logger.Error("Error walking directory", zap.Error(err))
			return err
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(path), f.format) {
			fileNames = append(fileNames, path)
		}
		return nil
	})
	if err != nil {
		f.logger.Error("Error listing files", zap.Error(err))
		return nil, err
	}

	f.logger.Info("Finished listing "+f.format+" files", zap.Int("count", len(fileNames)))
	return fileNames, nil
}
