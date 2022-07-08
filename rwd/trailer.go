package rwd

type Trailer struct {
	// Looks to have some text here as well.
	_               [280]byte
	DirectoryLength int32
	_               int32
}
