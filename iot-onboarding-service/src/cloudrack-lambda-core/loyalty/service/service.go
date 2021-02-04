package service

import (
	bookingModel "cloudrack-lambda-core/booking/model"
	core "cloudrack-lambda-core/core"
	loyaltyModel "cloudrack-lambda-core/loyalty/model"
)

type Service struct {
	Client core.HttpService
}

func Init(endpoint string) Service {
	return Service{Client: core.HttpInit(endpoint)}
}

func (s Service) FindUserProfile(booking bookingModel.Booking) (loyaltyModel.ResWrapper, error) {
	res := loyaltyModel.ResWrapper{}
	res2, err := s.Client.HttpGet(map[string]string{}, res, core.RestOptions{SubEndpoint: booking.Segments[0].Holder.LoyaltyId})
	res = res2.(loyaltyModel.ResWrapper)
	return res, err
}
