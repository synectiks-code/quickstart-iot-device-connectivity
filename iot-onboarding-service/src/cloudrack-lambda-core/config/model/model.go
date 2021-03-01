package model

import (
	core "cloudrack-lambda-core/core"
	"log"
	"reflect"
)

var SELLABLE_INVENTORY_TYPE_PER_DAY string = "per_day"
var SELLABLE_INVENTORY_TYPE_PER_NIGHT string = "per_night"
var SELLABLE_INVENTORY_TYPE_PER_UNIT string = "per_unit"
var TAG_PRICE_UPDATE_UNIT_PERCENT string = "percent"
var TAG_PRICE_UPDATE_UNIT_CURRENCY string = "currency"

type Hotel struct {
	Code              string               `json:"code"`              //unique identifier of the hotel
	Name              string               `json:"name"`              //hotel name
	Description       string               `json:"description"`       //Hotel description
	CurrencyCode      string               `json:"currencyCode"`      //Hotel curency code
	Lat               float64              `json:"lat"`               //Hotel latitude
	Lng               float64              `json:"lng"`               //Hotel longitude
	Timezone          HotelTimezone        `json:"timezone"`          //Hotel timezone
	DefaultMarketting HotelMarketing       `json:"defaultMarketting"` //Object containing marketting text
	Options           HotelOptions         `json:"options"`           //Object containing transactional options (open for search, reservation?)
	Pictures          []HotelPicture       `json:"pictures"`          //List of Picture object containing URL of image content and tags
	Buildings         []HotelBuilding      `json:"buildings"`         //List of hotel building (each buildings has floors and rooms)
	RoomTypes         []HotelRoomType      `json:"roomTypes"`         //List of Roome types (Queen, Suite, presidential suite...)
	Tags              []HotelRoomAttribute `json:"tags"`              //List of tags impacting room type prices (such as Sea view or minibar)
	Sellables         []HotelSellable      `json:"sellables"`         //List of sellable item (could be meals, excursion...)
	Rules             []HotelBusinessRule  `json:"businessRules"`     //List of Business rules defining prices and availability of rooms, tags sellable
	PendingChanges    []HotelConfigChange  `json:"pendingChanges"`    //List of Pending Changes to the hotel configuration record
	History           []HotelConfigChange  `json:"history"`           //List of All Changes to the hotel configuration record
}

type HotelTimezone struct {
	Id   string `json:"id"`   //ex: Europe/Paris
	Name string `json:"name"` //ex: Central European Standard Time
}

func (bom Hotel) ListUpdatedFields(oldBom Hotel) []string {
	v := reflect.ValueOf(bom)
	oldV := reflect.ValueOf(bom)
	fields := make([]string, 0)
	for i := 0; i < v.NumField(); i++ {
		log.Printf("HOTEL Diff: checking  %+v againstt %+v", v.Field(i), oldV.Field(i))
		if !reflect.DeepEqual(v.Field(i).Interface(), oldV.Field(i).Interface()) {
			fields = append(fields, v.Type().Field(i).Name)
		}
	}
	return fields
}

func (bom *Hotel) AddRoomTypePictures(picsByCode map[string][]HotelPicture) {
	log.Printf("Picture map to add %+v", picsByCode)

	for i, roomType := range bom.RoomTypes {
		log.Printf("Adding pic for room tyep %+v to : %+v", roomType, picsByCode[roomType.Code])
		roomType.Pictures = append(roomType.Pictures, picsByCode[roomType.Code]...)
		bom.RoomTypes[i] = roomType
	}
	log.Printf("Added pic to : %+v", bom.RoomTypes)
}

func (bom *Hotel) AddSellablePictures(picsByCode map[string][]HotelPicture) {
	for i, sell := range bom.Sellables {
		sell.Pictures = append(sell.Pictures, picsByCode[sell.Code]...)
		bom.Sellables[i] = sell
	}
}

//Populate hotelAttribute struct with all fields
func (bom *Hotel) AddAttributesToRooms(tagsByCode map[string]HotelRoomAttribute) {
	for i, building := range bom.Buildings {
		for j, floor := range building.Floors {
			for k, room := range floor.Rooms {
				for l, tag := range room.Attributes {
					fullTag := tagsByCode[tag.Code]
					if fullTag.Code != "" {
						room.Attributes[l] = tagsByCode[tag.Code]
					}
				}
				floor.Rooms[k] = room
			}
			building.Floors[j] = floor
		}
		bom.Buildings[i] = building
	}
}

type HotelBusinessRule struct {
	Id        string                     `json:"id"`
	AppliesTo []HotelBusinessObjectScope `json:"appliesTo"`
	Effect    []HotelBusinessRuleEffect  `json:"effect"`
	On        []HotelDateRange           `json:"on"`
	NotOn     []HotelDateRange           `json:"notOn"`
}

//Defines the scope of a business rule
type HotelBusinessObjectScope struct {
	Type string `json:"type"` //Room, RoomType, Attribute, Sellable
	Id   string `json:"id"`   //unique identifier of the business object
}

//Effect of a business rule on a hotel product
type HotelBusinessRuleEffect struct {
	EffectType      string `json:"effectType"` //availability, pricing, inventory
	Available       bool   `json:"available"`
	Bookable        bool   `json:"bookable"`
	PriceImpact     int64  `json:"priceImpact"`
	PriceUpdateUnit string `json:"priceUpdateUnit"`
}

//definition of the time frame a rule applies
type HotelDateRange struct {
	From   string   `json:"from"`   //From date
	To     string   `json:"from"`   //To date
	Dow    []string `json:"dow"`    // Day of weeks
	Within int64    `json:"within"` // with x seconds
	Of     string   `json:"of"`     // time
}

type HotelPicture struct {
	Id              string   `json:"id"`
	RawData         string   `json:"rawData"`
	Format          string   `json:"format"`
	Url             string   `json:"url"`
	Tags            []string `json:"tags"`
	Main            bool     `json:"main"`
	PictureItemCode string   `json:"pictureItemCode"` //code of the business object the picture belongs to
}

type HotelSellable struct {
	Code          string               `json:"code"`          //unique identifier of the sellable
	Category      string               `json:"category"`      //Category of the sellable (custom test such as "food and beverage")
	Name          string               `json:"name"`          //name of the sellable as it woudl be return in search response
	Quantity      int64                `json:"quantity"`      //Base quantity of the product at the hotel (us for per_unit inventory type)
	InventoryType string               `json:"inventoryType"` //inventoty type per_day / per_unit / per night
	PricePerUnit  float64              `json:"pricePerUnit"`  // Price per unit
	Description   string               `json:"description"`   //description
	OptionalTags  []HotelRoomAttribute `json:"optionalTags"`  //Optional tags
	Pictures      []HotelPicture       `json:"pictures"`
}

type HotelRoomType struct {
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	LowPrice    float64        `json:"lowPrice"`
	MedPrice    float64        `json:"medPrice"`
	HighPrice   float64        `json:"highPrice"`
	Description string         `json:"description"`
	Pictures    []HotelPicture `json:"pictures"`
}

type HotelRoomAttribute struct {
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	PriceImpact     float64 `json:"priceImpact"`
	PriceUpdateUnit string  `json:"priceUpdateUnit"` //TAG_PRICE_UPDATE_UNIT_<?>
	Category        string  `json:"category"`
}

type HotelOptions struct {
	Bookable     bool   `json:"bookable"`
	Shoppable    bool   `json:"shoppable"`
	CheckInTime  string `json:"checkInTime"`
	CheckOutTime string `json:"checkOutTime"`
}

type HotelMarketing struct {
	DefaultTagLine     string `json:"defaultTagLine"`
	DefaultDescription string `json:"defaultDescription"`
}

type HotelBuilding struct {
	Name   string       `json:"name"`
	Floors []HotelFloor `json:"floors"`
}

type HotelFloor struct {
	Num   int32       `json:"num"`
	Rooms []HotelRoom `json:"rooms"`
}

type HotelRoom struct {
	Number     int64                `json:"number"`
	Name       string               `json:"name"`
	Type       string               `json:"type"`
	Attributes []HotelRoomAttribute `json:"attributes"`
}

type HotelConfigChange struct {
	Id                   string `json:"id"`
	TimeStamp            string `json:"timeStamp"`
	EventName            string `json:"eventName"`
	ObjectName           string `json:"objectName"`
	FieldName            string `json:"fieldName"`
	OldValue             string `json:"oldValue"`
	NewValue             string `json:"newValue"`
	HumanReadableComment string `json:"humanReadableComment"`
}

type DynamoRecord struct {
	Code     string `json:"code"`     //mandatory for Dynamo db config table
	ItemType string `json:"itemType"` //mandatory for dynamo db config tabble
}

type DynamoHotelHistory struct {
	Code                 string `json:"code"`     //mandatory for Dynamo db config table
	ItemType             string `json:"itemType"` //mandatory for dynamo db config tabble
	TimeStamp            string `json:"timeStamp"`
	EventName            string `json:"eventName"`
	ObjectName           string `json:"objectName"`
	FieldName            string `json:"fieldName"`
	OldValue             string `json:"oldValue"`
	NewValue             string `json:"newValue"`
	HumanReadableComment string `json:"humanReadableComment"`
}

type DynamoHotelTimezone struct {
	Id   string `json:"id"`   //ex: Europe/Paris
	Name string `json:"name"` //ex: Central European Standard Time
}

type DynamoHotelRec struct {
	Code               string              `json:"code"`     //mandatory for Dynamo db config table
	ItemType           string              `json:"itemType"` //mandatory for dynamo db config tabble
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	CurrencyCode       string              `json:"currencyCode"`
	Lat                float64             `json:"lat"`
	Lng                float64             `json:"lng"`
	Timezone           DynamoHotelTimezone `json:"timezone"` //Hotel timezone
	DefaultTagLine     string              `json:"defaultTagLine"`
	DefaultDescription string              `json:"defaultDescription"`
	Bookable           string              `json:"bookable"`
	Shoppable          string              `json:"shoppable"`
	CheckInTime        string              `json:"checkInTime"`
	CheckOutTime       string              `json:"checkOuTime"`
	User               string              `json:"user",omitempty`
	LastUpdatedBy      string              `json:"lastUpdatedBy"`
	//PICTURE SPECIFIC ATTRIBUTES
	PictureId       string `json:"pictureId"`
	Tags            string `json:"tags"`
	Main            bool   `json:"main"`
	PictureItemCode string `json:"pictureItemCode"` //code of the business object the picture belongs to
	//ROOM SPECIFIC ATTRIBUTES
	Number       int64  `json:"number"`
	Type         string `json:"type"`
	Floor        int32  `json:"floor"`
	Building     int32  `json:"building"`
	BuildingName string `json:"buildingName"`
	Attributes   string `json:"attributes"`
	//ROOM TYPE, TAGS AND SELLABLE
	Quantity        int64   `json:"quantity"`
	Category        string  `json:"category"`
	LowPrice        float64 `json:"lowPrice"`
	MedPrice        float64 `json:"medPrice"`
	HighPrice       float64 `json:"highPrice"`
	PriceImpact     float64 `json:"priceImpact"`
	PriceUpdateUnit string  `json:"priceUpdateUnit"`
	InventoryType   string  `json:"inventoryType"`
	PricePerUnit    float64 `json:"pricePerUnit"`
	//HISTORY
	TimeStamp            string `json:"timeStamp"`
	EventName            string `json:"eventName"`
	ObjectName           string `json:"objectName"`
	FieldName            string `json:"fieldName"`
	OldValue             string `json:"oldValue"`
	NewValue             string `json:"newValue"`
	HumanReadableComment string `json:"humanReadableComment"`
}

//Use for write (this avoid having indexed attribute sent in request empty (like user))
type DynamoHotelPictureRec struct {
	Code     string `json:"code"`     //mandatory for Dynamo db config table
	ItemType string `json:"itemType"` //mandatory for dynamo db config tabble
	//PICTURE SPECIFIC ATTRIBUTES
	PictureId       string `json:"pictureId"`
	PictureItemCode string `json:"pictureItemCode"` //code of the business object the picture belongs to
	Tags            string `json:"tags"`
	Main            bool   `json:"main"`
	LastUpdatedBy   string `json:"lastUpdatedBy"`
}

type DynamoRoomRec struct {
	Code          string `json:"code"`     //mandatory for Dynamo db config table
	ItemType      string `json:"itemType"` //mandatory for dynamo db config tabble
	Number        int64  `json:"number"`
	Name          string `json:name"`
	Type          string `json:"type"`
	Floor         int32  `json:"floor"`
	Building      int32  `json:"building"`
	BuildingName  string `json:"buildingName"`
	Attributes    string `json:"attributes"`
	LastUpdatedBy string `json:"lastUpdatedBy"`
}

type DynamoSellableRec struct {
	Code          string  `json:"code"`
	ItemType      string  `json:"itemType"` //mandatory for dynamo db config tabble
	Category      string  `json:"category"`
	Name          string  `json:"name"`
	Quantity      int64   `json:"quantity"`
	InventoryType string  `json:"inventoryType"`
	PricePerUnit  float64 `json:"pricePerUnit"`
	Description   string  `json:"description"`
	LastUpdatedBy string  `json:"lastUpdatedBy"`
}

type DynamoRoomTypeRec struct {
	Code          string  `json:"code"`
	ItemType      string  `json:"itemType"` //mandatory for dynamo db config tabble
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	LowPrice      float64 `json:"lowPrice"`
	MedPrice      float64 `json:"medPrice"`
	HighPrice     float64 `json:"highPrice"`
	LastUpdatedBy string  `json:"lastUpdatedBy"`
}

type DynamoRoomAttributeRec struct {
	Code            string  `json:"code"`
	ItemType        string  `json:"itemType"` //mandatory for dynamo db config tabble
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Category        string  `json:"category"`
	PriceImpact     float64 `json:"priceImpact"`
	PriceUpdateUnit string  `json:"priceUpdateUnit"`
	LastUpdatedBy   string  `json:"lastUpdatedBy"`
}

//Implementing core.Batchable interface to allow this struc to be written in batch
func (drr DynamoRoomRec) GetPk() string {
	return drr.Code
}

//Requests and Respone to deserialize teh Json Payload into.
//ideally this should be a generic struct in teh core package with a specialised body but
//the design needs more thoughts so we will include it here for the moment
type RqWrapper struct {
	UserInfo    core.User `json:"userInfo"`
	SubFunction string    `json:"subFunction"`
	Id          string    `json:"id"` //ID to be used for GET request
	Request     Hotel     `json:"request"`
	StartDate   string    `json:"startDate"` //in case the config request involves a date range
	EndDate     string    `json:"endDate"`   //in case the config request involves a date range
}

type ResWrapper struct {
	Error       core.ResError      `json:"error"`
	SubFunction string             `json:"subFunction"`
	Response    []Hotel            `json:"response"`
	InvResponse map[string]Product `json:"inveResponse"`
}

//Inventory model for inventory display in config service
type Product struct {
	HotelCode string                          `json:"hotelCode"`
	Code      string                          `json:"code"`
	Type      string                          `json:"type"`
	Counts    map[string]map[string]Allotment `json:"counts"` //2019-05-14 => GINV => {CODE: GINV, BASECOUNT: 10, HELD: 2}
}

type Allotment struct {
	Code      string `json:"code"`
	BaseCount int    `json:"baseCount"`
	Held      int    `json:"held"`
}
