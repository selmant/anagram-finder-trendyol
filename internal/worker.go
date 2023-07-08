package internal

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/selmant/anagram-finder-trendyol/internal/input"
	"github.com/selmant/anagram-finder-trendyol/internal/storage"
	log "github.com/sirupsen/logrus"
)

const (
	ErrInputNotSet   = "input is not set"
	ErrStorageNotSet = "storage is not set"
)

type WorkerPool struct {
	workerCount int
	job         Job
}

func NewWorkerPool(workerCount int, job Job) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		job:         job,
	}
}

func (p *WorkerPool) Start(ctx context.Context) error {
	var doneJobCount atomic.Int32

	wg := &sync.WaitGroup{}
	wg.Add(p.workerCount)
	errChannel := make(chan error, 1)
	for i := 0; i < p.workerCount; i++ {
		workerID := i
		go func() {
			defer func() {
				wg.Done()

				if int(doneJobCount.Add(1)) == p.workerCount {
					log.Debug("All workers finished")
					close(errChannel)
				}
			}()
			log.Debugf("Worker started %d\n", workerID)
			err := p.job.Process(ctx)
			if err != nil {
				errChannel <- err
			}
			log.Debugf("Worker finished %d\n", workerID)
		}()
	}
	errs := make([]error, 0)
	for err := range errChannel {
		errs = append(errs, err)
	}
	wg.Wait()

	return errors.Join(errs...)
}

type Job interface {
	Process(ctx context.Context) error
}

type ReadAndMatchAnagramJob struct {
	storage storage.Storage
	input   input.DataReader
}

func NewReadAndMatchAnagramJob(storage storage.Storage, input input.DataReader) *ReadAndMatchAnagramJob {
	return &ReadAndMatchAnagramJob{
		storage: storage,
		input:   input,
	}
}

func (j *ReadAndMatchAnagramJob) Process(ctx context.Context) error {
	if j.input == nil {
		return errors.New(ErrInputNotSet)
	}

	if j.storage == nil {
		return errors.New(ErrStorageNotSet)
	}

	errs := make([]error, 0)
	for line := range j.input.Lines(ctx) {
		letterMap, err := NewAnagramLetterMap(line)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		err = j.storage.Store(ctx, letterMap.AnagramHash(), line)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	return errors.Join(errs...)
}
