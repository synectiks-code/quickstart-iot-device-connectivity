package model

import (
	core "cloudrack-lambda-core/core"
	"encoding/json"
	"time"
)

var INVENTORY_TYPE_ROOM string = "room"
var INVENTORY_TYPE_SELLABLE string = "sellable"
var INVENTORY_ALLOTMENT_GINV string = "GINV"
var INVENTORY_DB_TIMEFORMAT string = "20060102"

type ResWrapper struct {
	Error       core.ResError      `json:"error"`
	SubFunction string             `json:"subFunction"`
	Response    map[string]Product `json:"response"`
	Status      string             `json:"status"` // overall transaction status success/failure
}

func (p ResWrapper) Decode(dec json.Decoder) (error, core.CloudrackObject) {
	return dec.Decode(&p), p
}

type InventoryRequest struct {
	Actions []InventoryAction `json:"actions"`
}

func (invRq InventoryRequest) HotelCodes() []string {
	hotelCodes := []string{}
	hotelCodeMap := map[string]string{}
	for _, action := range invRq.Actions {
		hotelCodeMap[action.HotelCode] = action.HotelCode
	}
	for _, code := range hotelCodeMap {
		hotelCodes = append(hotelCodes, code)
	}
	return hotelCodes
}
func (invRq InventoryRequest) MinStartDate() string {
	minDate, _ := time.Parse(INVENTORY_DB_TIMEFORMAT, invRq.Actions[0].StartDate)
	for _, action := range invRq.Actions {
		date, _ := time.Parse(INVENTORY_DB_TIMEFORMAT, action.StartDate)
		if date.Before(minDate) {
			minDate = date
		}
	}
	return minDate.Format(INVENTORY_DB_TIMEFORMAT)
}
func (invRq InventoryRequest) MaxEndDate() string {
	maxDate, _ := time.Parse(INVENTORY_DB_TIMEFORMAT, invRq.Actions[0].EndDate)
	for _, action := range invRq.Actions {
		date, _ := time.Parse(INVENTORY_DB_TIMEFORMAT, action.EndDate)
		if date.After(maxDate) {
			maxDate = date
		}
	}
	return maxDate.Format(INVENTORY_DB_TIMEFORMAT)
}

type InventoryAction struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	HotelCode string `json:"hotelCode"`
	Code      string `json:"code"`
	Type      string `json:"type"`
	Action    string `json:"action"`  //"hold"/"release"
	Success   bool   `json:"success"` //result of inventory action  (true=SUCCESS/false=FAILURE)
}

//////////////////////////////////////////////
//*-----------------------------------------
// TYPE  CODE   ALOTMENT	COUNTTYPE    2019-05-14  2019-05-15  2019-05-16
// ROOM  DBL    GINV		BASECOUNT    10				9			10
// ROOM  DBL    GINV		HELD 	     2				1			0
// ROOM  DBL    A1			BASECOUNT    3				3			3
// ROOM  DBL    A1			HELD 	     2				1			0

type Product struct {
	HotelCode    string                          `json:"hotelCode"`
	Code         string                          `json:"code"`
	Type         string                          `json:"type"`
	Counts       map[string]map[string]Allotment `json:"counts"`       //2019-05-14 => GINV => {CODE: GINV, BASECOUNT: 10, HELD: 2}
	ActionStatus InventoryAction                 `json:"actionStatus"` //status of the action on inventory for thsi product
}

type Allotment struct {
	Code      string `json:"code"`
	BaseCount int    `json:"baseCount"`
	Held      int    `json:"held"`
}

type AuroraInventoryRecord struct {
	HotelCode     string
	ProductType   string
	ProductCode   string
	AllotmentCode string
	Status        string //BASECOUNT/HELD
	InventoryType string //per_unit / per_day... see config.Sellable model
	BaseCount     int
	Held          int //to be used for "per_unit" inventory types
	CountByDate   map[string]int
	HeldByDate    map[string]int
}
