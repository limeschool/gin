package gin

import (
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

func requestParams(c *Context) map[string]string {
	postParams := func(c *Context) string {
		data, _ := ioutil.ReadAll(c.Request.Body)
		return *(*string)(unsafe.Pointer(&data))
	}
	getParams := func(c *Context) string {
		return c.Request.URL.Query().Encode()
	}

	return map[string]string{
		"params": getParams(c),
		"body":   postParams(c),
	}
}

func panicErr() []string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller
	var arr []string
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		arr = append(arr, fmt.Sprintf("%s:%d", file, line))
	}
	return arr
}

func GormWhere(db *gorm.DB, tb string, src interface{}) *gorm.DB {
	m := ToSqlMap(tb, src)
	if m == nil {
		return db
	}
	for key, val := range m {
		if val == nil {
			db = db.Where(key)
		} else {
			db = db.Where(key, val)
		}
	}
	return db
}

func ToSqlMap(tb string, src interface{}) map[string]interface{} {
	if src == nil {
		return nil
	}
	m := make(map[string]interface{})
	value := reflect.ValueOf(src)
	tp := reflect.TypeOf(src)
	for value.Kind().String() == "ptr" {
		value = value.Elem()
		tp = tp.Elem()
	}
	if value.Kind().String() == "struct" {
		num := value.NumField()
		for i := 0; i < num; i++ {
			if isBlank(value.Field(i)) {
				continue
			}

			if value.Kind() == reflect.Slice {
				m[tb+"."+tp.Field(i).Name+" in ?"] = value.Field(i).Interface()
				continue
			}

			sqlTag := tp.Field(i).Tag.Get("sql")
			if sqlTag == "-" {
				continue
			}

			field := tp.Field(i).Tag.Get("json")
			if tp.Field(i).Tag.Get("field") != "" {
				field = tp.Field(i).Tag.Get("field")
			}

			if sqlTag == "" {
				m[tb+"."+field+" = ?"] = value.Field(i).Interface()
			} else {
				sqlTag = strings.ReplaceAll(sqlTag, "?", fmt.Sprint(value.Field(i).Interface()))
				m[tb+"."+field+" "+sqlTag] = nil
			}
		}
	}
	return m
}

func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}
