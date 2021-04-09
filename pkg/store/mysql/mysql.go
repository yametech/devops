package mysql

import (
	"fmt"
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

func (m *Mysql) Save(obj interface{}) error {
	tx := m.Db.Save(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (m *Mysql) Update(src, dst interface{}) error {
	tx := m.Db.Model(src).Updates(dst)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (m *Mysql) List(result interface{}) error {
	tx := m.Db.Find(result)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (m *Mysql) GetByFilter(filter map[string]interface{}, result interface{}) error {
	tx := m.Db.Where(filter).Find(result)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return tx.Error
	}
	return nil
}

func (m *Mysql) Del(obj interface{}) error {
	tx := m.Db.Delete(obj)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (m *Mysql) Apply(object interface{}) error {
	tx := m.Db.Create(object)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

//
//func (m *Mysql) Get(table string, filter map[string]interface{}, result interface{}) error {
//	tx := m.Db.Table(table).Where(filter).Find(result)
//	if tx.Error != nil {
//		return tx.Error
//	}
//	return nil
//}

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
