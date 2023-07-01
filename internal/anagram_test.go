package internal_test

import (
	"testing"

	"github.com/selmant/anagram-finder-trendyol/internal"
	"github.com/stretchr/testify/assert"
)

func TestAnagramMapCreated(t *testing.T) {
	assert := assert.New(t)

	word := "test"
	am, err := internal.NewAnagramLetterMap(word)
	if err != nil {
		t.Errorf("AnagramLetterMap creation failed: %v", err)
	}
	expected := internal.AnagramLetterMap{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 0, 0, 0, 0, 0, 0}
	assert.Equal(expected, am)
}

func TestAnagramMapCreatedWithSpace(t *testing.T) {
	assert := assert.New(t)

	word := "test test"
	_, err := internal.NewAnagramLetterMap(word)
	assert.NoError(err)
}

func TestAnagramMapCreatedWithEmptyWord(t *testing.T) {
	assert := assert.New(t)

	word := ""
	_, err := internal.NewAnagramLetterMap(word)
	assert.Error(err)
}

func TestAnagramMapCreatedWithNonAlphabeticWord(t *testing.T) {
	assert := assert.New(t)

	word := "123"
	_, err := internal.NewAnagramLetterMap(word)
	assert.Error(err)
}

func TestAnagramMapCreatedWithMixedCaseWord(t *testing.T) {
	assert := assert.New(t)

	word := "Test"
	_, err := internal.NewAnagramLetterMap(word)
	assert.Error(err)
}

func TestAnagramMapCreatedWithTurkishWord(t *testing.T) {
	assert := assert.New(t)

	word := "şamil"
	_, err := internal.NewAnagramLetterMap(word)
	assert.Error(err)
}

func TestAnagramMapCreatedWithNonAlphabeticCharacter(t *testing.T) {
	assert := assert.New(t)

	word := "test!"
	_, err := internal.NewAnagramLetterMap(word)
	assert.Error(err)
}

func TestAnagramMapCreatedWithSpaceAndNonAlphabeticCharacter(t *testing.T) {
	assert := assert.New(t)

	word := "test !"
	_, err := internal.NewAnagramLetterMap(word)
	assert.Error(err)
}

func TestAnagramGeneratedSameLettermapForSameWord(t *testing.T) {
	assert := assert.New(t)

	word := "test"
	am1, err := internal.NewAnagramLetterMap(word)
	if err != nil {
		t.Errorf("AnagramLetterMap creation failed: %v", err)
	}

	am2, err := internal.NewAnagramLetterMap(word)
	if err != nil {
		t.Errorf("AnagramLetterMap creation failed: %v", err)
	}

	assert.Equal(am1, am2)
	assert.True(am1.IsAnagram(am2))
}

func TestAnagramGeneratedDifferentLettermapForDifferentWord(t *testing.T) {
	assert := assert.New(t)

	word1 := "test"
	am1, err := internal.NewAnagramLetterMap(word1)
	if err != nil {
		t.Errorf("AnagramLetterMap creation failed: %v", err)
	}

	word2 := "anothertest"
	am2, err := internal.NewAnagramLetterMap(word2)
	if err != nil {
		t.Errorf("AnagramLetterMap creation failed: %v", err)
	}

	assert.NotEqual(am1, am2)
	assert.False(am1.IsAnagram(am2))
}

func TestAnagramGeneratedSameLettermapForSameWordWithSpace(t *testing.T) {
	assert := assert.New(t)

	word := "test test"
	am1, err := internal.NewAnagramLetterMap(word)
	if err != nil {
		t.Errorf("AnagramLetterMap creation failed: %v", err)
	}

	word2 := "testtest"
	am2, err := internal.NewAnagramLetterMap(word2)
	if err != nil {
		t.Errorf("AnagramLetterMap creation failed: %v", err)
	}

	assert.Equal(am1, am2)
	assert.True(am1.IsAnagram(am2))
}

func TestAnagramGeneratedSameLettermapForDifferentOrder(t *testing.T) {
	assert := assert.New(t)

	word1 := "test"
	am1, err := internal.NewAnagramLetterMap(word1)
	if err != nil {
		t.Errorf("AnagramLetterMap creation failed: %v", err)
	}

	word2 := "tset"
	am2, err := internal.NewAnagramLetterMap(word2)
	if err != nil {
		t.Errorf("AnagramLetterMap creation failed: %v", err)
	}

	assert.Equal(am1, am2)
	assert.True(am1.IsAnagram(am2))
}

func TestAnagramGeneratedSameLettermapForDifferentOrderWithSpace(t *testing.T) {
	assert := assert.New(t)

	word1 := "test test"
	am1, err := internal.NewAnagramLetterMap(word1)
	if err != nil {
		t.Errorf("AnagramLetterMap creation failed: %v", err)
	}

	word2 := "tsettest"
	am2, err := internal.NewAnagramLetterMap(word2)
	if err != nil {
		t.Errorf("AnagramLetterMap creation failed: %v", err)
	}

	assert.Equal(am1, am2)
	assert.True(am1.IsAnagram(am2))
}

func TestAnagramHashesSameForSameWord(t *testing.T) {
	assert := assert.New(t)

	word := "test"
	am1, err := internal.NewAnagramLetterMap(word)
	if err != nil {
		t.Errorf("AnagramLetterMap creation failed: %v", err)
	}

	am2, err := internal.NewAnagramLetterMap(word)
	if err != nil {
		t.Errorf("AnagramLetterMap creation failed: %v", err)
	}

	assert.Equal(am1.AnagramHash(), am2.AnagramHash())
}

func TestAnagramHashesDifferentForDifferentWord(t *testing.T) {
	assert := assert.New(t)

	word1 := "test"
	am1, err := internal.NewAnagramLetterMap(word1)
	if err != nil {
		t.Errorf("AnagramLetterMap creation failed: %v", err)
	}

	word2 := "anothertest"
	am2, err := internal.NewAnagramLetterMap(word2)
	if err != nil {
		t.Errorf("AnagramLetterMap creation failed: %v", err)
	}

	assert.NotEqual(am1.AnagramHash(), am2.AnagramHash())
}

func TestLettermapAndHashesConvertedToEachOther(t *testing.T) {
	assert := assert.New(t)

	word := "test"
	am, err := internal.NewAnagramLetterMap(word)
	if err != nil {
		t.Errorf("AnagramLetterMap creation failed: %v", err)
	}

	am2, err := internal.NewAnagramLetterMapFromHash(am.AnagramHash())
	if err != nil {
		t.Errorf("AnagramLetterMap creation failed: %v", err)
	}

	assert.Equal(am, am2)
	assert.Equal(am.AnagramHash(), am2.AnagramHash())
}

func TestLettermapFromHashFailsWithInvalidHash(t *testing.T) {
	assert := assert.New(t)

	_, err := internal.NewAnagramLetterMapFromHash("invalidhash")
	assert.Error(err)
}
