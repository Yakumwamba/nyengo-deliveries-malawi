package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"nyengo-deliveries/internal/models"
	"nyengo-deliveries/internal/repository"
)

type OrderService struct {
	repo        *repository.OrderRepository
	courierRepo *repository.CourierRepository
	pricing     *PricingService
}

func NewOrderService(repo *repository.OrderRepository, courierRepo *repository.CourierRepository, pricing *PricingService) *OrderService {
	return &OrderService{repo: repo, courierRepo: courierRepo, pricing: pricing}
}

func (s *OrderService) Create(ctx context.Context, courierID uuid.UUID, req *models.CreateOrderRequest) (*models.Order, error) {
	estimate, err := s.pricing.CalculateEstimate(&models.PriceEstimateRequest{
		PickupLatitude: req.PickupLatitude, PickupLongitude: req.PickupLongitude,
		DeliveryLatitude: req.DeliveryLatitude, DeliveryLongitude: req.DeliveryLongitude,
		PackageWeight: req.PackageWeight, IsFragile: req.IsFragile,
	})
	if err != nil {
		return nil, err
	}

	platformFee, earnings := s.pricing.CalculateCourierEarnings(estimate.TotalFare)

	order := &models.Order{
		CourierID: courierID, StoreID: req.StoreID, ExternalOrderID: req.ExternalOrderID,
		CustomerName: req.CustomerName, CustomerPhone: req.CustomerPhone, CustomerEmail: req.CustomerEmail,
		PickupAddress: req.PickupAddress, PickupLatitude: req.PickupLatitude, PickupLongitude: req.PickupLongitude,
		PickupNotes: req.PickupNotes, PickupContactName: req.PickupContactName, PickupContactPhone: req.PickupContactPhone,
		DeliveryAddress: req.DeliveryAddress, DeliveryLatitude: req.DeliveryLatitude, DeliveryLongitude: req.DeliveryLongitude,
		DeliveryNotes: req.DeliveryNotes, PackageDescription: req.PackageDescription,
		PackageSize: req.PackageSize, PackageWeight: req.PackageWeight,
		IsFragile: req.IsFragile, RequiresSignature: req.RequiresSignature,
		Distance: estimate.Distance, BaseFare: estimate.BaseFare, DistanceFare: estimate.DistanceFare,
		SurgeFare: estimate.SurgeFare, TotalFare: estimate.TotalFare,
		PlatformFee: platformFee, CourierEarnings: earnings,
		PaymentMethod: req.PaymentMethod, ScheduledPickup: req.ScheduledPickup,
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrderService) List(ctx context.Context, courierID uuid.UUID, filters *models.OrderListFilters) (*models.OrderListResponse, error) {
	return s.repo.List(ctx, courierID, filters)
}

func (s *OrderService) UpdateStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error {
	return s.repo.UpdateStatus(ctx, orderID, status)
}

func (s *OrderService) Accept(ctx context.Context, orderID uuid.UUID) error {
	return s.repo.UpdateStatus(ctx, orderID, models.OrderStatusAccepted)
}

func (s *OrderService) Decline(ctx context.Context, orderID uuid.UUID) error {
	return s.repo.UpdateStatus(ctx, orderID, models.OrderStatusDeclined)
}

func (s *OrderService) GetDailyStats(ctx context.Context, courierID uuid.UUID) (map[string]interface{}, error) {
	return s.repo.GetDailyStats(ctx, courierID, time.Now())
}

func (s *OrderService) GetMonthlyStats(ctx context.Context, courierID uuid.UUID) (map[string]interface{}, error) {
	now := time.Now()
	return s.repo.GetMonthlyStats(ctx, courierID, now.Year(), int(now.Month()))
}
