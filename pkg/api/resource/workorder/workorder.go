package workorder

import "github.com/yametech/devops/pkg/resource/workorder"

type Request struct {
	workorder.OrderType   `json:"order_type"`
	workorder.OrderStatus `json:"order_status"`
	Title                 string                 `json:"title"`
	Relation              string                 `json:"relation"`
	Attribute             map[string]interface{} `json:"attribute"`
	Apply                 map[string]interface{} `json:"apply"`
	Check                 map[string]interface{} `json:"check"`
	Result                map[string]interface{} `json:"result"`
	Extends               map[string]interface{} `json:"extends"`
}
