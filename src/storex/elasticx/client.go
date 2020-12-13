package elasticx

import (
	"context"
	"github.com/ivalue2333/pkg/src/storex/elasticx/es_wrapped"
	"sync"

	"github.com/olivere/elastic/v7"
)

func NewClient(ctx context.Context, uris []string) (*es_wrapped.WrappedClient, error) {
	c, err := es_wrapped.NewClient(ctx, elastic.SetURL(uris...))
	if err != nil {
		return nil, err
	}
	return c, nil
}

type clientMgr struct {
	lock    sync.RWMutex
	clients map[string]*es_wrapped.WrappedClient
}

func (mgr *clientMgr) NewClient(ctx context.Context, name string, uris ...string) error {
	c, err := NewClient(ctx, uris)
	if err != nil {
		return err
	}
	mgr.Add(name, c)
	return err
}

func (mgr *clientMgr) Add(name string, client *es_wrapped.WrappedClient) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	mgr.clients[name] = client
}

func (mgr *clientMgr) Delete(name string) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	delete(mgr.clients, name)
}

func (mgr *clientMgr) Get(name string) *es_wrapped.WrappedClient {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	client, ok := mgr.clients[name]
	if !ok {
		return nil
	}
	return client
}

var _clientsMgr *clientMgr

func init() {
	_clientsMgr = &clientMgr{}
	_clientsMgr.clients = map[string]*es_wrapped.WrappedClient{}
}

func ClientsMgr() *clientMgr {
	return _clientsMgr
}
