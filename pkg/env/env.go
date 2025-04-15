package env

import (
	"flag"
	"fmt"
	"strings"
)

// Profile 环境常量
type Profile int

// Profile 环境配置常量
const (
	Dev Profile = iota
	Fat
	Uat
	Pro
)

// String 将 Profile 转换为字符串，用来提供 toString 方法
func (p Profile) String() string {
	return [...]string{"dev", "fat", "uat", "pro"}[p]
}

// Environment 环境配置接口
type Environment interface {
	Value() Profile // 获取环境值
	ignore()
}

// environment 环境配置
type environment struct {
	value Profile
}

func (e environment) Value() Profile {
	return e.value
}

func (e environment) ignore() {
	panic("implement me")
}

var _ Environment = (*environment)(nil)

// 当前环境
var current Environment

// Current 获取当前环境
func Current() Profile {
	return current.Value()
}

// Init 初始化环境
func Init() {
	env := flag.String("env", "", "请输入运行环境:\n dev:开发环境\n fat:测试环境\n uat:预上线环境\n pro:正式环境\n")
	flag.Parse()

	switch strings.ToLower(strings.TrimSpace(*env)) {
	case "dev":
		current = &environment{value: Dev}
	case "fat":
		current = &environment{value: Fat}
	case "uat":
		current = &environment{value: Uat}
	case "pro":
		current = &environment{value: Pro}
	default:
		current = &environment{value: Dev}
		fmt.Println("Warning: '-env' cannot be found, or it is illegal. The default 'dev' will be used.")
	}
}
