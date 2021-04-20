package workorder

import (
	"fmt"
	"github.com/yametech/devops/pkg/core"
	"time"
)

type OrderStatus uint8

const (
	WaitCommit OrderStatus = iota // 待提交
	Checking                      // 审核中
	Rejected                      // 驳回
	Finish                        // 完成
	Running                       // 执行中
	Failed                        // 执行失败
	Withdrawn                     // 已撤回
)

type OrderType uint8

const (
	Resources OrderType = iota // 服务配置工单
	Namespace                  // 资源池配置工单
)

var mark map[int64]interface{}

func init() {
	mark = map[int64]interface{}{
		0: "RS",
		1: "NS",
	}
}

type WorkOrder struct {
	core.Metadata `json:"metadata"`
	Spec          WorkOrderSpec `json:"spec"`
}

type WorkOrderSpec struct {
	OrderType   `json:"order_type" bson:"order_type"`
	OrderStatus `json:"order_status" bson:"order_status"`
	Number      string                 `json:"number" bson:"number"`
	Title       string                 `json:"title" bson:"title"`
	Creator     string                 `json:"creator" bson:"creator"`
	UpdateTime  int64                  `json:"update_time" bson:"update_time"`
	Attribute   map[string]interface{} `json:"attribute" bson:"attribute"`
	Apply       map[string]interface{} `json:"apply" bson:"apply"`
	Check       map[string]interface{} `json:"check" bson:"check"`
	Result      map[string]interface{} `json:"result" bson:"result"`
}

func (w *WorkOrder) GenerateNumber()  {
	today := time.Now()
	w.Spec.Number = fmt.Sprintf("%d%s%d", today.Year(), today.Month(), today.Day())
	fmt.Println(w.Spec.Number)
}