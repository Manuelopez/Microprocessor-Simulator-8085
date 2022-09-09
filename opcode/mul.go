package opcode

const (
	MUL_AL_VAL = 0x47
	MUL_AL_M   = 0x48
	MUL_AL_AH  = 0x49
	MUL_AL_B   = 0x4A
	MUL_AL_C   = 0x4B
	MUL_AL_D   = 0x4C
	MUL_AL_E   = 0x4D
	MUL_AL_L   = 0x4E
	MUL_AL_H   = 0x4F
)

var MUL = map[string]byte{
	"VAL": MUL_AL_VAL,
	"M":   MUL_AL_M,
	"AH":  MUL_AL_AH,
	"B":   MUL_AL_B,
	"C":   MUL_AL_C,
	"D":   MUL_AL_D,
	"E":   MUL_AL_E,
	"L":   MUL_AL_L,
	"H":   MUL_AL_H,
}
