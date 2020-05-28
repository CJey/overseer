package overseer

//overseer listeners and connections allow graceful
//restarts by tracking when all connections from a listener
//have been closed

import (
	"net"
	"os"
	"time"
)

func newOverseerListener(l net.Listener) *overseerListener {
	return &overseerListener{
		Listener: l,
	}
}

//gracefully closing net.Listener
type overseerListener struct {
	net.Listener
}

func (l *overseerListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.(*net.TCPListener).AcceptTCP()
	if err != nil {
		return nil, err
	}
	conn.SetKeepAlive(true)                  // see http.tcpKeepAliveListener
	conn.SetKeepAlivePeriod(3 * time.Minute) // see http.tcpKeepAliveListener
	return conn, nil
}

//blocking wait for close
func (l *overseerListener) Close() error {
	return l.Listener.Close()
}

func (l *overseerListener) File() *os.File {
	// returns a dup(2) - FD_CLOEXEC flag *not* set
	tl := l.Listener.(*net.TCPListener)
	fl, _ := tl.File()
	return fl
}
