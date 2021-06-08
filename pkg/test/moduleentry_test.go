package test

import (
	"errors"
	"flag"
	"fmt"
	"github.com/yametech/devops/pkg/api"
	base3 "github.com/yametech/devops/pkg/api/action/base"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/base"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/store/mongo"
	"testing"
)

func TestCreateEntry(t *testing.T) {
	filter := map[string]interface{}{
		"spec.user": "weihengxing",
	}
	uuid := "6f2ac063797f484788f96b68727609f5"
	flag.Parse()
	store, _, _ := mongo.NewMongo("mongodb://10.200.10.46:27017/admin")
	baseService := service.NewBaseService(store)
	server := api.NewServer(baseService)
	b := base3.NewBaseServer("baseserver", server)
	privateModule := &base.PrivateModule{}
	if err := b.ModuleEntry.GetByFilter(common.DefaultNamespace, common.ModuleEntry, privateModule, filter); err != nil {
		fmt.Println(nil, errors.New("此用户数据在数据库中不存在！"))
		me := &base.PrivateModule{
			Metadata: core.Metadata{},
			Spec: base.PrivateModuleSpec{
				User:    "weihengxing",
				Modules: []string{uuid},
			},
		}
		_, err = b.ModuleEntry.Create(common.DefaultNamespace, common.ModuleEntry, me)
		if err != nil {
			fmt.Println(nil, err)
			return
		}
		module := &base.Module{}
		moduleSlice := make([]*base.Module, 0)
		err = b.ModuleEntry.GetByUUID(common.DefaultNamespace, common.AllModule, uuid, module)
		if err != nil {
			fmt.Println(nil, err)
			return
		}
		moduleSlice = append(moduleSlice, module)
		fmt.Println(moduleSlice, nil)
		return
	}
	modulates := &base.Module{}
	if err := b.ModuleEntry.GetByUUID(common.DefaultNamespace, common.AllModule, uuid, modulates); err != nil {
		fmt.Println(nil, errors.New("此模块的uuid在数据库中不存在！"))
		return
	}
	for _, v := range privateModule.Spec.Modules {
		if v == uuid {
			fmt.Println(nil, errors.New("此模块快捷入口已经存在！"))
			return
		}
	}
	privateModule.Spec.Modules = append(privateModule.Spec.Modules, uuid)
	_, judge, err := b.ModuleEntry.Apply(common.DefaultNamespace, common.ModuleEntry, privateModule.UUID, privateModule, true)
	if !judge && err != nil {
		fmt.Println(nil, err)
		return
	}
	moduleSlice := make([]*base.Module, 0)
	for _, v := range privateModule.Spec.Modules {
		module := &base.Module{}
		err := b.ModuleEntry.GetByUUID(common.DefaultNamespace, common.AllModule, v, module)
		if err != nil {
			fmt.Println(nil, err)
			return
		}
		moduleSlice = append(moduleSlice, module)
	}
	fmt.Println(moduleSlice, nil)
	return

}

func TestDeleteModule(t *testing.T) {
	filter := map[string]interface{}{
		"spec.user": "weihengxing",
	}

	uuid := "6f2ac063797f484788f96b68727609f5"
	flag.Parse()
	store, _, _ := mongo.NewMongo("mongodb://10.200.10.46:27017/admin")
	baseService := service.NewBaseService(store)
	server := api.NewServer(baseService)
	b := base3.NewBaseServer("baseserver", server)

	privateModule := &base.PrivateModule{}
	if err := b.ModuleEntry.GetByFilter(common.DefaultNamespace, common.ModuleEntry, privateModule, filter); err != nil {
		fmt.Println(nil, errors.New("此用户数据在数据库中不存在！"))
		return
	}
	fmt.Println(privateModule)
	modulates := &base.Module{}
	if err := b.ModuleEntry.GetByUUID(common.DefaultNamespace, common.AllModule, uuid, modulates); err != nil {
		fmt.Println(nil, errors.New("此模块的uuid在数据库中不存在！"))
		return
	}
	for k, v := range privateModule.Spec.Modules {
		if v == uuid {
			privateModule.Spec.Modules = append(privateModule.Spec.Modules[:k], privateModule.Spec.Modules[k+1:]...)
		}
	}
	_, judge, err := b.ModuleEntry.Apply(common.DefaultNamespace, common.ModuleEntry, privateModule.UUID, privateModule, true)
	if !judge && err != nil {
		fmt.Println(nil, err)
		return
	}
	moduleSlice := make([]*base.Module, 0)
	for _, v := range privateModule.Spec.Modules {
		module := &base.Module{}
		err := b.ModuleEntry.GetByUUID(common.DefaultNamespace, common.AllModule, v, module)
		if err != nil {
			fmt.Println(nil, err)
			return
		}
		moduleSlice = append(moduleSlice, module)
	}
	fmt.Println(moduleSlice, nil)
	return
}

func TestQueryModule(t *testing.T) {
	filter := map[string]interface{}{
		"spec.user": "weihengxing",
	}
	flag.Parse()
	store, _, _ := mongo.NewMongo("mongodb://10.200.10.46:27017/admin")
	baseService := service.NewBaseService(store)
	server := api.NewServer(baseService)
	b := base3.NewBaseServer("baseserver", server)

	privateModule := &base.PrivateModule{}
	if err := b.ModuleEntry.GetByFilter(common.DefaultNamespace, common.ModuleEntry, privateModule, filter); err != nil {
		fmt.Println(nil, errors.New("此用户数据在数据库中不存在！"))
		return
	}
	fmt.Println(privateModule)
	moduleSlice := make([]*base.Module, 0)
	for _, v := range privateModule.Spec.Modules {
		module := &base.Module{}
		err := b.ModuleEntry.GetByUUID(common.DefaultNamespace, common.AllModule, v, module)
		if err != nil {
			fmt.Println(nil, err)
			return
		}
		moduleSlice = append(moduleSlice, module)
	}
	fmt.Println(moduleSlice, nil)
	return
}
