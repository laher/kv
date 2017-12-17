package main

import (
	"flag"
	"fmt"
	"log"
	"os/user"
	"strings"
	"time"

	"github.com/boltdb/bolt"

	"github.com/laher/kv/internal"
)

var (
	usr, _ = user.Current()
	d      = flag.String("db", fmt.Sprintf("%s/.kv.db", usr.HomeDir), "database location")
	b      = flag.String("bck", "default", "bucket")
	ls     = flag.Bool("ls", false, "list entries")
	v      = flag.Bool("v", false, "verbose")
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
	if *ls {
		keys, err := kvdb.Ls()
		if err != nil {
			log.Fatal(err)
		}
		for _, key := range keys {
			if *v {
				val, err := kvdb.Get(key)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%s: %s\n", key, strings.TrimSpace(val))
			} else {
				fmt.Println(key)
			}
		}
		return
	}
	val, err := kvdb.Get(key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(val)
}
