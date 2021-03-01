package dbadapter

import (
	//core "cloudrack-lambda-core/core"
	model "cloudrack-lambda-core/avail/model"
	//"strconv"
	//"strings"
)

func BomToDynamoProductRecord(hotelCode string, bom model.Product) model.DynamoProductRecord{
	return model.DynamoProductRecord{ 
		Code : hotelCode,
	}
}

//determine whether a product is available for a given avail request based on inventory and config data
func (p model.Product) IsAvailable(rq modelAvailRequest) bool{
	//TO IMPLEMENT
	return true
}




