package mgopool

import (
	"os"

	"sync"

	"fmt"

	"github.com/globalsign/mgo"
)

func Example() {
	host := os.Getenv("MONGO_HOST")

	// create the initial mgo.Session, defer closing it
	initial, _ := mgo.Dial(host)
	defer initial.Close()

	// create a new leaky pool with size of 3
	pool := NewLeaky(initial, 3)
	defer pool.Close()

	wg := sync.WaitGroup{}

	// create a few workers to utilize the sessions in the pool
	for i := 0; i < 3; i++ {
		wg.Add(1)

		go func() {
			for j := 0; j < 2; j++ {
				// get a session from the pool
				session := pool.Get()

				// do something with that session, refresh it if there is an error
				if err := doSomething(session); err != nil {
					session.Refresh()
				}

				// return the session to the pool
				pool.Put(session)
			}
			wg.Done()
		}()
	}

	// wait for the workers to complete
	wg.Wait()

	// Output:
	// foobar
	// foobar
	// foobar
	// foobar
	// foobar
	// foobar
}

func doSomething(s *mgo.Session) error {
	fmt.Println("foobar")
	return s.Ping()
}
