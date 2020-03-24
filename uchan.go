package uchan

import (
	"time"
)

// UnboundedChannel provides a non-blocking, unbuffered channel construct. The channel uncouples
// sender and receiver without having a fixed-size channel buffer and both sides can go about
// their operation without dependence on each other. The `In` channel is used to send, while the
// `Out` channel is used to receive items.
type UnboundedChannel struct {

	// In is the channel to which senders can publish.
	In chan interface{}

	// Out is the channel to which receivers can subscribe.
	Out chan interface{}

	// queue contains the internal buffer of items.
	queue []interface{}

	// isConnected is a flag keeping the In connected to the Out channel. If this
	// is false, the goroutine managing the transport will shut down.
	isConnected bool
}

// Close shuts down the unbounded channel and is to be used instead of close(chan).
// This will block the caller until the whole queue is drained - ie. all items have
// been received via the Out channel.
func (c *UnboundedChannel) Close() {
	close(c.In)
	for len(c.queue) > 0 {
		time.Sleep(50 * time.Millisecond)
	}
	close(c.Out)
	c.Disconnect()
}

// Connect spawns the transport between In and Out using the intermediate queue. If the
// channels are already connected, it will do nothing. This cannot be called after a
// the Close() function has been used as channels will be permanently closed.
// This in fact, should never be called manually, unless the NewUnboundedChannel factory
// is not used.
func (c *UnboundedChannel) Connect() {
	if c.isConnected {
		return
	}
	c.isConnected = true
	next := func() interface{} {
		if len(c.queue) == 0 {
			return nil
		}
		return c.queue[0]
	}
	out := func() chan interface{} {
		if len(c.queue) == 0 {
			return nil
		}
		return c.Out
	}
	go func() {
		for c.isConnected {
			select {
			case item, ok := <-c.In:
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

// Disconnect sets the internal stop flag and terminates the transport between In and Out
// effectively.
func (c *UnboundedChannel) Disconnect() {
	c.isConnected = false
}

// NewUnboundedChannel creates and connects a new unbounded channel ready for receiving and sending
// items via In and Out.
func NewUnboundedChannel() *UnboundedChannel {
	uc := &UnboundedChannel{
		In:    make(chan interface{}),
		Out:   make(chan interface{}),
		queue: []interface{}{},
	}
	uc.Connect()
	return uc
}
