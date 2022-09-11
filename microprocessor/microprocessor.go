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

	a := m.B.GetValue()
	fmt.Println(util.BinaryToDecimal(a[:]))

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
  
   

	case opcode.MOV_AL_VAL:
		m.Al.SetLoad()
		m.Al.LoadValue(lbitsInst)
		m.Clock.Wait()

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
