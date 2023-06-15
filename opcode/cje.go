package opcode

const (
	CJE_AL_VAL_ADDR = 0xD4
	CJE_AH_VAL_ADDR = 0xD5
	CJE_B_VAL_ADDR  = 0xD6
	CJE_C_VAL_ADDR  = 0xD7
	CJE_D_VAL_ADDR  = 0xD8
	CJE_E_VAL_ADDR  = 0xD9
	CJE_L_VAL_ADDR  = 0xDA
	CJE_H_VAL_ADDR  = 0xDB
)

var CJE = map[string]byte{
	"AL": CJE_AL_VAL_ADDR,
	"AH": CJE_AH_VAL_ADDR,
	"B":  CJE_B_VAL_ADDR,
	"C":  CJE_C_VAL_ADDR,
	"D":  CJE_D_VAL_ADDR,
	"E":  CJE_E_VAL_ADDR,
	"L":  CJE_L_VAL_ADDR,
	"H":  CJE_H_VAL_ADDR,
}
