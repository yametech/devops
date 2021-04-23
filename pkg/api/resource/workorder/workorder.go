package workorder

import "github.com/yametech/devops/pkg/resource/workorder"

type Request struct {
	workorder.OrderType `json:"order_type"`
	Title               string                 `json:"title"`
	Attribute           map[string]interface{} `json:"attribute"`
	Apply               map[string]interface{} `json:"apply"`
	Check               map[string]interface{} `json:"check"`
	Result              map[string]interface{} `json:"result"`
}

