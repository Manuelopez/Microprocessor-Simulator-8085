package alu

import (
	"micp-sim/register"
	"micp-sim/stack"
	"micp-sim/util"
)

type Alu struct{
  Temp1 register.Register
  Temp2 register.Register
  Al *register.Register
  Stack *stack.Stack
  Carry bool
  Zero bool
  Comparison bool
	Mar *[2]register.Register
	Mbr *[2]register.Register
}

func New(al *register.Register, mar, mbr *[2]register.Register, stack *stack.Stack) Alu{
  return Alu{Al: al, Mar: mar, Mbr: mbr, Stack: stack}
}

func (a *Alu) Addition(s string){
  bT1 := a.Temp1.GetValue()
  bT2 := a.Temp2.GetValue()
  t1 := util.BinaryToDecimal(bT1[:])
  t2 := util.BinaryToDecimal(bT2[:])
  res := t1 + t2 

  if((int(t1) + int(t2)) > 255){
    a.Carry = true
  } else {
    a.Carry = false
  }
  if(s == "s"){
    a.Stack.Push(util.DecimalToBinary(res))
  } else{
    a.Al.SetLoad()
    a.Al.LoadValue(util.DecimalToBinary(res))
  }

}
