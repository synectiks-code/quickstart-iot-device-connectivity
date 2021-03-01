package service

import (
	bookingModel "cloudrack-lambda-core/booking/model"
	core "cloudrack-lambda-core/core"
	"log"
)

type Service struct {
	Client    core.HttpService
	ChannelID string
}

func Init(endpoint string, channelId string) Service {
	return Service{Client: core.HttpInit(endpoint), ChannelID: channelId}
}

func (s Service) BookAndCommit(booking bookingModel.Booking) (bookingModel.Booking, error) {
	log.Printf("[CORE][BOOKING] BookAndCommitRQ: %+v\n", booking)

	commitRes := bookingModel.ResWrapper{}
	res := bookingModel.ResWrapper{}
	res2, err := s.Client.HttpPost(bookingModel.RqWrapper{Request: []bookingModel.Booking{booking}}, res, core.RestOptions{Headers: map[string]string{bookingModel.BOOKING_CHANNEL_HEADER: s.ChannelID}})
	res = res2.(bookingModel.ResWrapper)
	if err != nil {
		return bookingModel.Booking{}, err
	}
	_, err = s.Client.HttpPost(bookingModel.RqWrapper{Request: []bookingModel.Booking{booking}}, commitRes, core.RestOptions{SubEndpoint: "commit/" + res.Bookings[0].HotelCode + "/" + res.Bookings[0].Id, Headers: map[string]string{bookingModel.BOOKING_CHANNEL_HEADER: s.ChannelID}})
	return res.Bookings[0], err
}
