package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	client "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"logagent/tailf"
	"strings"
	"time"
)

type EtcdClient struct {
	Client *client.Client
	keys   []string
}

func initEtcd() (etcdClient *EtcdClient, err error) {
	cli, err := client.New(client.Config{
		Endpoints:   []string{"127.0.0.1:2379", "127.0.0.1:22379", "127.0.0.1:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}
	etcdClient = &EtcdClient{
		Client: cli,
	}
	return
}

func (cli *EtcdClient) getEtcdConf(key string) (collectConf []tailf.Collect, err error) {
	if strings.HasSuffix(key, "/") == false {
		key = key + "/"
	}
	for _, ip := range localIPArray {
		etcdKey := fmt.Sprintf("%s%s", key, ip)
		cli.keys = append(cli.keys, etcdKey)
		//从Etcd获取value超时机制
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := cli.Client.Get(ctx, etcdKey)
		if err != nil {
			logs.Error("client get from etcd failed, err:%v", err)
			continue
		}
		cancel()
		for _, v := range resp.Kvs {
			if string(v.Key) == etcdKey {
				err = json.Unmarshal(v.Value, &collectConf)
				if err != nil {
					logs.Error("unmarshal failed, err:%v", err)
					continue
				}
			}
		}
	}
	return
}

func (cli *EtcdClient) initEtcdWatcher() (err error) {
	for _, key := range cli.keys {
		go watcher(key, cli.Client)
	}
	return nil
}

func watcher(key string, cli *client.Client) {
	for {
		rch := cli.Watch(context.Background(), key)
		for wresp := range rch {
			var collectConf []tailf.Collect
			var getConfSucc = true
			for _, ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s] 's config deleted", key)
					continue
				}
				//判断etcd的key是否更新并且etcd的key和当前Client的key相等
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err := json.Unmarshal(ev.Kv.Value, &collectConf)
					if err != nil {
						logs.Error("key [%s], Unmarshal[%s], err:%v ", err)
						getConfSucc = false
						continue
					}
				}
				logs.Debug("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
			if getConfSucc {
				//传入更新后的配置
				tailf.UpdateConfig(collectConf)
			}
		}
	}
}
