package rwd

// Section is a repeated component of the Trailer
type Section struct {
	// Zero terminated string
	NameZ [32]uint16
	// Offset (in bytes) into file where this section begins
	Offset   uint32
	Unknown1 uint32
	// Length (in bytes) of this section.
	Length   uint32
	Unknown3 uint32
	Unknown4 uint32
	Unknown5 uint32
	// AlternateLength appears to be a duplicate of the Length.
	AlternateLength uint32
	Unknown7        uint32
}

// Trailer represents the end of the RDW archive file.
type Trailer struct {
	Header Section
	Files  Section
	Footer Section
}
