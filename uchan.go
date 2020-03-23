package uchan


import (
	"time"
)

type UnboundedChannel struct {
	In chan interface{}
	Out chan interface{}

	queue []interface{}
	isConnected bool
}

func (c *UnboundedChannel) Close() {
	close(c.In)
	for len(c.queue) > 0 {
		time.Sleep(50 * time.Millisecond)
	}
	close(c.Out)
	c.Disconnect()
}

func (c *UnboundedChannel) Connect() {
	if c.isConnected {
		return
	}
	c.isConnected = true
	next := func () interface{} {
		if len(c.queue) == 0 {
			return nil
		}
		return c.queue[0]
	}
	out := func () chan interface{} {
		if len(c.queue) == 0 {
			return nil
		}
		return c.Out
	}
	go func() {
		for c.isConnected {
			select {
			case item, ok := <- c.In:
				if !ok {
					// do nothing
				} else {
					c.queue = append(c.queue, item)
				}
			case out() <- next():
				c.queue = c.queue[1:]
			}
		}
	}()

}

func (c *UnboundedChannel) Disconnect() {
	c.isConnected = false
}

func NewUnboundedChannel() *UnboundedChannel {
	uc := &UnboundedChannel{
		In: make(chan interface{}),
		Out: make(chan interface{}),
		queue: []interface{}{},
	}
	uc.Connect()
	return uc
}
