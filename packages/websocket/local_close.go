package websocket

import "errors"

// LocalCloseError reports explicit local closure or retirement during Reconnect.
// It does not report peer EOF or other transport failures: a consumer may still
// have a complete result received before that failure.
func (c *Connection[T]) LocalCloseError() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.localCloseErrorLocked()
}

func (c *Connection[T]) localCloseErrorLocked() error {
	if c.userClosed || c.reconnecting || errors.Is(c.err, ErrClosed) {
		return ErrClosed
	}
	select {
	case <-c.closeDone:
		return ErrClosed
	default:
		return nil
	}
}

// LocalCloseError reports lane detachment or explicit connection closure or
// retirement. A peer failure alone does not invalidate an already received result.
func (l *Lane[T]) LocalCloseError() error {
	c := l.connection
	c.mu.Lock()
	defer c.mu.Unlock()
	if l.queue.closed {
		return ErrLaneClosed
	}
	return c.localCloseErrorLocked()
}
