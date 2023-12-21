package ie

import "pgregory.net/rand"

var iePool = newIEPool(100, 32) // Default small buffer for applications who don't care about allocations

type pool struct {
	pool      []chan *IE
	nChannels int
	allocs    int64
	frees     int64
	gets      int64
	releases  int64
}

func newIEPool(bufferLen, nChannels int) *pool {
	sp := &pool{pool: make([]chan *IE, nChannels), nChannels: nChannels}
	for i := 0; i < nChannels; i++ {
		sp.pool[i] = make(chan *IE, bufferLen)
	}
	return sp
}

func InitIEPool(bufferLen, nChannels int) {
	iePool = newIEPool(bufferLen, nChannels)
}

func IEPoolStats() (allocs, frees, gets, releases int64) {
	return iePool.allocs, iePool.frees, iePool.gets, iePool.releases
}

func (p *pool) get() (c *IE) {
	p.gets++
	select {
	case c = <-p.pool[rand.Int()%p.nChannels]:
		// Try to fetch an allocated struct from the pool
	default:
		// Init a new struct if nothing available
		c = &IE{}
		p.allocs++
	}
	return c
}

func (p *pool) release(c *IE) (n *IE) {
	// Ignore nil releases
	if c == nil {
		return nil
	}

	// reset fields but preserve slice capacity
	c.Type = 0
	c.Length = 0
	c.instance = 0
	c.Payload = c.Payload[:0]
	for _, i := range c.ChildIEs {
		iePool.release(i)
	}
	ReleaseMultiParseContainer(c.ChildIEs)
	c.ChildIEs = nil

	p.releases++
	select {
	case p.pool[rand.Int()%p.nChannels] <- c:
		// Return c to the pool
	default:
		p.frees++
		// No space in pool, let c be garbage collected
	}
	return nil
}

var sPool = newIESlicePool(100, 4) // Default small buffer for applications who don't care about allocations

type slicePool struct {
	pool      []chan []*IE
	nChannels int
	allocs    int64
	frees     int64
	gets      int64
	releases  int64
}

func newIESlicePool(bufferLen, nChannels int) *slicePool {
	sp := &slicePool{pool: make([]chan []*IE, nChannels), nChannels: nChannels}
	for i := 0; i < nChannels; i++ {
		sp.pool[i] = make(chan []*IE, bufferLen)
	}
	return sp
}

func InitIESlicePool(bufferLen, nChannels int) {
	sPool = newIESlicePool(bufferLen, nChannels)
}

func IESlicePoolStats() (allocs, frees, gets, releases int64) {
	return sPool.allocs, sPool.frees, sPool.gets, sPool.releases
}

func (p *slicePool) get() (c []*IE) {
	p.gets++
	select {
	case c = <-p.pool[rand.Int()%p.nChannels]:
		// Try to fetch an allocated slice from a random part of the pool
	default:
		// Init a new slice if nothing available
		c = make([]*IE, 0, 8)
		p.allocs++
	}
	return c
}

func (p *slicePool) release(c []*IE) {
	// Ignore nil releases
	if c == nil {
		return
	}
	// reset length but preserve slice capacity
	c = c[:0]

	p.releases++
	select {
	case p.pool[rand.Int()%p.nChannels] <- c:
		// Return c to a random part of the pool
	default:
		p.frees++
		// No space in pool, let c be garbage collected
	}
}

func Release(i *IE) *IE {
	iePool.release(i)
	return nil
}

func ReleaseSlice(s []*IE) []*IE {
	for _, i := range s {
		iePool.release(i)
	}
	return s[:0]
}

func ReleaseMultiParseContainer(s []*IE) {
	sPool.release(s)
}

func ReleaseMultiIEsAndContainer(s []*IE) {
	for _, i := range s {
		iePool.release(i)
	}
	sPool.release(s)
}
