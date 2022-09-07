package stack

import (
	"micp-sim/register"
	"micp-sim/util"
)

type Stack struct {
	Mem    [8][8][8]byte
	Sp     register.Register
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
	return Stack{Top: top, Bottom: bottom, Sp: sp}
}

func (s *Stack) Push(value byte) {

	//completeAddress := s.Mar[util.LOW_BITS].GetValue()
	//decimalCompleteAddress := util.BinaryToDecimal(completeAddress[:])
	//hbitsAddress := decimalCompleteAddress >> 3
	//lbitsAddress := decimalCompleteAddress & 0x7
	//valueBits := s.Mbr[util.LOW_BITS].GetValue()

	spB := s.Sp.GetValue()
	sp := util.BinaryToDecimal(spB[:])
	topB := s.Top.GetValue()
	top := util.BinaryToDecimal(topB[:])
	if sp < top {
		sp++
	}

	hbitsAddress := sp >> 3
	lbitsAddress := sp & 0x7
	valueBits := util.DecimalToBinary(int(value))
	for i := 0; i < 8; i++ {
		s.Mem[hbitsAddress][lbitsAddress][i] = valueBits[i]
	}

	s.Sp.SetLoad()
	s.Sp.LoadValue(util.DecimalToBinary(sp))
}

func (s *Stack) Pop() {
	//completeAddress := s.Mar[util.LOW_BITS].GetValue()
	//decimalCompleteAddress := util.BinaryToDecimal(completeAddress[:])
	//hbitsAddress := decimalCompleteAddress >> 3
	//lbitsAddress := decimalCompleteAddress & 0x7
	//valueBits := s.Mem[hbitsAddress][lbitsAddress]

	spB := s.Sp.GetValue()
	sp := util.BinaryToDecimal(spB[:])
	bottomB := s.Bottom.GetValue()
	bottom := util.BinaryToDecimal(bottomB[:])

	hbitsAddress := sp >> 3
	lbitsAddress := sp & 0x7
	valueBits := s.Mem[hbitsAddress][lbitsAddress]

	s.Mbr[util.LOW_BITS].SetLoad()
	s.Mbr[util.LOW_BITS].LoadValue(valueBits)
	if sp > bottom {
		sp--
	}

  s.Sp.SetLoad()
  s.Sp.LoadValue(util.DecimalToBinary(sp))
}
