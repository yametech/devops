package appservice

import (
	"github.com/yametech/devops/pkg/core"
)

type AppConfig struct {
	core.Metadata `json:"metadata"`
	Spec          AppConfigSpec `json:"spec"`
}

type AppConfigSpec struct {
	App    string                 `json:"app" bson:"app"`
	Config map[string]interface{} `json:"config" bson:"config"`
}

func (ap *AppConfig) Clone() core.IObject {
	result := &AppConfig{}
	core.Clone(ap, result)
	return result
}

type HistoryStatus uint8

const (
	Success HistoryStatus = iota
	Failed
)

type AppConfigHistory struct {
	core.Metadata `json:"metadata"`
	Spec          AppConfigHistorySpec `json:"spec"`
}

type AppConfigHistorySpec struct {
	App           string                 `json:"app" bson:"app"`
	History       map[string]interface{} `json:"history" bson:"history"`
	HistoryStatus `json:"history_status" bson:"history_status"`
}

func (a *AppConfigHistory) Clone() core.IObject {
	result := &AppConfigHistory{}
	core.Clone(a, result)
	return result
}
