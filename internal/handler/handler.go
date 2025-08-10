package handler

import (
    coreService "abet/internal/service"
)

type Handler struct {
    *coreService.Service
}

func NewHandler(service *coreService.Service) *Handler {
    return &Handler{service}
}
