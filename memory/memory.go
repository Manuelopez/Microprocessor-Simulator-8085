package memory

import (
	"micp-sim/register"
	"micp-sim/util"
)



type Memory struct {
	Mem [255][255][16]byte
	Mbr *[2]register.Register
	Mar *[2]register.Register
}

func New(mbr, mar *[2]register.Register) Memory {
	return Memory{Mbr: mbr, Mar: mar}
}

func (m *Memory) Write() {
	hbitsAdressBinary := m.Mar[util.HIGH_BITS].GetValue()
	lbitsAdressBinary := m.Mar[util.LOW_BITS].GetValue()
	address1 := util.BinaryToDecimal(hbitsAdressBinary[:])
	address2 := util.BinaryToDecimal(lbitsAdressBinary[:])
	hbitsValue := m.Mbr[util.HIGH_BITS].GetValue()
	lbitsValue := m.Mbr[util.LOW_BITS].GetValue()

	for i := 0; i < 16; i++ {
		if i < 8 {
			m.Mem[address1][address2][i] = hbitsValue[i]
		}else{
      m.Mem[address1][address2][i] = lbitsValue[i-8]
    }
	}

}

func (m *Memory) Read() {
	hbits := m.Mar[util.HIGH_BITS].GetValue()
	lbits := m.Mar[util.LOW_BITS].GetValue()

	address1 := util.BinaryToDecimal(hbits[:])
	address2 := util.BinaryToDecimal(lbits[:])

	data := m.Mem[address1][address2]

	mbrHbits := [8]byte{}
	mbrLbits := [8]byte{}

	for i := 0; i < len(mbrHbits); i++ {
		mbrHbits[i] = data[i]
	}

	for i := 0; i < 8; i++ {
		mbrLbits[i] = data[i+8]
	}

	m.Mbr[util.HIGH_BITS].SetLoad()
	m.Mbr[util.HIGH_BITS].LoadValue(mbrHbits)
	m.Mbr[util.LOW_BITS].SetLoad()
	m.Mbr[util.LOW_BITS].LoadValue(mbrLbits)


}
