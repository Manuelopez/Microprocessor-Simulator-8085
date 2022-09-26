package alu

import (
	"fmt"
	"micp-sim/register"
	"micp-sim/stack"
	"micp-sim/util"
)

type Alu struct {
	Temp1      register.Register     `json:"temp1"`
	Temp2      register.Register     `json:"temp2"`
	Al         *register.Register    `json:"al"`
	Stack      *stack.Stack          `json:"stack"`
	Carry      bool                  `json:"carry"`
	Zero       bool                  `json:"zero"`
	Comparison register.Register     `json:"comparison"`
	Mar        *[2]register.Register `json:"mar"`
	Mbr        *[2]register.Register `json:"mbr"`
}

func New(al *register.Register, mar, mbr *[2]register.Register, stack *stack.Stack) Alu {
	return Alu{Al: al, Mar: mar, Mbr: mbr, Stack: stack}
}

func (a *Alu) Addition(s string) {
	bT1 := a.Temp1.GetValue()
	bT2 := a.Temp2.GetValue()
	t1 := util.BinaryToDecimal(bT1[:])
	t2 := util.BinaryToDecimal(bT2[:])
	res := t1 + t2
	byteRes := byte(res)
	if res > 255 {
		a.Carry = true
	} else {
		a.Carry = false
	}
	if s == "s" {

		a.Stack.Push(util.DecimalToBinary(int(byteRes)))
	} else {
		a.Al.SetLoad()
		a.Al.LoadValue(util.DecimalToBinary(int(byteRes)))
	}
}

func (a *Alu) Multiplication(s string) {
	bT1 := a.Temp1.GetValue()
	bT2 := a.Temp2.GetValue()
	t1 := util.BinaryToDecimal(bT1[:])
	t2 := util.BinaryToDecimal(bT2[:])
	res := t1 * t2

	byteRes := byte(res)
	if res > 255 {
		a.Carry = true
	} else {
		a.Carry = false
	}
	if s == "s" {
		a.Stack.Push(util.DecimalToBinary(int(byteRes)))
	} else {
		a.Al.SetLoad()
		a.Al.LoadValue(util.DecimalToBinary(int(byteRes)))
	}
}

func (a *Alu) Subtraction(s string) {
	bT1 := a.Temp1.GetValue()
	bT2 := a.Temp2.GetValue()
	t1 := util.BinaryToDecimal(bT1[:])
	t2 := util.BinaryToDecimal(bT2[:])
	res := t1 - t2

	byteRes := byte(res)
	if res < 0 {
		a.Carry = true
	} else {
		a.Carry = false
	}
	if s == "s" {
		a.Stack.Push(util.DecimalToBinary(int(byteRes)))
	} else {
		a.Al.SetLoad()
		a.Al.LoadValue(util.DecimalToBinary(int(byteRes)))
	}
}

func (a *Alu) Division(s string) {
	bT1 := a.Temp1.GetValue()
	bT2 := a.Temp2.GetValue()
	t1 := util.BinaryToDecimal(bT1[:])
	t2 := util.BinaryToDecimal(bT2[:])
	if t2 == 0 {
		fmt.Println("0 Division")
		return
	}
	res := t1 / t2

	byteRes := byte(res)

	if s == "s" {
		a.Stack.Push(util.DecimalToBinary(int(byteRes)))
	} else {
		a.Al.SetLoad()
		a.Al.LoadValue(util.DecimalToBinary(int(byteRes)))
	}
}

func (a *Alu) Xor(s string) {
	bT1 := a.Temp1.GetValue()
	bT2 := a.Temp2.GetValue()
	t1 := util.BinaryToDecimal(bT1[:])
	t2 := util.BinaryToDecimal(bT2[:])
	res := t1 ^ t2

	byteRes := byte(res)

	if s == "s" {
		a.Stack.Push(util.DecimalToBinary(int(byteRes)))
	} else {
		a.Al.SetLoad()
		a.Al.LoadValue(util.DecimalToBinary(int(byteRes)))
	}
}

func (a *Alu) Orl(s string) {
	bT1 := a.Temp1.GetValue()
	bT2 := a.Temp2.GetValue()
	t1 := util.BinaryToDecimal(bT1[:])
	t2 := util.BinaryToDecimal(bT2[:])
	res := t1 | t2

	byteRes := byte(res)

	if s == "s" {
		a.Stack.Push(util.DecimalToBinary(int(byteRes)))
	} else {
		a.Al.SetLoad()
		a.Al.LoadValue(util.DecimalToBinary(int(byteRes)))
	}
}

func (a *Alu) And(s string) {
	bT1 := a.Temp1.GetValue()
	bT2 := a.Temp2.GetValue()
	t1 := util.BinaryToDecimal(bT1[:])
	t2 := util.BinaryToDecimal(bT2[:])
	res := t1 & t2

	byteRes := byte(res)

	if s == "s" {
		a.Stack.Push(util.DecimalToBinary(int(byteRes)))
	} else {
		a.Al.SetLoad()
		a.Al.LoadValue(util.DecimalToBinary(int(byteRes)))
	}
}

func (a *Alu) Increment(value [8]byte) [8]byte {
	t1 := util.BinaryToDecimal(value[:])
	res := t1 + 1

	if res == 0 {
		a.Zero = true
	} else {
		a.Zero = false
	}

	if res > 255 {
		a.Carry = true
	} else {
		a.Carry = false
	}

	byteRes := byte(res)

	binaryByteRes := util.DecimalToBinary(int(byteRes))
	return binaryByteRes

}

func (a *Alu) Decrement(value [8]byte) [8]byte {
	t1 := util.BinaryToDecimal(value[:])
	res := t1 - 1

	if res == 0 {
		a.Zero = true
	} else {
		a.Zero = false
	}
	if res < 0 {
		a.Carry = true
	} else {
		a.Carry = false
	}

	byteRes := byte(res)

	binaryByteRes := util.DecimalToBinary(int(byteRes))
	return binaryByteRes
}

func (a *Alu) Not(value [8]byte) [8]byte {
	t1 := util.BinaryToDecimal(value[:])

	res := ^t1
	byteRes := byte(res)

	return util.DecimalToBinary(int(byteRes))
}

func (a *Alu) Equal() {
	bT1 := a.Temp1.GetValue()
	bT2 := a.Temp2.GetValue()
	t1 := util.BinaryToDecimal(bT1[:])
	t2 := util.BinaryToDecimal(bT2[:])
	res := t1 - t2

	a.Comparison.SetLoad()
	if res > 0 {
		a.Comparison.LoadValue(util.DecimalToBinary(1))
	} else if res == 0 {
		a.Comparison.LoadValue(util.DecimalToBinary(0))
	} else {
		a.Comparison.LoadValue(util.DecimalToBinary(255))
	}

}
