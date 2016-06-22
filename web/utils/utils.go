package utils

import (
	"encoding/json"
	"net/http"
)

// Header return json
func Header(w http.ResponseWriter) {
	w.Header().Set("Content-type", "application/json")
}

// Ret 返回json
type Ret struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// NewRet 构造返回
func NewRet() Ret {
	return Ret{0, "", nil}
}

const (
	// ErrMethod 请求方法错误
	ErrMethod = -1
	// ErrServer 服务错误
	ErrServer = -2
	// ErrParam 参数错误
	ErrParam = -3
)

// CheckErr 检查异常
func CheckErr(err error, ret *Ret) bool {
	if err != nil {
		ret.Code = ErrServer
		ret.Msg = err.Error()
		return true
	}
	return false
}

// Output 数据返回
func Output(w http.ResponseWriter, ret Ret) {
	d, _ := json.Marshal(ret)
	w.Write(d)
}
