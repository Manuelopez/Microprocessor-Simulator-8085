package main

import (
	"micp-sim/microprocessor"
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
    go micp.Start(inst)



    return ""
}

func subtract(this js.Value, i []js.Value) any {
	js.Global().Set("output", js.ValueOf(i[0].Int()+i[1].Int()))
	sum := i[0].Int() - i[1].Int()
	println(sum)
	return ""
}


func syntaxChequer(this js.Value, i []js.Value) any {
    instString := i[0].String()
    savedPcVariables := make(map[string]int)
    insts := strings.Split(instString, "\n")
    errors := [][]string{}
    for _, inst := range insts{
        if strings.TrimSpace(inst) == ""{
            continue
        }
        error := microprocessor.SyntaxHilighter(inst, savedPcVariables)
        errors = append(errors, error)
    }


	errosJs := js.Global().Get("Array").New(len(errors))
    for i, x := range errors{
	    errorJs := js.Global().Get("Array").New(len(x))
        for j, err := range x{
            errorJs.SetIndex(j, err)
        }
		errosJs.SetIndex(i, errorJs)
    }
    return errosJs
}

func registerCallbacks() {
	js.Global().Set("add", js.FuncOf(add))
	js.Global().Set("subtract", js.FuncOf(subtract))
	js.Global().Set("sendInst", js.FuncOf(sendInst))
	js.Global().Set("syntaxChequer", js.FuncOf(syntaxChequer))
}

func main() {
	c := make(chan struct{}, 0)
	println("Go webassembly initializes")
	registerCallbacks()


	<-c
}
