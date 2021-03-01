package service

import (
	avail "cloudrack-lambda-core/avail/model"
	cfg "cloudrack-lambda-core/config/model"
	core "cloudrack-lambda-core/core"
	"strconv"
)

type Service struct {
	Client core.HttpService
}

func Init(endpoint string) Service {
	return Service{Client: core.HttpInit(endpoint)}
}

func (s Service) UpdateHotelList(hotel cfg.Hotel) error {
	_, err := s.Client.HttpPut(hotel.Code, hotel)
	return err
}

func (s Service) FindHotelByName(hotelName string) (avail.ResWrapper, error) {
	params := map[string]string{
		"name": hotelName,
	}
	res := avail.ResWrapper{}
	res2, err := s.Client.HttpGet(params, res, core.RestOptions{SubEndpoint: "config"})
	res = res2.(avail.ResWrapper)
	return res, err
}

func (s Service) FindAvailability(avilRq avail.AvailRequest) (avail.ResWrapper, error) {
	params := map[string]string{
		"hotel":     avilRq.Hotels[0],
		"startDate": avilRq.StartDate,
		"nNights":   strconv.Itoa(avilRq.NNights),
		"nGuests":   strconv.Itoa(avilRq.NGuests),
	}
	res := avail.ResWrapper{}
	res2, err := s.Client.HttpGet(params, res, core.RestOptions{SubEndpoint: "single"})
	res = res2.(avail.ResWrapper)
	return res, err
}
