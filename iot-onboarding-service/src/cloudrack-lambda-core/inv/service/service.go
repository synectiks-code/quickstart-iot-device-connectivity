package service

import (
	cfg "cloudrack-lambda-core/config/model"
	core "cloudrack-lambda-core/core"
	inv "cloudrack-lambda-core/inv/model"
)

type Service struct {
	Client core.HttpService
}

func Init(endpoint string) Service {
	return Service{Client: core.HttpInit(endpoint)}
}

func (s Service) UpdateInventory(hotel cfg.Hotel) error {
	_, err := s.Client.HttpPut("", hotel)
	return err
}

func (s Service) PerformInventoryOperation() {

}

func (s Service) GetHotelInventoryForDateRange(hotelCode string, startDate string, endDate string) (inv.ResWrapper, error) {
	params := map[string]string{
		"hotelCode": hotelCode,
		"startDate": startDate,
		"endDate":   endDate,
	}
	res := inv.ResWrapper{}
	res2, err := s.Client.HttpGet(params, res, core.RestOptions{})
	res = res2.(inv.ResWrapper)
	return res, err
}
