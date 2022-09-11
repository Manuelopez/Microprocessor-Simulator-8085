package stack

import (
	"errors"
	"micp-sim/register"
	"micp-sim/util"
)

type Stack struct {
	Mem    [8][8]register.Register
	Sp     *register.Register
	Top    register.Register
	Bottom register.Register

	Mar *[2]register.Register
	Mbr *[2]register.Register
}

func New() Stack {
	sp := register.Register{}
	sp.SetLoad()
	sp.LoadValue(util.DecimalToBinary(0))
	top := register.Register{}
	top.SetLoad()
	top.LoadValue(util.DecimalToBinary(63))
	bottom := register.Register{}
	bottom.SetLoad()
	bottom.LoadValue(util.DecimalToBinary(0))
	return Stack{Top: top, Bottom: bottom, Sp: nil}
}

func (s *Stack) Push(value [8]byte) {

	//completeAddress := s.Mar[util.LOW_BITS].GetValue()
	//decimalCompleteAddress := util.BinaryToDecimal(completeAddress[:])
	//hbitsAddress := decimalCompleteAddress >> 3
	//lbitsAddress := decimalCompleteAddress & 0x7
	//valueBits := s.Mbr[util.LOW_BITS].GetValue()
	if s.Sp == nil {
		sp := register.Register{}
		sp.SetLoad()
		sp.LoadValue(util.DecimalToBinary(0))
		s.Sp = &sp

		s.Mem[0][0].SetLoad()
		s.Mem[0][0].LoadValue(value)

	} else {
		spB := s.Sp.GetValue()
		sp := util.BinaryToDecimal(spB[:])
		topB := s.Top.GetValue()
		top := util.BinaryToDecimal(topB[:])
		if sp == top {
			return
		}
		sp++
		hbitsAddress := sp >> 3
		lbitsAddress := sp & 0x7

		s.Mem[hbitsAddress][lbitsAddress].SetLoad()
		s.Mem[hbitsAddress][lbitsAddress].LoadValue(value)

		s.Sp.SetLoad()
		s.Sp.LoadValue(util.DecimalToBinary(sp))

	}
}

func (s *Stack) Pop() ([8]byte, error) {
	//completeAddress := s.Mar[util.LOW_BITS].GetValue()
	//decimalCompleteAddress := util.BinaryToDecimal(completeAddress[:])
	//hbitsAddress := decimalCompleteAddress >> 3
	//lbitsAddress := decimalCompleteAddress & 0x7
	//valueBits := s.Mem[hbitsAddress][lbitsAddress]

	if s.Sp == nil {
		return [8]byte{}, errors.New("Empty Stack")
	}
	spB := s.Sp.GetValue()
	sp := util.BinaryToDecimal(spB[:])

	hbitsAddress := sp >> 3
	lbitsAddress := sp & 0x7

	value := s.Mem[hbitsAddress][lbitsAddress].GetValue()
	sp--
	if sp == -1 {
		s.Sp = nil
	} else {
		s.Sp.SetLoad()
		s.Sp.LoadValue(util.DecimalToBinary(sp))
	}

  return value, nil

}
