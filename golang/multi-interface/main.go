package multiinterface

import "context"

type Interface1 interface {
	Method1()
	Method2(ctx context.Context) string
	Method3(ctx context.Context) string
}

type Interface2 interface {
	Method2(ctx context.Context) Interface1
}

type implementationInterface1 struct{}

func (i *implementationInterface1) Method1() {
	println("Method1")
}

func (i *implementationInterface1) Method2(ctx context.Context) string {
	println("Method2")
	return "Method2"
}

func (i *implementationInterface1) Method3(ctx context.Context) string {
	println("Method3")
	return "Method3"
}
