// Licensed to Apache Software Foundation(ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation(ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package schema

import (
	"dtstack.com/dtstack/easymatrix/matrix/encrypt/aes"
	"fmt"
	"regexp"
	"strings"
)

func (c ConfigMap) AddOne(i int) int { return i + 1 }

func (c ConfigMap) AesEncryptByPassword(adminPass string) {
	for key, configItem := range c {
		if regexp.MustCompile(`(?i).*password.*`).Match([]byte(key)) {
			defaultValue, _ := aes.AesEncryptByPassword(fmt.Sprintf("%s", *(configItem.(VisualConfig).Default.(*string))), adminPass)
			value, _ := aes.AesEncryptByPassword(fmt.Sprintf("%s", *(configItem.(VisualConfig).Value.(*string))), adminPass)
			c[key] = VisualConfig{
				Default: defaultValue,
				Desc:    configItem.(VisualConfig).Desc,
				Type:    configItem.(VisualConfig).Type,
				Value:   value,
			}
		}
	}

}
func (c ConfigMap) Join(field interface{}, s ...interface{}) (string, error) {
	switch f := field.(type) {
	case string:
		if vc, ok := c[f]; ok {
			field = vc
			break
		}
		return "", fmt.Errorf("can't found config `%v`", f)
	case VisualConfig:
	default:
		return "", fmt.Errorf("the config `%v` type `%T` not support", field, f)
	}

	switch v := field.(VisualConfig).Value.(type) {
	case *ServiceAddrStruct:
		if v == nil {
			return "", fmt.Errorf("the config `%v` is nil", field)
		}
		all := ""
		for _, ss := range s {
			switch ss_ := ss.(type) {
			case string:
				all += ss_
			case VisualConfig:
				all += ss_.String()
			default:
				return "", fmt.Errorf("the s `%v` type `%T` not support", s, ss_)
			}
		}
		return strings.Join(v.IP, all), nil
	default:
		return "", fmt.Errorf("the config `%v` type `%T` not support", field, v)
	}
}

func (c ConfigMap) JoinHost(field interface{}, s ...interface{}) (string, error) {
	switch f := field.(type) {
	case string:
		if vc, ok := c[f]; ok {
			field = vc
			break
		}
		return "", fmt.Errorf("can't found config `%v`", f)
	case VisualConfig:
	default:
		return "", fmt.Errorf("the config `%v` type `%T` not support", field, f)
	}

	switch v := field.(VisualConfig).Value.(type) {
	case *ServiceAddrStruct:
		if v == nil {
			return "", fmt.Errorf("the config `%v` is nil", field)
		}
		all := ""
		for _, ss := range s {
			switch ss_ := ss.(type) {
			case string:
				all += ss_
			case VisualConfig:
				all += ss_.String()
			default:
				return "", fmt.Errorf("the s `%v` type `%T` not support", s, ss_)
			}
		}
		return strings.Join(v.Host, all), nil
	default:
		return "", fmt.Errorf("the config `%v` type `%T` not support", field, v)
	}
}

func (c ConfigMap) Joinx(field, sep interface{}, s ...interface{}) (string, error) {
	switch f := field.(type) {
	case string:
		if vc, ok := c[f]; ok {
			field = vc
			break
		}
		return "", fmt.Errorf("can't found config `%v`", f)
	case VisualConfig:
	default:
		return "", fmt.Errorf("the config `%v` type `%T` not support", field, f)
	}

	switch v := field.(VisualConfig).Value.(type) {
	case *ServiceAddrStruct:
		if v == nil {
			return "", fmt.Errorf("the config `%v` is nil", field)
		}
		all := ""
		for _, ss := range s {
			switch ss_ := ss.(type) {
			case string:
				all += ss_
			case VisualConfig:
				all += ss_.String()
			default:
				return "", fmt.Errorf("the s `%v` type `%T` not support", s, ss_)
			}
		}
		switch sep_ := sep.(type) {
		case string:
		case VisualConfig:
			sep = sep_.String()
		default:
			return "", fmt.Errorf("the sep `%v` type `%T` not support", sep, sep_)
		}
		return strings.Join(v.IP, all+sep.(string)) + all, nil
	default:
		return "", fmt.Errorf("the config `%v` type `%T` not support", field, v)
	}
}

func (c ConfigMap) JoinxHost(field, sep interface{}, s ...interface{}) (string, error) {
	switch f := field.(type) {
	case string:
		if vc, ok := c[f]; ok {
			field = vc
			break
		}
		return "", fmt.Errorf("can't found config `%v`", f)
	case VisualConfig:
	default:
		return "", fmt.Errorf("the config `%v` type `%T` not support", field, f)
	}

	switch v := field.(VisualConfig).Value.(type) {
	case *ServiceAddrStruct:
		if v == nil {
			return "", fmt.Errorf("the config `%v` is nil", field)
		}
		all := ""
		for _, ss := range s {
			switch ss_ := ss.(type) {
			case string:
				all += ss_
			case VisualConfig:
				all += ss_.String()
			default:
				return "", fmt.Errorf("the s `%v` type `%T` not support", s, ss_)
			}
		}
		switch sep_ := sep.(type) {
		case string:
		case VisualConfig:
			sep = sep_.String()
		default:
			return "", fmt.Errorf("the sep `%v` type `%T` not support", sep, sep_)
		}
		return strings.Join(v.Host, all+sep.(string)) + all, nil
	default:
		return "", fmt.Errorf("the config `%v` type `%T` not support", field, v)
	}
}

func (c ConfigMap) LastSegIP(field interface{}) (string, error) {
	switch f := field.(type) {
	case string:
		if vc, ok := c[f]; ok {
			field = vc
			break
		}
		return "", fmt.Errorf("can't found config `%v`", f)
	case VisualConfig:
	default:
		return "", fmt.Errorf("the config `%v` type `%T` not support", field, f)
	}

	switch v := field.(VisualConfig).Value.(type) {
	case *ServiceAddrStruct:
		if v == nil {
			return "", fmt.Errorf("the config `%v` is nil", field)
		}
		return strings.Split(v.getIP(), ".")[3], nil
	default:
		return "", fmt.Errorf("the config `%v` type `%T` not support", field, v)
	}
}

func (c ConfigMap) NodeCount(field interface{}) (int, error) {
	switch f := field.(type) {
	case string:
		if vc, ok := c[f]; ok {
			field = vc
			break
		}
		return 0, fmt.Errorf("can't found config `%v`", f)
	case VisualConfig:
	default:
		return 0, fmt.Errorf("the config `%v` type `%T` not support", field, f)
	}

	switch v := field.(VisualConfig).Value.(type) {
	case *ServiceAddrStruct:
		if v == nil {
			return 0, fmt.Errorf("the config `%v` is nil", field)
		}
		return len(v.IP), nil
	default:
		return 0, fmt.Errorf("the config `%v` type `%T` not support", field, v)
	}
}

// start from 0
func (c ConfigMap) NodeIndex(field interface{}) (uint, error) {
	switch f := field.(type) {
	case string:
		if vc, ok := c[f]; ok {
			field = vc
			break
		}
		return 0, fmt.Errorf("can't found config `%v`", f)
	case VisualConfig:
	default:
		return 0, fmt.Errorf("the config `%v` type `%T` not support", field, f)
	}

	switch v := field.(VisualConfig).Value.(type) {
	case *ServiceAddrStruct:
		if v == nil {
			return 0, fmt.Errorf("the config `%v` is nil", field)
		}
		return v.SingleIndex, nil
	default:
		return 0, fmt.Errorf("the config `%v` type `%T` not support", field, v)
	}
}

// start from 1
func (c ConfigMap) NodeID(field interface{}) (uint, error) {
	switch f := field.(type) {
	case string:
		if vc, ok := c[f]; ok {
			field = vc
			break
		}
		return 0, fmt.Errorf("can't found config `%v`", f)
	case VisualConfig:
	default:
		return 0, fmt.Errorf("the config `%v` type `%T` not support", field, f)
	}

	switch v := field.(VisualConfig).Value.(type) {
	case *ServiceAddrStruct:
		if v == nil {
			return 0, fmt.Errorf("the config `%v` is nil", field)
		}
		return v.NodeId, nil
	default:
		return 0, fmt.Errorf("the config `%v` type `%T` not support", field, v)
	}
}

func (c ConfigMap) IPList(field interface{}) ([]string, error) {
	switch f := field.(type) {
	case string:
		if vc, ok := c[f]; ok {
			field = vc
			break
		}
		return nil, fmt.Errorf("can't found config `%v`", f)
	case VisualConfig:
	default:
		return nil, fmt.Errorf("the config `%v` type `%T` not support", field, f)
	}

	switch v := field.(VisualConfig).Value.(type) {
	case *ServiceAddrStruct:
		if v == nil {
			return nil, fmt.Errorf("the config `%v` is nil", field)
		}
		return v.IP, nil
	default:
		return nil, fmt.Errorf("the config `%v` type `%T` not support", field, v)
	}
}

func (c ConfigMap) HostList(field interface{}) ([]string, error) {
	switch f := field.(type) {
	case string:
		if vc, ok := c[f]; ok {
			field = vc
			break
		}
		return nil, fmt.Errorf("can't found config `%v`", f)
	case VisualConfig:
	default:
		return nil, fmt.Errorf("the config `%v` type `%T` not support", field, f)
	}

	switch v := field.(VisualConfig).Value.(type) {
	case *ServiceAddrStruct:
		if v == nil {
			return nil, fmt.Errorf("the config `%v` is nil", field)
		}
		return v.Host, nil
	default:
		return nil, fmt.Errorf("the config `%v` type `%T` not support", field, v)
	}
}

func (c ConfigMap) Hostname(field interface{}) (string, error) {
	switch f := field.(type) {
	case string:
		if vc, ok := c[f]; ok {
			field = vc
			break
		}
		return "", fmt.Errorf("can't found config `%v`", f)
	case VisualConfig:
	default:
		return "", fmt.Errorf("the config `%v` type `%T` not support", field, f)
	}

	switch v := field.(VisualConfig).Value.(type) {
	case *ServiceAddrStruct:
		if v == nil {
			return "", fmt.Errorf("the config `%v` is nil", field)
		}
		return v.getHost(), nil
	default:
		return "", fmt.Errorf("the config `%v` type `%T` not support", field, v)
	}
}

func (c ConfigMap) GetIpByNodeID(field interface{}, nodeID uint) (string, error) {
	switch f := field.(type) {
	case string:
		if vc, ok := c[f]; ok {
			field = vc
			break
		}
		return "", fmt.Errorf("can't found config `%v`", f)
	case VisualConfig:
	default:
		return "", fmt.Errorf("the config `%v` type `%T` not support", field, f)
	}

	switch v := field.(VisualConfig).Value.(type) {
	case *ServiceAddrStruct:
		if v == nil {
			return "", fmt.Errorf("the config `%v` is nil", field)
		}
		return v.getIpByNodeID(nodeID)
	default:
		return "", fmt.Errorf("the config `%v` type `%T` not support", field, v)
	}
}

func (c ConfigMap) GetHostByNodeID(field interface{}, nodeID uint) (string, error) {
	switch f := field.(type) {
	case string:
		if vc, ok := c[f]; ok {
			field = vc
			break
		}
		return "", fmt.Errorf("can't found config `%v`", f)
	case VisualConfig:
	default:
		return "", fmt.Errorf("the config `%v` type `%T` not support", field, f)
	}

	switch v := field.(VisualConfig).Value.(type) {
	case *ServiceAddrStruct:
		if v == nil {
			return "", fmt.Errorf("the config `%v` is nil", field)
		}
		return v.getHostByNodeID(nodeID)
	default:
		return "", fmt.Errorf("the config `%v` type `%T` not support", field, v)
	}
}
