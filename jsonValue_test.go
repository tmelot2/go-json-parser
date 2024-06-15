package main

import (
	// "fmt"
	"testing"
)

func TestSomething(t *testing.T) {
	// TODO
}


// JsonValue tests
//
// JSON i was testing stuff with
// {
// 	"a": "b",
// 	"b": 1,
// 	"c": 2.222,
// 	"d": {
// 		"d1": "d",
// 		"d2": 1,
// 		"d3": 2.222
// 	},
// 	"e": [
// 		[1,2]
// 	]
// }
// j := NewJsonValue(jsonResult)
// s, sErr := j.GetString("a")
// fmt.Println(s, sErr)

// i, iErr := j.GetInt("b")
// fmt.Println(i, iErr)

// f, fErr := j.GetFloat("c")
// fmt.Println(f, fErr)

// o, oErr := j.GetObject("d")
// fmt.Println(o, oErr)
// fmt.Println(o.GetString("d1"))
// fmt.Println(o.GetInt("d2"))
// fmt.Println(o.GetFloat("d3"))

// a, _ := j.GetArray("e")
// for _, aVal := range a {
// 	fmt.Println(aVal)
// 	fmt.Println(aVal.GetObject(""))
// 	z, _ := aVal.GetArray("")
// 	for _, zVal := range z {
// 		fmt.Println(zVal.GetInt(""))
// 	}
