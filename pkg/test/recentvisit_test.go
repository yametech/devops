package test

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/yametech/devops/pkg/api"
	base3 "github.com/yametech/devops/pkg/api/action/base"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/resource/base"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/store/mongo"
	"testing"
)

func TestCreateRecent(t *testing.T) {
	user := ""
	uuid := "7ff083d39ea14d39ad51124d7b2cc7b3"
	state1 := false
	state2 := false
	flag.Parse()
	store, _, _ := mongo.NewMongo("mongodb://10.200.10.46:27017/admin")
	baseService := service.NewBaseService(store)
	server := api.NewServer(baseService)
	b := base3.NewBaseServer("baseserver", server)
	privateModule := &base.PrivateModule{}
	filter := map[string]interface{}{}
	filter["spec.user"] = user
	err := b.RecentVisit.GetByFilter(common.DefaultNamespace, common.RecentVisit, privateModule, filter)
	if err != nil {
		fmt.Println(nil, errors.New("没有找到当前用户！"))
	}
	modulates := &base.Module{}
	if err := b.RecentVisit.GetByUUID(common.DefaultNamespace, common.AllModule, uuid, modulates); err != nil {
		fmt.Println(nil, errors.New("此模块的uuid在数据库中不存在！"))
	}
	if len(privateModule.Spec.Modules) < 6 {
		for k, _ := range privateModule.Spec.Modules {
			if privateModule.Spec.Modules[k] == uuid {
				privateModule.Spec.Modules = append(privateModule.Spec.Modules[:k], privateModule.Spec.Modules[k+1:]...)
				privateModule.Spec.Modules = append(privateModule.Spec.Modules, uuid)
				privateModule.Spec.User = user
				_, judge, err := b.RecentVisit.Apply(common.DefaultNamespace, common.RecentVisit, privateModule.UUID, privateModule, true)
				if !judge && err != nil {
					fmt.Println(nil, errors.New("最近访问更新失败！"))
				}
				state1 = true
			}
		}
		if state1 == false {
			fmt.Println(666)
			privateModule.Spec.Modules = append(privateModule.Spec.Modules, uuid)
			_, judge, err := b.RecentVisit.Apply(common.DefaultNamespace, common.RecentVisit, privateModule.UUID, privateModule, true)
			if !judge && err != nil {
				fmt.Println(nil, errors.New("最近访问更新失败！"))
			}
		}

	} else {
		for k, _ := range privateModule.Spec.Modules {
			if privateModule.Spec.Modules[k] == uuid {
				privateModule.Spec.Modules = append(privateModule.Spec.Modules[:k], privateModule.Spec.Modules[k+1:]...)
				privateModule.Spec.Modules = append(privateModule.Spec.Modules, uuid)
				_, judge, err := b.RecentVisit.Apply(common.DefaultNamespace, common.RecentVisit, privateModule.UUID, privateModule, true)
				if !judge && err != nil {
					fmt.Println(nil, errors.New("最近访问更新失败！"))
				}
				state2 = true
			}
		}
		if state2 == false {
			privateModule.Spec.Modules = append(privateModule.Spec.Modules[1:], uuid)
			_, judge, err := b.RecentVisit.Apply(common.DefaultNamespace, common.RecentVisit, privateModule.UUID, privateModule, true)
			if !judge && err != nil {
				fmt.Println(nil, errors.New("最近访问更新失败！"))
			}
		}
	}
}
