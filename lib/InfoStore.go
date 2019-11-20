package InfoStore

import (
	"log"
	"fmt"
	bolt "go.etcd.io/bbolt"
)

type InfoPair struct {
	Key string
	Value string
}


func Write(server string) {


db, err := bolt.Open("info.db", 0666, nil)
if err != nil {
  log.Fatal(err)
}
defer db.Close()

db.Update(func(tx *bolt.Tx) error {

	_, err := tx.CreateBucketIfNotExists([]byte("agentInfo"))
	if err != nil {
		return err
	}


	bucket := tx.Bucket([]byte("agentInfo"))
	error := bucket.Put([]byte("server"), []byte(server))
	



		v := bucket.Get([]byte("server"))
		fmt.Printf("The answer is: %s\n", v)
		return error

})

}