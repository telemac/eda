package event

import (
	"github.com/fatih/camelcase"
	"reflect"
	"strings"
)

func GetTypeName(event Eventer) string {
	evType := reflect.TypeOf(event).String()
	lastIndex := strings.LastIndex(evType, ".")
	if lastIndex == -1 {
		return evType
	}
	return evType[lastIndex+1:]
}

func GetTypeNameCamelCase(event Eventer) string {
	evType := GetTypeName(event)
	splitted := camelcase.Split(evType)
	for i, s := range splitted {
		splitted[i] = strings.ToLower(s)
	}
	return strings.Join(splitted, ".")
}
