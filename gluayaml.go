package gluayaml

import (
	_ "errors"
	"fmt"
	"github.com/yuin/gopher-lua"
	"gopkg.in/yaml.v2"
	_ "strconv"
)

func Loader(L *lua.LState) int {
	tb := L.NewTable()
	L.SetFuncs(tb, map[string]lua.LGFunction{
		"parse": apiParse,
		"dump":  apiDump,
	})
	L.Push(tb)

	return 1
}

// apiParse parses yaml formatted text to the table.
func apiParse(L *lua.LState) int {
	str := L.CheckString(1)

	var value interface{}
	err := yaml.Unmarshal([]byte(str), &value)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(fromYAML(L, value))
	return 1
}

// apiDump dumps the yaml formatted text.
func apiDump(L *lua.LState) int {
	L.RaiseError("unimplemented function")
	return 0

	//	value := L.CheckAny(1)
	//	visited := make(map[*lua.LTable]bool)
	//	data, err := toYAML(value, visited)
	//	if err != nil {
	//		L.Push(lua.LNil)
	//		L.Push(lua.LString(err.Error()))
	//		return 2
	//	}
	//	L.Push(lua.LString(string(data)))
	//	return 1
}

//var (
//	errFunction = errors.New("cannot encode function to YAML")
//	errChannel  = errors.New("cannot encode channel to YAML")
//	errState    = errors.New("cannot encode state to YAML")
//	errUserData = errors.New("cannot encode userdata to YAML")
//	errNested   = errors.New("cannot encode recursively nested tables to YAML")
//)
//
//type yamlValue struct {
//	lua.LValue
//	visited map[*lua.LTable]bool
//}
//
//func (y yamlValue) MarshalYAML() (interface{}, error) {
//	fmt.Println(y.LValue)
//	return toYAML(y.LValue, y.visited)
//}
//
//func toYAML(value lua.LValue, visited map[*lua.LTable]bool) (data []byte, err error) {
//	switch converted := value.(type) {
//	case lua.LBool:
//		data, err = yaml.Marshal(converted)
//	case lua.LChannel:
//		err = errChannel
//	case lua.LNumber:
//		data, err = yaml.Marshal(converted)
//	case *lua.LFunction:
//		err = errFunction
//	case *lua.LNilType:
//		data, err = yaml.Marshal(converted)
//	case *lua.LState:
//		err = errState
//	case lua.LString:
//		data, err = yaml.Marshal(converted)
//	case *lua.LTable:
//		var arr []yamlValue
//		var obj map[string]yamlValue
//
//		if visited[converted] {
//			panic(errNested)
//			return
//		}
//		visited[converted] = true
//
//		converted.ForEach(func(k lua.LValue, v lua.LValue) {
//			i, numberKey := k.(lua.LNumber)
//			if numberKey && obj == nil {
//				index := int(i) - 1
//				if index != len(arr) {
//					// map out of order; convert to map
//					obj = make(map[string]yamlValue)
//					for i, value := range arr {
//						obj[strconv.Itoa(i+1)] = value
//					}
//					obj[strconv.Itoa(index+1)] = yamlValue{v, visited}
//					return
//				}
//				arr = append(arr, yamlValue{v, visited})
//				return
//			}
//			if obj == nil {
//				obj = make(map[string]yamlValue)
//				for i, value := range arr {
//					obj[strconv.Itoa(i+1)] = value
//				}
//			}
//			obj[k.String()] = yamlValue{v, visited}
//		})
//		if obj != nil {
//			data, err = yaml.Marshal(obj)
//		} else {
//			data, err = yaml.Marshal(arr)
//		}
//	case *lua.LUserData:
//		err = errUserData
//	}
//	return
//}

func fromYAML(L *lua.LState, value interface{}) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case int:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case []interface{}:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(fromYAML(L, item))
		}
		return arr
	case map[interface{}]interface{}:
		tbl := L.CreateTable(0, len(converted))
		for key, item := range converted {
			if s, ok := key.(string); ok {
				tbl.RawSetH(lua.LString(s), fromYAML(L, item))
			}
		}
		return tbl
	default:
		panic(fmt.Sprintf("unexpected type: %T\n", converted))
	}

	return lua.LNil
}
