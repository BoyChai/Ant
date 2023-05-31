package Ant

import (
	"errors"
	"fmt"
	"regexp"
)

var rules = rulesMap{
	// 表单相关
	"email": `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, // 邮箱
	// 网络相关
	"network_transport_protocol_common": `(?i)^(UDP|TCP)$`,                                           // 常用的网络传输协议 (TCP|UDP)
	"port":                              `^([1-9]|[1-9]\d{1,3}|[1-5]\d{4}|[6][0-5][0-5][0-3][0-5])$`, // 端口 (1-65535)
	"addr_ipv4":                         `^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`,                           // ipv4地址
	"addr_ipv6":                         `^(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`,                // 标准的ipv6地址

}

type rulesMap map[string]string

var custom customRule

type customRule map[string]func(value string) error

func init() {
	custom = make(map[string]func(value string) error)
}

// 获取规则对象
func (r *rulesMap) getRule(name string) (*regexp.Regexp, error) {
	if rules[name] != "" {
		regexpObj, err := regexp.Compile(rules[name])
		if err != nil {
			return nil, errors.New(fmt.Sprint("Rule compilation error, please contact the developer Err:", err))
		}
		return regexpObj, nil
	}
	return nil, errors.New("unknown rule")
}

func NewCustomRule(ruleName string, fun func(value string) error) {
	custom[ruleName] = fun
}

func (c *customRule) getRule(name string) func(value string) error {
	return custom[name]
}
