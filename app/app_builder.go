package app

import (
	"errors"

	"github.com/selmant/anagram-finder-trendyol/app/config"
	"github.com/selmant/anagram-finder-trendyol/internal/input"
	"github.com/selmant/anagram-finder-trendyol/internal/storage"
)

const (
	ErrConfigNotSet          = "config is not set"
	ErrStorageFactoryNotSet  = "storage factory is not set"
	ErrReaderFactoryNotSet   = "reader factory is not set"
	ErrFailedToCreateStorage = "failed to create storage"
	ErrFailedToCreateReader  = "failed to create reader"
)

type AnagramApplicationBuilder struct {
	cfg            *config.Config
	storageFactory storage.Factory
	readerFactory  input.Factory
}

func NewAnagramApplicationBuilder() *AnagramApplicationBuilder {
	return &AnagramApplicationBuilder{}
}

func (b *AnagramApplicationBuilder) WithConfig(cfg *config.Config) *AnagramApplicationBuilder {
	b.cfg = cfg
	return b
}

func (b *AnagramApplicationBuilder) WithStorageFactory(factory storage.Factory) *AnagramApplicationBuilder {
	b.storageFactory = factory
	return b
}

func (b *AnagramApplicationBuilder) WithReaderFactory(factory input.Factory) *AnagramApplicationBuilder {
	b.readerFactory = factory
	return b
}

func (b *AnagramApplicationBuilder) Build() (*AnagramApplication, error) {
	if b.cfg == nil {
		return nil, errors.New(ErrConfigNotSet)
	}
	config.GlobalConfig = b.cfg

	if b.storageFactory == nil {
		return nil, errors.New(ErrStorageFactoryNotSet)
	}

	if b.readerFactory == nil {
		return nil, errors.New(ErrReaderFactoryNotSet)
	}

	storage := b.storageFactory.CreateStorage(b.cfg)
	reader := b.readerFactory.CreateReader(b.cfg)

	if storage == nil {
		return nil, errors.New(ErrFailedToCreateStorage)
	}

	if reader == nil {
		return nil, errors.New(ErrFailedToCreateReader)
	}

	return &AnagramApplication{
		Input:          reader,
		AnagramStorage: storage,
	}, nil
}
