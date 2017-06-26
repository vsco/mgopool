package mgopool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeaky_Get(t *testing.T) {
	orig := session(t)
	defer orig.Close()
	p := NewLeaky(orig, 1)
	defer p.Close()

	s := p.Get()
	assert.NotNil(t, s)

	p.Put(s)
	next := p.Get()
	assert.Equal(t, s, next)
}

func TestLeaky_Put(t *testing.T) {
	orig := session(t)
	defer orig.Close()
	p := NewLeaky(orig, 1)
	defer p.Close()

	leak := p.Get()

	p.Put(p.Get())
	p.Put(leak)

	actual := p.Get()
	assert.NotEqual(t, leak, actual)

	p.Put(nil)
	actual = p.Get()

	assert.NotNil(t, actual)
}

func TestLeak_Close(t *testing.T) {
	orig := session(t)
	defer orig.Close()
	p := NewLeaky(orig, 1)

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

func TestLeaky_Used(t *testing.T) {
	orig := session(t)
	defer orig.Close()

	p := NewLeaky(orig, 1)
	defer p.Close()
	assert.Equal(t, 0, p.Used())

	leak := p.Get()
	assert.Equal(t, 1, p.Used())

	p.Put(p.Get())
	p.Put(leak)
	assert.Equal(t, 0, p.Used())
}
