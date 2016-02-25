// Package mgopool provides buffer implementations around mgo.Session (v2).
//
// With mgo, an initial session is created with the server, and then either a copy or clone of that initial session is
// used for each unit of work (HTTP request handler, worker loop, etc.), which is closed upon completion. This allows
// mgo to efficiently manage a pool of sockets between those sessions; however, this results in the underlying sockets
// to constantly log in and out of the Mongo cluster.
//
// mgopool avoids this issue by creating a Pool of copies from the initial session, storing them in a free list. Work
// units request an existing copy (or one is created if it doesn't exist) and then return it to the pool on completion.
// Both a leaky and capped version of the Pool is available, depending on need.
package mgopool

import "gopkg.in/mgo.v2"

// The Pool interface describes the mechanisms for retrieving and releasing mgo.Sessions. The pool should be used in
// place of calling Copy/Clone and Close on mgo.Session directly.
type Pool interface {
	// Get returns a *mgo.Session from the Pool. If there are no free sessions, a new one is copied from the initial
	// Session. If the Pool is capped and all sessions have been retrieved, Get will block until a session is returned to
	// with Put.
	Get() *mgo.Session

	// Put releases an *mgo.Session back into the Pool. Only healthy sessions (or nil) should be returned to the pool. If
	// a session is errors, Refresh should be called before releasing it to the Pool.
	Put(*mgo.Session)

	// Close drains the Pool, closing all held sessions. The initial session is not closed by this method. Calling Get or
	// Put on a closed pool will result in a panic.
	Close()
}
