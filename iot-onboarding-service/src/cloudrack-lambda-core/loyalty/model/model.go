package model

import (
	core "cloudrack-lambda-core/core"
	"encoding/json"
)

type RqWrapper struct {
	Id      string  `json:"id"` //ID to be used for GET request
	Profile Profile `json:"profile"`
}

type ResWrapper struct {
	Error    core.ResError `json:"error"`
	Profiles []Profile     `json:"profiles"`
	Status   string        `json:"status"` // overall transaction status success/failure
}

func (p ResWrapper) Decode(dec json.Decoder) (error, core.CloudrackObject) {
	return dec.Decode(&p), p
}

type Profile struct {
	ID        string  `json:"id"`
	Points    float64 `json:"points"`
	Email     string  `json:"email"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Status    string  `json:"status"`
}
