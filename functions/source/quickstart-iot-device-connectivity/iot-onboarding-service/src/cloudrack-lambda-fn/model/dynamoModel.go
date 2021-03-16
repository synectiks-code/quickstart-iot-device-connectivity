package model

type DynamoRecord struct {
	DeviceGroup    string `json:"deviceGroup"`
	SerialNumber   string `json:"serialNumber"`
	DeviceName     string `json:"deviceName"`
	ThingID        string `json:"thingId"`
	CertificateArn string `json:"certificateArn"`
	CertificateID  string `json:"certificateId"`
	CertificatePem string `json:"certificatePem"`
	PrivateKey     string `json:"privateKey"`
	PublicKey      string `json:"publicKey"`
}

type GenericDynamoRecord struct {
	DeviceGroup  string `json:"deviceGroup"`
	SerialNumber string `json:"serialNumber"`
}
