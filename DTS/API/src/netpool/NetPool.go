package netpool

import (
	"net"
)

const MaxConnections = 10000

type Error string

func (e Error) Error() string {
	return string(e)
}

var ErrMaxConn = Error("Maximum connections reached")

type Netpool struct {
	name  string
	conns int
	free  []net.Conn
}

func NewNetpool(name string) *Netpool {
	return &Netpool{
    		name: name,
	}
}

func (n *Netpool) Open() (conn net.Conn, err error) {
	if n.conns >= MaxConnections && len(n.free) == 0 {
		return nil, ErrMaxConn
	}

	if len(n.free) > 0 {
		// return the first free connection in the pool
		conn = n.free[0]
		n.free = n.free[1:]
	} else {
		addr, err := net.ResolveTCPAddr("tcp", n.name)
	if err != nil {
		return nil, err
	}
	conn, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	n.conns += 1
	}
	return conn, err
}

func (n *Netpool) Close(conn net.Conn) error {
	n.free = append(n.free, conn)
	return nil
}
