package adapter

import (
	iot "cloudrack-lambda-core/iot"
	model "cloudrack-lambda-fn/model"
)

var DEVICE_GROUP = "rigado_quickstart"

//This function creates a DynamoDB records with hte apropriate primary
func RequestToDynamo(rq model.Request, device iot.Device) model.DynamoRecord {
	dynRec := model.DynamoRecord{}
	dynRec.DeviceGroup = DEVICE_GROUP
	dynRec.SerialNumber = rq.SerialNumber
	dynRec.DeviceName = device.Name
	dynRec.ThingID = device.ID
	dynRec.CertificateArn = device.CertificateArn
	dynRec.CertificateID = device.CertificateId
	dynRec.CertificatePem = device.CertificatePem
	dynRec.PrivateKey = device.PrivateKey
	dynRec.PublicKey = device.PublicKey
	return dynRec
}

//this create a record with only PK and SK for retreive purpose
func CreateGenericDynamoRecord(rq model.Request) model.GenericDynamoRecord {
	dynRec := model.GenericDynamoRecord{}
	dynRec.DeviceGroup = DEVICE_GROUP
	dynRec.SerialNumber = rq.SerialNumber
	return dynRec
}
func DynamoToResponse(dynReq model.DynamoRecord, mqttEndpoint string) model.Response {
	res := model.Response{}
	res.SerialNumber = dynReq.SerialNumber
	res.DeviceName = dynReq.DeviceName
	res.ThingID = dynReq.ThingID
	res.Credential = model.DeviceCredentials{
		CertificateArn: dynReq.CertificateArn,
		CertificateId:  dynReq.CertificateID,
		CertificatePem: dynReq.CertificatePem,
		PrivateKey:     dynReq.PrivateKey,
		PubilKey:       dynReq.PublicKey,
	}
	res.MqttEndpoint = mqttEndpoint
	return res
}

func CreateResponse(rq model.Request, device iot.Device, mqttEndpoint string) model.Response {
	res := model.Response{}
	res.SerialNumber = rq.SerialNumber
	res.DeviceName = device.Name
	res.ThingID = device.ID
	res.Credential = model.DeviceCredentials{
		CertificateArn: device.CertificateArn,
		CertificateId:  device.CertificateId,
		CertificatePem: device.CertificatePem,
		PrivateKey:     device.PrivateKey,
		PubilKey:       device.PublicKey,
	}
	res.MqttEndpoint = mqttEndpoint
	return res
}

func ResponseToIotDevice(res model.Response) iot.Device {
	device := iot.Device{
		CertificateArn: res.Credential.CertificateArn,
		Name:           res.DeviceName,
		CertificateId:  res.Credential.CertificateId,
	}
	return device
}
