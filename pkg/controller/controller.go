package main

type Controller interface {
	Run() error
	Stop() error
}
