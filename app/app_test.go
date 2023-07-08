package app_test

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	mock_input "github.com/selmant/anagram-finder-trendyol/.mock/mock_internal/input"
	mock_storage "github.com/selmant/anagram-finder-trendyol/.mock/mock_internal/storage"
	"github.com/selmant/anagram-finder-trendyol/app"
	"github.com/selmant/anagram-finder-trendyol/app/config"
	"github.com/selmant/anagram-finder-trendyol/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AppTestSuite struct {
	suite.Suite
	ctrl             *gomock.Controller
	mockInputReader  *mock_input.MockDataReader
	mockAnagramStore *mock_storage.MockStorage
}

func (suite *AppTestSuite) SetupSuite() {
	config.GlobalConfig = &config.Config{}
	config.GlobalConfig.WorkerPoolSize = 1
	config.GlobalConfig.WordsChannelSize = 8
}

func (suite *AppTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockInputReader = mock_input.NewMockDataReader(suite.ctrl)
	suite.mockAnagramStore = mock_storage.NewMockStorage(suite.ctrl)
}

func (suite *AppTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *AppTestSuite) TestAppPrintAnagrams() {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	mockInputReader := mock_input.NewMockDataReader(suite.ctrl)
	mockAnagramStore := mock_storage.NewMockStorage(suite.ctrl)

	app := app.AnagramApplication{
		Input:          mockInputReader,
		AnagramStorage: mockAnagramStore,
	}

	mockAnagramStore.EXPECT().AllAnagrams(gomock.Any()).Return(chanAnagrams(
		storage.AnagramResult{HashKey: "abc", Anagrams: []string{"abc", "acb"}, Error: nil},
		storage.AnagramResult{HashKey: "bac", Anagrams: []string{"bac", "bca"}, Error: nil},
	))

	err := app.PrintAnagrams(context.Background())
	assert.NoError(suite.T(), err)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	assert.Equal(suite.T(), "abc, acb\nbac, bca\n", string(out))
}

func (suite *AppTestSuite) TestAppPrintAnagramsError() {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	mockInputReader := mock_input.NewMockDataReader(suite.ctrl)
	mockAnagramStore := mock_storage.NewMockStorage(suite.ctrl)

	app := app.AnagramApplication{
		Input:          mockInputReader,
		AnagramStorage: mockAnagramStore,
	}

	mockAnagramStore.EXPECT().AllAnagrams(gomock.Any()).Return(chanAnagrams(
		storage.AnagramResult{Error: assert.AnError},
	))

	err := app.PrintAnagrams(context.Background())
	assert.Error(suite.T(), err)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	assert.Equal(suite.T(), "", string(out))
}

func (suite *AppTestSuite) TestAppPrintAnagramsEmpty() {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	mockInputReader := mock_input.NewMockDataReader(suite.ctrl)
	mockAnagramStore := mock_storage.NewMockStorage(suite.ctrl)

	app := app.AnagramApplication{
		Input:          mockInputReader,
		AnagramStorage: mockAnagramStore,
	}

	mockAnagramStore.EXPECT().AllAnagrams(gomock.Any()).Return(chanAnagrams())

	err := app.PrintAnagrams(context.Background())
	assert.NoError(suite.T(), err)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	assert.Equal(suite.T(), "", string(out))
}

func (suite *AppTestSuite) TestAppPrintAnagramsEmptyAnagrams() {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	mockInputReader := mock_input.NewMockDataReader(suite.ctrl)
	mockAnagramStore := mock_storage.NewMockStorage(suite.ctrl)

	app := app.AnagramApplication{
		Input:          mockInputReader,
		AnagramStorage: mockAnagramStore,
	}

	mockAnagramStore.EXPECT().AllAnagrams(gomock.Any()).Return(chanAnagrams(
		storage.AnagramResult{HashKey: "abc", Anagrams: []string{}, Error: nil},
		storage.AnagramResult{HashKey: "bac", Anagrams: []string{}, Error: nil},
	))

	err := app.PrintAnagrams(context.Background())
	assert.NoError(suite.T(), err)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	assert.Equal(suite.T(), "", string(out))
}

func (suite *AppTestSuite) TestAppHashAndStore() {
	mockInputReader := mock_input.NewMockDataReader(suite.ctrl)
	mockAnagramStore := mock_storage.NewMockStorage(suite.ctrl)

	app := app.AnagramApplication{
		Input:          mockInputReader,
		AnagramStorage: mockAnagramStore,
	}

	mockInputReader.EXPECT().Lines(gomock.Any()).Return(chanLines("abc", "acb", "bac", "bca"))
	for _, line := range []string{"abc", "acb", "bac", "bca"} {
		mockAnagramStore.EXPECT().Store(gomock.Any(), gomock.Any(), line).Return(nil)
	}

	err := app.HashAndStore(context.Background())
	assert.NoError(suite.T(), err)
}

func chanLines(lines ...string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, line := range lines {
			ch <- line
		}
	}()
	return ch
}

func chanAnagrams(anagrams ...storage.AnagramResult) <-chan storage.AnagramResult {
	ch := make(chan storage.AnagramResult)
	go func() {
		defer close(ch)
		for _, anagram := range anagrams {
			ch <- anagram
		}
	}()
	return ch
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}
