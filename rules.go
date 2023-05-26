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
	"network_transport_protocol_common": `(?i)^(UDP|TCP)$`,                            // 常用的网络传输协议 (TCP|UDP)
	"addr_ipv4":                         `^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`,            // ipv4地址
	"addr_ipv6":                         `^(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`, // 标准的ipv6地址
}

type rulesMap map[string]string

// 获取规则对象
func (r *rulesMap) getRule(name string) (*regexp.Regexp, error) {
	if rules[name] != "" {
		regexpObj, err := regexp.Compile(rules[name])
		if err != nil {
			return nil, errors.New(fmt.Sprint("规则编译错误,请联系开发者. err:", err))
		}
		return regexpObj, nil
	}
	return nil, errors.New("未知的规则")
}
