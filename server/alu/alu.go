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

	result, carry := a.AdditionLogic(bT1, bT2)
	if carry == true {
		a.Carry = true
	} else {
		a.Carry = false
	}

	if s == "s" {
		a.Stack.Push(result)
	} else {
		a.Al.SetLoad()
		a.Al.LoadValue(result)
	}
}

func (a *Alu) AdditionLogic(valA, valB [8]byte) ([8]byte, bool) {
	sum := false
	carry := false
	result := [8]byte{}
	for i := 7; i >= 0; i-- {

		if i == 7 {
			sum, carry = a.HalfAdder(valA[i] == 1, valB[i] == 1)
		} else {
			sum, carry = a.FullAdder(valA[i] == 1, valB[i] == 1, carry)
		}

		if sum == true {
			result[i] = 1
		} else {
			result[i] = 0
		}
	}

	return result, carry
}

func (a *Alu) ComplementLogic(val [8]byte) [8]byte {
	result := [8]byte{}
	for i := 0; i < 8; i++ {
		byte := a.NotGate(val[i] == 1)
		if byte {
			result[i] = 1
		} else {
			result[i] = 0
		}
	}

	return result
}

func (a *Alu) Multiplication(s string) {
	bT1 := a.Temp1.GetValue()
	bT2 := a.Temp2.GetValue()
	// t1 := util.BinaryToDecimal(bT1[:])
	// t2 := util.BinaryToDecimal(bT2[:])
	// res := t1 * t2
	//
	// byteRes := byte(res)
	// if res > 255 {
	// 	a.Carry = true
	// } else {
	// 	a.Carry = false
	// }
	// if s == "s" {
	// 	a.Stack.Push(util.DecimalToBinary(int(byteRes)))
	// } else {
	// 	a.Al.SetLoad()
	// 	a.Al.LoadValue(util.DecimalToBinary(int(byteRes)))
	// }

    rest, carry := a.MultiplicationLogic(bT1, bT2)

	a.Carry = carry
	if s == "s" {
		a.Stack.Push(rest)
	} else {
		a.Al.SetLoad()
		a.Al.LoadValue(rest)
	}

}

func (a *Alu) MultiplicationLogic(m, q [8]byte) ([8]byte, bool){
    al := [8]byte{}
    carry := false
    for i := 0; i < 8; i++{
        if q[7] == 1{
            carry2 := false
            al, carry2 = a.AdditionLogic(m, al)
            carry = a.OrGate(carry, carry2)
        }

        q, _ = a.ShiftRight(q)
        al, _ = a.ShiftLeft(al)
    }

    return al, carry 
}

func (a *Alu) Subtraction(s string) {
	bT1 := a.Temp1.GetValue()
	bT2 := a.Temp2.GetValue()

	res, carry := a.SubtractionLogic(bT1, bT2)
	if carry {
		a.Carry = true
	} else {
		a.Carry = false
	}
	if s == "s" {
		a.Stack.Push(res)
    } else {
		a.Al.SetLoad()
		a.Al.LoadValue(res)
	}
}

func (a *Alu) SubtractionLogic(valA, valB [8]byte) ([8]byte, bool){
	sum := false
	carry := false
	result := [8]byte{}
	for i := 7; i >= 0; i-- {

		if i == 7 {
			sum, carry = a.HalfSubtracter(valA[i] == 1, valB[i] == 1)
		} else {
			sum, carry = a.FullSubtracter(valA[i] == 1, valB[i] == 1, carry)
		}

		if sum == true {
			result[i] = 1
		} else {
			result[i] = 0
		}
	}

	return result, carry
}

func (a *Alu) TwosComplementLogic(val [8]byte) ([8]byte, bool){
	complement := a.ComplementLogic(val)
    twosComplement, carry := a.AdditionLogic(complement, [8]byte{0, 0, 0, 0, 0, 0, 0, 1})

    return twosComplement, carry
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

func (a *Alu) AndGate(b1 bool, b2 bool) bool {
	return b1 && b2
}

func (a *Alu) OrGate(b1 bool, b2 bool) bool {
	return b1 || b2
}

func (a *Alu) NotGate(b1 bool) bool {
	return !b1
}

func (a *Alu) XorGate(b1, b2 bool) bool {
	return (b1 || b2) && !(b1 && b2)
}

func (a *Alu) HalfAdder(b1, b2 bool) (bool, bool) {
	sum := a.XorGate(b1, b2)
	carry := a.AndGate(b1, b2)
	return sum, carry
}

func (a *Alu) FullAdder(b1, b2, c bool) (bool, bool) {
	sum, carry1 := a.HalfAdder(b1, b2)
	sum, carry2 := a.HalfAdder(sum, c)
	carry := a.OrGate(carry1, carry2)
	return sum, carry
}

func (a *Alu) HalfSubtracter(b1, b2 bool) (bool, bool){
    sub := a.XorGate(b1, b2)
    carry := a.AndGate(a.NotGate(b1), b2)
    return sub, carry
}

func (a *Alu) ShiftRight(b [8]byte) ([8]byte, byte){
	disposed := b[7]
	for i := 7; i > 0; i-- {
		b[i] = b[i-1]
		if i == 1 {
			b[i-1] = 0
		}
	}
	return b, disposed
}

func (a *Alu) ShiftLeft(b [8]byte) ([8]byte, byte){
    disposed := b[0]

    for i := 0; i < 7; i++{
        b[i] = b[i+1]
    }

    b[7] = 0

    return b, disposed
}

func (a *Alu) FullSubtracter(b1,b2, c bool) (bool, bool){
    sub, carry1 := a.HalfSubtracter(b1, b2)
    sub, carry2 := a.HalfSubtracter(sub, c)
    carry := a.OrGate(carry1, carry2)

    return sub, carry
}
