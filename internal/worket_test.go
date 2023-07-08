package internal_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mock_internal "github.com/selmant/anagram-finder-trendyol/.mock/mock_internal"
	mock_input "github.com/selmant/anagram-finder-trendyol/.mock/mock_internal/input"
	mock_storage "github.com/selmant/anagram-finder-trendyol/.mock/mock_internal/storage"
	"github.com/selmant/anagram-finder-trendyol/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WorkerPoolTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	mockJob         *mock_internal.MockJob
	mockStorage     *mock_storage.MockStorage
	mockInputReader *mock_input.MockDataReader
}

func (suite *WorkerPoolTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockJob = mock_internal.NewMockJob(suite.ctrl)
	suite.mockStorage = mock_storage.NewMockStorage(suite.ctrl)
	suite.mockInputReader = mock_input.NewMockDataReader(suite.ctrl)
}

func (suite *WorkerPoolTestSuite) BeforeTest(_, _ string) {
	suite.ctrl.Finish()
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockJob = mock_internal.NewMockJob(suite.ctrl)
	suite.mockStorage = mock_storage.NewMockStorage(suite.ctrl)
	suite.mockInputReader = mock_input.NewMockDataReader(suite.ctrl)
}

func (suite *WorkerPoolTestSuite) TestWorkerPoolStart() {
	ctx := context.Background()

	workerCount := 2
	suite.mockJob.EXPECT().Process(gomock.Any()).Return(nil).Times(workerCount)

	wp := internal.NewWorkerPool(workerCount, suite.mockJob)
	err := wp.Start(ctx)
	assert.NoError(suite.T(), err)
}

func (suite *WorkerPoolTestSuite) TestWorkePoolStartWithError() {
	ctx := context.Background()

	workerCount := 2
	suite.mockJob.EXPECT().Process(gomock.Any()).Return(assert.AnError).Times(workerCount)

	wp := internal.NewWorkerPool(workerCount, suite.mockJob)
	err := wp.Start(ctx)
	assert.Error(suite.T(), err)
	asd := assert.AnError.Error()
	assert.Equal(suite.T(), asd+"\n"+asd, err.Error())
}

func (suite *WorkerPoolTestSuite) TestReadAndMatchAnagram() {
	ctx := context.Background()

	workerCount := 2
	linesChannel := make(chan string)
	suite.mockInputReader.EXPECT().Lines(gomock.Any()).Return(linesChannel).Times(workerCount)

	lines := []string{"abc", "def", "bac", "fed"}
	go func() {
		for _, line := range lines {
			linesChannel <- line
		}
		close(linesChannel)
	}()

	for _, line := range lines {
		suite.mockStorage.EXPECT().Store(gomock.Any(), gomock.Any(), line).Return(nil)
	}

	job := internal.NewReadAndMatchAnagramJob(suite.mockStorage, suite.mockInputReader)
	wp := internal.NewWorkerPool(workerCount, job)
	err := wp.Start(ctx)
	assert.NoError(suite.T(), err)
}

func (suite *WorkerPoolTestSuite) TestReadAndMatchAnagramWithStorageError() {
	ctx := context.Background()

	workerCount := 3
	linesChannel := make(chan string)

	lines := []string{"abc", "def", "bac", "def"}
	suite.mockInputReader.EXPECT().Lines(gomock.Any()).Return(linesChannel).Times(workerCount)
	go func() {
		for _, line := range lines {
			linesChannel <- line
		}
		close(linesChannel)
	}()

	for _, line := range lines {
		if line == "def" {
			suite.mockStorage.EXPECT().Store(gomock.Any(), gomock.Any(), "def").Return(assert.AnError)
			continue
		}
		suite.mockStorage.EXPECT().Store(gomock.Any(), gomock.Any(), line).Return(nil)
	}

	job := internal.NewReadAndMatchAnagramJob(suite.mockStorage, suite.mockInputReader)
	wp := internal.NewWorkerPool(workerCount, job)
	err := wp.Start(ctx)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), assert.AnError.Error()+"\n"+assert.AnError.Error(), err.Error())
}

func (suite *WorkerPoolTestSuite) TestReadAndMatchAnagramWithWrongInput() {
	ctx := context.Background()

	workerCount := 2
	linesChannel := make(chan string)
	suite.mockInputReader.EXPECT().Lines(gomock.Any()).Return(linesChannel).Times(workerCount)

	lines := []string{"abc", "def", "bac", "...123sdas"}
	go func() {
		for _, line := range lines {
			linesChannel <- line
		}
		close(linesChannel)
	}()

	for _, line := range lines {
		if line == "...123sdas" {
			continue
		}
		suite.mockStorage.EXPECT().Store(gomock.Any(), gomock.Any(), line).Return(nil)
	}

	job := internal.NewReadAndMatchAnagramJob(suite.mockStorage, suite.mockInputReader)
	wp := internal.NewWorkerPool(workerCount, job)
	err := wp.Start(ctx)
	assert.Error(suite.T(), err)
}

func TestWorkerPoolTestSuite(t *testing.T) {
	suite.Run(t, new(WorkerPoolTestSuite))
}
