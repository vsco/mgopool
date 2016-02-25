package mgopool

import (
	"os"
	"testing"

	"gopkg.in/mgo.v2"
)

func session(t *testing.T) *mgo.Session {
	host := os.Getenv("MONGO_HOST")

	sess, err := mgo.Dial(host)
	if err != nil {
		t.Fatalf("%v", err)
	}

	return sess
}
