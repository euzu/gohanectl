package utils

import (
	"encoding/json"
	"fmt"
	"gohanectl/hanectl/model"
	"regexp"
	"strings"
)

func IsBlank(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}

func IsNotBlank(value string) bool {
	return !IsBlank(value)
}

// Sprintf support named format
var re = regexp.MustCompile("%\\{([a-zA-Z0-9_]+)\\}[.0-9]*[xsvTtbcdoqXxUeEfFgGp]")

func reformat(f string) (string, []string) {
	i := re.FindAllStringSubmatchIndex(f, -1)

	ord := []string{}
	pair := []int{0}
	for _, v := range i {
		ord = append(ord, f[v[2]:v[3]])
		pair = append(pair, v[2]-1)
		pair = append(pair, v[3]+1)
	}
	pair = append(pair, len(f))
	plen := len(pair)

	out := ""
	for n := 0; n < plen; n += 2 {
		out += f[pair[n]:pair[n+1]]
	}

	return out, ord
}

func parse(format string, params model.Dictionary) (string, []interface{}) {
	f, n := reformat(format)
	p := make([]interface{}, len(n))
	for i, v := range n {
		p[i] = params[v]
	}
	return f, p
}

// Named format
func Sprintf(format string, params model.Dictionary) string {
	f, p := parse(format, params)
	return fmt.Sprintf(f, p...)
}

func ContainsStr(list []string, item string) bool {
	for _, x := range list {
		if strings.Compare(x, item) == 0 {
			return true
		}
	}
	return false
}

func Json(value interface{}) string {
	jString, _ := json.Marshal(value)
	return string(jString)
}

func mapValue(value interface{}) interface{} {
	if assertedTMap, okMap := value.(map[interface{}]interface{}); okMap {
		return MapToStringKey(assertedTMap)
	} else if assertedToSlice, okSlice := value.([]interface{}); okSlice {
		var rslice []interface{}
		for _, x := range assertedToSlice {
			if assertedElement, okElem := x.(map[interface{}]interface{}); okElem {
				rslice = append(rslice, MapToStringKey(assertedElement))
			} else {
				rslice = append(rslice, mapValue(x))
			}
		}
		return rslice
	} else {
		return value
	}
}

func MapToStringKey(source map[interface{}]interface{}) model.Dictionary {
	if source == nil {
		return nil
	}
	result := make(model.Dictionary)
	for k, v := range source {
		key := fmt.Sprintf("%v", k)
		result[key] = mapValue(v)
	}
	return result
}

func MapContentToStringKey(source model.Dictionary) model.Dictionary {
	for k, v := range source {
		if dict, ok := v.(map[interface{}]interface{}); ok {
			source[k] = MapToStringKey(dict)
		}
	}
	return source
}
