package alu

import (
	"micp-sim/register"
	"micp-sim/stack"
	"micp-sim/util"
	"testing"
)

func TestAddition(t *testing.T){
	mar := &[2]register.Register{register.New(), register.New()}
	mbr := &[2]register.Register{register.New(), register.New()}


	al := register.New()
	stack := stack.New()
	a := New(&al, mar, mbr, &stack)

    valA := 23
    valB := 255
    a.Temp1.SetLoad() 
    a.Temp1.LoadValue(util.DecimalToBinary(valA))
    a.Temp2.SetLoad() 
    a.Temp2.LoadValue(util.DecimalToBinary(valB))

    a.Addition("")
    res := a.Al.GetValue()
    resDec := util.BinaryToDecimal(res[:])
    if byte(resDec) != byte(valA + valB){
        t.Fatalf("Val A = %v, Val B = %v, expected Sum = %v, acutal sum = %v", valA, valB, byte(valA + valB), resDec)
    }
}
