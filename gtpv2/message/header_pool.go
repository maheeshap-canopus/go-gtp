package message

var headerPool = &pool{pool: make(chan *Header, 100)} // Default small buffer for applications who don't care about allocations

type pool struct {
	pool     chan *Header
	allocs   int64
	frees    int64
	gets     int64
	releases int64
}

func InitHeaderPool(bufferLen int) {
	headerPool = &pool{pool: make(chan *Header, bufferLen)}
}

func HeaderPoolStats() (allocs, frees, gets, releases int64) {
	return headerPool.allocs, headerPool.frees, headerPool.gets, headerPool.releases
}

func (p *pool) get() (c *Header) {
	p.gets++
	select {
	case c = <-p.pool:
		// Try to fetch an allocated struct from the pool
	default:
		// Init a new struct if nothing available
		c = &Header{}
		p.allocs++
	}
	return c
}

func (p *pool) release(c *Header) (n *Header) {
	// Ignore nil releases
	if c == nil {
		return nil
	}
	*c = Header{} // Reset fields
	p.releases++
	select {
	case p.pool <- c:
		// Return c to the pool
	default:
		p.frees++
		// No space in pool, let c be garbage collected
	}
	return nil
}

func ReleaseHeader(h *Header) *Header {
	headerPool.release(h)
	return nil
}
