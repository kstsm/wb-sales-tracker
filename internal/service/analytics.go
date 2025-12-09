package service

import (
	"context"

	"github.com/kstsm/wb-sales-tracker/internal/dto"
)

func (s *Service) GetAnalytics(ctx context.Context, req dto.AnalyticsRequest) (*dto.AnalyticsResponse, error) {
	return s.repo.GetAnalytics(ctx, req)
}
