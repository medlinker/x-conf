package models

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"x-conf/client/goclient"
)

// Mode 用户权限
type Mode int8

const (
	// Super 管理员
	Super Mode = 1
	// Normal 普通用户
	Normal Mode = 2

	prefix = "/users"
)

// User model
type User struct {
	Name string `json:"name"`
	Pass string `json:"-"`
	Mode `json:"mode"`
}

// IsSuper 判断用户是否是超级管理员
func (u User) IsSuper() bool {
	if u.Mode == Super {
		return true
	}
	return false
}

// EncrytPass 加密密码
func EncrytPass(originPass string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(originPass)))
}

// ValidPass 用户验证
func (u User) ValidPass() bool {
	resp, err := goclient.Get(makeKey(u.Name), nil)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if EncrytPass(u.Pass) == resp.Node.Value {
		return true
	}
	fmt.Println(2)
	return false
}

func makeKey(key string) string {
	return fmt.Sprint(prefix, "/", key)
}

func (u User) String() string {
	jsonStr, _ := json.Marshal(u)
	return string(jsonStr)
}
