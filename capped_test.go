package mgopool

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	mgo "gopkg.in/mgo.v2"
)

func TestCapped_Get(t *testing.T) {
	orig := session(t)
	defer orig.Close()
	p := NewCapped(orig, 1)
	defer p.Close()

	expected := p.Get()
	var actual *mgo.Session

	done := make(chan struct{})
	go func() {
		actual = p.Get()
		close(done)
	}()

	p.Put(expected)

	select {
	case <-done:
	case <-time.NewTimer(time.Millisecond * 10).C:
		t.Fatal("capped bucket never unblocked")
	}

	assert.Equal(t, expected, actual)
}

func TestCapped_Put(t *testing.T) {
	orig := session(t)
	defer orig.Close()
	p := NewCapped(orig, 1)
	defer p.Close()

	expected := p.Get()

	p.Put(expected)
	p.Put(orig.Copy())

	actual := p.Get()
	assert.Equal(t, expected, actual)

	done := make(chan struct{})
	go func() {
		actual = p.Get()
		close(done)
	}()

	p.Put(expected)

	select {
	case <-done:
	case <-time.NewTimer(time.Millisecond * 10).C:
		t.Fatal("capped bucket never unblocked")
	}

	assert.Equal(t, expected, actual)
}

func TestCapped_Close(t *testing.T) {
	orig := session(t)
	defer orig.Close()
	p := NewCapped(orig, 1)

	s := p.Get()

	p.Put(s)
	p.Close()

	assert.Panics(t, func() {
		p.Get()
	})

	assert.Panics(t, func() {
		p.Put(s)
	})

	assert.Panics(t, func() {
		s.Copy()
	})

	assert.NotPanics(t, func() {
		orig.Copy().Close()
	})
}

func TestCapped_Used(t *testing.T) {
	orig := session(t)
	defer orig.Close()

	p := NewCapped(orig, 1)
	defer p.Close()
	assert.Equal(t, 0, p.Used())

	leak := p.Get()
	assert.Equal(t, 1, p.Used())

	p.Put(leak)
	assert.Equal(t, 0, p.Used())
}
