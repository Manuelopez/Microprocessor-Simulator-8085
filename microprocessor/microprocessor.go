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

	a := m.Al.GetValue()
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
	case opcode.MOV_AL_VAL:
		m.Al.SetLoad()
		m.Al.LoadValue(lbitsInst)
    m.Clock.Wait()
	case opcode.ADD_VAL:
		m.Alu.Temp1.SetLoad()
		m.Alu.Temp1.LoadValue(m.Al.GetValue())
    m.Clock.Wait()
		m.Alu.Temp2.SetLoad()
		m.Alu.Temp2.LoadValue(lbitsInst)
    m.Clock.Wait()
		m.Alu.Addition("")
    m.Clock.Wait()
	default:
		fmt.Println("OPERATION NOT IMPLEMENTED")
	}

	return false
}

func (m *MicroProcessor) ReadInstructon() {
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
	m.Mar[util.HIGH_BITS].SetLoad()
	m.Mar[util.HIGH_BITS].LoadValue(util.DecimalToBinary(hbits))
	m.Mar[util.LOW_BITS].SetLoad()
	m.Mar[util.LOW_BITS].LoadValue(util.DecimalToBinary(lbits))

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

	pc++
	hbitsPcD := pc >> 8
	lbitsPcD := pc & 0xFF

	hbitsPcBN, lbitsPcBN := util.DecimalToBinary(hbitsPcD), util.DecimalToBinary(lbitsPcD)
	m.Pc[util.HIGH_BITS].SetLoad()
	m.Pc[util.HIGH_BITS].LoadValue(hbitsPcBN)
	m.Pc[util.LOW_BITS].SetLoad()
	m.Pc[util.LOW_BITS].LoadValue(lbitsPcBN)

	m.Ir[util.HIGH_BITS].SetLoad()
	m.Ir[util.HIGH_BITS].LoadValue(hbitsMbr)
	m.Ir[util.LOW_BITS].SetLoad()
	m.Ir[util.LOW_BITS].LoadValue(lbitsMbr)
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

		if isMemoryOp {
			//TODO MEMORY OPERATION
		}

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
			/* TODO MEMORY AS OPERAND1 */
			mapedOp1, ok := opcode.MOV[operand1]
			checkOperand(ok, operand1)
			val, err := strconv.Atoi(operand2)
			if err == nil {
				check8BitOverflow(val)
				code, _ := mapedOp1["VAL"]
				return util.DecimalToBinary(int(code)), util.DecimalToBinary(val), false
			} else {
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
		panic(fmt.Sprintf("%v OVERFLOW"))
	}
}

func checkOperand(ok bool, operand string) {
	if !ok {
		panic(operand + " NOT A VALID OPERAND")
	}
}
