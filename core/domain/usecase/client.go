package usecase

import (
	"log"

	"github.com/ellofae/deanery-gateway/core/domain"
	"github.com/ellofae/deanery-gateway/pkg/logger"
)

type ClientUsecase struct {
	logger *log.Logger
}

func NewClientUsecase() domain.IClientUsecase {
	return &ClientUsecase{
		logger: logger.GetLogger(),
	}
}
