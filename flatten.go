package flatten

import (
	"encoding/json"
	"errors"
	"regexp"
	"fmt"
	"strings"
	"strconv"
	"sort"
	"github.com/imdario/mergo"
	"gopkg.in/mgo.v2/bson"
)

type SeparatorStyle struct {
	Before string
	Middle string
	After  string
}


var (
	DotStyle = SeparatorStyle{Middle: "."}

	PathStyle = SeparatorStyle{Middle: "/"}

	RailsStyle = SeparatorStyle{Before: "[", After: "]"}

	UnderscoreStyle = SeparatorStyle{Middle: "_"}
)

var NotValidInputError = errors.New("Not a valid input: map or slice")
var NotValidJsonInputError = errors.New("Not a valid input, must be a map")
var isJsonMap = regexp.MustCompile(`^\s*\{`)

func Flatten(nested map[string]interface{}, prefix string, style SeparatorStyle) (map[string]interface{}, error) {
	flatmap := make(map[string]interface{})

	err := flatten(true, flatmap, nested, prefix, style)
	if err != nil {
		return nil, err
	}

	return flatmap, nil
}

func UnFlatten(nested map[string]interface{}, prefix string, style SeparatorStyle) (map[string]interface{}, error) {
	flatmap := make(map[string]interface{})

	err := unflatten(true, flatmap, nested, prefix, style)
	if err != nil {
		return nil, err
	}

	return flatmap, nil
}


func FlattenString(nestedstr, prefix string, style SeparatorStyle) (string, error) {
	if !isJsonMap.MatchString(nestedstr) {
		return "", NotValidJsonInputError
	}

	var nested map[string]interface{}
	err := json.Unmarshal([]byte(nestedstr), &nested)
	if err != nil {
		return "", err
	}

	flatmap, err := Flatten(nested, prefix, style)
	if err != nil {
		return "", err
	}

	flatb, err := json.Marshal(&flatmap)
	if err != nil {
		return "", err
	}

	return string(flatb), nil
}

func UnFlattenString(nestedstr, prefix string, style SeparatorStyle) (string, error) {
	if !isJsonMap.MatchString(nestedstr) {
		return "", NotValidJsonInputError
	}

	var nested map[string]interface{}
	err := json.Unmarshal([]byte(nestedstr), &nested)
	if err != nil {
		return "", err
	}

	flatmap, err := UnFlatten(nested, prefix, style)
	if err != nil {
		return "", err
	}

	flatb, err := json.Marshal(&flatmap)
	if err != nil {
		return "", err
	}

	return string(flatb), nil
}

func flatten(top bool, flatMap map[string]interface{}, nested interface{}, prefix string, style SeparatorStyle) error {
	assign := func(newKey string, v interface{}) error {
		switch v.(type) {
		case map[string]interface{}:
			if err := flatten(false, flatMap, v, newKey, style); err != nil {
				return err
			}
		case []interface{}:
			if err := flatten(false, flatMap, v, newKey, style); err != nil {
				return err
			}
		default:
			flatMap[newKey] = v
		}

		return nil
	}

	switch nested.(type) {
	case map[string]interface{}:
		for k, v := range nested.(map[string]interface{}) {
			newKey := enkey(top, prefix, k, style)
			assign(newKey, v)
		}

	case []interface{}:
		for i, v := range nested.([]interface{}) {
			newKey := enkey(top, prefix, strconv.Itoa(i), style)
			assign(newKey, v)
		}
	default:
		return NotValidInputError
	}

	return nil
}

func getKeys(m map[string]interface{}) []string {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率较高
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}


func unflatten(top bool, flatMap map[string]interface{}, nested interface{}, prefix string, style SeparatorStyle) error {
	uf := func (k string, v interface{}) (n interface{}) {
		n = v

		keys := strings.Split(k, ".")

		for i := len(keys) - 1; i >= 0; i-- {
			temp := make(map[string]interface{})
			temp[keys[i]] = n
			n = temp
		}

		return
	}


	switch nested.(type) {
	case map[string]interface{}:
		abc := nested.(map[string]interface{})
		keys := getKeys(abc)
		sort.Strings(keys)
		for _, v := range  keys {
                        temp := uf(v, abc[v]).(map[string]interface{})
			err := mergo.Merge(&flatMap, temp)
			if err != nil {
				fmt.Println(err)
			}
		}
	case []interface{}:

		fmt.Println(nested.([]interface{}))
		//for i, v := range nested.([]interface{}) {
		//	newKey := enkey(top, prefix, strconv.Itoa(i), style)
		//	assign(newKey, v)
		//}
	default:
		return NotValidInputError
	}

	//return nil
	return nil
}


func enkey(top bool, prefix, subkey string, style SeparatorStyle) string {
	key := prefix

	if top {
		key += subkey
	} else {
		key += style.Before + style.Middle + subkey + style.After
	}

	return key
}


func unenkey(top bool, prefix, subkey string, style SeparatorStyle) string {
	key := prefix

	if top {
		key += subkey
	} else {
		key += style.Before + style.Middle + subkey + style.After
	}

	return key
}

func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func getKeys2(m map[string][]interface{}) []string {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率较高
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func test5 (nested interface{}, key string, nested2 map[string][]interface{}, nested3 []string) {
	switch nested.(type) {
	case map[string]interface{}:
		for k, v := range nested.(map[string]interface{}) {
			newKey := enkey2(key, ".", k)
			cc := strings.Trim(newKey, ".")
			fmt.Println(cc)
			for _,v2 := range nested3 {
				if strings.Split(v2, ".")[0] == cc {

					test5(v, newKey, nested2, nested3)
					break

				} else if v2  == cc {
					test5(v, newKey, nested2, nested3)
					break

				}

			}

			//if _, ok := nested3[]; ok {
			//	test5(v, newKey, nested2, nested3)
			//}
			//fmt.Println(k)
			//fmt.Println(v)
		}
	case []interface{}:
		global_len := strconv.Itoa(len(getKeys2(nested2))-1)


		for k, v := range nested.([]interface{}) {
			entry := DeepCopy(nested2[global_len])

			if k > 0 {
				nested2[strconv.Itoa(len(getKeys2(nested2)))] = entry.([]interface{})
			}
			test5(v, key, nested2, nested3)

		}
		var qq []interface{}
		for _,v := range nested2 {
			for _,v2 := range v {
				qq = append(qq, v2)
			}

		}
		nested2["0"] = qq

		for _,v := range getKeys2(nested2) {
                    if v != "0" {
			delete(nested2,v)
		}

		}

	default:
		global_len := strconv.Itoa(len(getKeys2(nested2))-1)
		for _,v := range nested2[global_len] {
			cc := strings.Trim(key, ".")
			v.(map[string]interface{})[cc] = nested
		}
	}
}

func DeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = DeepCopy(v)
		}

		return newSlice
	} else if valueMap, ok := value.(bson.M); ok {
		newMap := make(bson.M)
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}
	}
	return value
}


func enkey2(prefix string, subkey string, style string) string {
	//fmt.Println(prefix)
	//fmt.Println(subkey)
	//fmt.Println(style)
	//key := prefix

	//if top {
	//	key += subkey
	//} else {
	//	key += style.Before + style.Middle + subkey + style.After
	//}
	//key := ""
	//if style == "" {
	//	key  = prefix
	//} else {
	//}
	key := strings.Trim(prefix, ".") + subkey + style
	//key += style.Before + style.Middle + subkey + style.After

	return key
}


func FlattenPreserveListsString(json2 string) string {
	flat, err := FlattenString(json2, "", DotStyle)
	if err!= nil {
		fmt.Println(err)
	}
	//fmt.Println(string(flat))

	var nested map[string]interface{}

	err = json.Unmarshal([]byte(flat), &nested)
	if err != nil {
		fmt.Println(err)
	}


	nested2 := make(map[string]interface{})
	keys := getKeys(nested)
	for _,v := range keys {
		vv := strings.Split(v, ".")
		c := ""
		for _,v2 := range vv {
			//fmt.Println(IsNum(v2))
			if IsNum(v2) == false {
				if c == "" {
					c = v2
				} else {
					c = c + "." +v2
				}
			}
		}
		nested2[string(c)] = ""
		//fmt.Println(c)
	}
	//fmt.Println(keys)
	nested3 := make(map[string][]interface{})
	//nested3["0"] = []interface{}
	nested3["0"] = append(nested3["0"], nested2)

	var nested10 map[string]interface{}

	err = json.Unmarshal([]byte(json2), &nested10)
	if err != nil {
		fmt.Println(err)
	}

	//nested4 := make(map[string]interface{})
	//nested4["age"] = ""

	nested4 :=[]string{"age","name.first"}
	test5(nested10, "", nested3, nested4)

	q := nested3["0"]

	flatb, err := json.Marshal(q)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(string(flatb))
	return string(flatb)

}
func FlattenPreserveLists(json2 string, test ...string) *[]interface{} {

	flat, err := FlattenString(json2, "", DotStyle)
	if err!= nil {
		fmt.Println(err)
	}
	//fmt.Println(string(flat))

	var nested map[string]interface{}

	err = json.Unmarshal([]byte(flat), &nested)
	if err != nil {
		fmt.Println(err)
	}


	nested2 := make(map[string]interface{})
	//keys := getKeys(nested)
	//for _,v := range keys {
	//	vv := strings.Split(v, ".")
	//	c := ""
	//	for _,v2 := range vv {
	//		//fmt.Println(IsNum(v2))
	//		if IsNum(v2) == false {
	//			if c == "" {
	//				c = v2
	//			} else {
	//				c = c + "." +v2
	//			}
	//		}
	//	}
	//	nested2[string(c)] = ""
	//	//fmt.Println(c)
	//}
	//fmt.Println(keys)
	nested3 := make(map[string][]interface{})
	//nested3["0"] = []interface{}
	//nested2["age"] = ""
	nested3["0"] = append(nested3["0"], nested2)
	//map[0:[map[age: name.first: name.last: qqqq.first: qqqq.last: qqqq.test.test:]]]

	var nested10 map[string]interface{}

	err = json.Unmarshal([]byte(json2), &nested10)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(nested3)

	//nested4 :=[]string{"age","name.first"}
	//nested4["age"] = ""
	//nested4["name.first"] = ""


	test5(nested10, "", nested3, test)
	q := nested3["0"]
	return &q
}


//const json2 = `{"name":[{"first":"Janet","last":"Prichard"},{"first":"Janet","last":"Prichard22"}],"age":47,"qqqq":[{"first":"Janet","last":"Prichard"},{"first":"Janet","last":"Prichard22"}]}`
//const json2 = `{"name":[{"first":"Janet","last":"Prichard"},{"first":"Janet","last":"Prichard22"}],"age":47,"qqqq":[{"first":"Janet","last":"Prichard","test":[{"test":"test"}]},{"first":"Janet","last":"Prichard22","test":[{"test":"test"}]}]}`
//const json2 = `{"name":[{"first":"Janet","last":"Prichard"},{"first":"Janet","last":"Prichard22"}],"age":47,"qqqq":[{"first":"Janet","last":"Prichard","test":[{"test":"test2","test5":"test5"},{"test":"test2","test5":"test5"}]},{"first":"Janet","last":"Prichard22","test":[{"test":"test2"}]}]}`

const json2 = `{"F0001":"上海仓库","F0003":"344","F0004":[{"C0002":"橘子","C0003":"0002","C0004":"5","C0005":"10","C0006":"50"},{"C0002":"橘子","C0003":"0002","C0004":"19","C0005":"10","C0006":"190"},{"C0002":"西瓜","C0003":"0001","C0004":"3","C0005":"5","C0006":"15"}],"F0005":"3","F0006":"255"}`

func main() {
	//nested4 :=[]string{"age","name.first"}
	nested4 :=[]string{"F0001","F0004.C0002", "F0004.C0006", "F0004.C0003"}
	flatb := FlattenPreserveLists(json2, nested4...)

	fmt.Println(flatb)
	//flatb2 := FlattenPreserveListsString(json2)

	//fmt.Println(flatb2)
}
