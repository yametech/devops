package workorder

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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

var Mark = map[OrderType]interface{}{
	Resources: "RS",
	Namespace: "NS",
}

type WorkOrder struct {
	core.Metadata `json:"metadata"`
	Spec          Spec `json:"spec"`
}

type Spec struct {
	OrderType   `json:"order_type" bson:"order_type"`
	OrderStatus `json:"order_status" bson:"order_status"`
	Number      string                 `json:"number" bson:"number"`
	Title       string                 `json:"title" bson:"title"`
	Creator     string                 `json:"creator" bson:"creator"`
	Attribute   map[string]interface{} `json:"attribute" bson:"attribute"`
	Apply       map[string]interface{} `json:"apply" bson:"apply"`
	Check       map[string]interface{} `json:"check" bson:"check"`
	Result      map[string]interface{} `json:"result" bson:"result"`
}

func (w *WorkOrder) GenerateNumber() error {
	today := time.Now().Format("20060102")
	if mark, ok := Mark[w.Spec.OrderType]; ok{
		w.Spec.Number = fmt.Sprintf("%s%s-%s", mark, today, uuid.New().String())
	}

	fmt.Println(w.Spec.Number)
	return errors.New("Have no this type workorder")
}

func (w *WorkOrder) Clone() core.IObject {
	result := &WorkOrder{}
	core.Clone(w, result)
	return result
}
