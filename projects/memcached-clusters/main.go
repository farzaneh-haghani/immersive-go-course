package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// (Sharding)      go run . --mcrouter=11211 --memcacheds=11212,11213
// (Replicated)    go run . --mcrouter=11211 --memcacheds=11212,11213,11214

func main() {
	countOfKeys := 1000

	leaderPort := flag.String("mcrouter", "", "Port of leader")
	followersPort := flag.String("memcacheds", "", "Port of followers")
	flag.Parse()
	flag.Args()

	listOfFollowersPort := strings.Split(*followersPort, ",")

	for i := 0; i < countOfKeys; i++ {
		myItem := memcache.Item{
			Key:        fmt.Sprint(i),
			Value:      []byte(fmt.Sprint(i)),
			Expiration: 60,
		}

		address := fmt.Sprintf("localhost:%s", *leaderPort)
		mc := memcache.New(address)
		if err := mc.Set(&memcache.Item{Key: myItem.Key, Value: myItem.Value, Expiration: myItem.Expiration}); err != nil {
			fmt.Fprintf(os.Stderr, "Can't set a value: %s", err)
			os.Exit(52) //The server did not reply anything, which here is considered an error.
		}
	}
	for i := 0; i < 2; i++ { // Just for testing TTL
		countOfKeysInFollowers := []int{}
		replicated := true
		sharding := false

		for _, eachFollower := range listOfFollowersPort {
			count := 0
			address := fmt.Sprintf("localhost:%s", eachFollower)
			mc := memcache.New(address)

			for i := 0; i < countOfKeys; i++ {
				r, err := mc.Get(fmt.Sprint(i))

				if err != nil { // memcache: cache miss
					replicated = false
				} else if string(r.Value) == fmt.Sprint(i) {
					sharding = true
					count++
				} else {
					replicated = false
				}
			}
			countOfKeysInFollowers = append(countOfKeysInFollowers, count)
		}

		if replicated {
			fmt.Println("replicated")
		} else if sharding {
			isGood := true

			for i, countOfKeysInEachFollower := range countOfKeysInFollowers {
				percent := float64(countOfKeysInEachFollower) / float64(countOfKeys) * 100
				fmt.Printf("%.2f %% in server %d\n", percent, i)
				if math.Abs(percent-float64(100/len(listOfFollowersPort))) > 10 {
					isGood = false
				}
			}
			if isGood {
				fmt.Println("Good Sharded")
			} else {
				fmt.Println("Bad Sharded")
			}

		} else {
			fmt.Println("Expired")
		}
		replicated = true
		sharding = false

		time.Sleep(60 * time.Second) // Just for testing TTL
	}
}
