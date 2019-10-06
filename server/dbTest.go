package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/tecbot/gorocksdb"
)

var db *gorocksdb.DB
var ro *gorocksdb.ReadOptions
var wo *gorocksdb.WriteOptions

func openDB(path_to_db string) {
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(gorocksdb.NewLRUCache(3 << 30))
	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	db1, err := gorocksdb.OpenDb(opts, path_to_db)
	if err != nil {
		fmt.Println("Error in openDB")
	}
	db = db1
	ro = gorocksdb.NewDefaultReadOptions()
	wo = gorocksdb.NewDefaultWriteOptions()
}

func updateUserList(user_id string) {
	key := "U" + user_id
	val, err := db.Get(ro, []byte(key))
	if err != nil {
		fmt.Println("Error in updateUserList")
	}
	if val == nil {
		db.Put(wo, []byte(key), []byte(""))
	}
}

func updateMoveList(user_id string, x string, y string, color string) {

	move_key := "M" + strconv.Itoa(rand.Int())
	move_id_err := db.Put(wo, []byte(move_key), []byte(x+","+y+","+color+","+strconv.FormatInt(int64(time.Now().Unix()), 10)))
	if move_id_err != nil {
		fmt.Println("Error in generating move id")
	}
}

func main() {
	openDB("./rocksdb1")

}
