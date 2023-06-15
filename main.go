package main

import (
	"micp-sim/microprocessor"
	"micp-sim/util"
	"strconv"
	"strings"
	"syscall/js"
)


func add(this js.Value, i []js.Value) any {
	value1 := js.Global().Get("document").Call("getElementById", i[0].String()).Get("value").String()
	value2 := js.Global().Get("document").Call("getElementById", i[1].String()).Get("value").String()

	int1, _ := strconv.Atoi(value1)
	int2, _ := strconv.Atoi(value2)
	js.Global().Set("output", int1+int2)
	sum := int1 + int2
	println(sum)

	js.Global().Get("document").Call("getElementById", i[2].String()).Set("value", sum)
	return ""
}

func sendInst(this js.Value, i []js.Value) any{
	value1 := js.Global().Get("document").Call("getElementById", "inst").Get("value").String()
    inst := strings.Split(value1, "\n")

    micp := microprocessor.New(1)
    micp.Start(inst)
    al := micp.Al.GetValue()

    
    println(util.BinaryToDecimal(al[:])) 

    return ""
}

func subtract(this js.Value, i []js.Value) any {
	js.Global().Set("output", js.ValueOf(i[0].Int()+i[1].Int()))
	sum := i[0].Int() - i[1].Int()
	println(sum)
	return ""
}

func Hello(this js.Value, i []js.Value) any {
	println("yes i got called")
	return ""
}

func registerCallbacks() {
	js.Global().Set("add", js.FuncOf(add))
	js.Global().Set("subtract", js.FuncOf(subtract))
	js.Global().Set("Hello", js.FuncOf(Hello))
	js.Global().Set("sendInst", js.FuncOf(sendInst))
}

func main() {
	c := make(chan struct{}, 0)
	println("Go webassembly initializes")
	registerCallbacks()


	<-c
}
