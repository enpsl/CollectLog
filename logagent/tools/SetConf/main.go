package main

import (
	"context"
	"encoding/json"
	"fmt"
	client "github.com/coreos/etcd/clientv3"
	"logagent/tailf"
	"time"
)

const (
	EtcdKey = "logagent/192.168.0.114"
)

func SetLogConfToEtcd() {
	cli, err := client.New(client.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}

	fmt.Println("connect succ")
	defer cli.Close()
	var logConfArr []tailf.Collect
	logConfArr = append(
		logConfArr,
		tailf.Collect{
			LogPath: "/tmp/logagent.log",
			Topic:   "nginx_log",
		},
	)
	logConfArr = append(
		logConfArr,
		tailf.Collect{
			LogPath: "/tmp/nginx_err.log",
			Topic:   "nginx_log_err",
		},
	)
	data, err := json.Marshal(logConfArr)
	if err != nil {
		fmt.Println("json failed, ", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, EtcdKey, string(data))
	//cli.Delete(ctx, EtcdKey)
	cancel()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}
}

func main() {
	SetLogConfToEtcd()
}
