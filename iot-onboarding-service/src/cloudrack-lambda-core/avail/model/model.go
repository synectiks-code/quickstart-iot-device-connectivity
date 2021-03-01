package model

import (
	cfg "cloudrack-lambda-core/config/model"
	core "cloudrack-lambda-core/core"
	"encoding/json"
	"time"
)

var PRODUCT_STATUS_AVAILABLE string = "available"
var PRODUCT_STATUS_UNAVAILABLE string = "unavailable"
var AVAILABILITY_DB_TIMEFORMAT string = "20060102"

//Requests and Respone to deserialize teh Json Payload into.
//ideally this should be a generic struct in teh core package with a specialised body but
//the design needs more thoughts so we will include it here for the moment
type RqWrapper struct {
	SubFunction string       `json:"subFunction"`
	Id          string       `json:"id"` //ID to be used for GET request
	Request     AvailRequest `json:"request"`
}

type ResWrapper struct {
	Error       core.ResError `json:"error"`
	SubFunction string        `json:"subFunction"`
	Response    AvailResponse `json:"response"`
}

func (res *ResWrapper) AddError(err error) {
	res.Error = core.BuildResError(err)
}

func (p ResWrapper) Decode(dec json.Decoder) (error, core.CloudrackObject) {
	return dec.Decode(&p), p
}

type AvailRequest struct {
	Hotels      []string `json:"hotels"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	TextSearch  string   `json:"textSearch"`
	NNights     int      `json:"nNights"`
	NGuests     int      `json:"nGuests"`
	NRooms      int      `json:"nRooms"`
	Lat         float64  `json:"lat"`
	Lng         float64  `json:"lng"`
	Tags        []string `json:"tags"`
	PromoCode   string   `json:"promoCode"`
	GroupCode   string   `json:"groupCode"`
	CorporateId string   `json:"corporateId"`
	LoyaltyId   string   `json:"loyaltyId"`
	Name        string   `json:"name"`
}

type AvailResponse struct {
	SingleAvail       []Product         `json:"singleAvail"`
	MultiAvail        []HotelAvail      `json:"multiAvail"`
	AdditionalResults AdditionalResults `json:"additionalResults"`
	//TODO: to build dedicated config model for avail.
	//do not reuse the config model iself to avoid create
	//model dependency
	Config    cfg.Hotel       `json:"config"`
	HotelList []string        `json:"hotelList"`
	Pricing   PricingResponse `json:"pricingResponse"`
}

type AdditionalResults struct {
	KendraResults map[string]string `json:"kendraResults"`
}

//////////////////////////////
//Pricing
//////////////////////////////

type PricingRequest struct {
	Segments []PricingRequestSegment `json:"segments"`
}

type PricingRequestSegment struct {
	SegmentId string           `json:"segmentId"`
	Hotels    []string         `json:"hotels"`
	StartDate string           `json:"startDate"`
	EndDate   string           `json:"endDate"`
	NNights   int              `json:"nNights"`
	NGuests   int              `json:"nGuests"`
	Products  []BookingProduct `json:"products"`
}

type PricingResponse struct {
	PricePerSegment map[string]SegmentPrice `json:"pricePerSegment"`
}

type BookingProduct struct {
	Id           string   `json:"id"`
	RoomTypeCode string   `json:"roomTypeCode"`
	TagCodes     []string `json:"tagCodes"`
	ProductCodes []string `json:"productCodes"`
}

type SegmentPrice struct {
	PricePerNight        map[string]map[string]ProductPrice `json:"pricePerNight"`
	TaxPerNight          map[string]map[string]BookingTax   `json:"taxPerNight"`
	PricePerStay         map[string]ProductPrice            `json:"pricePerStay"`
	TaxPerStay           []BookingTax                       `json:"taxPerStay"`
	BusinessRules        []BusinessRule                     `json:"businessRules"`
	TotalPricePerNight   map[string]float64                 `json:"totalPricePerNight"`
	TotalPricePerProduct map[string]float64                 `json:"totalPricePerProduct"`
	Total                float64                            `json:"total"`
}

type ProductPrice struct {
	Price float64  `json:"roomTypePrice"`
	Tags  []string `json:"tagCodes"`
}

type BookingTax struct {
	TaxCode   string  `json:"taxCode"`
	TaxLabel  string  `json:"taxLabel"`
	TaxAmount float64 `json:"taxAmount"`
	Included  bool    `json:"included"`
}

type BusinessRule struct {
}

type HotelAvail struct {
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Pictures    []Picture `json:"pictures"`
	StartPrice  Price     `json:"startPrice"`
	Tags        []Tag     `json:"tags"`
	Amenities   []Amenity `json:"amenities"`
	Status      string    `json:"status"` //status of the hotel (avail/ not avail)

}

type Product struct {
	Code           string   `json:"code"`
	Name           string   `json:"name"`
	RoomType       RoomType `json:"roomType"`
	TagList        []Tag    `json:"tags"`
	AddOns         []AddOn  `json:"addOn"`
	Price          Price    `json:"price"`
	Status         string   `json:"status"`         //status of the product (avail/ not avail)
	NumProductLeft int      `json:"numProductLeft"` //number of products left in allotment
}

type Amenity struct {
}

type Price struct {
	BasePrice     float64            `json:"baseRoomTypePrice"`
	IncludedTaxes map[string]float64 `json:"includedTaxes"`
	ExcludedTaxes map[string]float64 `json:"excludedTaxes"`
	AddOnPrices   map[string]Price   `json:"addOnPrices"`
	Currency      string             `json:"currency"`
}

type AddOn struct {
	Code          string    `json:"code"`
	Category      string    `json:"category"`
	Name          string    `json:"name"`
	InventoryType string    `json:"inventoryType"`
	PricePerUnit  float64   `json:"pricePerUnit"`
	MaxQuantity   int64     `json:"maxQuantity"`
	Description   string    `json:"description"`
	Pictures      []Picture `json:"pictures"`
	Status        string    `json:"status"`
}

type RoomType struct {
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Pictures    []Picture `json:"pictures"`
}

type Tag struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type Picture struct {
	Id      string   `json:"id"`
	RawData string   `json:"rawData"`
	Format  string   `json:"format"`
	Url     string   `json:"url"`
	Tags    []string `json:"tags"`
}

type DynamoMultiAvailRecord struct {
	Code      string `json:"code"`
	Geohash   string `json:"geohash"`
	HotelCode string `json:"hotelCode"`
}

type DynamoProductRecord struct {
	Code                    string             `json:"code"`
	Name                    string             `json:"name"`
	RoomTypeCode            string             `json:"roomTypeCode"`
	TagCodeList             []string           `json:"code"`
	AddOnCodeList           []string           `json:"code"`
	IncludedTaxes           map[string]float64 `json:"taxes"`
	ExcludedTaxes           map[string]float64 `json:"taxes"`
	AddOnPriceList          []string           `json:"addOnPrices"`
	AddOnincludedTaxeAmount map[string]float64 `json:"addOnPrices"`
	AddOnExcludedTaxeAmount map[string]float64 `json:"addOnPrices"`
	Currency                string             `json:"currency"`
	Status                  string             `json:"status"`
}

//////////////////////////////
//Regrets
//////////////////////////////

type RegretRequest struct {
	UtcTimestamp string `json:"utcTimestamp"`
	Hotel        string `json:"hotel"`
	LoyaltyID    string `json:"loyaltyId"`
	StartDate    string `json:"startDate"`
	NNights      int    `json:"conNightsde"`
	NGuests      int    `json:"nGuests"`
	NRooms       int    `json:"nRooms"`
}

func (rq AvailRequest) BuildDateRange() []string {
	startDate, _ := time.Parse(AVAILABILITY_DB_TIMEFORMAT, rq.StartDate)
	dateRange := make([]string, 0, 0)
	for i := 0; i < rq.NNights; i++ {
		t := startDate.AddDate(0, 0, i)
		dateRange = append(dateRange, t.Format(AVAILABILITY_DB_TIMEFORMAT))
	}
	return dateRange
}
