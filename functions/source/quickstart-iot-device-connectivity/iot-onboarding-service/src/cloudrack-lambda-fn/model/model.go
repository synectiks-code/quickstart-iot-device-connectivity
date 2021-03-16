package model

import (
	core "cloudrack-lambda-core/core"
)

type Request struct {
	SerialNumber string `json:"serialNumber"`
}

type Response struct {
	SerialNumber string            `json:"serialNumber"`
	DeviceName   string            `json:"deviceName"`
	ThingID      string            `json:"thingId"`
	Credential   DeviceCredentials `json:"credential"`
	MqttEndpoint string            `json:"mqttEndpoint"`
	Error        core.ResError     `json:"error"`
}

type DeviceCredentials struct {
	CertificateArn string `json:"certificateArn"`
	CertificateId  string `json:"certificateId"`
	CertificatePem string `json:"certificatePem"`
	PrivateKey     string `json:"privateKey"`
	PubilKey       string `json:"publicKey"`
}
