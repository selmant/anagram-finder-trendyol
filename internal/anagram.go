package internal

type AnagramLetterMap [26]uint8

func NewAnagramLetterMap(word string) AnagramLetterMap {
	var w AnagramLetterMap
	for _, c := range word {
		if c < 'a' || c > 'z' {
			continue
		}
		w[c-'a']++
	}
	return w
}

func (w *AnagramLetterMap) AnagramHash() string {
	return string(w[:])
}
