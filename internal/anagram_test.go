package internal_test

import (
	"testing"

	"github.com/selmant/anagram-finder-trendyol/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AnagramLetterMapTestSuite struct {
	suite.Suite
}

func (suite *AnagramLetterMapTestSuite) TestAnagramMapCreated() {
	assert := assert.New(suite.T())

	word := "test"
	am, err := internal.NewAnagramLetterMap(word)
	assert.NoError(err)

	expected := internal.AnagramLetterMap{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 0, 0, 0, 0, 0, 0}
	assert.Equal(expected, am)
}

func (suite *AnagramLetterMapTestSuite) TestAnagramMapCreatedWithSpace() {
	assert := assert.New(suite.T())

	word := "test test"
	_, err := internal.NewAnagramLetterMap(word)
	assert.NoError(err)
}

func (suite *AnagramLetterMapTestSuite) TestAnagramMapCreatedWithEmptyWord() {
	assert := assert.New(suite.T())

	word := ""
	_, err := internal.NewAnagramLetterMap(word)
	assert.Error(err)
}

func (suite *AnagramLetterMapTestSuite) TestAnagramMapCreatedWithNonAlphabeticWord() {
	assert := assert.New(suite.T())

	word := "123"
	_, err := internal.NewAnagramLetterMap(word)
	assert.Error(err)
}

func (suite *AnagramLetterMapTestSuite) TestAnagramMapCreatedWithMixedCaseWord() {
	assert := assert.New(suite.T())

	word := "Test"
	_, err := internal.NewAnagramLetterMap(word)
	assert.Error(err)
}

func (suite *AnagramLetterMapTestSuite) TestAnagramMapCreatedWithTurkishWord() {
	assert := assert.New(suite.T())

	word := "şamil"
	_, err := internal.NewAnagramLetterMap(word)
	assert.Error(err)
}

func (suite *AnagramLetterMapTestSuite) TestAnagramMapCreatedWithNonAlphabeticCharacter() {
	assert := assert.New(suite.T())

	word := "test!"
	_, err := internal.NewAnagramLetterMap(word)
	assert.Error(err)
}

func (suite *AnagramLetterMapTestSuite) TestAnagramMapCreatedWithSpaceAndNonAlphabeticCharacter() {
	assert := assert.New(suite.T())

	word := "test !"
	_, err := internal.NewAnagramLetterMap(word)
	assert.Error(err)
}

func (suite *AnagramLetterMapTestSuite) TestAnagramGeneratedSameLettermapForSameWord() {
	assert := assert.New(suite.T())

	word := "test"
	am1, err := internal.NewAnagramLetterMap(word)
	assert.NoError(err)

	am2, err := internal.NewAnagramLetterMap(word)
	assert.NoError(err)

	assert.Equal(am1, am2)
	assert.True(am1.IsAnagram(am2))
}

func (suite *AnagramLetterMapTestSuite) TestAnagramGeneratedDifferentLettermapForDifferentWord() {
	assert := assert.New(suite.T())

	word1 := "test"
	am1, err := internal.NewAnagramLetterMap(word1)
	assert.NoError(err)

	word2 := "anothertest"
	am2, err := internal.NewAnagramLetterMap(word2)
	assert.NoError(err)

	assert.NotEqual(am1, am2)
	assert.False(am1.IsAnagram(am2))
}

func (suite *AnagramLetterMapTestSuite) TestAnagramGeneratedSameLettermapForSameWordWithSpace() {
	assert := assert.New(suite.T())

	word1 := "test test"
	am1, err := internal.NewAnagramLetterMap(word1)
	assert.NoError(err)

	word2 := "testtest"
	am2, err := internal.NewAnagramLetterMap(word2)
	assert.NoError(err)

	assert.Equal(am1, am2)
	assert.True(am1.IsAnagram(am2))
}

func (suite *AnagramLetterMapTestSuite) TestAnagramGeneratedSameLettermapForDifferentOrder() {
	assert := assert.New(suite.T())

	word1 := "test"
	am1, err := internal.NewAnagramLetterMap(word1)
	assert.NoError(err)

	word2 := "tset"
	am2, err := internal.NewAnagramLetterMap(word2)
	assert.NoError(err)

	assert.Equal(am1, am2)
	assert.True(am1.IsAnagram(am2))
}

func (suite *AnagramLetterMapTestSuite) TestAnagramGeneratedSameLettermapForDifferentOrderWithSpace() {
	assert := assert.New(suite.T())

	word1 := "test test"
	am1, err := internal.NewAnagramLetterMap(word1)
	assert.NoError(err)

	word2 := "tsettest"
	am2, err := internal.NewAnagramLetterMap(word2)
	assert.NoError(err)

	assert.Equal(am1, am2)
	assert.True(am1.IsAnagram(am2))
}

func (suite *AnagramLetterMapTestSuite) TestAnagramHashesSameForSameWord() {
	assert := assert.New(suite.T())

	word := "test"
	am1, err := internal.NewAnagramLetterMap(word)
	assert.NoError(err)

	am2, err := internal.NewAnagramLetterMap(word)
	assert.NoError(err)

	assert.Equal(am1.AnagramHash(), am2.AnagramHash())
}

func (suite *AnagramLetterMapTestSuite) TestAnagramHashesDifferentForDifferentWord() {
	assert := assert.New(suite.T())

	word1 := "test"
	am1, err := internal.NewAnagramLetterMap(word1)
	assert.NoError(err)

	word2 := "anothertest"
	am2, err := internal.NewAnagramLetterMap(word2)
	assert.NoError(err)

	assert.NotEqual(am1.AnagramHash(), am2.AnagramHash())
}

func (suite *AnagramLetterMapTestSuite) TestLettermapAndHashesConvertedToEachOther() {
	assert := assert.New(suite.T())

	word := "test"
	am, err := internal.NewAnagramLetterMap(word)
	assert.NoError(err)

	am2, err := internal.NewAnagramLetterMapFromHash(am.AnagramHash())
	assert.NoError(err)

	assert.Equal(am, am2)
	assert.Equal(am.AnagramHash(), am2.AnagramHash())
}

func (suite *AnagramLetterMapTestSuite) TestLettermapFromHashFailsWithInvalidHash() {
	assert := assert.New(suite.T())

	_, err := internal.NewAnagramLetterMapFromHash("invalidhash")
	assert.Error(err)
}

func TestAnagramLetterMapTestSuite(t *testing.T) {
	suite.Run(t, new(AnagramLetterMapTestSuite))
}
