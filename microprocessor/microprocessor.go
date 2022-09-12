package microprocessor

import (
	"bufio"
	"fmt"
	"micp-sim/alu"
	"micp-sim/clock"
	"micp-sim/memory"
	"micp-sim/opcode"
	"micp-sim/register"
	"micp-sim/stack"
	"micp-sim/util"
	"os"
	"strconv"
	"strings"
)

var MEMORY_ADDRESS_FOR_OPERATION uint16 = 0

type MicroProcessor struct {
	// AH HIGHT BITS AL LOW BITS
	Al, Ah, B, C, D, E, L, H *register.Register
	*memory.Memory
	*stack.Stack
	*clock.Clock
	*alu.Alu
	Ir  *[2]register.Register
	St  register.Register
	Mar *[2]register.Register
	Mbr *[2]register.Register
	Pc  *[2]register.Register
}

func New(freq float64) MicroProcessor {

	mar := &[2]register.Register{register.New(), register.New()}
	mbr := &[2]register.Register{register.New(), register.New()}
	pc := &[2]register.Register{register.New(), register.New()}

	al := register.New()
	ah := register.New()
	b := register.New()
	c := register.New()
	d := register.New()
	e := register.New()
	l := register.New()
	h := register.New()

	memory := memory.New(mbr, mar)
	clock := clock.New(freq)
	stack := stack.New()
	alu := alu.New(&al, mar, mbr, &stack)
	return MicroProcessor{
		Al:     &al,
		Ah:     &ah,
		B:      &b,
		C:      &c,
		D:      &d,
		E:      &e,
		L:      &l,
		H:      &h,
		Memory: &memory,
		Mar:    mar,
		Mbr:    mbr,
		Clock:  &clock,
		Stack:  &stack,
		Alu:    &alu,
		Ir:     &[2]register.Register{register.New(), register.New()},
		St:     register.New(),
		Pc:     pc,
	}
}

func (m *MicroProcessor) Start() {
	go m.Clock.TurnOn()
	m.LoadInstructions("./instructions.txt")
	for {
		m.ReadInstructon()
		endProgram := m.Execute()
		if endProgram {
			m.Clock.TurnOff()
			break
		}
	}

	//	a := m.Al.GetValue()
	al := m.Al.GetValue()
	fmt.Println(util.BinaryToDecimal(al[:]))

}

func (m *MicroProcessor) Test() {

	m.LoadInstructions("./instructions.txt")
	for {
		m.ReadInstructon()
		endProgram := m.Execute()
		if endProgram {
			break
		}
	}

	a := m.Al.GetValue()
	fmt.Println(util.BinaryToDecimal(a[:]))

}

func (m *MicroProcessor) Execute() bool {
	hbitsInst := m.Ir[util.HIGH_BITS].GetValue()
	lbitsInst := m.Ir[util.LOW_BITS].GetValue()
	code := util.BinaryToDecimal(hbitsInst[:])

	switch code {
	case opcode.BEGIN:
		fmt.Println("PROGRAM STARTED")
		m.Clock.Wait()
	case opcode.END:
		fmt.Println("PROGRAM ENDED")
		m.Clock.Wait()
		return true

	case opcode.ADD_VAL:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Clock.Wait()
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Clock.Wait()
		m.Alu.Addition("")
		m.Clock.Wait()
	case opcode.ADD_AH:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.Ah.GetValue())
		m.Alu.Addition("")
	case opcode.ADD_B:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.B.GetValue())
		m.Alu.Addition("")
	case opcode.ADD_C:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.C.GetValue())
		m.Alu.Addition("")
	case opcode.ADD_D:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.D.GetValue())
		m.Alu.Addition("")
	case opcode.ADD_E:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.E.GetValue())
		m.Alu.Addition("")
	case opcode.ADD_L:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.L.GetValue())
		m.Alu.Addition("")
	case opcode.ADD_H:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.H.GetValue())
		m.Alu.Addition("")
	case opcode.ADD_S:
		val1, err := m.Stack.Pop()
		check(err)
		val2, err := m.Stack.Pop()
		check(err)
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(val1)
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(val2)
		m.Alu.Addition("s")
	case opcode.ADD_MEM:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.LoadMarWithPC()
		m.Read()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Memory.Read()

		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Alu.Addition("")
		m.IncreasePc()

		// Sta
	case opcode.STA:
		m.LoadMarWithPC()
		m.Memory.Read()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Memory.Read()
		m.Mbr[util.HIGH_BITS].SetLoad()
		m.Mbr[util.HIGH_BITS].LoadValue(m.Ah.GetValue())
		m.Mbr[util.LOW_BITS].SetLoad()
		m.Mbr[util.LOW_BITS].LoadValue(m.Al.GetValue())
		m.Write()

		//-------- MOV AL
	case opcode.MOV_AL_AH:
		m.Al.SetLoad()
		m.Al.LoadValue(m.Ah.GetValue())
	case opcode.MOV_AL_B:
		m.Al.SetLoad()
		m.Al.LoadValue(m.B.GetValue())
	case opcode.MOV_AL_C:
		m.Al.SetLoad()
		m.Al.LoadValue(m.C.GetValue())
	case opcode.MOV_AL_D:
		m.Al.SetLoad()
		m.Al.LoadValue(m.D.GetValue())
	case opcode.MOV_AL_E:
		m.Al.SetLoad()
		m.Al.LoadValue(m.E.GetValue())
	case opcode.MOV_AL_L:
		m.Al.SetLoad()
		m.Al.LoadValue(m.L.GetValue())
	case opcode.MOV_AL_H:
		m.Al.SetLoad()
		m.Al.LoadValue(m.H.GetValue())

		//------ MOV AH
	case opcode.MOV_AH_AL:
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.Al.GetValue())
	case opcode.MOV_AH_B:
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.B.GetValue())
	case opcode.MOV_AH_C:
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.C.GetValue())
	case opcode.MOV_AH_D:
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.D.GetValue())
	case opcode.MOV_AH_E:
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.E.GetValue())
	case opcode.MOV_AH_L:
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.L.GetValue())
	case opcode.MOV_AH_H:
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.H.GetValue())

		// ------ MOV VAL
	case opcode.MOV_AL_VAL:
		m.Al.SetLoad()
		m.Al.LoadValue(lbitsInst)
		m.Clock.Wait()
	case opcode.MOV_AH_VAL:
		m.Ah.SetLoad()
		m.Ah.LoadValue(lbitsInst)
		m.Clock.Wait()
	case opcode.MOV_B_VAL:
		m.B.SetLoad()
		m.B.LoadValue(lbitsInst)
		m.Clock.Wait()
	case opcode.MOV_C_VAL:
		m.C.SetLoad()
		m.C.LoadValue(lbitsInst)
		m.Clock.Wait()
	case opcode.MOV_D_VAL:
		m.D.SetLoad()
		m.D.LoadValue(lbitsInst)
		m.Clock.Wait()
	case opcode.MOV_E_VAL:
		m.E.SetLoad()
		m.E.LoadValue(lbitsInst)
		m.Clock.Wait()
	case opcode.MOV_H_VAL:
		m.H.SetLoad()
		m.H.LoadValue(lbitsInst)
		m.Clock.Wait()
	case opcode.MOV_L_VAL:
		m.L.SetLoad()
		m.L.LoadValue(lbitsInst)
		m.Clock.Wait()

		// ------- MOV B
	case opcode.MOV_B_AL:
		m.B.SetLoad()
		m.B.LoadValue(m.Al.GetValue())
	case opcode.MOV_B_AH:
		m.B.SetLoad()
		m.B.LoadValue(m.Ah.GetValue())
	case opcode.MOV_B_C:
		m.B.SetLoad()
		m.B.LoadValue(m.C.GetValue())
	case opcode.MOV_B_D:
		m.B.SetLoad()
		m.B.LoadValue(m.D.GetValue())
	case opcode.MOV_B_E:
		m.B.SetLoad()
		m.B.LoadValue(m.E.GetValue())
	case opcode.MOV_B_H:
		m.B.SetLoad()
		m.B.LoadValue(m.H.GetValue())
	case opcode.MOV_B_L:
		m.B.SetLoad()
		m.B.LoadValue(m.L.GetValue())
	case opcode.MOV_B_M:
		m.B.SetLoad()
		m.LoadMarWithPC()
		m.Memory.Read()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Memory.Read()
		m.B.LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.IncreasePc()

		// -------- MOV C
	case opcode.MOV_C_AL:
		m.C.SetLoad()
		m.C.LoadValue(m.Al.GetValue())
	case opcode.MOV_C_AH:
		m.C.SetLoad()
		m.C.LoadValue(m.Ah.GetValue())
	case opcode.MOV_C_B:
		m.C.SetLoad()
		m.C.LoadValue(m.B.GetValue())
	case opcode.MOV_C_D:
		m.C.SetLoad()
		m.C.LoadValue(m.D.GetValue())
	case opcode.MOV_C_E:
		m.C.SetLoad()
		m.C.LoadValue(m.E.GetValue())
	case opcode.MOV_C_H:
		m.C.SetLoad()
		m.C.LoadValue(m.H.GetValue())
	case opcode.MOV_C_L:
		m.C.SetLoad()
		m.C.LoadValue(m.L.GetValue())
	case opcode.MOV_C_M:
		m.LoadMarWithPC()
		m.Memory.Read()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Memory.Read()
		m.C.SetLoad()
		m.C.LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.IncreasePc()

		// -------- MOV D
	case opcode.MOV_D_AL:
		m.D.SetLoad()
		m.D.LoadValue(m.Al.GetValue())
	case opcode.MOV_D_AH:
		m.D.SetLoad()
		m.D.LoadValue(m.Ah.GetValue())
	case opcode.MOV_D_B:
		m.D.SetLoad()
		m.D.LoadValue(m.B.GetValue())
	case opcode.MOV_D_C:
		m.D.SetLoad()
		m.D.LoadValue(m.C.GetValue())
	case opcode.MOV_D_E:
		m.D.SetLoad()
		m.D.LoadValue(m.E.GetValue())
	case opcode.MOV_D_L:
		m.D.SetLoad()
		m.D.LoadValue(m.L.GetValue())
	case opcode.MOV_D_H:
		m.D.SetLoad()
		m.D.LoadValue(m.H.GetValue())
	case opcode.MOV_D_M:
		m.LoadMarWithPC()
		m.Memory.Read()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Read()
		m.D.SetLoad()
		m.D.LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.IncreasePc()

		// --------MOV E
	case opcode.MOV_E_AL:
		m.E.SetLoad()
		m.E.LoadValue(m.Al.GetValue())
	case opcode.MOV_E_AH:
		m.E.SetLoad()
		m.E.LoadValue(m.Ah.GetValue())
	case opcode.MOV_E_B:
		m.E.SetLoad()
		m.E.LoadValue(m.B.GetValue())
	case opcode.MOV_E_C:
		m.E.SetLoad()
		m.E.LoadValue(m.C.GetValue())
	case opcode.MOV_E_D:
		m.E.SetLoad()
		m.E.LoadValue(m.D.GetValue())
	case opcode.MOV_E_H:
		m.E.SetLoad()
		m.E.LoadValue(m.H.GetValue())
	case opcode.MOV_E_L:
		m.E.SetLoad()
		m.E.LoadValue(m.L.GetValue())
	case opcode.MOV_E_M:
		m.LoadMarWithPC()
		m.Read()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.E.SetLoad()
		m.E.LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.IncreasePc()

		// --------- H
	case opcode.MOV_H_AL:
		m.H.SetLoad()
		m.H.LoadValue(m.Al.GetValue())
	case opcode.MOV_H_AH:
		m.H.SetLoad()
		m.H.LoadValue(m.Ah.GetValue())
	case opcode.MOV_H_B:
		m.H.SetLoad()
		m.H.LoadValue(m.B.GetValue())
	case opcode.MOV_H_C:
		m.H.SetLoad()
		m.H.LoadValue(m.C.GetValue())
	case opcode.MOV_H_D:
		m.H.SetLoad()
		m.H.LoadValue(m.D.GetValue())
	case opcode.MOV_H_E:
		m.H.SetLoad()
		m.H.LoadValue(m.E.GetValue())
	case opcode.MOV_H_L:
		m.H.SetLoad()
		m.H.LoadValue(m.L.GetValue())

	// --------- MOV  L
	case opcode.MOV_L_AL:
		m.L.SetLoad()
		m.L.LoadValue(m.Al.GetValue())
	case opcode.MOV_L_AH:
		m.L.SetLoad()
		m.L.LoadValue(m.Ah.GetValue())
	case opcode.MOV_L_B:
		m.L.SetLoad()
		m.L.LoadValue(m.B.GetValue())
	case opcode.MOV_L_C:
		m.L.SetLoad()
		m.L.LoadValue(m.C.GetValue())
	case opcode.MOV_L_D:
		m.L.SetLoad()
		m.L.LoadValue(m.D.GetValue())
	case opcode.MOV_L_E:
		m.L.SetLoad()
		m.L.LoadValue(m.E.GetValue())
	case opcode.MOV_L_H:
		m.L.SetLoad()
		m.L.LoadValue(m.H.GetValue())

		// --------- MOV Memory
	case opcode.MOV_M_AL:
		alBinary := m.Al.GetValue()
		m.LoadMarWithPC()
		m.Memory.Read()
		hbitsAddress, lbitsAddress := m.Mbr[util.HIGH_BITS].GetValue(), m.Mbr[util.LOW_BITS].GetValue()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(hbitsAddress)
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(lbitsAddress)
		m.Mbr[util.HIGH_BITS].SetLoad()
		m.Mbr[util.HIGH_BITS].LoadValue(util.DecimalToBinary(0))
		m.Mbr[util.LOW_BITS].SetLoad()
		m.Mbr[util.LOW_BITS].LoadValue(alBinary)
		m.Memory.Write()
		m.IncreasePc()
	case opcode.MOV_M_AH:
		m.LoadMarWithPC()
		m.Memory.Read()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Mbr[util.HIGH_BITS].SetLoad()
		m.Mbr[util.HIGH_BITS].LoadValue(util.DecimalToBinary(0))
		m.Mbr[util.LOW_BITS].SetLoad()
		m.Mbr[util.LOW_BITS].LoadValue(m.Ah.GetValue())
		m.Memory.Write()
		m.IncreasePc()
	case opcode.MOV_M_B:
		m.LoadMarWithPC()
		m.Memory.Read()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Mbr[util.HIGH_BITS].SetLoad()
		m.Mbr[util.HIGH_BITS].LoadValue(util.DecimalToBinary(0))
		m.Mbr[util.LOW_BITS].SetLoad()
		m.Mbr[util.LOW_BITS].LoadValue(m.B.GetValue())
		m.Memory.Write()
		m.IncreasePc()
	case opcode.MOV_M_C:
		m.LoadMarWithPC()
		m.Memory.Read()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Mbr[util.HIGH_BITS].SetLoad()
		m.Mbr[util.HIGH_BITS].LoadValue(util.DecimalToBinary(0))
		m.Mbr[util.LOW_BITS].SetLoad()
		m.Mbr[util.LOW_BITS].LoadValue(m.C.GetValue())
		m.Memory.Write()
		m.IncreasePc()
	case opcode.MOV_M_D:
		m.LoadMarWithPC()
		m.Memory.Read()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Mbr[util.HIGH_BITS].SetLoad()
		m.Mbr[util.HIGH_BITS].LoadValue(util.DecimalToBinary(0))
		m.Mbr[util.LOW_BITS].SetLoad()
		m.Mbr[util.LOW_BITS].LoadValue(m.D.GetValue())
		m.Memory.Write()
		m.IncreasePc()
	case opcode.MOV_M_E:
		m.LoadMarWithPC()
		m.Memory.Read()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Mbr[util.HIGH_BITS].SetLoad()
		m.Mbr[util.HIGH_BITS].LoadValue(util.DecimalToBinary(0))
		m.Mbr[util.LOW_BITS].SetLoad()
		m.Mbr[util.LOW_BITS].LoadValue(m.E.GetValue())
		m.Memory.Write()
		m.IncreasePc()

		// ---------- MOV AX
	case opcode.MOV_AX_M:
		m.LoadMarWithPC()
		m.Memory.Read()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Memory.Read()
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Al.SetLoad()
		m.Al.LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.IncreasePc()
		//	MUL_AL_VAL = 0x47
	case opcode.MUL_AL_VAL:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Multiplication("")

		//MUL_AL_M   = 0x48
	case opcode.MUL_AL_M:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.LoadMarWithPC()
		m.Read()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Memory.Read()

		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Alu.Multiplication("")
		m.IncreasePc()

		//MUL_AL_AH  = 0x49
	case opcode.MUL_AL_AH:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.Ah.GetValue())
		m.Alu.Multiplication("")

		//MUL_AL_B   = 0x4A
	case opcode.MUL_AL_B:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.B.GetValue())
		m.Alu.Multiplication("")

		//MUL_AL_C   = 0x4B
	case opcode.MUL_AL_C:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.C.GetValue())
		m.Alu.Multiplication("")

		//MUL_AL_D   = 0x4C
	case opcode.MUL_AL_D:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.D.GetValue())
		m.Alu.Multiplication("")

		//MUL_AL_E   =n 0x4D
	case opcode.MUL_AL_E:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.E.GetValue())
		m.Alu.Multiplication("")

		//MUL_AL_L   = 0x4E
	case opcode.MUL_AL_L:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.L.GetValue())
		m.Alu.Multiplication("")

		//MUL_AL_H   = 0x4F
	case opcode.MUL_AL_H:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.H.GetValue())
		m.Alu.Multiplication("")

		//MUL_S      = 0xF2
	case opcode.MUL_S:
		t1, err := m.Stack.Pop()
		check(err)
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(t1)
		t2, err := m.Stack.Pop()
		check(err)
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(t2)
		m.Alu.Multiplication("s")

	//	SUB_AL_VAL = 0x3E
	case opcode.SUB_AL_VAL:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Subtraction("")

		//SUB_AL_M   = 0x3F
	case opcode.SUB_AL_M:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.LoadMarWithPC()
		m.Read()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Memory.Read()

		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Alu.Subtraction("")
		m.IncreasePc()

	//SUB_AL_AH  = 0x40
	case opcode.SUB_AL_AH:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.Ah.GetValue())
		m.Alu.Subtraction("")

		//SUB_AL_B   = 0x41
	case opcode.SUB_AL_B:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.B.GetValue())
		m.Alu.Subtraction("")

		//SUB_AL_C   = 0x42
	case opcode.SUB_AL_C:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.C.GetValue())
		m.Alu.Subtraction("")

	//SUB_AL_D   = 0x43
	case opcode.SUB_AL_D:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.D.GetValue())
		m.Alu.Subtraction("")

	//SUB_AL_E   = 0x44
	case opcode.SUB_AL_E:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.E.GetValue())
		m.Alu.Subtraction("")

	//SUB_AL_L   = 0x45
	case opcode.SUB_AL_L:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.L.GetValue())
		m.Alu.Subtraction("")

	//SUB_AL_H   = 0x46
	case opcode.SUB_AL_H:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.H.GetValue())
		m.Alu.Subtraction("")

	//SUB_S      = 0xF1
	case opcode.SUB_S:
		t1, err := m.Stack.Pop()
		check(err)
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(t1)
		t2, err := m.Stack.Pop()
		check(err)
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(t2)
		m.Alu.Subtraction("s")

		//DIV_AL_VAL = 0x50
	case opcode.DIV_AL_VAL:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Division("")

		//DIV_AL_M   = 0x51
	case opcode.DIV_AL_M:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.LoadMarWithPC()
		m.Read()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Memory.Read()

		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Alu.Division("")
		m.IncreasePc()

		//DIV_AL_AH  = 0x52
	case opcode.DIV_AL_AH:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.Ah.GetValue())
		m.Alu.Division("")

		//DIV_AL_B   = 0x53
	case opcode.DIV_AL_B:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.B.GetValue())
		m.Alu.Division("")

		//DIV_AL_C   = 0x54
	case opcode.DIV_AL_C:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.C.GetValue())
		m.Alu.Division("")

		//DIV_AL_D   = 0x55
	case opcode.DIV_AL_D:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.D.GetValue())
		m.Alu.Division("")

		//DIV_AL_E   = 0x56
	case opcode.DIV_AL_E:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.E.GetValue())
		m.Alu.Division("")

		//DIV_AL_L   = 0x57
	case opcode.DIV_AL_L:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.Al.GetValue())
		m.Alu.Division("")

		//DIV_AL_H   = 0x58
	case opcode.DIV_AL_H:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.H.GetValue())
		m.Alu.Division("")

		//DIV_S      = 0xF3
	case opcode.DIV_S:
		t1, err := m.Stack.Pop()
		check(err)
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(t1)
		t2, err := m.Stack.Pop()
		check(err)
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(t2)
		m.Alu.Division("s")

	// XOR_AL_H   = 0x06
	case opcode.XOR_AL_H:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.H.GetValue())
		m.Alu.Xor("")
	case opcode.XOR_AL_M:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.LoadMarWithPC()
		m.Read()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Memory.Read()

		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Alu.Xor("")
		m.IncreasePc()
	case opcode.XOR_AL_VAL:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Xor("")
	case opcode.XOR_AL_AH:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.Ah.GetValue())
		m.Alu.Xor("")
	case opcode.XOR_AL_B:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.B.GetValue())
		m.Alu.Xor("")
	case opcode.XOR_AL_C:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.C.GetValue())
		m.Alu.Xor("")
	case opcode.XOR_AL_D:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.D.GetValue())
		m.Alu.Xor("")
	case opcode.XOR_AL_E:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.E.GetValue())
		m.Alu.Xor("")
	case opcode.XOR_AL_L:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.L.GetValue())
		m.Alu.Xor("")
	case opcode.XOR_S:
		t1, err := m.Stack.Pop()
		check(err)
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(t1)
		t2, err := m.Stack.Pop()
		check(err)
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(t2)
		m.Alu.Xor("s")

	//	ORL_AL_AH  = 0x1D
	case opcode.ORL_AL_AH:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.Ah.GetValue())
		m.Alu.Orl("")
		//ORL_AL_B   = 0x1E
	case opcode.ORL_AL_B:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.B.GetValue())
		m.Alu.Orl("")
		//ORL_AL_C   = 0x1F
	case opcode.ORL_AL_C:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.C.GetValue())
		m.Alu.Orl("")
		//ORL_AL_D   = 0x20
	case opcode.ORL_AL_D:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.D.GetValue())
		m.Alu.Orl("")
		//ORL_AL_E   = 0x21
	case opcode.ORL_AL_E:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.E.GetValue())
		m.Alu.Orl("")
		//ORL_AL_L   = 0x22
	case opcode.ORL_AL_L:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.L.GetValue())
		m.Alu.Orl("")
		//ORL_AL_H   = 0x23
	case opcode.ORL_AL_H:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.H.GetValue())
		m.Alu.Orl("")
		//ORL_AL_M   = 0x2B
	case opcode.ORL_AL_M:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.LoadMarWithPC()
		m.Read()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Memory.Read()

		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Alu.Orl("")
		m.IncreasePc()

	//ORL_AL_VAL = 0x2D
	case opcode.ORL_AL_VAL:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Orl("")

		//ORL_S      = 0xF4
	case opcode.ORL_S:
		t1, err := m.Stack.Pop()
		check(err)
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(t1)
		t2, err := m.Stack.Pop()
		check(err)
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(t2)
		m.Alu.Orl("s")

		//	AND_AL_AH  = 0x24
	case opcode.AND_AL_AH:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.Ah.GetValue())
		m.Alu.And("")

		//AND_AL_B   = 0x25
	case opcode.AND_AL_B:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.B.GetValue())
		m.Alu.And("")
		//AND_AL_C   = 0x26
	case opcode.AND_AL_C:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.C.GetValue())
		m.Alu.And("")
		//AND_AL_D   = 0x27
	case opcode.AND_AL_D:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.D.GetValue())
		m.Alu.And("")
		//AND_AL_E   = 0x28
	case opcode.AND_AL_E:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.E.GetValue())
		m.Alu.And("")
		//AND_AL_L   = 0x29
	case opcode.AND_AL_L:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.L.GetValue())
		m.Alu.And("")
		//AND_AL_H   = 0x2A
	case opcode.AND_AL_H:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.H.GetValue())
		m.Alu.And("")
		//AND_AL_M   = 0x2C
	case opcode.AND_AL_M:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.LoadMarWithPC()
		m.Read()
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Memory.Read()

		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		m.Alu.And("")
		m.IncreasePc()
		//AND_AL_VAL = 0x2E
	case opcode.AND_AL_VAL:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.And("")
		//AND_S      = 0xF5
	case opcode.AND_S:
		t1, err := m.Stack.Pop()
		check(err)
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(t1)
		t2, err := m.Stack.Pop()
		check(err)
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(t2)
		m.Alu.And("s")

	//	INC_AL = 0x64
	case opcode.INC_AL:
		m.Al.SetLoad()
		m.Al.LoadValue(m.Alu.Increment(m.Al.GetValue()))
		//INC_AH = 0x65
	case opcode.INC_AH:
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.Alu.Increment(m.Ah.GetValue()))

		//INC_B  = 0x66
	case opcode.INC_B:
		m.B.SetLoad()
		m.B.LoadValue(m.Alu.Increment(m.B.GetValue()))

		//INC_C  = 0x67
	case opcode.INC_C:
		m.C.SetLoad()
		m.C.LoadValue(m.Alu.Increment(m.C.GetValue()))

		//INC_D  = 0x68
	case opcode.INC_D:
		m.D.SetLoad()
		m.D.LoadValue(m.Alu.Increment(m.D.GetValue()))

		//INC_E  = 0x69
	case opcode.INC_E:
		m.E.SetLoad()
		m.E.LoadValue(m.Alu.Increment(m.E.GetValue()))

		//INC_H  = 0x6A
	case opcode.INC_H:
		m.H.SetLoad()
		m.H.LoadValue(m.Alu.Increment(m.H.GetValue()))

		//INC_L  = 0x6B
	case opcode.INC_L:
		m.L.SetLoad()
		m.L.LoadValue(m.Alu.Increment(m.L.GetValue()))

	//	DEC_AL = 0x6C
	case opcode.DEC_AL:
		m.Al.SetLoad()
		m.Al.LoadValue(m.Alu.Decrement(m.Al.GetValue()))

	//	DEC_AH = 0x6D
	case opcode.DEC_AH:
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.Alu.Decrement(m.Ah.GetValue()))

	//	DEC_B  = 0x6E
	case opcode.DEC_B:
		m.B.SetLoad()
		m.B.LoadValue(m.Alu.Decrement(m.B.GetValue()))

		//DEC_C  = 0x6F
	case opcode.DEC_C:
		m.C.SetLoad()
		m.C.LoadValue(m.Alu.Decrement(m.C.GetValue()))

		//DEC_D  = 0x70
	case opcode.DEC_D:
		m.D.SetLoad()
		m.D.LoadValue(m.Alu.Decrement(m.D.GetValue()))

		//DEC_E  = 0x71
	case opcode.DEC_E:
		m.E.SetLoad()
		m.E.LoadValue(m.Alu.Decrement(m.E.GetValue()))

		//DEC_H  = 0x72
	case opcode.DEC_H:
		m.H.SetLoad()
		m.H.LoadValue(m.Alu.Decrement(m.H.GetValue()))

		//DEC_L  = 0x73
	case opcode.DEC_L:
		m.L.SetLoad()
		m.L.LoadValue(m.Alu.Decrement(m.L.GetValue()))

//NOT_AL = 0x15
	//NOT_AH = 0x16
	//NOT_B  = 0x17
	//NOT_C  = 0x18
	//NOT_D  = 0x19
	//NOT_E  = 0x1A

	default:
		fmt.Println("OPERATION NOT IMPLEMENTED")
	}

	return false
}

func (m *MicroProcessor) LoadMarWithPC() {
	hbits, lbits := m.GetPcBinary()
	m.Mar[util.HIGH_BITS].SetLoad()
	m.Mar[util.HIGH_BITS].LoadValue(hbits)
	m.Mar[util.LOW_BITS].SetLoad()
	m.Mar[util.LOW_BITS].LoadValue(lbits)

}

func (m *MicroProcessor) ReadInstructon() {

	m.LoadMarWithPC()

	m.Memory.Read()
	hbitsMbr := m.Mbr[util.HIGH_BITS].GetValue()
	lbitsMbr := m.Mbr[util.LOW_BITS].GetValue()
	/*
		inst := [16]byte{}
		for i, _ := range hbitsMbr {
			inst[i] = hbitsMbr[i]
		}
		for i, _ := range lbitsMbr {
			inst[i+8] = lbitsMbr[i]
		}
	*/

	m.IncreasePc()
	m.Ir[util.HIGH_BITS].SetLoad()
	m.Ir[util.HIGH_BITS].LoadValue(hbitsMbr)
	m.Ir[util.LOW_BITS].SetLoad()
	m.Ir[util.LOW_BITS].LoadValue(lbitsMbr)
}

func (m MicroProcessor) GetPcBinary() ([8]byte, [8]byte) {
	hbitsPcB := m.Pc[util.HIGH_BITS].GetValue()
	lbitsPcB := m.Pc[util.LOW_BITS].GetValue()
	a := [16]byte{}
	for i, _ := range hbitsPcB {
		a[i] = hbitsPcB[i]
	}
	for i, _ := range lbitsPcB {
		a[i+8] = lbitsPcB[i]
	}
	pc := util.BinaryToDecimal(a[:])

	hbits := pc >> 8
	lbits := pc & 0xFF

	return util.DecimalToBinary(hbits), util.DecimalToBinary(lbits)

}

func (m *MicroProcessor) IncreasePc() {
	hbitsPcB := m.Pc[util.HIGH_BITS].GetValue()
	lbitsPcB := m.Pc[util.LOW_BITS].GetValue()
	a := [16]byte{}
	for i, _ := range hbitsPcB {
		a[i] = hbitsPcB[i]
	}
	for i, _ := range lbitsPcB {
		a[i+8] = lbitsPcB[i]
	}
	pc := util.BinaryToDecimal(a[:])

	pc++
	hbitsPcD := pc >> 8
	lbitsPcD := pc & 0xFF

	hbitsPcBN, lbitsPcBN := util.DecimalToBinary(hbitsPcD), util.DecimalToBinary(lbitsPcD)
	m.Pc[util.HIGH_BITS].SetLoad()
	m.Pc[util.HIGH_BITS].LoadValue(hbitsPcBN)
	m.Pc[util.LOW_BITS].SetLoad()
	m.Pc[util.LOW_BITS].LoadValue(lbitsPcBN)

}

func (m *MicroProcessor) LoadInstructions(filePath string) {
	file, err := os.Open(filePath)
	pc := 0
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		instructions := []byte(text)
		hbitsValue, lbitsValue, isMemoryOp := Assembler(string(instructions))

		hbits := pc >> 8
		lbits := pc & 0xFF
		m.Mar[util.HIGH_BITS].SetLoad()
		m.Mar[util.HIGH_BITS].LoadValue(util.DecimalToBinary(hbits))
		m.Mar[util.LOW_BITS].SetLoad()
		m.Mar[util.LOW_BITS].LoadValue(util.DecimalToBinary(lbits))
		m.Mbr[util.HIGH_BITS].SetLoad()
		m.Mbr[util.HIGH_BITS].LoadValue(hbitsValue)
		m.Mbr[util.LOW_BITS].SetLoad()
		m.Mbr[util.LOW_BITS].LoadValue(lbitsValue)
		m.Memory.Write()
		pc++

		if isMemoryOp {

			hbitsValueMemory, lbitsValueMemory := util.DecimalToBinary16(int(MEMORY_ADDRESS_FOR_OPERATION))

			hbits := pc >> 8
			lbits := pc & 0xFF
			m.Mar[util.HIGH_BITS].SetLoad()
			m.Mar[util.HIGH_BITS].LoadValue(util.DecimalToBinary(hbits))
			m.Mar[util.LOW_BITS].SetLoad()
			m.Mar[util.LOW_BITS].LoadValue(util.DecimalToBinary(lbits))
			m.Mbr[util.HIGH_BITS].SetLoad()
			m.Mbr[util.HIGH_BITS].LoadValue(hbitsValueMemory)
			m.Mbr[util.LOW_BITS].SetLoad()
			m.Mbr[util.LOW_BITS].LoadValue(lbitsValueMemory)
			m.Memory.Write()
			pc++

		}
		// m.Clock.Wait()

	}

}

func Assembler(instructions string) ([8]byte, [8]byte, bool) {
	instructions = strings.ToUpper(instructions)
	splitInstructions := strings.Split(instructions, " ")
	if len(splitInstructions) == 1 {
		switch splitInstructions[0] {
		case "BEGIN":
			return util.DecimalToBinary(opcode.BEGIN), util.DecimalToBinary(opcode.BEGIN), false
		case "END":
			return util.DecimalToBinary(opcode.END), util.DecimalToBinary(opcode.NOTHING), false
		default:
			panic(splitInstructions[0] + "is not an operation")
		}
	} else if len(splitInstructions) == 2 {
		operation := splitInstructions[0]
		operand1 := splitInstructions[1]
		switch operation {
		case "INC":
			code, ok := opcode.INC[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false
		case "DEC":
			code, ok := opcode.DEC[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false

		case "ADD":
			code, ok := opcode.ADD[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false
		case "SUB":
			code, ok := opcode.SUB[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false
		case "MUL":
			code, ok := opcode.MUL[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false
		case "DIV":
			code, ok := opcode.DIV[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false
		case "ORL":
			code, ok := opcode.ORL[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false
		case "AND":
			code, ok := opcode.AND[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false
		case "XOR":
			code, ok := opcode.XOR[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false

		}
		// TODO IMPLEMENT ALU STACK
		fmt.Println("NOT IMPLEMENTED")
	} else if len(splitInstructions) == 3 {
		operation := splitInstructions[0]
		operand1 := splitInstructions[1]
		operand2 := splitInstructions[2]

		switch operation {
		case "ADD":
			if operand1 != "AL" {
				panic("ADDITION CAN ONLY BE DONE WITH AL AS OPERAND 1 " + operand1 + " NOT ALLOWED")
			}
			if operand2[0] == 'M' {
				mOp := strings.Replace(operand2, "M0X", "", -1)
				n, err := strconv.ParseUint(mOp, 16, 64)
				check(err)
				MEMORY_ADDRESS_FOR_OPERATION = uint16(n)
				code, ok := opcode.ADD["M"]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true
			}
			val, err := strconv.Atoi(operand2)
			if err == nil {
				byteVal := byte(val)
				code, _ := opcode.ADD["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), false
			} else {
				code, ok := opcode.ADD[operand2]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false
			}
		case "MOV":
			if operand1[0] == 'M' {
				mOp := strings.Replace(operand1, "M0X", "", -1)
				n, err := strconv.ParseUint(mOp, 16, 64)
				check(err)
				MEMORY_ADDRESS_FOR_OPERATION = uint16(n)
				code, ok := opcode.MOV["M"][operand2]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true
			}
			mapedOp1, ok := opcode.MOV[operand1]
			checkOperand(ok, operand1)
			val, err := strconv.Atoi(operand2)
			if err == nil {
				byteVal := byte(val)
				code, _ := mapedOp1["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), false
			} else {
				if operand2[0] == 'M' {
					mOp := strings.Replace(operand2, "M0X", "", -1)
					n, err := strconv.ParseUint(mOp, 16, 64)
					check(err)
					MEMORY_ADDRESS_FOR_OPERATION = uint16(n)
					code, ok := mapedOp1["M"]
					checkOperand(ok, operand2)
					return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true
				}
				code, ok := mapedOp1[operand2]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false
			}
		case "MUL":
			if operand1 != "AL" {
				panic("ADDITION CAN ONLY BE DONE WITH AL AS OPERAND 1 " + operand1 + " NOT ALLOWED")
			}
			if operand2[0] == 'M' {
				mOp := strings.Replace(operand2, "M0X", "", -1)
				n, err := strconv.ParseUint(mOp, 16, 64)
				check(err)
				MEMORY_ADDRESS_FOR_OPERATION = uint16(n)
				code, ok := opcode.MUL["M"]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true
			}
			val, err := strconv.Atoi(operand2)
			if err == nil {
				byteVal := byte(val)
				code, _ := opcode.MUL["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), false
			} else {
				code, ok := opcode.MUL[operand2]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false
			}
		case "AND":
			if operand1 != "AL" {
				panic("ADDITION CAN ONLY BE DONE WITH AL AS OPERAND 1 " + operand1 + " NOT ALLOWED")
			}
			if operand2[0] == 'M' {
				mOp := strings.Replace(operand2, "M0X", "", -1)
				n, err := strconv.ParseUint(mOp, 16, 64)
				check(err)
				MEMORY_ADDRESS_FOR_OPERATION = uint16(n)
				code, ok := opcode.AND["M"]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true
			}
			val, err := strconv.Atoi(operand2)
			if err == nil {
				byteVal := byte(val)
				code, _ := opcode.AND["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), false
			} else {
				code, ok := opcode.AND[operand2]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false
			}

		case "XOR":
			if operand1 != "AL" {
				panic("ADDITION CAN ONLY BE DONE WITH AL AS OPERAND 1 " + operand1 + " NOT ALLOWED")
			}
			if operand2[0] == 'M' {
				mOp := strings.Replace(operand2, "M0X", "", -1)
				n, err := strconv.ParseUint(mOp, 16, 64)
				check(err)
				MEMORY_ADDRESS_FOR_OPERATION = uint16(n)
				code, ok := opcode.XOR["M"]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true
			}
			val, err := strconv.Atoi(operand2)
			if err == nil {
				byteVal := byte(val)
				code, _ := opcode.XOR["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), false
			} else {
				code, ok := opcode.XOR[operand2]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false
			}

		case "ORL":
			if operand1 != "AL" {
				panic("ADDITION CAN ONLY BE DONE WITH AL AS OPERAND 1 " + operand1 + " NOT ALLOWED")
			}
			if operand2[0] == 'M' {
				mOp := strings.Replace(operand2, "M0X", "", -1)
				n, err := strconv.ParseUint(mOp, 16, 64)
				check(err)
				MEMORY_ADDRESS_FOR_OPERATION = uint16(n)
				code, ok := opcode.ORL["M"]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true
			}
			val, err := strconv.Atoi(operand2)
			if err == nil {
				byteVal := byte(val)
				code, _ := opcode.ORL["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), false
			} else {
				code, ok := opcode.ORL[operand2]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false
			}

		case "SUB":
			if operand1 != "AL" {
				panic("ADDITION CAN ONLY BE DONE WITH AL AS OPERAND 1 " + operand1 + " NOT ALLOWED")
			}
			if operand2[0] == 'M' {
				mOp := strings.Replace(operand2, "M0X", "", -1)
				n, err := strconv.ParseUint(mOp, 16, 64)
				check(err)
				MEMORY_ADDRESS_FOR_OPERATION = uint16(n)
				code, ok := opcode.SUB["M"]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true
			}
			val, err := strconv.Atoi(operand2)
			if err == nil {
				byteVal := byte(val)
				code, _ := opcode.SUB["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), false
			} else {
				code, ok := opcode.SUB[operand2]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false
			}

		case "DIV":
			if operand1 != "AL" {
				panic("ADDITION CAN ONLY BE DONE WITH AL AS OPERAND 1 " + operand1 + " NOT ALLOWED")
			}
			if operand2[0] == 'M' {
				mOp := strings.Replace(operand2, "M0X", "", -1)
				n, err := strconv.ParseUint(mOp, 16, 64)
				check(err)
				MEMORY_ADDRESS_FOR_OPERATION = uint16(n)
				code, ok := opcode.DIV["M"]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true
			}
			val, err := strconv.Atoi(operand2)
			if err == nil {
				byteVal := byte(val)
				code, _ := opcode.DIV["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), false
			} else {
				code, ok := opcode.DIV[operand2]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false
			}

		case "STA":
			if operand1[0] == 'M' {
				mOp := strings.Replace(operand1, "M0X", "", -1)
				n, err := strconv.ParseUint(mOp, 16, 64)
				fmt.Println(operand1, operand2)
				check(err)
				MEMORY_ADDRESS_FOR_OPERATION = uint16(n)
				if operand2 != "AX" {
					checkOperand(false, operand2)
				}
				code := opcode.STA
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true
			} else {
				checkOperand(false, operand1)
			}

		default:
			panic(operation + " NOT A VALID OPERATION")

		}
	}

	panic("WORNG COMMAND")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func check8BitOverflow(val int) {
	if val > 255 {
		panic(fmt.Sprintf("%v OVERFLOW", val))
	}
}

func checkOperand(ok bool, operand string) {
	if !ok {
		panic(operand + " NOT A VALID OPERAND")
	}
}
