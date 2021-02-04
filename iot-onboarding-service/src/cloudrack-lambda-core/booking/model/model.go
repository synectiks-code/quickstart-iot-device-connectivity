package model

import (
	core "cloudrack-lambda-core/core"
	"encoding/json"
)

var STATUS_PENDING string = "pending"
var STATUS_CONFIRMED string = "confirmed"
var STATUS_CANCELLED string = "cancelled"
var STATUS_IGNORED string = "ignored"
var BOOKING_CHANNEL_HEADER = "cloudrack-booking-channel-id"

type RqWrapper struct {
	Request []Booking `json:"request"`
}

type ResWrapper struct {
	Errors   []ResError `json:"errors"`
	Bookings []Booking  `json:"bookings"`
	Status   string     `json:"status"` // overall transaction status success/failure
}

type ResError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (p ResWrapper) Decode(dec json.Decoder) (error, core.CloudrackObject) {
	return dec.Decode(&p), p
}

type Booking struct {
	ObjectVersion      int                       `json:"objectVersion"`      //unique identifier of the hotel
	LastUpdated        string                    `json:"lastUpdated"`        //unique identifier of the hotel
	Id                 string                    `json:"id"`                 //unique identifier of the hotel
	QldbDocumentId     string                    `json:"qldbDocumentId"`     //unique identifier of the hotel
	CreationChannelId  string                    `json:"creationChannelId"`  //unique identifier of the hotel
	HotelCode          string                    `json:"hotelCode"`          //unique identifier of the hotel
	NNights            int                       `json:"nNights"`            //unique identifier of the hotel
	NGuests            int                       `json:"nGuests"`            //unique identifier of the hotel
	StartDate          string                    `json:"startDate"`          //unique identifier of the hotel
	Comments           []string                  `json:"comments"`           //unique identifier of the hotel
	Holder             BookingGuestInfo          `json:"holder"`             //unique identifier of the hotel
	PaymentInformation BookingPaymentInformation `json:"paymentInformation"` //unique identifier of the hotel
	ContextId          string                    `json:"contextId"`          //unique identifier of the hotel
	GroupId            string                    `json:"groupId"`            //unique identifier of the hotel
	Status             string                    `json:"status"`             //unique identifier of the hotel
	Currency           string                    `json:"currency"`           //unique identifier of the hotel
	CancelReason       BookingCancellationReason `json:"cancelReason"`       //unique identifier of the hotel
	Segments           []BookingSegment          `json:"segments"`           //unique identifier of the hotel
}

type BookingCancellationReason struct {
	Reason  string `json:"reason"`
	Comment string `json:"comment"`
}

type BookingGuestInfo struct {
	FirstName       string                  `json:"firstName"`       //unique identifier of the hotel
	MiddleName      string                  `json:"middleName"`      //unique identifier of the hotel
	LastName        string                  `json:"lastName"`        //unique identifier of the hotel
	Email           string                  `json:"email"`           //unique identifier of the hotel
	Phone           string                  `json:"phone"`           //unique identifier of the hotel
	LoyaltyId       string                  `json:"loyaltyId"`       //unique identifier of the hotel
	CompanyCode     string                  `json:"companyCode"`     //unique identifier of the hotel
	Mot             string                  `json:"mot"`             //unique identifier of the hotel
	MotTicketId     string                  `json:"motTicketId"`     //unique identifier of the hotel
	Segment         string                  `json:"segment"`         //unique identifier of the hotel
	SubSegment      string                  `json:"subSegment"`      //unique identifier of the hotel
	IdentityProofs  []BookingIdentityProof  `json:"identityProofs"`  //unique identifier of the hotel
	LoyaltyPrograms []BookingLoyaltyProgram `json:"loyaltyPrograms"` //unique identifier of the hotel
	Address         BookingLocation         `json:"address"`         //unique identifier of the hotel

}

type BookingLoyaltyProgram struct {
	ProgramId        string `json:"programId"`
	ProgramName      string `json:"programName"`
	LoyaltyNumber    string `json:"loyaltyNumber"`
	LoyaltyStatus    string `json:"loyaltyStatus"`
	LoyaltyProfileId string `json:"loyaltyProfileId"`
}

type BookingIdentityProof struct {
	proofId     string `json:"proofId"`
	proofTypeId string `json:"proofTypeId"`
	proofdata   string `json:"proofdata"`
}

type BookingPaymentInformation struct {
	PaymentType string                       `json:"paymentType"` //unique identifier of the hotel
	CcInfo      BookingCreditCardInformation `json:"ccInfo"`      //unique identifier of the hotel
}

type BookingCreditCardInformation struct {
	Token      string          `json:"token"`      //unique identifier of the hotel
	Expiration string          `json:"expiration"` //unique identifier of the hotel
	Name       string          `json:"name"`       //unique identifier of the hotel
	Address    BookingLocation `json:"address"`    //unique identifier of the hotel
}

type BookingLocation struct {
	Lines     []string `json:"lines"`     //unique identifier of the hotel
	City      string   `json:"city"`      //unique identifier of the hotel
	Country   string   `json:"country"`   //unique identifier of the hotel
	Latitude  float64  `json:"latitude"`  //unique identifier of the hotel
	Longitude float64  `json:"longitude"` //unique identifier of the hotel
	ZipCode   string   `json:"zipcode"`   //unique identifier of the hotel
}

type BookingSegment struct {
	Id                 string                    `json:"id"`                 //unique identifier of the hotel
	HotelCode          string                    `json:"hotelCode"`          //unique identifier of the hotel
	NNights            int                       `json:"nNights"`            //unique identifier of the hotel
	NGuests            int                       `json:"nGuests"`            //unique identifier of the hotel
	StartDate          string                    `json:"startDate"`          //unique identifier of the hotel
	Products           []BookingProduct          `json:"products"`           //unique identifier of the hotel
	Comments           []string                  `json:"comments"`           //unique identifier of the hotel
	Holder             BookingGuestInfo          `json:"holder"`             //unique identifier of the hotel
	AdditionalGuests   []BookingGuestInfo        `json:"additionalGuests"`   //unique identifier of the hotel
	PaymentInformation BookingPaymentInformation `json:"paymentInformation"` //unique identifier of the hotel
	ContextId          string                    `json:"contextId"`          //unique identifier of the hotel
	GroupId            string                    `json:"groupId"`            //unique identifier of the hotel
	Status             string                    `json:"status"`             //unique identifier of the hotel
	Price              BookingPrice              `json:"price"`              //unique identifier of the hotel

}

type BookingPrice struct {
	PricePerNight        map[string]map[string]BookingAmount `json:"pricePerNight"`        //unique identifier of the hotel
	TaxePerNight         map[string]map[string]BookingAmount `json:"taxePerNight"`         //unique identifier of the hotel
	PricePerStay         map[string]BookingAmount            `json:"pricePerStay"`         //unique identifier of the hotel
	TaxPerStay           map[string]BookingAmount            `json:"taxPerStay"`           //unique identifier of the hotel
	BusinessRules        []BookingBusinessRule               `json:"businessRules"`        //unique identifier of the hotel
	TotalPricePerNight   map[string]float64                  `json:"totalPricePerNight"`   //unique identifier of the hotel
	TotalPricePerProduct map[string]float64                  `json:"totalPricePerProduct"` //unique identifier of the hotel
	Total                float64                             `json:"total"`                //unique identifier of the hotel
}

type BookingBusinessRule struct {
	Id        string                    `json:"pricePerNight"`
	AppliesTo BookingBusinessRuleScope  `json:"appliesTo"`
	Effect    BookingBusinessRuleEffect `json:"effect"`
	On        []BookingDateRange        `json:"on"`
	NotOn     []BookingDateRange        `json:"notOn"`
}
type BookingBusinessRuleScope struct {
	Selectors BookingBusinessObjectSelector `json:"selectors"`
}
type BookingBusinessObjectSelector struct {
	ObjectType string `json:"objectType"`
	Field      string `json:"field"`
	Op         string `json:"op"`
	Value      string `json:"value"`
}

type BookingBusinessRuleEffect struct {
	Type            string  `json:"type"`
	Available       bool    `json:"available"`
	PriceImpact     float64 `json:"priceImpact"`
	PriceUpdateUnit string  `json:"priceUpdateUnit"`
}
type BookingDateRange struct {
	from string   `json:"from"`
	to   string   `json:"to"`
	dow  []string `json:"dow"`
}

type BookingAmount struct {
	RoomTypePrice float64  `json:"roomTypePrice"` //unique identifier of the hotel
	TagCodes      []string `json:"tagCodes"`      //unique identifier of the hotel
}

type BookingProduct struct {
	Id                   string            `json:"id"`                   //unique identifier of the hotel
	RoomTypeCode         string            `json:"roomTypeCode"`         //unique identifier of the hotel
	TagCodes             []string          `json:"tagCodes"`             //unique identifier of the hotel
	ProductCodes         []string          `json:"productCodes"`         //unique identifier of the hotel
	RoomTypeName         string            `json:"roomTypeName"`         //unique identifier of the hotel
	RoomTypeDescription  string            `json:"roomTypeDescription"`  //unique identifier of the hotel
	TagNames             map[string]string `json:"tagNames"`             //unique identifier of the hotel
	TagDescriptions      map[string]string `json:"tagDescriptions"`      //unique identifier of the hotel
	SellableNames        map[string]string `json:"sellableNames"`        //unique identifier of the hotel
	SellableDescriptions map[string]string `json:"sellableDescriptions"` //unique identifier of the hotel

}
