package base

import apiResource "github.com/yametech/devops/pkg/api/resource"

func (u *User) ListAr() (interface{}, error) {
	data := make([]apiResource.RespArtifact, 0)
	_, err := u.IService.List("artifacts", 0, 1, true, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
