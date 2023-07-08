package input_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/selmant/anagram-finder-trendyol/app/config"
	"github.com/selmant/anagram-finder-trendyol/internal/input"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FileReaderSuite struct {
	suite.Suite
	tempFile string
}

func (suite *FileReaderSuite) SetupSuite() {
	config.GlobalConfig = &config.Config{
		WordsChannelSize: 8,
		WorkerPoolSize:   16,
	}
}

func (suite *FileReaderSuite) SetupTest() {
	f, err := os.CreateTemp("", "test")
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.tempFile = f.Name()
}

func (suite *FileReaderSuite) TearDownTest() {
	os.Remove(suite.tempFile)
}

func (suite *FileReaderSuite) TestFileReaderFileNotFound() {
	ctx := context.Background()
	fr := input.NewFileReader("unexistedfile.txt")

	err := fr.Prepare(ctx)
	assert.Error(suite.T(), err)
}

func (suite *FileReaderSuite) TestFileReaderFileFound() {
	ctx := context.Background()
	fr := input.NewFileReader(suite.tempFile)

	err := fr.Prepare(ctx)
	assert.NoError(suite.T(), err)
}

func (suite *FileReaderSuite) TestFileReaderReadSingleLine() {
	ctx := context.Background()
	fr := input.NewFileReader(suite.tempFile)

	err := fr.Prepare(ctx)
	assert.NoError(suite.T(), err)

	expectedWord := "test"
	err = os.WriteFile(suite.tempFile, []byte(expectedWord), 0644)
	if err != nil {
		suite.T().Fatal(err)
	}

	data := <-fr.Lines(ctx)
	assert.Equal(suite.T(), expectedWord, data)
}

func (suite *FileReaderSuite) TestFileReaderReadMultipleLines() {
	ctx := context.Background()
	fr := input.NewFileReader(suite.tempFile)

	err := fr.Prepare(ctx)
	assert.NoError(suite.T(), err)

	expectedWords := []string{"test1", "test2", "test3"}
	err = os.WriteFile(suite.tempFile, []byte("test1\ntest2\ntest3\n"), 0644)
	if err != nil {
		suite.T().Fatal(err)
	}

	count := 0
	for word := range fr.Lines(ctx) {
		assert.Equal(suite.T(), expectedWords[count], word)
		count++
	}

	assert.Equal(suite.T(), len(expectedWords), count)
}

func (suite *FileReaderSuite) TestFileReaderConcurrentRead() {
	ctx := context.Background()
	fr := input.NewFileReader(suite.tempFile)

	err := fr.Prepare(ctx)
	assert.NoError(suite.T(), err)

	expectedWord := "test"
	err = os.WriteFile(suite.tempFile, []byte("test\ntest\ntest\ntest\ntest\ntest\ntest\ntest\ntest\ntest\n"), 0644)
	if err != nil {
		suite.T().Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	count1 := 0
	count2 := 0
	go func() {
		for word := range fr.Lines(ctx) {
			time.Sleep(1 * time.Millisecond)
			count1++
			assert.Equal(suite.T(), expectedWord, word)
		}
		assert.Greater(suite.T(), count1, 0)
		wg.Done()
	}()

	go func() {
		for word := range fr.Lines(ctx) {
			time.Sleep(1 * time.Millisecond)
			count2++
			assert.Equal(suite.T(), expectedWord, word)
		}
		assert.Greater(suite.T(), count2, 0)
		wg.Done()
	}()

	wg.Wait()
	assert.Equal(suite.T(), 10, count1+count2)
}

func TestFileReaderSuite(t *testing.T) {
	suite.Run(t, new(FileReaderSuite))
}
