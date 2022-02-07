package entity

import "github.com/kamva/mgm/v3"

type APIKey struct {
	mgm.DefaultModel `bson:",inline"`
	Key              string `json:"key" bson:"key"`
	EnableAt         string `json:"enableAt" bson:"enableAt"`
}
