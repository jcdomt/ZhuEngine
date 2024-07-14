package router

import (
	"fmt"
	"reflect"

	"github.com/kataras/iris/v12"
)

func InitRouter(app *iris.Application) {
	fmt.Printf("%+v", getRouterFuncs("app"))
}

func getRouterFuncs(packageName string) []string {
	var funcs []string
	pkg := reflect.ValueOf(packageName)
	t := pkg.Type()
	for i := 0; i < t.NumMethod(); i++ {
		funcs = append(funcs, t.Method(i).Name)
	}
	return funcs
}
