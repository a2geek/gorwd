package rwd

// Header is the known definition for the RDW starting sequence.
type Header struct {
	Magic [4]byte
	// 3 int32 with 2,3,2
	// 1 int16 with length?
	// "K2" in 16-bit values followed by 4 zeros
	// 1 int32 of calculated values (varies by file)?
	Data [26]byte
}
