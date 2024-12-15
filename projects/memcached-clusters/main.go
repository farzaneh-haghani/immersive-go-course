package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type item struct {
	key   string
	value string
	ttl   int32
}

func main() {

	myItem := item{
		key:   "0",
		value: "4",
		ttl:   60,
	}
	replicated := true
	sharding := false

	mcLeaderPortString := flag.String("mcrouter", "", "Port of leader")
	mcFollowersPortString := flag.String("memcacheds", "", "Ports of followers")
	flag.Parse()
	flag.Args()

	mcFollowersSlice := strings.Split(*mcFollowersPortString, ",")

	address := fmt.Sprintf("localhost:%s", *mcLeaderPortString)
	mc := memcache.New(address)
	if err := mc.Set(&memcache.Item{Key: myItem.key, Value: []byte(myItem.value), Expiration: myItem.ttl}); err != nil {
		fmt.Printf("Can't set a value: %s",err)
		return
	}

	for i := 0; i < 2; i++ { // Just for testing TTL

		for _, v := range mcFollowersSlice {
			address := fmt.Sprintf("localhost:%s", v)
			mc := memcache.New(address)

			r, err := mc.Get(myItem.key)

			if err != nil {
				replicated = false
			} else if myItem.value == string(r.Value) {
				sharding = true
			} else if myItem.value != string(r.Value) {
				replicated = false
			}
		}
		if replicated {
			fmt.Println("replicated")
		} else if sharding {
			fmt.Println("sharded")
		} else {
			fmt.Println("non")
		}
		replicated = true
		sharding = false

		time.Sleep(60 * time.Second) // Just for testing TTL
	}
}
