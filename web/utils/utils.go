package utils

import (
	"encoding/json"
	"net/http"
)

// Filter 拦截器
func Filter(w http.ResponseWriter, r *http.Request) {
	/*
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Max-Age", "1728000")
		w.Header().Set("Content-type", "application/json")
	*/
}

// Header 设置返回header
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
	// ErrUser 用户名或密码错误
	ErrUser = -4
)

var (
	// Envs env
	Envs = [4]string{"prod", "qa", "dev", "local"}
)

// CheckErr 检查异常
func CheckErr(err error, ret *Ret) bool {
	if err != nil {
		if ret != nil {
			ret.Code = ErrServer
			ret.Msg = err.Error()
		}
		return true
	}
	return false
}

// SetMethodErr 设置
func SetMethodErr(ret *Ret) {
	ret.Code = ErrMethod
	ret.Msg = "request method error"
}

// CheckParamsErr 参数空
func CheckParamsErr(ret *Ret, params ...string) bool {
	for _, p := range params {
		if p == "" {
			if ret != nil {
				ret.Code = ErrParam
				ret.Msg = "params error"
			}
			return true
		}
	}
	return false
}

// Output 数据返回
func Output(w http.ResponseWriter, ret Ret) {
	d, _ := json.Marshal(ret)
	w.Write(d)
}

// CheckSessFromCookie cookie中检测session
func CheckSessFromCookie(r *http.Request) (bool, *Session) {
	c, err := r.Cookie("SESSIONID")
	if err != nil || !CheckSess(c.Value) {
		return false, nil
	}
	return true, SessMap[c.Value]
}
