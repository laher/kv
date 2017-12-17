package main

import (
	"flag"
	"fmt"
	"log"
	"os/user"
	"time"

	"github.com/boltdb/bolt"

	"github.com/laher/kv/internal"
)

var (
	b      = flag.String("bck", "default", "bucket")
	usr, _ = user.Current()
	d      = flag.String("db", fmt.Sprintf("%s/.kv.db", usr.HomeDir), "database location")
)

func main() {
	flag.Parse()
	db, err := bolt.Open(*d, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	kvdb, err := kv.NewKeyValue(db, *b)
	if err != nil {
		log.Fatal(err)
	}
	key := "default"
	if len(flag.Args()) > 0 {
		key = flag.Arg(0)
	}
	val, err := kvdb.Get(key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(val)
}
