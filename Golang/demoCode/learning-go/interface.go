package main

type LogicProvider struct {
}

type Logic interface {
	Process(data string) string
}

type Client struct {
	L Logic
}
