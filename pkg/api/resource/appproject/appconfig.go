package appproject

type AppConfigRequest struct {
	App                   string `json:"app"`
	Config                map[string]interface{} `json:"config"`
}
