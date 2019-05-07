package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/context"

	"github.com/coreos/etcd/clientv3"
)

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 2 * time.Second
	endpoints      = []string{"localhost:2379"}
)

func main() {
	etcdHost := "127.0.0.1:2379"
	keys := 20000

	fmt.Println("connecting to etcd - " + etcdHost)

	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://" + etcdHost},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("connected to etcd - " + etcdHost)

	defer etcd.Close()

	watchChan := etcd.Watch(context.Background(), "user/", clientv3.WithPrefix())

	data := make([]int64, keys)

	var recv int32
	go func() {
		for watchResp := range watchChan {
			go func(resp clientv3.WatchResponse) {
				for _, event := range resp.Events {
					i, _ := strconv.Atoi(strings.TrimPrefix(string(event.Kv.Key), "user/"))
					data[i] = time.Now().UnixNano() - data[i]
				}
			}(watchResp)
			atomic.AddInt32(&recv, int32(len(watchResp.Events)))
		}
	}()

	go func() {
		var wg sync.WaitGroup
		wg.Add(keys)
		for i := 0; i < keys; i++ {
			go func(index int) {
				defer wg.Done()
				now := time.Now()
				k := strconv.Itoa(index)
				v := now.String()

				data[index] = now.UnixNano()

				_, err := etcd.Put(context.Background(), "user/"+k, v)
				if err != nil {
					fmt.Printf("Put[%v] failed, %v\n", k, err)
				}
			}(i)
		}
		wg.Wait()
		fmt.Println("Put done")
	}()

	for int(atomic.LoadInt32(&recv)) != keys {
		time.Sleep(time.Second)
		fmt.Printf("recv: %d\n", atomic.LoadInt32(&recv))
	}
}
