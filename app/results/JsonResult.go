package results

import (
	"encoding/json"
	"github.com/revel/revel"
	"net/http"
)

type JSONResult struct {
	obj      interface{}
	callback string
}

func (r JSONResult) Apply(req *revel.Request, resp *revel.Response) {
	var b []byte
	var err error
	if revel.Config.BoolDefault("results.pretty", false) {
		b, err = json.MarshalIndent(r.obj, "", "  ")
	} else {
		b, err = json.Marshal(r.obj)
	}

	if err != nil {
		revel.ErrorResult{Error: err}.Apply(req, resp)
		return
	}

	if r.callback == "" {
		resp.WriteHeader(http.StatusOK, "application/json; charset=utf-8")
		if _, err = resp.GetWriter().Write(b); err != nil {
			revel.AppLog.Error("Apply: Response write failed:", "error", err)
		}
		return
	}

	resp.WriteHeader(http.StatusOK, "application/javascript; charset=utf-8")
	if _, err = resp.GetWriter().Write([]byte(r.callback + "(")); err != nil {
		revel.AppLog.Error("Apply: Response write failed", "error", err)
	}
	if _, err = resp.GetWriter().Write(b); err != nil {
		revel.AppLog.Error("Apply: Response write failed", "error", err)
	}
	if _, err = resp.GetWriter().Write([]byte(");")); err != nil {
		revel.AppLog.Error("Apply: Response write failed", "error", err)
	}
}
