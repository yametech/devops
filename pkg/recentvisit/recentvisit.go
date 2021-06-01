package recentvisit

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yametech/devops/pkg/api"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/resource/base"
	"github.com/yametech/devops/pkg/service"
)

func RecentVisit(service service.IService) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.Query("uuid")
		user := ""
		filter := map[string]interface{}{"spec.user": user}
		privateModule := &base.PrivateModule{}
		err := service.GetByFilter(common.DefaultNamespace, common.RecentVisit, privateModule, filter)
		if err != nil {
			log.Println("没有找到当前用户！")
			return
		}
		if len(privateModule.Spec.Modules) < 6 {
			privateModule.Spec.Modules = append(privateModule.Spec.Modules, uuid)
			privateModule.Spec.User = user
			_, judge, _err := service.Apply(common.DefaultNamespace, common.RecentVisit, privateModule.UUID, privateModule, true)
			if !judge && _err != nil {
				api.ResponseError(c, _err)
				return
			}
		} else {
			privateModule.Spec.Modules = append(privateModule.Spec.Modules[1:], uuid)
			_, judge, _err := service.Apply(common.DefaultNamespace, common.RecentVisit, privateModule.UUID, privateModule, true)
			if !judge && _err != nil {
				api.ResponseError(c, _err)
				return
			}
		}
	}
}
