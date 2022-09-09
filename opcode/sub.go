package opcode

const (
	SUB_AL_VAL = 0x3E
	SUB_AL_M   = 0x3F
	SUB_AL_AH  = 0x40
	SUB_AL_B   = 0x41
	SUB_AL_C   = 0x42
	SUB_AL_D   = 0x43
	SUB_AL_E   = 0x44
	SUB_AL_L   = 0x45
	SUB_AL_H   = 0x46
)

var SUB = map[string]byte{
	"VAL": SUB_AL_VAL,
	"M":   SUB_AL_M,
	"AH":  SUB_AL_AH,
	"B":   SUB_AL_B,
	"C":   SUB_AL_C,
	"D":   SUB_AL_D,
	"E":   SUB_AL_E,
	"L":   SUB_AL_L,
	"H":   SUB_AL_H,
}
