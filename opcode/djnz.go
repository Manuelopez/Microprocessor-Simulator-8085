package opcode

const (
	DJNZ_AL_ADDR = 0xCA
	DJNZ_AH_ADDR = 0xCB
	DJNZ_B_ADDR  = 0xCC
	DJNZ_C_ADDR  = 0xCD
	DJNZ_D_ADDR  = 0xCE
	DJNZ_E_ADDR  = 0xCF
	DJNZ_L_ADDR  = 0xD0
	DJNZ_H_ADDR  = 0xD1
)

var DJNZ = map[string]byte{
	"AL": DJNZ_AL_ADDR,
	"AH": DJNZ_AH_ADDR,
	"B":  DJNZ_B_ADDR,
	"C":  DJNZ_C_ADDR,
	"D":  DJNZ_D_ADDR,
	"E":  DJNZ_E_ADDR,
	"L":  DJNZ_L_ADDR,
	"H":  DJNZ_H_ADDR,
}
