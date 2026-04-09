package editor

// IsBinary checks the first 8KB for null bytes — a strong indicator of binary content.
func IsBinary(data []byte) bool {
	size := 8192
	if len(data) < size {
		size = len(data)
	}
	for _, b := range data[:size] {
		if b == 0 {
			return true
		}
	}
	return false
}
