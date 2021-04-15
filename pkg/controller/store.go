package main

import "github.com/yametech/devops/pkg/store"

type StorageInterface interface {
	UpdateStatus(run FlowRun) error
}

type StorageImpl struct {
	store.IKVStore
}

func (s StorageImpl) UpdateStatus(FlowRun) error {
	panic("implement me")
}

var _ StorageInterface = &StorageImpl{}
