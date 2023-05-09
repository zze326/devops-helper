package results

import (
	"fmt"
	"github.com/revel/revel"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type PageData struct {
	Total int `json:"total"`
	Items any `json:"items"`
}

func JsonOk() revel.Result {
	return JSONResult{obj: Response{Code: 1, Msg: "ok", Data: nil}}
}

func JsonOkMsg(msg string) revel.Result {
	return JSONResult{obj: Response{Code: 1, Msg: msg, Data: nil}}
}

func JsonError(err error) revel.Result {
	if err == nil {
		return JsonOk()
	}
	return JSONResult{obj: Response{Code: 0, Msg: err.Error(), Data: nil}}
}

func JsonErrorMsg(errMsg string) revel.Result {
	return JSONResult{obj: Response{Code: 0, Msg: errMsg, Data: nil}}
}

func JsonErrorMsgf(errMsg string, args ...any) revel.Result {
	return JSONResult{obj: Response{Code: 0, Msg: fmt.Sprintf(errMsg, args...), Data: nil}}
}

func JsonOkData(data interface{}) revel.Result {
	return JSONResult{obj: Response{Code: 1, Msg: "ok", Data: data}}
}
