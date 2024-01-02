package utils

import (
	"fmt"
	"reflect"
	S "strings"
)

func StringifyAnyArgs(args ...any) string {
	var sb S.Builder
	// var pargs []any
	for _, pp := range args {
		// fmt.Printf("G_2 AnyArg %d \n", i)
		if p, ok := pp.(*string); ok {
			s := *p // s is a string
			if s == "" {
				s = "ϕ" // "nil/null"
			}
			// Ignore return values (int, error)
			sb.WriteString("\"" + s + "\" ")
		} else {
			v := reflect.ValueOf(pp)
			if v.Kind() == reflect.Ptr {
				v.Elem() // or v.Indirect() ?
			}
			sb.WriteString(fmt.Sprintf("<%T>|%+v| ", pp, pp))
		}
	}
	return sb.String()
}
