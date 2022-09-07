package microprocessor

import (
	"bufio"
	"fmt"
	"micp-sim/alu"
	"micp-sim/clock"
	"micp-sim/memory"
	"micp-sim/register"
	"micp-sim/stack"
	"micp-sim/util"
	"os"
	"strconv"
	"strings"
)

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
	for {
		if _, ts := m.Clock.GetState(); ts == clock.State0To1 {
			fmt.Println("a")
		} else {
			fmt.Println("b")
		}
	}
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

  a:= m.Al.GetValue()
  fmt.Println(util.BinaryToDecimal(a[:]))

}

func (m *MicroProcessor) Execute() bool {
	hbitsInst := m.Ir[util.HIGH_BITS].GetValue()
	lbitsInst := m.Ir[util.LOW_BITS].GetValue()
	instB := []byte{}
	for _, v := range hbitsInst {
		if v == 0 {
			continue
		}
		instB = append(instB, v)
	}
	for _, v := range lbitsInst {
		if v == 0 {
			continue
		}

		instB = append(instB, v)
	}

	instruction := string(instB)
	splitInstruction := strings.Split(instruction, " ")
	operation := splitInstruction[0]
	var operand1 string
	var operand2 string
	if len(splitInstruction) == 3 {
		operand1 = splitInstruction[1]
		operand2 = splitInstruction[2]
	}

	switch operation {
	case "END":
		return true
	case "BEGIN":
		fmt.Println("PROGRAM STARTED")
	case "MOV":
		switch operand1 {
		case "B":
			switch operand2 {
			default:
				value, err := strconv.Atoi(operand2)
				check(err)
				m.B.SetLoad()
				m.B.LoadValue(util.DecimalToBinary(value))
			}
		}
	case "ADD":
		switch operand1 {
		case "B":
			switch operand2 {
			default:
				value, err := strconv.Atoi(operand2)
				check(err)
				m.Alu.Temp1.SetLoad()
				m.Alu.Temp1.LoadValue(m.B.GetValue())
				m.Alu.Temp2.SetLoad()
				m.Alu.Temp2.LoadValue(util.DecimalToBinary(value))
        m.Alu.Addition("")

			}
		}

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
		hbitsValue := [8]byte{}

		lastIndex := 0
		for i, x := range instructions {

			lastIndex = i
			if i == 8 {
				break
			}
			hbitsValue[i] = x

		}

		lbitsValue := [8]byte{}
		if len(instructions) > 8 {
			for i := lastIndex; i < len(instructions); i++ {
				lbitsValue[i-8] = instructions[i]
			}
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

	}

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
