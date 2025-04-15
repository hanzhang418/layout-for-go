package env

import (
	"flag"
	"fmt"
	"strings"
)

// Environment 环境配置接口
type Environment interface {
	Value() string // 获取环境值
	ignore()
}

// environment 环境配置
type environment struct {
	value string
}

func (e environment) Value() string {
	return e.value
}

func (e environment) ignore() {
	panic("implement me")
}

var _ Environment = (*environment)(nil)

// 初始化一些环境
var (
	dev     Environment = &environment{value: "dev"} // 开发环境
	fat     Environment = &environment{value: "fat"} // 测试环境
	uat     Environment = &environment{value: "uat"} // 预上线环境
	pro     Environment = &environment{value: "pro"} // 生产环境
	current Environment                              // 当前环境
)

// Current 获取当前环境
func Current() Environment {
	return current
}

// Init 初始化环境
func Init() {
	env := flag.String("env", "", "请输入运行环境:\n dev:开发环境\n fat:测试环境\n uat:预上线环境\n pro:正式环境\n")
	flag.Parse()

	switch strings.ToLower(strings.TrimSpace(*env)) {
	case "dev":
		current = dev
	case "fat":
		current = fat
	case "uat":
		current = uat
	case "pro":
		current = pro
	default:
		current = dev
		fmt.Println("Warning: '-env' cannot be found, or it is illegal. The default 'dev' will be used.")
	}
}
