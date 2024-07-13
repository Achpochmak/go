package service

import (
	"HOMEWORK-1/internal/module"
	"context"
	"log"

	"HOMEWORK-1/internal/models"
	"HOMEWORK-1/pkg/api/proto/pvz/v1/pvz/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	packagingMap = map[string]models.Packaging{
		"bag":  models.NewBag(),
		"box":  models.NewBox(),
		"film": models.NewFilm(),
		"none": models.NewNoPackaging(),
	}
)

type PVZService struct {
	Module module.Module
	pvz.UnimplementedPVZServer
}

func (t *PVZService) AddOrder(ctx context.Context, req *pvz.AddOrderRequest) (*emptypb.Empty, error) {
	log.Printf("[PVZService.AddOrder] %v", req)
	if errValidate := req.ValidateAll(); errValidate != nil {
		return nil, status.Error(codes.InvalidArgument, errValidate.Error())
	}
	order, errToDomain := pvzToDomain(req)
	if errToDomain != nil {
		return nil, status.Error(codes.InvalidArgument, errToDomain.Error())
	}
	err := t.Module.AddOrder(ctx, order)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func pvzToDomain(req *pvz.AddOrderRequest) (models.Order, error) {
	packaging, ok := packagingMap[req.GetOrder().GetPackaging()]
	if !ok {
		packaging = models.NewNoPackaging()
	}
	return models.Order{
		ID:          models.ID(req.GetOrder().GetId()),
		IDReceiver:  models.ID(req.GetOrder().GetIdReceiver()),
		StorageTime: req.GetOrder().GetStorageTime().AsTime(),
		WeightKg:    req.GetOrder().GetWeight(),
		Packaging:   packaging,
	}, nil
}

func (t *PVZService) DeleteOrder(ctx context.Context, req *pvz.DeleteOrderRequest) (*pvz.DeleteOrderResponse, error) {
	err := t.Module.DeleteOrder(ctx, models.ID(req.GetID().GetId()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pvz.DeleteOrderResponse{}, nil
}

func (t *PVZService) ListOrder(ctx context.Context, req *pvz.ListOrderRequest) (*pvz.ListOrderResponse, error) {
	list, err := t.Module.ListOrder(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	listPVZ := make([]*pvz.OrderInfo, 0, len(list))
	for _, order := range list {
		listPVZ = append(listPVZ, &pvz.OrderInfo{
			Order: &pvz.Order{
				Id:         order.ID,
				IdReceiver: order.IDReceiver,
				Weight:     order.WeightKg,
				Packaging:  order.Packaging.GetName(), // Assuming Packaging has a String() method
				// Add other fields as necessary
			},
		})
	}
	return &pvz.ListOrderResponse{
		List: listPVZ,
	}, nil
}

func (t *PVZService) FindOrder(ctx context.Context, req *pvz.FindOrderRequest) (*pvz.FindOrderResponse, error) {
	order, err := t.Module.GetOrderByID(ctx, models.ID(req.GetId().GetId()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pvz.FindOrderResponse{
		Order: &pvz.OrderInfo{
			Order: pvz.Order{
				Id:         order.ID,
				IdReceiver: order.IDReceiver,
				Weight:     order.WeightKg,
				Packaging:  order.Packaging.GetName(), // Assuming Packaging has a String() method
			},
		},
	}, nil
}
