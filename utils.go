package fatsecret

// padLeft left-pads the string with pad up to len runes
// len may be exceeded if
func padLeft(str string, length int, pad string) string {
	return times(pad, length-len(str)) + str
}

// padRight right-pads the string with pad up to len runes
func padRight(str string, length int, pad string) string {
	return str + times(pad, length-len(str))
}

func times(str string, n int) (out string) {
	for i := 0; i < n; i++ {
		out += str
	}
	return
}
