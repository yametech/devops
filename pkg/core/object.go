package core

import (
	"github.com/yametech/devops-zpk-server/pkg/utils"
	"gorm.io/datatypes"
	"time"
)

type IObject interface {
	GetUUID() string
	GetKind() string
}

type Metadata struct {
	Name        string `json:"name" bson:"name" gorm:"size:255;not null"`
	Kind        string `json:"kind"  bson:"kind" gorm:"default:'';size:255"`
	UUID        string `json:"uuid" bson:"uuid" gorm:"size:255;not null;primaryKey"`
	Version     int64  `json:"version" bson:"version" gorm:"autoUpdateTime:nano"`
	IsDelete    bool   `json:"is_delete" bson:"is_delete" gorm:"type:bool;default:false"`
	CreatedTime int64  `json:"created_time" bson:"created_time" gorm:"autoCreateTime"`

	Labels datatypes.JSON `json:"labels" bson:"labels" gorm:"type:json"`
}

func (m *Metadata) GetKind() string {
	return m.Kind
}

func (m *Metadata) GenerateVersion() IObject {
	m.CreatedTime = time.Now().Unix()
	if m.UUID == "" {
		m.UUID = utils.NewSUID().String()
	}
	return m
}

func (m *Metadata) GetUUID() string {
	return m.UUID
}
