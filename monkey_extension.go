package monkeyExtensions

import (
	"fmt"
	"reflect"

	"github.com/bouk/monkey"
)

func PatchInstanceMethodFlexible(target reflect.Type, methodName string, replacement interface{}) {
	m, ok := target.MethodByName(methodName)
	if !ok {
		panic(fmt.Sprintf("unknown method %s", methodName))
	}

	replacementInputLen := reflect.TypeOf(replacement).NumIn()
	if replacementInputLen > m.Type.NumIn() {
		panic(fmt.Sprintf("replacement functoin has too many input parameters: %d, replaced function: %d", replacementInputLen, m.Type.NumIn()))
	}

	replacementWrapper := reflect.MakeFunc(m.Type, func(args []reflect.Value) []reflect.Value {
		inputsForReplacement := make([]reflect.Value, 0, replacementInputLen)
		for i := 0; i < cap(inputsForReplacement); i++ {
			elem := args[i].Convert(reflect.TypeOf(replacement).In(i))
			inputsForReplacement = append(inputsForReplacement, elem)
		}

		return reflect.ValueOf(replacement).Call(inputsForReplacement)
	}).Interface()

	monkey.PatchInstanceMethod(target, methodName, replacementWrapper)
}
