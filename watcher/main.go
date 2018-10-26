package main

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
)

var ErrNoChange = errors.New("key watcher: value has not been changed.")

type KvResult struct {
	Key   string
	Value string
}

type KeyWatcher struct {
	indexMap map[string]uint64
	agent    *api.Client
	env      string
}

func NewKeyWatcher(agentAddr string, env string) (*KeyWatcher, error) {
	rst := &KeyWatcher{
		env:      env,
		indexMap: make(map[string]uint64),
	}

	cfg := api.DefaultConfig()
	cfg.Address = agentAddr

	cli, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	rst.agent = cli

	return rst, nil
}

func (w *KeyWatcher) WatchKey(key string) <-chan *KvResult {
	rst := make(chan *KvResult)

	go func(ch chan<- *KvResult) {
		for {
			kv, err := w.doWatchKey(key)
			if err != nil {
				if err == ErrNoChange {
					log.Printf("value not changed")
					continue
				}
				log.Printf("watch error, %v", err)
				close(rst)
				return
			}
			rst <- kv
		}
	}(rst)
	return rst
}

func (w *KeyWatcher) doWatchKey(key string) (*KvResult, error) {
	idx, ok := w.indexMap[key]
	var pair *api.KVPair
	var meta *api.QueryMeta
	var err error
	if !ok {
		pair, meta, err = w.agent.KV().Get(key, nil)
	} else {
		option := &api.QueryOptions{WaitIndex: idx}
		pair, meta, err = w.agent.KV().Get(key, option)
	}

	if err != nil {
		return nil, err
	}

	if ok && idx == meta.LastIndex {
		return nil, ErrNoChange
	}

	if pair == nil {
		return nil, fmt.Errorf("key %s not found", key)
	}

	w.indexMap[key] = meta.LastIndex
	return &KvResult{
		Key:   key,
		Value: string(pair.Value),
	}, nil
}

func main() {
	watcher, err := NewKeyWatcher("192.168.2.251:8500", "")
	if err != nil {
		log.Printf("NewKeyWatcher failed, %v", err)
		return
	}

	ch := watcher.WatchKey("WATCH_TEST")

	for c := range ch {
		log.Printf("receive check result. key=%s, value=%s", c.Key, c.Value)
	}
}
