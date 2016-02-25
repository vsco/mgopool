package mgopool

import (
	"gopkg.in/mgo.v2"
)

type leaky struct {
	base     *mgo.Session
	freeList chan *mgo.Session
}

// NewLeaky creates a leaky Pool of sessions copied from initial. A maximum of size sessions will be held in the free
// list, but the Pool will not block calls to Get. Releasing excess sessions to the pool are automatically Closed and
// removed.
func NewLeaky(initial *mgo.Session, size int) Pool {
	return &leaky{
		base:     initial.Clone(),
		freeList: make(chan *mgo.Session, size),
	}
}

func (p *leaky) Get() *mgo.Session {
	select {
	case s, closed := <-p.freeList:
		if s == nil || !closed {
			panic("pool has been closed")
		}
		return s
	default:
		return p.base.Copy()
	}
}

func (p *leaky) Put(s *mgo.Session) {
	if s == nil {
		return
	}

	select {
	case p.freeList <- s:
	default:
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

var _ Pool = &leaky{}
