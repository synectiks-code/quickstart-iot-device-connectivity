package service

import (
	bookingModel "cloudrack-lambda-core/booking/model"
	core "cloudrack-lambda-core/core"
	engModel "cloudrack-lambda-core/engagement/model"
)

type Service struct {
	Client core.HttpService
}

func Init(endpoint string) Service {
	return Service{Client: core.HttpInit(endpoint)}
}

func (s Service) CreateEmailTemplate(booking bookingModel.Booking) error {
	_, err := s.Client.HttpPost(booking, engModel.ResWrapper{}, core.RestOptions{SubEndpoint: "template"})
	return err
}

func (s Service) SendConfirmationMessage(booking bookingModel.Booking) error {
	_, err := s.Client.HttpPost(booking, engModel.ResWrapper{}, core.RestOptions{SubEndpoint: "confirmation"})
	return err
}

func (s Service) UpdateEnpoints(booking bookingModel.Booking) error {
	_, err := s.Client.HttpPost(booking, engModel.ResWrapper{}, core.RestOptions{SubEndpoint: "enpoint"})
	return err
}

func (s Service) UpdateSegmentsAndJourney(booking bookingModel.Booking) error {
	_, err := s.Client.HttpPost(booking, engModel.ResWrapper{}, core.RestOptions{SubEndpoint: "journey"})
	return err
}
