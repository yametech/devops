package resource

type RequestUser struct {
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Password string `json:"password"`
	Username string `json:"username"`
	NickName string `json:"nick_name"`

}
