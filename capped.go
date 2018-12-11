package mgopool

import mgo "github.com/globalsign/mgo"

type capped struct {
	pool   Pool
	leases chan struct{}
}

// NewCapped creates a capped Pool of sessions copied from initial. Get on a capped Pool will block after size sessions
// have been retrieved until one is Put back in.
func NewCapped(initial *mgo.Session, size int) Pool {
	return &capped{
		pool:   NewLeaky(initial, size),
		leases: make(chan struct{}, size),
	}
}

func (p *capped) Get() *mgo.Session {
	p.leases <- struct{}{}
	return p.pool.Get()
}

func (p *capped) Put(s *mgo.Session) {
	select {
	case <-p.leases:
	default:
		//if leases is empty, continue anyways
	}

	p.pool.Put(s)
}

func (p *capped) Close() {
	close(p.leases)
	p.pool.Close()
}

func (p *capped) Used() int {
	return p.pool.Used()
}

var _ Pool = &capped{}
