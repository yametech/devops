package mysql

import (
	"fmt"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource"
	"github.com/yametech/devops/pkg/store"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Mysql struct {
	Uri string
	Db  *gorm.DB
}

func (m *Mysql) List(db, table string, result interface{}) error {
	tx := m.Db.Table(table).Find(&result)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (m *Mysql) GetByFilter(db, table string, filter map[string]interface{}, result interface{}) error {
	panic("implement me")
}

func (m *Mysql) Del(db, table string, object interface{}) error {
	panic("implement me")
}

func (m *Mysql) Apply(db, table string, object core.IObject) error {
	tx := m.Db.Create(object)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

var _ store.IStore = &Mysql{}

func Setup(uri, user, pw, database string, errC chan<- error) (*Mysql, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pw, uri, database)
	//	dsn := "root:123456@tcp(127.0.0.1:3306)/ccmose?charset=utf8mb4&parseTime=True&loc=Local"
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := conn.DB()
	if err != nil {
		return nil, err
	}
	err = conn.AutoMigrate(&resource.User{}, &resource.Artifact{})

	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxIdleTime(time.Second * 600)

	go func() {
		for {
			time.Sleep(1 * time.Second)

			if err := sqlDB.Ping(); err != nil {
				_ = sqlDB.Close()
				errC <- err
			}
		}
	}()
	return &Mysql{Uri: uri, Db: conn}, nil
}
