package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"nyengo-deliveries/internal/repository"
)

type AnalyticsService struct {
	orderRepo   *repository.OrderRepository
	courierRepo *repository.CourierRepository
}

func NewAnalyticsService(orderRepo *repository.OrderRepository, courierRepo *repository.CourierRepository) *AnalyticsService {
	return &AnalyticsService{orderRepo: orderRepo, courierRepo: courierRepo}
}

type DashboardMetrics struct {
	Today     *DayMetrics   `json:"today"`
	ThisWeek  *WeekMetrics  `json:"thisWeek"`
	ThisMonth *MonthMetrics `json:"thisMonth"`
}

type DayMetrics struct {
	TotalOrders   int     `json:"totalOrders"`
	Completed     int     `json:"completed"`
	Pending       int     `json:"pending"`
	TotalRevenue  float64 `json:"totalRevenue"`
	TotalEarnings float64 `json:"totalEarnings"`
}

type WeekMetrics struct {
	TotalOrders   int     `json:"totalOrders"`
	TotalRevenue  float64 `json:"totalRevenue"`
	TotalEarnings float64 `json:"totalEarnings"`
	AvgPerDay     float64 `json:"avgPerDay"`
}

type MonthMetrics struct {
	TotalOrders   int     `json:"totalOrders"`
	Completed     int     `json:"completed"`
	TotalRevenue  float64 `json:"totalRevenue"`
	TotalEarnings float64 `json:"totalEarnings"`
	AvgRating     float64 `json:"avgRating"`
}

func (s *AnalyticsService) GetDashboardMetrics(ctx context.Context, courierID uuid.UUID) (*DashboardMetrics, error) {
	todayStats, err := s.orderRepo.GetDailyStats(ctx, courierID, time.Now())
	if err != nil {
		return nil, err
	}

	now := time.Now()
	monthStats, err := s.orderRepo.GetMonthlyStats(ctx, courierID, now.Year(), int(now.Month()))
	if err != nil {
		return nil, err
	}

	return &DashboardMetrics{
		Today: &DayMetrics{
			TotalOrders:   todayStats["totalOrders"].(int),
			Completed:     todayStats["completed"].(int),
			Pending:       todayStats["pending"].(int),
			TotalRevenue:  todayStats["totalRevenue"].(float64),
			TotalEarnings: todayStats["totalEarnings"].(float64),
		},
		ThisMonth: &MonthMetrics{
			TotalOrders:   monthStats["totalOrders"].(int),
			Completed:     monthStats["completed"].(int),
			TotalRevenue:  monthStats["totalRevenue"].(float64),
			TotalEarnings: monthStats["totalEarnings"].(float64),
			AvgRating:     monthStats["avgRating"].(float64),
		},
	}, nil
}
