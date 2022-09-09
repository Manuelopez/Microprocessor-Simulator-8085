package opcode

const (
	ORL_AL_AH = 0x1D
	ORL_AL_B  = 0x1E
	ORL_AL_C  = 0x1F
	ORL_AL_D  = 0x20
	ORL_AL_E  = 0x21
	ORL_AL_L  = 0x22
	ORL_AL_H  = 0x23
)

var ORL = map[string]byte{
	"AH": ORL_AL_AH,
	"B":  ORL_AL_B,
	"C":  ORL_AL_C,
	"D":  ORL_AL_D,
	"E":  ORL_AL_E,
	"L":  ORL_AL_L,
	"H":  ORL_AL_H,
}
