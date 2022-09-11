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
  a := []byte{}
  ah := m.Ah.GetValue()
  al := m.Al.GetValue()
  a = append(a, ah[:]...)
  a = append(a, al[:]...)
	fmt.Println(util.BinaryToDecimal(a))

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

	case opcode.MOV_AL_AH:
		m.Al.SetLoad()
		m.Al.LoadValue(m.Ah.GetValue())

	case opcode.MOV_AL_B:
		m.Al.SetLoad()
		m.Al.LoadValue(m.B.GetValue())
		//	MOV_AL_C = 0x84
	case opcode.MOV_AL_C:
		m.Al.SetLoad()
		m.Al.LoadValue(m.C.GetValue())
		//	MOV_AL_D = 0x85
	case opcode.MOV_AL_D:
		m.Al.SetLoad()
		m.Al.LoadValue(m.D.GetValue())
		//	MOV_AL_E = 0x86
	case opcode.MOV_AL_E:
		m.Al.SetLoad()
		m.Al.LoadValue(m.E.GetValue())
		//	MOV_AL_L = 0x87
	case opcode.MOV_AL_L:
		m.Al.SetLoad()
		m.Al.LoadValue(m.L.GetValue())
		//	MOV_AL_H = 0x88
	case opcode.MOV_AL_H:
		m.Al.SetLoad()
		m.Al.LoadValue(m.H.GetValue())
		//------ AH

		//MOV_AH_AL = 0x82
	case opcode.MOV_AH_AL:
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.Al.GetValue())
		//MOV_AH_B = 0x89
	case opcode.MOV_AH_B:
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.B.GetValue())
		//	MOV_AH_C = 0x8A
	case opcode.MOV_AH_C:
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.C.GetValue())
	//	MOV_AH_D = 0x8B
	case opcode.MOV_AH_D:
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.D.GetValue())
		//	MOV_AH_E = 0x8C
	case opcode.MOV_AH_E:
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.E.GetValue())
		//	MOV_AH_L = 0x8D
	case opcode.MOV_AH_L:
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.L.GetValue())
		//	MOV_AH_H = 0x8E
	case opcode.MOV_AH_H:
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.H.GetValue())

		// ------ VAL
		//	MOV_AL_VAL = 0x8F
	case opcode.MOV_AL_VAL:
		m.Al.SetLoad()
		m.Al.LoadValue(lbitsInst)
		m.Clock.Wait()

		//MOV_AH_VAL = 0x90
	case opcode.MOV_AH_VAL:
		m.Ah.SetLoad()
		m.Ah.LoadValue(lbitsInst)
		m.Clock.Wait()

		//	MOV_B_VAL  = 0x91
	case opcode.MOV_B_VAL:
		m.B.SetLoad()
		m.B.LoadValue(lbitsInst)
		m.Clock.Wait()

		//MOV_C_VAL  = 0x92
	case opcode.MOV_C_VAL:
		m.C.SetLoad()
		m.C.LoadValue(lbitsInst)
		m.Clock.Wait()

		//MOV_D_VAL  = 0x93
	case opcode.MOV_D_VAL:
		m.D.SetLoad()
		m.D.LoadValue(lbitsInst)
		m.Clock.Wait()

		//MOV_E_VAL  = 0x94
	case opcode.MOV_E_VAL:
		m.E.SetLoad()
		m.E.LoadValue(lbitsInst)
		m.Clock.Wait()

		//MOV_H_VAL  = 0x95
	case opcode.MOV_H_VAL:
		m.H.SetLoad()
		m.H.LoadValue(lbitsInst)
		m.Clock.Wait()

		//MOV_L_VAL  = 0x96
	case opcode.MOV_L_VAL:
		m.L.SetLoad()
		m.L.LoadValue(lbitsInst)
		m.Clock.Wait()

		// ------- B
		//MOV_B_AL = 0x97
	case opcode.MOV_B_AL:
		m.B.SetLoad()
		m.B.LoadValue(m.Al.GetValue())

		//MOV_B_AH = 0x98
	case opcode.MOV_B_AH:
		m.B.SetLoad()
		m.B.LoadValue(m.Ah.GetValue())

		//MOV_B_C  = 0x99
	case opcode.MOV_B_C:
		m.B.SetLoad()
		m.B.LoadValue(m.C.GetValue())

		//MOV_B_D  = 0x9A
	case opcode.MOV_B_D:
		m.B.SetLoad()
		m.B.LoadValue(m.D.GetValue())

		//MOV_B_E  = 0x9B
	case opcode.MOV_B_E:
		m.B.SetLoad()
		m.B.LoadValue(m.E.GetValue())

		//MOV_B_H  = 0x9C
	case opcode.MOV_B_H:
		m.B.SetLoad()
		m.B.LoadValue(m.H.GetValue())

		//MOV_B_L  = 0x9D
	case opcode.MOV_B_L:
		m.B.SetLoad()
		m.B.LoadValue(m.L.GetValue())

		//MOV_B_M  = 0x11
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

		// -------- C
		//MOV_C_AL = 0x9E
	case opcode.MOV_C_AL:
		m.C.SetLoad()
		m.C.LoadValue(m.Al.GetValue())
	//MOV_C_AH = 0x9F
	case opcode.MOV_C_AH:
		m.C.SetLoad()
		m.C.LoadValue(m.Ah.GetValue())

		//MOV_C_B  = 0xA0
	case opcode.MOV_C_B:
		m.C.SetLoad()
		m.C.LoadValue(m.B.GetValue())

		//MOV_C_D  = 0xA1
	case opcode.MOV_C_D:
		m.C.SetLoad()
		m.C.LoadValue(m.D.GetValue())

		//MOV_C_E  = 0xA2
	case opcode.MOV_C_E:
		m.C.SetLoad()
		m.C.LoadValue(m.E.GetValue())

		//MOV_C_H  = 0xA3
	case opcode.MOV_C_H:
		m.C.SetLoad()
		m.C.LoadValue(m.H.GetValue())

		//MOV_C_L  = 0xA4
	case opcode.MOV_C_L:
		m.C.SetLoad()
		m.C.LoadValue(m.L.GetValue())

		////MOV_C_M  = 0x12
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

		// -------- D
		//MOV_D_AL = 0xA5
	case opcode.MOV_D_AL:
		m.D.SetLoad()
		m.D.LoadValue(m.Al.GetValue())
		//MOV_D_AH = 0xA6
	case opcode.MOV_D_AH:
		m.D.SetLoad()
		m.D.LoadValue(m.Ah.GetValue())

		//MOV_D_B  = 0xA7
	case opcode.MOV_D_B:
		m.D.SetLoad()
		m.D.LoadValue(m.B.GetValue())

	//MOV_D_C  = 0xA8
	case opcode.MOV_D_C:
		m.D.SetLoad()
		m.D.LoadValue(m.C.GetValue())

		//MOV_D_E  = 0xA9
	case opcode.MOV_D_E:
		m.D.SetLoad()
		m.D.LoadValue(m.E.GetValue())

		//MOV_D_L  = 0xAA
	case opcode.MOV_D_L:
		m.D.SetLoad()
		m.D.LoadValue(m.L.GetValue())

		//MOV_D_H  = 0xAB
	case opcode.MOV_D_H:
		m.D.SetLoad()
		m.D.LoadValue(m.H.GetValue())

		//MOV_D_M  = 0x13
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

		// -------- E
		//MOV_E_AL = 0xAC
	case opcode.MOV_E_AL:
		m.E.SetLoad()
		m.E.LoadValue(m.Al.GetValue())
		//MOV_E_AH = 0xAD
	case opcode.MOV_E_AH:
		m.E.SetLoad()
		m.E.LoadValue(m.Ah.GetValue())

		//MOV_E_B  = 0xAE
	case opcode.MOV_E_B:
		m.E.SetLoad()
		m.E.LoadValue(m.B.GetValue())

		//MOV_E_C  = 0xAF
	case opcode.MOV_E_C:
		m.E.SetLoad()
		m.E.LoadValue(m.C.GetValue())

		//MOV_E_D  = 0xB0
	case opcode.MOV_E_D:
		m.E.SetLoad()
		m.E.LoadValue(m.D.GetValue())

		//MOV_E_H  = 0xB1
	case opcode.MOV_E_H:
		m.E.SetLoad()
		m.E.LoadValue(m.H.GetValue())

		//MOV_E_L  = 0xB2
	case opcode.MOV_E_L:
		m.E.SetLoad()
		m.E.LoadValue(m.L.GetValue())

		//MOV_E_M  = 0x14
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
		//MOV_H_AL = 0xB3
	case opcode.MOV_H_AL:
		m.H.SetLoad()
		m.H.LoadValue(m.Al.GetValue())
		//MOV_H_AH = 0xB4
	case opcode.MOV_H_AH:
		m.H.SetLoad()
		m.H.LoadValue(m.Ah.GetValue())

		//MOV_H_B  = 0xB5
	case opcode.MOV_H_B:
		m.H.SetLoad()
		m.H.LoadValue(m.B.GetValue())

		//MOV_H_C  = 0xB6
	case opcode.MOV_H_C:
		m.H.SetLoad()
		m.H.LoadValue(m.C.GetValue())

		//MOV_H_D  = 0xB7
	case opcode.MOV_H_D:
		m.H.SetLoad()
		m.H.LoadValue(m.D.GetValue())

		//MOV_H_E  = 0xB8
	case opcode.MOV_H_E:
		m.H.SetLoad()
		m.H.LoadValue(m.E.GetValue())

		//MOV_H_L  = 0xB9
	case opcode.MOV_H_L:
		m.H.SetLoad()
		m.H.LoadValue(m.L.GetValue())

	// --------- L
	//MOV_L_AL = 0xBA
	case opcode.MOV_L_AL:
		m.L.SetLoad()
		m.L.LoadValue(m.Al.GetValue())

	//MOV_L_AH = 0xBB
	case opcode.MOV_L_AH:
		m.L.SetLoad()
		m.L.LoadValue(m.Ah.GetValue())

	//MOV_L_B  = 0xBC
	case opcode.MOV_L_B:
		m.L.SetLoad()
		m.L.LoadValue(m.B.GetValue())

	//MOV_L_C  = 0xBD
	case opcode.MOV_L_C:
		m.L.SetLoad()
		m.L.LoadValue(m.C.GetValue())

	//MOV_L_D  = 0xBE
	case opcode.MOV_L_D:
		m.L.SetLoad()
		m.L.LoadValue(m.D.GetValue())

		//MOV_L_E  = 0xBF
	case opcode.MOV_L_E:
		m.L.SetLoad()
		m.L.LoadValue(m.E.GetValue())

		//MOV_L_H  = 0xC1
	case opcode.MOV_L_H:
		m.L.SetLoad()
		m.L.LoadValue(m.H.GetValue())

		// --------- Memory
		//MOV_M_AL = 0x0A
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

		//MOV_M_AH = 0x0B
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

		//MOV_M_B  = 0x0C
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

		//MOV_M_C  = 0x0D
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

		//MOV_M_D  = 0x0E
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
		//MOV_M_E  = 0x0F

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

		// ---------- AX
		//MOV_AX_M = 0x10
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
		}
	} else if len(splitInstructions) == 2 {
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
			val, err := strconv.Atoi(operand2)
			if err == nil {
				check8BitOverflow(val)
				code, _ := opcode.ADD["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(val), false
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
				check8BitOverflow(val)
				code, _ := mapedOp1["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(val), false
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
