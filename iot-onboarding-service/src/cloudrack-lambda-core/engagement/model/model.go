package model

import (
	core "cloudrack-lambda-core/core"
	"encoding/json"
)

type ResWrapper struct {
	Error       core.ResError `json:"error"`
	SubFunction string        `json:"subFunction"`
	Status      string        `json:"status"` // overall transaction status success/failure
}

func (p ResWrapper) Decode(dec json.Decoder) (error, core.CloudrackObject) {
	return dec.Decode(&p), p
}

type EmailTemplateData struct {
	BookingNumber   string
	HotelName       string
	StartDate       string
	Currency        string
	NGuests         string
	NNights         string
	Price           string
	HolderFirstName string
	HolderLastName  string
	Status          string
}
type DynamoSegmentRecord struct {
	HotelCode   string `json:"hotelCode"`
	SegmentCode string `json:"segmentCode"`
	SegmentId   string `json:"segmentId"`
}

type BookingGuestInfo struct {
	FirstName       string `json:"firstName"`       //unique identifier of the hotel
	MiddleName      string `json:"middleName"`      //unique identifier of the hotel
	LastName        string `json:"lastName"`        //unique identifier of the hotel
	Email           string `json:"email"`           //unique identifier of the hotel
	Phone           string `json:"phone"`           //unique identifier of the hotel
	LoyaltyId       string `json:"loyaltyId"`       //unique identifier of the hotel
	CompanyCode     string `json:"companyCode"`     //unique identifier of the hotel
	Mot             string `json:"mot"`             //unique identifier of the hotel
	MotTicketId     string `json:"motTicketId"`     //unique identifier of the hotel
	Segment         string `json:"segment"`         //unique identifier of the hotel
	SubSegment      string `json:"subSegment"`      //unique identifier of the hotel
	IdentityProofs  string `json:"identityProofs"`  //unique identifier of the hotel
	LoyaltyPrograms string `json:"loyaltyPrograms"` //unique identifier of the hotel

}
