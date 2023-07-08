package internal

import (
	"unicode"

	"github.com/pkg/errors"
)

const EnglishLetterCount = 26

const (
	ErrLowerCaseLetters       = "word must contain only lowercase letters"
	ErrHashMustBe26Characters = "hash must be 26 characters long"
	ErrWordMustNotBeEmpty     = "word must not be empty"
)

var ErrorWordMustNotBeEmpty = errors.New(ErrWordMustNotBeEmpty)

type AnagramLetterMap [EnglishLetterCount]uint8

func NewAnagramLetterMap(word string) (AnagramLetterMap, error) {
	var w AnagramLetterMap
	if len(word) == 0 {
		return w, errors.New(ErrWordMustNotBeEmpty)
	}
	for _, c := range word {
		if unicode.IsSpace(c) {
			continue
		}
		if c < 'a' || c > 'z' {
			return AnagramLetterMap{}, ErrorWordMustNotBeEmpty
		}
		w[c-'a']++
	}
	return w, nil
}

func NewAnagramLetterMapFromHash(hash string) (AnagramLetterMap, error) {
	if len(hash) != EnglishLetterCount {
		return AnagramLetterMap{}, errors.New(ErrHashMustBe26Characters)
	}
	var w AnagramLetterMap
	for i, c := range hash {
		w[i] = uint8(c)
	}
	return w, nil
}

func (w *AnagramLetterMap) AnagramHash() string {
	return string(w[:])
}

func (w *AnagramLetterMap) IsAnagram(other AnagramLetterMap) bool {
	for i := 0; i < 26; i++ {
		if w[i] != other[i] {
			return false
		}
	}
	return true
}
