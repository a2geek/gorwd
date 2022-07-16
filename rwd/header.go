package rwd

// Header is the known definition for the RDW starting sequence.
type Header struct {
	Magic      [4]byte
	Value      [3]uint32
	NameLength uint16
	Name       [4]uint16
	Unknown    uint32
}
