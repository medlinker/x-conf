package models

// Mode 配置模式：文件上传或配置平台手动配置
type Mode int8

const (
	// FileConfig 文件上传配置
	FileConfig Mode = 1
	// HandleConfig 配置平台手动配置
	HandleConfig Mode = 2
)

// Config 配置对象
type Config struct {
	projName string
	key      string
	value    interface{}
	mode     Mode
}
