package opcode

const (
	// ----- AL
	MOV_AL_AH = 0x81

	MOV_AH_AL = 0x82 // AH

	MOV_AL_B = 0x83
	MOV_AL_C = 0x84
	MOV_AL_D = 0x85
	MOV_AL_E = 0x86
	MOV_AL_L = 0x87
	MOV_AL_H = 0x88
	//------ AH
	MOV_AH_B = 0x89
	MOV_AH_C = 0x8A
	MOV_AH_D = 0x8B
	MOV_AH_E = 0x8C
	MOV_AH_L = 0x8D
	MOV_AH_H = 0x8E

	// ------ VAL
	MOV_AL_VAL = 0x8F
	MOV_AH_VAL = 0x90
	MOV_B_VAL  = 0x91
	MOV_C_VAL  = 0x92
	MOV_D_VAL  = 0x93
	MOV_E_VAL  = 0x94
	MOV_H_VAL  = 0x95
	MOV_L_VAL  = 0x96

	// ------- B
	MOV_B_AL = 0x97
	MOV_B_AH = 0x98
	MOV_B_C  = 0x99
	MOV_B_D  = 0x9A
	MOV_B_E  = 0x9B
	MOV_B_H  = 0x9C
	MOV_B_L  = 0x9D
	MOV_B_M  = 0x11

	// -------- C
	MOV_C_AL = 0x9E
	MOV_C_AH = 0x9F
	MOV_C_B  = 0xA0
	MOV_C_D  = 0xA1
	MOV_C_E  = 0xA2
	MOV_C_H  = 0xA3
	MOV_C_L  = 0xA4
	MOV_C_M  = 0x12

	// -------- D
	MOV_D_AL = 0xA5
	MOV_D_AH = 0xA6
	MOV_D_B  = 0xA7
	MOV_D_C  = 0xA8
	MOV_D_E  = 0xA9
	MOV_D_L  = 0xAA
	MOV_D_H  = 0xAB
	MOV_D_M  = 0x13

	// -------- E
	MOV_E_AL = 0xAC
	MOV_E_AH = 0xAD
	MOV_E_B  = 0xAE
	MOV_E_C  = 0xAF
	MOV_E_D  = 0xB0
	MOV_E_H  = 0xB1
	MOV_E_L  = 0xB2
	MOV_E_M  = 0x14

	// --------- H
	MOV_H_AL = 0xB3
	MOV_H_AH = 0xB4
	MOV_H_B  = 0xB5
	MOV_H_C  = 0xB6
	MOV_H_D  = 0xB7
	MOV_H_E  = 0xB8
	MOV_H_L  = 0xB9

	// --------- L
	MOV_L_AL = 0xBA
	MOV_L_AH = 0xBB
	MOV_L_B  = 0xBC
	MOV_L_C  = 0xBD
	MOV_L_D  = 0xBE
	MOV_L_E  = 0xBF
	MOV_L_H  = 0xC1

	// --------- Memory
	MOV_M_AL = 0x0A
	MOV_M_AH = 0x0B
	MOV_M_B  = 0x0C
	MOV_M_C  = 0x0D
	MOV_M_D  = 0x0E
	MOV_M_E  = 0x0F

	// ---------- AX
	MOV_AX_M = 0x10
)

var MOV = map[string]map[string]byte{
	"AX": map[string]byte{
		"M": MOV_AX_M,
	},
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
	"AH": map[string]byte{
		"AL":  MOV_AH_AL,
		"B":   MOV_AH_B,
		"C":   MOV_AH_C,
		"D":   MOV_AH_D,
		"E":   MOV_AH_E,
		"L":   MOV_AH_L,
		"H":   MOV_AH_H,
		"VAL": MOV_AH_VAL,
	},
	"B": map[string]byte{
		"AL":  MOV_B_AL,
		"AH":  MOV_B_AH,
		"C":   MOV_B_C,
		"D":   MOV_B_D,
		"E":   MOV_B_E,
		"L":   MOV_B_L,
		"H":   MOV_B_H,
		"M":   MOV_B_M,
		"VAL": MOV_B_VAL,
	},
	"C": map[string]byte{
		"AL":  MOV_C_AL,
		"AH":  MOV_C_AH,
		"B":   MOV_C_B,
		"D":   MOV_C_D,
		"E":   MOV_C_E,
		"L":   MOV_C_L,
		"H":   MOV_C_H,
		"M":   MOV_C_M,
		"VAL": MOV_C_VAL,
	},
	"D": map[string]byte{
		"AL":  MOV_D_AL,
		"AH":  MOV_D_AH,
		"B":   MOV_D_B,
		"C":   MOV_D_C,
		"E":   MOV_D_E,
		"L":   MOV_D_L,
		"H":   MOV_D_H,
		"M":   MOV_D_M,
		"VAL": MOV_D_VAL,
	},
	"E": map[string]byte{
		"AL":  MOV_E_AL,
		"AH":  MOV_E_AH,
		"B":   MOV_E_B,
		"C":   MOV_E_C,
		"D":   MOV_E_D,
		"L":   MOV_E_L,
		"H":   MOV_E_H,
		"M":   MOV_E_M,
		"VAL": MOV_E_VAL,
	},
	"H": map[string]byte{
		"AL":  MOV_H_AL,
		"AH":  MOV_H_AH,
		"B":   MOV_H_B,
		"C":   MOV_H_C,
		"D":   MOV_H_D,
		"E":   MOV_H_E,
		"L":   MOV_H_L,
		"VAL": MOV_H_VAL,
	},
	"L": map[string]byte{
		"AL":  MOV_L_AL,
		"AH":  MOV_L_AH,
		"B":   MOV_L_B,
		"C":   MOV_L_C,
		"D":   MOV_L_D,
		"E":   MOV_L_E,
		"H":   MOV_L_H,
		"VAL": MOV_L_VAL,
	},
	"M": map[string]byte{
		"AL": MOV_M_AL,
		"AH": MOV_M_AH,
		"B":  MOV_M_B,
		"C":  MOV_M_C,
		"D":  MOV_M_D,
		"E":  MOV_M_E,
	},
}
