package ocr

import (
	"net/http"
)

type Service struct {
	// OCR service dependencies
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ProcessUpload(r *http.Request) (interface{}, error) {
	// TODO: Move OCR logic from document_processor.go here
	return map[string]interface{}{
		"status": "processing",
	}, nil
}
