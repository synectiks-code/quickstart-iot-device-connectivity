package usercase

import (
	db "cloudrack-lambda-core/db"
	iot "cloudrack-lambda-core/iot"
	adapter "cloudrack-lambda-fn/adapter"
	model "cloudrack-lambda-fn/model"
	"errors"
	"log"
)

//This function creates an IOT Thing from a device serial number.
//It generates cedentials (certificate + key pair)
//it returns a struct with the credentials, device serial number and IOT core
func Create(rq model.Request, dynSvc db.DBConfig, iotCore iot.Config, mqttEndpoint string) (model.Response, error) {
	//0-Checking if device with identical serial number exists. if yes we retuurn the Response
	log.Printf("[ONBOARDING][CREATE] checking if device has already been onboarded")
	res, err0 := Retrieve(rq, dynSvc, mqttEndpoint)
	//if the device has already been onboarded, we return the data directly, this allows teh service
	//to be idempotent which facilitates retry strategy
	if err0 == nil {
		log.Printf("[ONBOARDING][CREATE] device has already been onboarded")
		return res, nil
	}
	//0-Create the device using IOT Core SDK
	//TODO: make the device name prefix customizable in query
	log.Printf("[ONBOARDING][CREATE] Device does not exists, creating it")
	iotCredentals, err := iotCore.CreateDevice(rq.SerialNumber)
	if err != nil {
		return model.Response{}, err
	}
	//1-Create record in DynamoDB
	log.Printf("[ONBOARDING][CREATE] Creation success, Adding to DynamoDB")
	dynamoRecord := adapter.RequestToDynamo(rq, iotCredentals)
	_, err2 := dynSvc.Save(dynamoRecord)
	if err2 != nil {
		return model.Response{}, err2
	}
	//2-Build response
	log.Printf("[ONBOARDING][CREATE] Returning response")
	response := adapter.CreateResponse(rq, iotCredentals, mqttEndpoint)
	return response, nil
}

//Delete the device
func Delete(rq model.Request, dynSvc db.DBConfig, iotCore iot.Config) (model.Response, error) {
	res, err := Retrieve(rq, dynSvc, "")
	if err != nil {
		return model.Response{}, err
	}
	device := adapter.ResponseToIotDevice(res)
	err = iotCore.DeleteDevice(device)
	if err == nil {
		genRec := adapter.CreateGenericDynamoRecord(rq)
		dynSvc.Delete(genRec)
	}
	return model.Response{}, err
}

func Retrieve(rq model.Request, dynSvc db.DBConfig, mqttEndpoint string) (model.Response, error) {
	//this create a record with only PK and SK for retreive purpose
	log.Printf("[ONBOARDING][RETREIVE] Retreiveing device %s", rq.SerialNumber)
	genRec := adapter.CreateGenericDynamoRecord(rq)
	dynRec := model.DynamoRecord{}
	log.Printf("[ONBOARDING][RETREIVE] checking existence of record PK:%s and SK:%s", genRec.DeviceGroup, genRec.SerialNumber)
	err := dynSvc.Get(genRec.DeviceGroup, genRec.SerialNumber, &dynRec)
	if err != nil {
		log.Printf("[ONBOARDING][RETREIVE] error while getting record PK:%s and SK:%s. Msg: %v", genRec.DeviceGroup, genRec.SerialNumber, err)
		return model.Response{}, err
	}
	if dynRec.SerialNumber == "" {
		log.Printf("[ONBOARDING][RETREIVE] Record PK:%s and SK:%s not found", genRec.DeviceGroup, genRec.SerialNumber)
		return model.Response{}, errors.New("No Device found with serial Numbber " + rq.SerialNumber)
	}
	response := adapter.DynamoToResponse(dynRec, mqttEndpoint)
	return response, nil

}
