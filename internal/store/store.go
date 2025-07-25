package store

import (
	"context"

	"github.com/guionardo/go-dev-monitor/internal/logging"
	"github.com/guionardo/go-dev-monitor/internal/repository"
)

type (
	DataStore struct {
		dataChan chan postData
		store    Store
	}
	postData struct {
		hostname   string
		repository *repository.Local
	}
	Store interface {
		Post(hostName string, repository *repository.Local) error
		GetSummary() (map[string][]*repository.Local, error)
		BeginPosts(hostName string) error
	}
)

func (ds *DataStore) BeginPosts(hostname string) {
	if err := ds.store.BeginPosts(hostname); err != nil {
		logging.Error("DataStore.BeginPosts", err)
	}

}

func New(queueSize int, storeFolder string) (*DataStore, error) {
	store, err := NewSqliteStore(storeFolder)
	if err != nil {
		return nil, err
	}

	return &DataStore{
		dataChan: make(chan postData, queueSize),
		store:    store,
	}, nil

}

func (ds *DataStore) Post(hostName string, repositoryData *repository.Local) error {
	go func() {
		_ = ds.store.Post(hostName, repositoryData)
	}()
	return nil
}

func (ds *DataStore) Get() (map[string][]*repository.Local, error) {
	return ds.store.GetSummary()
}

func (ds *DataStore) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-ds.dataChan:
				_ = ds.store.Post(data.hostname, data.repository)
			}
		}
	}()
}
