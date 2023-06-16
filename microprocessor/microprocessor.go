package microprocessor

import (
	"fmt"
	"micp-sim/alu"
	"micp-sim/clock"
	"micp-sim/memory"
	"micp-sim/opcode"
	"micp-sim/register"
	"micp-sim/stack"
	"micp-sim/util"
	"strconv"
	"strings"
	"syscall/js"
)

var MEMORY_ADDRESS_FOR_OPERATION uint16 = 0
var savedPcVariables = make(map[string]int)

type MicroProcessor struct {
	// AH HIGHT BITS AL LOW BITS
	Al             *register.Register `json:"al"`
	Ah             *register.Register `json:"ah"`
	B              *register.Register `json:"b"`
	C              *register.Register `json:"c"`
	D              *register.Register `json:"d"`
	E              *register.Register `json:"e"`
	L              *register.Register `json:"l"`
	H              *register.Register `json:"h"`
	*memory.Memory `json:"memory"`
	*stack.Stack   `json:"stak"`
	*clock.Clock   `json:"clock"`
	*alu.Alu       `json:"alu"`
	Ir             *[2]register.Register `json:"ir"`
	Mar            *[2]register.Register `json:"mar"`
	Mbr            *[2]register.Register `json:"mbr"`
	Pc             *[2]register.Register `json:"pc"`
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
		Pc:     pc,
	}
}

func (m *MicroProcessor) Start(instructions []string) {
	go m.Clock.TurnOn()

	m.LoadInstructions(instructions)
	m.Pc[0].SetLoad()
	m.Pc[0].LoadValue(util.DecimalToBinary(0))
	m.Pc[1].SetLoad()
	m.Pc[1].LoadValue(util.DecimalToBinary(0))
	for {
		m.ReadInstructon()
		endProgram := m.Execute()

        update(m)
		if endProgram {
			m.Clock.TurnOff()
			break
		}
        

		m.Clock.Wait()

	}

}

func update(m *MicroProcessor) {
	alBinary := m.Al.GetValue()
	al := util.BinaryToDecimal(alBinary[:])
	// Ah
	ahBinary := m.Ah.GetValue()
	ah := util.BinaryToDecimal(ahBinary[:])

	// B              *register.Register `json:"b"`
	bBinary := m.B.GetValue()
	b := util.BinaryToDecimal(bBinary[:])
	// C              *register.Register `json:"c"`
	cBinary := m.C.GetValue()
	c := util.BinaryToDecimal(cBinary[:])
	// D              *register.Register `json:"d"`
	dBinary := m.D.GetValue()
	d := util.BinaryToDecimal(dBinary[:])
	// E              *register.Register `json:"e"`
	eBinary := m.E.GetValue()
	e := util.BinaryToDecimal(eBinary[:])
	// L              *register.Register `json:"l"`
	lBinary := m.L.GetValue()
	l := util.BinaryToDecimal(lBinary[:])
	// H              *register.Register `json:"h"`
	hBinary := m.H.GetValue()
	h := util.BinaryToDecimal(hBinary[:])

    js.Global().Call("updateRegisters", al, ah, b, c, d, e, l, h)
}

func (m *MicroProcessor) Execute() bool {
	hbitsInst := m.Ir[util.HIGH_BITS].GetValue()
	lbitsInst := m.Ir[util.LOW_BITS].GetValue()
	code := util.BinaryToDecimal(hbitsInst[:])

	switch code {
	case opcode.BEGIN:
		fmt.Println("PROGRAM STARTED")

	case opcode.END:
		fmt.Println("PROGRAM ENDED")
		return true

	case opcode.ADD_VAL:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Addition("")
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
		m.IncreasePc()

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
	case opcode.MOV_AH_VAL:
		m.Ah.SetLoad()
		m.Ah.LoadValue(lbitsInst)
	case opcode.MOV_B_VAL:
		m.B.SetLoad()
		m.B.LoadValue(lbitsInst)
	case opcode.MOV_C_VAL:
		m.C.SetLoad()
		m.C.LoadValue(lbitsInst)
	case opcode.MOV_D_VAL:
		m.D.SetLoad()
		m.D.LoadValue(lbitsInst)
	case opcode.MOV_E_VAL:
		m.E.SetLoad()
		m.E.LoadValue(lbitsInst)
	case opcode.MOV_H_VAL:
		m.H.SetLoad()
		m.H.LoadValue(lbitsInst)
	case opcode.MOV_L_VAL:
		m.L.SetLoad()
		m.L.LoadValue(lbitsInst)

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
	case opcode.NOT_AL:
		m.Al.SetLoad()
		m.Al.LoadValue(m.Alu.Not(m.Al.GetValue()))
		//NOT_AH = 0x16
	case opcode.NOT_AH:
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.Alu.Not(m.Ah.GetValue()))

		//NOT_B  = 0x17
	case opcode.NOT_B:
		m.B.SetLoad()
		m.B.LoadValue(m.Alu.Not(m.B.GetValue()))

		//NOT_C  = 0x18
	case opcode.NOT_C:
		m.C.SetLoad()
		m.C.LoadValue(m.Alu.Not(m.C.GetValue()))

		//NOT_D  = 0x19
	case opcode.NOT_D:
		m.D.SetLoad()
		m.D.LoadValue(m.Alu.Not(m.D.GetValue()))

		//NOT_E  = 0x1A
	case opcode.NOT_E:
		m.E.SetLoad()
		m.E.LoadValue(m.Alu.Not(m.E.GetValue()))

	case opcode.CJNE_AL_VAL_ADDR:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Equal()
		c := m.Alu.Comparison.GetValue()
		cv := util.BinaryToDecimal(c[:])
		if cv != 0 {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}

		//CJNE_AH_VAL_ADDR = 0xC3
	case opcode.CJNE_AH_VAL_ADDR:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Ah.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Equal()
		c := m.Alu.Comparison.GetValue()
		cv := util.BinaryToDecimal(c[:])
		if cv != 0 {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}

	//CJNE_B_VAL_ADDR  = 0xC4
	case opcode.CJNE_B_VAL_ADDR:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.B.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Equal()
		c := m.Alu.Comparison.GetValue()
		cv := util.BinaryToDecimal(c[:])
		if cv != 0 {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}

	//CJNE_C_VAL_ADDR  = 0xC5
	case opcode.CJNE_C_VAL_ADDR:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.C.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Equal()
		c := m.Alu.Comparison.GetValue()
		cv := util.BinaryToDecimal(c[:])
		if cv != 0 {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}

	//CJNE_D_VAL_ADDR  = 0xC6
	case opcode.CJNE_D_VAL_ADDR:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Equal()
		c := m.Alu.Comparison.GetValue()
		cv := util.BinaryToDecimal(c[:])
		if cv != 0 {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}

	//CJNE_E_VAL_ADDR  = 0xC7
	case opcode.CJNE_E_VAL_ADDR:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Equal()
		c := m.Alu.Comparison.GetValue()
		cv := util.BinaryToDecimal(c[:])
		if cv != 0 {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}

	//CJNE_L_VAL_ADDR  = 0xC8
	case opcode.CJNE_L_VAL_ADDR:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Equal()
		c := m.Alu.Comparison.GetValue()
		cv := util.BinaryToDecimal(c[:])
		if cv != 0 {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}

	//CJNE_H_VAL_ADDR  = 0xC9
	case opcode.CJNE_H_VAL_ADDR:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Equal()
		c := m.Alu.Comparison.GetValue()
		cv := util.BinaryToDecimal(c[:])
		if cv != 0 {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}

	case opcode.CJE_AL_VAL_ADDR:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Equal()
		c := m.Alu.Comparison.GetValue()
		cv := util.BinaryToDecimal(c[:])
		if cv == 0 {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}
	//	CJE_AH_VAL_ADDR = 0xD5
	case opcode.CJE_AH_VAL_ADDR:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Ah.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Equal()
		c := m.Alu.Comparison.GetValue()
		cv := util.BinaryToDecimal(c[:])
		if cv == 0 {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}
		////CJE_B_VAL_ADDR  = 0xD6
	case opcode.CJE_B_VAL_ADDR:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.B.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Equal()
		c := m.Alu.Comparison.GetValue()
		cv := util.BinaryToDecimal(c[:])
		if cv == 0 {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}
		//CJE_C_VAL_ADDR  = 0xD7
	case opcode.CJE_C_VAL_ADDR:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.C.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Equal()
		c := m.Alu.Comparison.GetValue()
		cv := util.BinaryToDecimal(c[:])
		if cv == 0 {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}
		//CJE_D_VAL_ADDR  = 0xD8
	case opcode.CJE_D_VAL_ADDR:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.D.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Equal()
		c := m.Alu.Comparison.GetValue()
		cv := util.BinaryToDecimal(c[:])
		if cv == 0 {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}
		//CJE_E_VAL_ADDR  = 0xD9
	case opcode.CJE_E_VAL_ADDR:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.E.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Equal()
		c := m.Alu.Comparison.GetValue()
		cv := util.BinaryToDecimal(c[:])
		if cv == 0 {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}
		//CJE_L_VAL_ADDR  = 0xDA
	case opcode.CJE_L_VAL_ADDR:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.L.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Equal()
		c := m.Alu.Comparison.GetValue()
		cv := util.BinaryToDecimal(c[:])
		if cv == 0 {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}
		//CJE_H_VAL_ADDR  = 0xDB
	case opcode.CJE_H_VAL_ADDR:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.H.GetValue())
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
		m.Alu.Equal()
		c := m.Alu.Comparison.GetValue()
		cv := util.BinaryToDecimal(c[:])
		if cv == 0 {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}

	//	DJNZ_AL_ADDR = 0xCA
	case opcode.DJNZ_AL_ADDR:
		m.Al.SetLoad()
		m.Al.LoadValue(m.Alu.Decrement(m.Al.GetValue()))
		if m.Zero == false {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}
		//DJNZ_AH_ADDR = 0xCB
	case opcode.DJNZ_AH_ADDR:
		m.Ah.SetLoad()
		m.Ah.LoadValue(m.Alu.Decrement(m.Ah.GetValue()))
		if m.Zero == false {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}

	//DJNZ_B_ADDR  = 0xCC
	case opcode.DJNZ_B_ADDR:
		m.B.SetLoad()
		m.B.LoadValue(m.Alu.Decrement(m.B.GetValue()))
		if m.Zero == false {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}

	//DJNZ_C_ADDR  = 0xCD
	case opcode.DJNZ_C_ADDR:
		m.C.SetLoad()
		m.C.LoadValue(m.Alu.Decrement(m.C.GetValue()))
		if m.Zero == false {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}

	//DJNZ_D_ADDR  = 0xCE
	case opcode.DJNZ_D_ADDR:
		m.D.SetLoad()
		m.D.LoadValue(m.Alu.Decrement(m.D.GetValue()))
		if m.Zero == false {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}

	//DJNZ_E_ADDR  = 0xCF
	case opcode.DJNZ_E_ADDR:
		m.E.SetLoad()
		m.E.LoadValue(m.Alu.Decrement(m.E.GetValue()))
		if m.Zero == false {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}

	//DJNZ_L_ADDR  = 0xD0
	case opcode.DJNZ_L_ADDR:
		m.L.SetLoad()
		m.L.LoadValue(m.Alu.Decrement(m.L.GetValue()))
		if m.Zero == false {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}

	//DJNZ_H_ADDR  = 0xD1
	case opcode.DJNZ_H_ADDR:
		m.H.SetLoad()
		m.H.LoadValue(m.Alu.Decrement(m.H.GetValue()))
		if m.Zero == false {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}

	case opcode.JC:
		if m.Alu.Carry == true {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())

		} else {
			m.IncreasePc()
		}

	case opcode.JNC:
		if m.Alu.Carry == false {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())

		} else {
			m.IncreasePc()
		}
	case opcode.JZ:
		a := m.Al.GetValue()
		al := util.BinaryToDecimal(a[:])
		if al == 0 {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}
	case opcode.JNZ:
		a := m.Al.GetValue()
		al := util.BinaryToDecimal(a[:])
		if al != 0 {
			m.LoadMarWithPC()
			m.Memory.Read()
			m.Pc[util.HIGH_BITS].SetLoad()
			m.Pc[util.HIGH_BITS].LoadValue(m.Mbr[util.HIGH_BITS].GetValue())
			m.Pc[util.LOW_BITS].SetLoad()
			m.Pc[util.LOW_BITS].LoadValue(m.Mbr[util.LOW_BITS].GetValue())
		} else {
			m.IncreasePc()
		}

	//	PUSH_VAL = 0xDE
	case opcode.PUSH_VAL:
		m.Push(lbitsInst)

		//PUSH_AL  = 0xDF
	case opcode.PUSH_AL:
		m.Push(m.Al.GetValue())
		//PUSH_AH  = 0xE0
	case opcode.PUSH_AH:
		m.Push(m.Ah.GetValue())
		//PUSH_B   = 0xE1
	case opcode.PUSH_B:
		m.Push(m.B.GetValue())
		//PUSH_C   = 0xE2
	case opcode.PUSH_C:
		m.Push(m.C.GetValue())
		//PUSH_D   = 0xE3
	case opcode.PUSH_D:
		m.Push(m.D.GetValue())
		//PUSH_E   = 0xE4
	case opcode.PUSH_E:
		m.Push(m.E.GetValue())
		//PUSH_L   = 0xE5
	case opcode.PUSH_L:
		m.Push(m.L.GetValue())
		//PUSH_H   = 0xE6
	case opcode.PUSH_H:
		m.Push(m.H.GetValue())

		//	POP_AL = 0xE7
	case opcode.POP_AL:
		m.Al.SetLoad()
		val, err := m.Stack.Pop()
		if err == nil {
			m.Al.LoadValue(val)
		}
	//POP_AH = 0xE8
	case opcode.POP_AH:
		m.Ah.SetLoad()
		val, err := m.Stack.Pop()
		if err == nil {
			m.Ah.LoadValue(val)
		}

	//POP_B  = 0xE9
	case opcode.POP_B:
		m.B.SetLoad()
		val, err := m.Stack.Pop()
		if err == nil {
			m.B.LoadValue(val)
		}

	//POP_C  = 0xEA
	case opcode.POP_C:
		m.C.SetLoad()
		val, err := m.Stack.Pop()
		if err == nil {
			m.C.LoadValue(val)
		}

	//POP_D  = 0xEB
	case opcode.POP_D:
		m.D.SetLoad()
		val, err := m.Stack.Pop()
		if err == nil {
			m.D.LoadValue(val)
		}

	//POP_E  = 0xEC
	case opcode.POP_E:
		m.E.SetLoad()
		val, err := m.Stack.Pop()
		if err == nil {
			m.E.LoadValue(val)
		}

	//POP_L  = 0xED
	case opcode.POP_L:
		m.L.SetLoad()
		val, err := m.Stack.Pop()
		if err == nil {
			m.L.LoadValue(val)
		}

	//POP_H  = 0xEE
	case opcode.POP_H:
		m.H.SetLoad()
		val, err := m.Stack.Pop()
		if err == nil {
			m.H.LoadValue(val)
		}

	//	CLR_CA = 0xF0
	case opcode.CLR_CA:
		m.Alu.Carry = false
		//CLR_AL = 0xF7
	case opcode.CLR_AL:
		m.Al.SetLoad()
		m.Al.LoadValue(util.DecimalToBinary(0))
		//CLR_AH = 0xF8
	case opcode.CLR_AH:
		m.Ah.SetLoad()
		m.Ah.LoadValue(util.DecimalToBinary(0))

		//CLR_B  = 0xF9
	case opcode.CLR_B:
		m.B.SetLoad()
		m.B.LoadValue(util.DecimalToBinary(0))

		//CLR_C  = 0xFA
	case opcode.CLR_C:
		m.C.SetLoad()
		m.C.LoadValue(util.DecimalToBinary(0))

		//CLR_D  = 0xFB
	case opcode.CLR_D:
		m.D.SetLoad()
		m.D.LoadValue(util.DecimalToBinary(0))

		//CLR_E  = 0xFC
	case opcode.CLR_E:
		m.E.SetLoad()
		m.E.LoadValue(util.DecimalToBinary(0))

		//CLR_H  = 0xFD
	case opcode.CLR_H:
		m.H.SetLoad()
		m.H.LoadValue(util.DecimalToBinary(0))

		//CLR_L  = 0xFE
	case opcode.CLR_L:
		m.L.SetLoad()
		m.L.LoadValue(util.DecimalToBinary(0))

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

func (m *MicroProcessor) LoadInstructions(instructions []string) {
	pc := 0
	for _, inst := range instructions {
		text := inst
		instructions := []byte(text)
		hbitsValue, lbitsValue, isMemoryOp, savePC, variable := Assembler(string(instructions))

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
		if savePC {
			// fmt.Println(pc)
			// fmt.Println(variable)
			savedPcVariables[variable] = pc
		}
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

func Assembler(instructions string) ([8]byte, [8]byte, bool, bool, string) {
	instructions = strings.ToUpper(instructions)
	instructions = strings.TrimSpace(instructions)
	splitInstructions := strings.Split(instructions, " ")
	allOperations := [...]string{
		"BEGIN",
		"END",
		"INC",
		"DEC",
		"ADD",
		"SUB",
		"MUL",
		"DIV",
		"ORL",
		"AND",
		"XOR",
		"NOT",
		"JC",
		"JNC",
		"JZ",
		"JNZ",
		"PUSH",
		"POP",
		"CLR",
		"MOV",
		"STA",
		"DJNZ",
		"CJNE",
		"CJE",
	}
	savePC := false
	variable := ""
	if len(splitInstructions) > 2 {
		found := false
		for _, op := range allOperations {
			if op == splitInstructions[0] {
				found = true
			}
		}

		if !found {
			// then is variable to save operation code
			savePC = true
			variable = splitInstructions[0]
			splitInstructions = splitInstructions[1:]
		}
	}

	if len(splitInstructions) == 1 {
		switch splitInstructions[0] {
		case "BEGIN":
			return util.DecimalToBinary(opcode.BEGIN), util.DecimalToBinary(opcode.BEGIN), false, savePC, variable
		case "END":
			return util.DecimalToBinary(opcode.END), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
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
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
		case "DEC":
			code, ok := opcode.DEC[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable

		case "ADD":
			code, ok := opcode.ADD[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
		case "SUB":
			code, ok := opcode.SUB[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
		case "MUL":
			code, ok := opcode.MUL[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
		case "DIV":
			code, ok := opcode.DIV[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
		case "ORL":
			code, ok := opcode.ORL[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
		case "AND":
			code, ok := opcode.AND[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
		case "XOR":
			code, ok := opcode.XOR[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
		case "NOT":
			code, ok := opcode.NOT[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
		case "JC":
			mOp, ok := savedPcVariables[operand1]
			checkOperand(ok, operand1)
			MEMORY_ADDRESS_FOR_OPERATION = uint16(mOp)
			code := opcode.JC
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true, savePC, variable
		case "JNC":
			mOp, ok := savedPcVariables[operand1]
			checkOperand(ok, operand1)
			MEMORY_ADDRESS_FOR_OPERATION = uint16(mOp)
			code := opcode.JNC
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true, savePC, variable
		case "JZ":
			if operand1[0] == 'M' {
				mOp := strings.Replace(operand1, "M0X", "", -1)
				n, err := strconv.ParseUint(mOp, 16, 64)
				check(err)
				MEMORY_ADDRESS_FOR_OPERATION = uint16(n)
				code := opcode.JZ
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true, savePC, variable
			} else {
				checkOperand(false, operand1)
			}
		case "JNZ":
			if operand1[0] == 'M' {
				mOp := strings.Replace(operand1, "M0X", "", -1)
				n, err := strconv.ParseUint(mOp, 16, 64)
				check(err)
				MEMORY_ADDRESS_FOR_OPERATION = uint16(n)
				code := opcode.JNZ
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true, savePC, variable
			} else {
				checkOperand(false, operand1)
			}

		case "PUSH":
			val, err := strconv.Atoi(operand1)
			if err == nil {
				byteVal := byte(val)
				code, _ := opcode.PUSH["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), false, savePC, variable
			} else {
				code, ok := opcode.PUSH[operand1]
				checkOperand(ok, operand1)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
			}

		case "POP":
			val, err := strconv.Atoi(operand1)
			if err == nil {
				byteVal := byte(val)
				code, _ := opcode.POP["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), false, savePC, variable
			} else {
				code, ok := opcode.POP[operand1]
				checkOperand(ok, operand1)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
			}

		case "CLR":
			code, ok := opcode.CLR[operand1]
			checkOperand(ok, operand1)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable

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
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true, savePC, variable
			}
			val, err := strconv.Atoi(operand2)
			if err == nil {
				byteVal := byte(val)
				code, _ := opcode.ADD["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), false, savePC, variable
			} else {
				code, ok := opcode.ADD[operand2]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
			}
		case "MOV":
			if operand1[0] == 'M' {
				mOp := strings.Replace(operand1, "M0X", "", -1)
				n, err := strconv.ParseUint(mOp, 16, 64)
				check(err)
				MEMORY_ADDRESS_FOR_OPERATION = uint16(n)
				code, ok := opcode.MOV["M"][operand2]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true, savePC, variable
			}
			mapedOp1, ok := opcode.MOV[operand1]
			checkOperand(ok, operand1)
			val, err := strconv.Atoi(operand2)
			if err == nil {
				byteVal := byte(val)
				code, _ := mapedOp1["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), false, savePC, variable
			} else {
				if operand2[0] == 'M' {
					mOp := strings.Replace(operand2, "M0X", "", -1)
					n, err := strconv.ParseUint(mOp, 16, 64)
					check(err)
					MEMORY_ADDRESS_FOR_OPERATION = uint16(n)
					code, ok := mapedOp1["M"]
					checkOperand(ok, operand2)
					return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true, savePC, variable
				}
				code, ok := mapedOp1[operand2]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
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
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true, savePC, variable
			}
			val, err := strconv.Atoi(operand2)
			if err == nil {
				byteVal := byte(val)
				code, _ := opcode.MUL["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), false, savePC, variable
			} else {
				code, ok := opcode.MUL[operand2]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
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
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true, savePC, variable
			}
			val, err := strconv.Atoi(operand2)
			if err == nil {
				byteVal := byte(val)
				code, _ := opcode.AND["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), false, savePC, variable
			} else {
				code, ok := opcode.AND[operand2]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
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
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true, savePC, variable
			}
			val, err := strconv.Atoi(operand2)
			if err == nil {
				byteVal := byte(val)
				code, _ := opcode.XOR["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), false, savePC, variable
			} else {
				code, ok := opcode.XOR[operand2]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
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
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true, savePC, variable
			}
			val, err := strconv.Atoi(operand2)
			if err == nil {
				byteVal := byte(val)
				code, _ := opcode.ORL["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), false, savePC, variable
			} else {
				code, ok := opcode.ORL[operand2]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
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
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true, savePC, variable
			}
			val, err := strconv.Atoi(operand2)
			if err == nil {
				byteVal := byte(val)
				code, _ := opcode.SUB["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), false, savePC, variable
			} else {
				code, ok := opcode.SUB[operand2]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
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
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true, savePC, variable
			}
			val, err := strconv.Atoi(operand2)
			if err == nil {
				byteVal := byte(val)
				code, _ := opcode.DIV["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), false, savePC, variable
			} else {
				code, ok := opcode.DIV[operand2]
				checkOperand(ok, operand2)
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), false, savePC, variable
			}

		case "STA":
			if operand1[0] == 'M' {
				mOp := strings.Replace(operand1, "M0X", "", -1)
				n, err := strconv.ParseUint(mOp, 16, 64)
				check(err)
				MEMORY_ADDRESS_FOR_OPERATION = uint16(n)
				if operand2 != "AX" {
					checkOperand(false, operand2)
				}
				code := opcode.STA
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true, savePC, variable
			} else {
				checkOperand(false, operand1)
			}
		case "DJNZ":
			code, ok := opcode.DJNZ[operand1]
			checkOperand(ok, operand1)
			mOp, ok := savedPcVariables[operand2]
			checkOperand(ok, operand2)
			MEMORY_ADDRESS_FOR_OPERATION = uint16(mOp)

			return util.DecimalToBinary(int(code)), util.DecimalToBinary(opcode.NOTHING), true, savePC, variable

		default:
			panic(operation + " NOT A VALID OPERATION")

		}
	} else if len(splitInstructions) == 4 {
		operation := splitInstructions[0]
		operand1 := splitInstructions[1]
		operand2 := splitInstructions[2]
		operand3 := splitInstructions[3]
		switch operation {
		case "CJNE":
			code, ok := opcode.CJNE[operand1]
			checkOperand(ok, operand1)
			val, err := strconv.Atoi(operand2)
			check(err)
			byteVal := byte(val)
			mOp, ok := savedPcVariables[operand3]
			checkOperand(ok, operand3)
			MEMORY_ADDRESS_FOR_OPERATION = uint16(mOp)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), true, savePC, variable
		case "CJE":
			code, ok := opcode.CJE[operand1]
			checkOperand(ok, operand1)
			val, err := strconv.Atoi(operand2)
			check(err)
			byteVal := byte(val)
			mOp, ok := savedPcVariables[operand3]
			checkOperand(ok, operand3)
			MEMORY_ADDRESS_FOR_OPERATION = uint16(mOp)
			return util.DecimalToBinary(int(code)), util.DecimalToBinary(int(byteVal)), true, savePC, variable

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
