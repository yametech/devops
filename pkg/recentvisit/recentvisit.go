package recentvisit

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/resource/base"
	"github.com/yametech/devops/pkg/service"
	"log"
)

func RecentVisit(service service.IService) gin.HandlerFunc {
	return func(g *gin.Context) {
		uuid := g.Query("uuid")
		filter := map[string]interface{}{}
		user := ""
		filter["spec.user"] = user
		privateModule := &base.PrivateModule{}
		err := service.GetByFilter(common.DefaultNamespace, common.RecentVisit, privateModule, filter)
		if err != nil {
			log.Println("没有找到当前用户！")
		}
		if len(privateModule.Spec.Modules) < 6 {
			privateModule.Spec.Modules = append(privateModule.Spec.Modules, uuid)
			privateModule.Spec.User = user
			_, judge, err := service.Apply(common.DefaultNamespace, common.RecentVisit, privateModule.UUID, privateModule, true)
			if !judge && err != nil {
				api.ResponseError(g, err)
			}
		} else {
			privateModule.Spec.Modules = append(privateModule.Spec.Modules[1:], uuid)
			_, judge, err := service.Apply(common.DefaultNamespace, common.RecentVisit, privateModule.UUID, privateModule, true)
			if !judge && err != nil {
				api.ResponseError(g, err)
			}
		}
	}
}
