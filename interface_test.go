package mgopool

import (
	"os"
	"testing"

	"github.com/globalsign/mgo"
)

func session(t *testing.T) *mgo.Session {
	host := os.Getenv("MONGO_HOST")

	sess, err := mgo.Dial(host)
	if err != nil {
		t.Fatalf("%v", err)
	}

	return sess
}
