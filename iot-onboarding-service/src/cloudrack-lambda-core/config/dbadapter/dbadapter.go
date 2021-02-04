package dbadapter

import (
	model "cloudrack-lambda-core/config/model"
	core "cloudrack-lambda-core/core"
	"strconv"
	"strings"
)

//“cfg”-<”general” | “history” | nested_object_type>-<“picture”>-<id>
var ITEM_TYPE_CONFIG_PREFIX string = "cfg"
var ITEM_TYPE_CONFIG_GENERAL string = ITEM_TYPE_CONFIG_PREFIX + "-general"
var ITEM_TYPE_CONFIG_HISTORY string = ITEM_TYPE_CONFIG_PREFIX + "-history"
var ITEM_TYPE_CONFIG_GENERAL_PICTURE string = ITEM_TYPE_CONFIG_GENERAL + "-picture"
var ITEM_TYPE_CONFIG_INVENTORY_ROOM string = ITEM_TYPE_CONFIG_PREFIX + "-inventory-room"
var ITEM_TYPE_CONFIG_INVENTORY_SELLABLE string = ITEM_TYPE_CONFIG_PREFIX + "-inventory-sellable"
var ITEM_TYPE_CONFIG_INVENTORY_SELLABLE_PICTURE string = ITEM_TYPE_CONFIG_INVENTORY_SELLABLE + "-picture"
var ITEM_TYPE_CONFIG_INVENTORY_ROOM_TYPE string = ITEM_TYPE_CONFIG_INVENTORY_ROOM + "-type"
var ITEM_TYPE_CONFIG_INVENTORY_ROOM_TYPE_PICTURE string = ITEM_TYPE_CONFIG_INVENTORY_ROOM_TYPE + "-picture"
var ITEM_TYPE_CONFIG_INVENTORY_ROOM_ATTRIBUTE string = ITEM_TYPE_CONFIG_INVENTORY_ROOM + "-attribute"

func IdToDynamo(id string) model.DynamoHotelRec {
	return model.DynamoHotelRec{Code: id, ItemType: ITEM_TYPE_CONFIG_GENERAL}
}

func DynamoPendingChangesForDelete(changes []model.DynamoHotelHistory) []model.DynamoRecord {
	recs := make([]model.DynamoRecord, 0, 0)
	for _, change := range changes {
		recs = append(recs, model.DynamoRecord{Code: change.Code, ItemType: change.ItemType})
	}
	return recs
}

func DynamoToBom(dbr model.DynamoHotelRec) model.Hotel {
	bom := model.Hotel{Code: dbr.Code,
		Name:         dbr.Name,
		Description:  dbr.Description,
		CurrencyCode: dbr.CurrencyCode,
		Lat:          dbr.Lat,
		Lng:          dbr.Lng,
		Timezone:     model.HotelTimezone{Id: dbr.Timezone.Id, Name: dbr.Timezone.Name},
		DefaultMarketting: model.HotelMarketing{
			DefaultTagLine:     dbr.DefaultTagLine,
			DefaultDescription: dbr.DefaultDescription},
		Options: model.HotelOptions{
			CheckInTime:  dbr.CheckInTime,
			CheckOutTime: dbr.CheckOutTime,
		}}

	bom.Options.Bookable, _ = strconv.ParseBool(dbr.Bookable)
	bom.Options.Shoppable, _ = strconv.ParseBool(dbr.Shoppable)

	return bom
}

func DynamoPicToBom(dbr model.DynamoHotelRec) model.HotelPicture {
	bom := model.HotelPicture{
		Id:              dbr.PictureId,
		Tags:            strings.Split(dbr.Tags, "|"),
		PictureItemCode: dbr.PictureItemCode,
		Main:            dbr.Main}
	return bom
}

func DynamoListToBom(dbrList []model.DynamoHotelRec) model.Hotel {
	var bom model.Hotel
	var pics []model.HotelPicture = make([]model.HotelPicture, 0, 0)
	var hotelConfigChanges []model.HotelConfigChange = make([]model.HotelConfigChange, 0, 0)
	invDbr := make([]model.DynamoHotelRec, 0, 0)
	roomTypePictures := make(map[string][]model.HotelPicture)
	sellablePictures := make(map[string][]model.HotelPicture)
	tagsByCode := make(map[string]model.HotelRoomAttribute)

	for _, dbr := range dbrList {
		if strings.HasPrefix(dbr.ItemType, ITEM_TYPE_CONFIG_INVENTORY_SELLABLE_PICTURE) {
			sellablePictures[dbr.PictureItemCode] = append(roomTypePictures[dbr.PictureItemCode], DynamoPicToBom(dbr))
		} else if strings.HasPrefix(dbr.ItemType, ITEM_TYPE_CONFIG_INVENTORY_ROOM_TYPE_PICTURE) {
			roomTypePictures[dbr.PictureItemCode] = append(sellablePictures[dbr.PictureItemCode], DynamoPicToBom(dbr))
		} else if strings.HasPrefix(dbr.ItemType, ITEM_TYPE_CONFIG_INVENTORY_ROOM_TYPE) {
			bom.RoomTypes = append(bom.RoomTypes, DynamoToBomRoomType(dbr))
		} else if strings.HasPrefix(dbr.ItemType, ITEM_TYPE_CONFIG_INVENTORY_ROOM_ATTRIBUTE) {
			tag := DynamoToBomTag(dbr)
			bom.Tags = append(bom.Tags, tag)
			tagsByCode[tag.Code] = tag
		} else if strings.HasPrefix(dbr.ItemType, ITEM_TYPE_CONFIG_INVENTORY_ROOM) {
			invDbr = append(invDbr, dbr)
		} else if strings.HasPrefix(dbr.ItemType, ITEM_TYPE_CONFIG_INVENTORY_SELLABLE) {
			bom.Sellables = append(bom.Sellables, DynamoToBomSellable(dbr))
		} else if strings.HasPrefix(dbr.ItemType, ITEM_TYPE_CONFIG_GENERAL_PICTURE) {
			pics = append(pics, DynamoPicToBom(dbr))
		} else if dbr.ItemType == ITEM_TYPE_CONFIG_GENERAL {
			bom = DynamoToBom(dbr)
		} else if strings.HasPrefix(dbr.ItemType, ITEM_TYPE_CONFIG_HISTORY) {
			hotelConfigChanges = append(hotelConfigChanges, DynamoToBomConfigChange(dbr))
		}

	}

	bom.Pictures = pics
	bom.Buildings = DynamoRoomListToBom(invDbr)
	bom.AddRoomTypePictures(roomTypePictures)
	bom.AddSellablePictures(sellablePictures)
	bom.AddAttributesToRooms(tagsByCode)
	bom.PendingChanges = hotelConfigChanges
	return bom
}

func BomToDynamo(bom model.Hotel, user core.User) model.DynamoHotelRec {
	return model.DynamoHotelRec{Code: bom.Code,
		Name:               bom.Name,
		Description:        bom.Description,
		CurrencyCode:       bom.CurrencyCode,
		Lat:                bom.Lat,
		Lng:                bom.Lng,
		Timezone:           model.DynamoHotelTimezone{Id: bom.Timezone.Id, Name: bom.Timezone.Name},
		DefaultTagLine:     bom.DefaultMarketting.DefaultTagLine,
		DefaultDescription: bom.DefaultMarketting.DefaultDescription,
		Bookable:           strconv.FormatBool(bom.Options.Bookable),
		Shoppable:          strconv.FormatBool(bom.Options.Shoppable),
		CheckInTime:        bom.Options.CheckInTime,
		CheckOutTime:       bom.Options.CheckOutTime,
		User:               user.Username,
		LastUpdatedBy:      user.Username,
		ItemType:           ITEM_TYPE_CONFIG_GENERAL}
}

func BomToDynamoPic(hotelCode string, user string, pic model.HotelPicture) model.DynamoHotelPictureRec {
	return model.DynamoHotelPictureRec{
		Code:            hotelCode,
		PictureItemCode: hotelCode,
		PictureId:       pic.Id,
		ItemType:        ITEM_TYPE_CONFIG_GENERAL_PICTURE + "-" + pic.Id,
		Tags:            strings.Join(pic.Tags, "|"),
		Main:            pic.Main,
		LastUpdatedBy:   user,
	}
}

func BomRoomTypeToDynamoPic(hotelCode string, roomTypeCode string, user string, pic model.HotelPicture) model.DynamoHotelPictureRec {
	return model.DynamoHotelPictureRec{
		Code:            hotelCode,
		ItemType:        ITEM_TYPE_CONFIG_INVENTORY_ROOM_TYPE_PICTURE + "-" + pic.Id,
		PictureId:       pic.Id,
		PictureItemCode: roomTypeCode,
		Tags:            strings.Join(pic.Tags, "|"),
		Main:            pic.Main,
		LastUpdatedBy:   user,
	}
}

func BomSellableToDynamoPic(hotelCode string, sellableCode string, user string, pic model.HotelPicture) model.DynamoHotelPictureRec {
	return model.DynamoHotelPictureRec{
		Code:            hotelCode,
		ItemType:        ITEM_TYPE_CONFIG_INVENTORY_SELLABLE_PICTURE + "-" + pic.Id,
		PictureId:       pic.Id,
		PictureItemCode: sellableCode,
		Tags:            strings.Join(pic.Tags, "|"),
		Main:            pic.Main,
		LastUpdatedBy:   user,
	}
}

func BomRoomToDynamoList(bom model.Hotel, user core.User) []model.DynamoRoomRec {
	dynRooms := make([]model.DynamoRoomRec, 0, 0)
	for i, building := range bom.Buildings {
		for j, floor := range building.Floors {
			for k, bomRoom := range floor.Rooms {
				dynRooms = append(dynRooms, model.DynamoRoomRec{
					Code:          bom.Code,
					ItemType:      ITEM_TYPE_CONFIG_INVENTORY_ROOM + "-" + BuildRoomId(i, j, k),
					Number:        bomRoom.Number,
					Name:          bomRoom.Name,
					Type:          bomRoom.Type,
					Floor:         floor.Num,
					Building:      int32(i),
					BuildingName:  building.Name,
					Attributes:    attributesToDynamo(bomRoom.Attributes),
					LastUpdatedBy: user.Username,
				})
			}
		}
	}
	return dynRooms
}

func BuildRoomId(building int, floor int, room int) string {
	return strconv.FormatInt(int64(building), 10) + "-" + strconv.FormatInt(int64(floor), 10) + "-" + strconv.FormatInt(int64(room), 10)
}

func DynamoRoomListToBom(roomList []model.DynamoHotelRec) []model.HotelBuilding {
	buildings := make([]model.HotelBuilding, 0, 0)
	for _, room := range roomList {
		nBuildings := int32(len(buildings))
		if nBuildings <= room.Building {
			//extending building array to fit size
			buildings = append(buildings, make([]model.HotelBuilding, room.Building+1-nBuildings, room.Building+1-nBuildings)...)
		}
		buildings[room.Building].Name = room.BuildingName
		nFloors := int32(len(buildings[room.Building].Floors))
		if nFloors < room.Floor {
			//extending floor array to fit size
			buildings[room.Building].Floors = append(buildings[room.Building].Floors, make([]model.HotelFloor, room.Floor-nFloors, room.Floor-nFloors)...)
		}
		//CAREFUL: floor num = floor position is array +1 contrary to buillding
		buildings[room.Building].Floors[room.Floor-1].Num = room.Floor
		buildings[room.Building].Floors[room.Floor-1].Rooms = append(buildings[room.Building].Floors[room.Floor-1].Rooms, model.HotelRoom{
			Number:     room.Number,
			Name:       room.Name,
			Type:       room.Type,
			Attributes: dynamoToAttributes(room.Attributes)})
	}
	return buildings

}

func attributesToDynamo(attributes []model.HotelRoomAttribute) string {
	codes := make([]string, len(attributes), len(attributes))
	for i, attr := range attributes {
		codes[i] = attr.Code
	}
	return strings.Join(codes, "|")
}

func BomSellableToDynamo(bom model.Hotel, user core.User) model.DynamoSellableRec {
	return model.DynamoSellableRec{
		Code:          bom.Code,
		Name:          bom.Sellables[0].Name,
		Description:   bom.Sellables[0].Description,
		Category:      bom.Sellables[0].Category,
		Quantity:      bom.Sellables[0].Quantity,
		InventoryType: bom.Sellables[0].InventoryType,
		PricePerUnit:  bom.Sellables[0].PricePerUnit,
		ItemType:      ITEM_TYPE_CONFIG_INVENTORY_SELLABLE + "-" + bom.Sellables[0].Code,
		LastUpdatedBy: user.Username,
	}
}

func BomRoomTypeToDynamo(bom model.Hotel, user core.User) model.DynamoRoomTypeRec {
	return model.DynamoRoomTypeRec{
		Code:          bom.Code,
		Name:          bom.RoomTypes[0].Name,
		Description:   bom.RoomTypes[0].Description,
		LowPrice:      bom.RoomTypes[0].LowPrice,
		MedPrice:      bom.RoomTypes[0].MedPrice,
		HighPrice:     bom.RoomTypes[0].HighPrice,
		ItemType:      ITEM_TYPE_CONFIG_INVENTORY_ROOM_TYPE + "-" + bom.RoomTypes[0].Code,
		LastUpdatedBy: user.Username,
	}
}

func BomTagToDynamo(bom model.Hotel, user core.User) model.DynamoRoomAttributeRec {
	return model.DynamoRoomAttributeRec{
		Code:            bom.Code,
		Name:            bom.Tags[0].Name,
		Description:     bom.Tags[0].Description,
		Category:        bom.Tags[0].Category,
		PriceImpact:     bom.Tags[0].PriceImpact,
		PriceUpdateUnit: bom.Tags[0].PriceUpdateUnit,
		ItemType:        ITEM_TYPE_CONFIG_INVENTORY_ROOM_ATTRIBUTE + "-" + bom.Tags[0].Code,
		LastUpdatedBy:   user.Username,
	}
}

//Bom to generic Dynamo reccord to be used for delete
func BomRoomTypeToDynamoRecord(bom model.Hotel) model.DynamoRecord {
	return model.DynamoRecord{
		Code:     bom.Code,
		ItemType: ITEM_TYPE_CONFIG_INVENTORY_ROOM_TYPE + "-" + bom.RoomTypes[0].Code}
}

//Bom to generic Dynamo reccord to be used for delete
func BomTagToDynamoRecord(bom model.Hotel) model.DynamoRecord {
	return model.DynamoRecord{
		Code:     bom.Code,
		ItemType: ITEM_TYPE_CONFIG_INVENTORY_ROOM_ATTRIBUTE + "-" + bom.Tags[0].Code}
}

//Bom to generic Dynamo reccord to be used for delete
func BomSellableToDynamoRecord(bom model.Hotel) model.DynamoRecord {
	return model.DynamoRecord{
		Code:     bom.Code,
		ItemType: ITEM_TYPE_CONFIG_INVENTORY_SELLABLE + "-" + bom.Sellables[0].Code}
}

//Bom to generic Dynamo reccord to be used for delete
func BomRoomTypePictureToDynamoRecord(hotelCode string, bom model.HotelPicture) model.DynamoRecord {
	return model.DynamoRecord{
		Code:     hotelCode,
		ItemType: ITEM_TYPE_CONFIG_INVENTORY_ROOM_TYPE_PICTURE + "-" + bom.Id}
}

//Bom to generic Dynamo reccord to be used for delete
func BomSellablePictureToDynamoRecord(hotelCode string, bom model.HotelPicture) model.DynamoRecord {
	return model.DynamoRecord{
		Code:     hotelCode,
		ItemType: ITEM_TYPE_CONFIG_INVENTORY_SELLABLE_PICTURE + "-" + bom.Id}
}

//Bom to generic Dynamo reccord to be used for delete
func BomPictureToDynamoRecord(hotelCode string, bom model.HotelPicture) model.DynamoRecord {
	return model.DynamoRecord{
		Code:     hotelCode,
		ItemType: ITEM_TYPE_CONFIG_GENERAL_PICTURE + "-" + bom.Id}
}

func DynamoToBomSellable(dbr model.DynamoHotelRec) model.HotelSellable {
	return model.HotelSellable{
		Code:          getDbrCode(dbr.ItemType),
		Name:          dbr.Name,
		Description:   dbr.Description,
		InventoryType: dbr.InventoryType,
		PricePerUnit:  dbr.PricePerUnit,
		Category:      dbr.Category,
		Quantity:      dbr.Quantity}
}

func DynamoToBomRoomType(dbr model.DynamoHotelRec) model.HotelRoomType {
	return model.HotelRoomType{
		Code:        getDbrCode(dbr.ItemType),
		Name:        dbr.Name,
		LowPrice:    dbr.LowPrice,
		MedPrice:    dbr.MedPrice,
		HighPrice:   dbr.HighPrice,
		Description: dbr.Description}
}

func DynamoToBomTag(dbr model.DynamoHotelRec) model.HotelRoomAttribute {
	return model.HotelRoomAttribute{
		Code:            getDbrCode(dbr.ItemType),
		Name:            dbr.Name,
		PriceImpact:     dbr.PriceImpact,
		PriceUpdateUnit: dbr.PriceUpdateUnit,
		Description:     dbr.Description,
		Category:        dbr.Category}
}

func DynamoToBomConfigChange(histRec model.DynamoHotelRec) model.HotelConfigChange {
	return model.HotelConfigChange{
		TimeStamp:            histRec.TimeStamp,
		EventName:            histRec.EventName,
		ObjectName:           histRec.ObjectName,
		FieldName:            histRec.FieldName,
		OldValue:             histRec.OldValue,
		NewValue:             histRec.NewValue,
		HumanReadableComment: histRec.HumanReadableComment,
	}
}

//returns code form bynamo dbitem by splitinting itemType field and getting teh last element following our datamodel convention
func getDbrCode(itemType string) string {
	return strings.Split(itemType, "-")[len(strings.Split(itemType, "-"))-1]
}

func dynamoToAttributes(attributeString string) []model.HotelRoomAttribute {
	attributes := make([]model.HotelRoomAttribute, 0, 0)
	for _, code := range strings.Split(attributeString, "|") {
		attributes = append(attributes, model.HotelRoomAttribute{Code: code})
	}
	return attributes
}
