package mgopool

import (
	"sync/atomic"

	mgo "gopkg.in/mgo.v2"
)

type leaky struct {
	base     *mgo.Session
	freeList chan *mgo.Session
	used     int32
}

// NewLeaky creates a leaky Pool of sessions copied from initial. A maximum of size sessions will be held in the free
// list, but the Pool will not block calls to Get. Releasing excess sessions to the pool are automatically Closed and
// removed.
func NewLeaky(initial *mgo.Session, size int) Pool {
	return &leaky{
		base:     initial.Clone(),
		freeList: make(chan *mgo.Session, size),
		used:     0,
	}

}

func (p *leaky) Get() *mgo.Session {
	atomic.AddInt32(&p.used, 1)
	select {
	case s, more := <-p.freeList:
		if s == nil || !more {
			panic("pool has been closed")
		}
		return s
	default:
		// if freeList is empty, Copy a new session
		return p.base.Copy()
	}
}

func (p *leaky) Put(s *mgo.Session) {
	atomic.AddInt32(&p.used, -1)
	if s == nil {
		return
	}

	select {
	case p.freeList <- s:
	default:
		// if freeList is full, close and discard the session
		s.Close()
	}
}

func (p *leaky) Close() {
	close(p.freeList)
	for s := range p.freeList {
		s.Close()
	}
	p.base.Close()
}

func (p *leaky) Used() int {
	return int(p.used)
}

var _ Pool = &leaky{}
