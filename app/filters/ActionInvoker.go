package filters

import (
	"github.com/asaskevich/govalidator"
	"github.com/revel/revel"
	"github.com/tidwall/gjson"
	"github.com/zze326/devops-helper/app/results"
	"github.com/zze326/devops-helper/app/utils"
	"io"
	"reflect"
)

var (
	websocketType = reflect.TypeOf((*revel.ServerWebSocket)(nil)).Elem()
)

func ActionInvoker(c *revel.Controller, _ []revel.Filter) {
	if c.Request.Method != "GET" && !utils.InSlice[string](c.Request.Method, []string{"DELETE", "OPTIONS"}) && c.Request.ContentType == "application/json" {
		if jsonBody := string(c.Params.JSON); govalidator.IsJSON(jsonBody) {
			gjson.Parse(jsonBody).ForEach(func(key, value gjson.Result) bool {
				c.Params.Add(key.String(), value.String())
				return true
			})
		}
	}

	// Instantiate the method.
	methodValue := reflect.ValueOf(c.AppController).MethodByName(c.MethodType.Name)

	// Collect the values for the method's arguments.
	var methodArgs []reflect.Value
	for _, arg := range c.MethodType.Args {
		// If they accept a websocket connection, treat that arg specially.
		var boundArg reflect.Value
		if arg.Type.Implements(websocketType) {
			boundArg = reflect.ValueOf(c.Request.WebSocket)
		} else {
			boundArg = revel.Bind(c.Params, arg.Name, arg.Type)
			// #756 - If the argument is a closer, defer a Close call,
			// so we don't risk on leaks.
			if closer, ok := boundArg.Interface().(io.Closer); ok {
				defer func() {
					_ = closer.Close()
				}()
			}
			// 判断 boundArg 是否是结构体指针
			if boundArg.Kind() == reflect.Struct || (boundArg.Kind() == reflect.Ptr && boundArg.Elem().Kind() == reflect.Struct) {
				if _, err := govalidator.ValidateStruct(boundArg.Interface()); err != nil {
					c.Result = results.JsonError(err)
					return
				}
			}
		}
		methodArgs = append(methodArgs, boundArg)
	}

	var resultValue reflect.Value
	if methodValue.Type().IsVariadic() {
		resultValue = methodValue.CallSlice(methodArgs)[0]
	} else {
		resultValue = methodValue.Call(methodArgs)[0]
	}
	if resultValue.Kind() == reflect.Interface && !resultValue.IsNil() {
		c.Result = resultValue.Interface().(revel.Result)
	}
}
