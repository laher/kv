package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"time"

	"github.com/boltdb/bolt"

	"github.com/laher/kv/internal"
)

var (
	b      = flag.String("bck", "default", "bucket")
	usr, _ = user.Current()
	d      = flag.String("db", fmt.Sprintf("%s/.kv.db", usr.HomeDir), "database location")
	del    = flag.Bool("del", false, "delete key")
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
	if *del {
		err = kvdb.Del(key)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	//take from cli if available. Otherwise stdin
	if len(flag.Args()) > 1 {
		val := flag.Arg(1)
		err = kvdb.Set(key, val)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	kvdb.Set(key, string(bytes))

}
