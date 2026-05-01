package common

import (
	"fmt"
	"net"
)

func RemoteIPFromConn(conn fmt.Stringer) string {
	if host, _, err := net.SplitHostPort(conn.String()); err == nil {
		return host
	}
	return conn.String()
}
