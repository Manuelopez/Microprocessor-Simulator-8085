package opcode

const (
	MOV_AL_AH  = 0x81
	MOV_AH_AL  = 0x82
	MOV_AL_B   = 0x83
	MOV_AL_C   = 0x84
	MOV_AL_D   = 0x85
	MOV_AL_E   = 0x86
	MOV_AL_L   = 0x87
	MOV_AL_H   = 0x88
	MOV_AL_VAL = 0x8F

	MOV_AH_B = 0x89
)

var MOV = map[string]map[string]byte{
	"AL": map[string]byte{
		"AH":  MOV_AH_AL,
		"B":   MOV_AL_B,
		"C":   MOV_AL_C,
		"D":   MOV_AL_D,
		"E":   MOV_AL_E,
		"L":   MOV_AL_L,
		"H":   MOV_AL_H,
		"VAL": MOV_AL_VAL,
	},
}
