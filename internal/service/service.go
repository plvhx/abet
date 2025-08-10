package service

import (
    "abet/pkg"
    coreRepository "abet/internal/repository"
)

type Service struct {
    *pkg.Options
    *coreRepository.Repository
}

func NewService(options *pkg.Options, repository *coreRepository.Repository) *Service {
    return &Service{options, repository}
}
