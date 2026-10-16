// Package xmpp implements the XMPP client and stream handlers.
package xmpp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/FacileStudio/Chat/internal/config"
	"mellium.im/sasl"
	"mellium.im/xmpp/dial"
	"mellium.im/xmpp/jid"
)

// CreateTLSConfig returns TLS configuration based on the user config.
func CreateTLSConfig(cfg *config.Config, domain string) *tls.Config {
	return &tls.Config{
		ServerName:         domain,
		InsecureSkipVerify: cfg.InsecureTLS,
		MinVersion:         tls.VersionTLS12,
	}
}

// BuildMechanisms returns SASL authentication mechanisms.
func BuildMechanisms(hasUser bool) []sasl.Mechanism {
	if !hasUser {
		return []sasl.Mechanism{sasl.Anonymous}
	}
	return []sasl.Mechanism{
		sasl.ScramSha256Plus,
		sasl.ScramSha1Plus,
		sasl.ScramSha256,
		sasl.ScramSha1,
		sasl.Plain,
	}
}

// DialConn establishes a network connection to the XMPP server.
func DialConn(ctx context.Context, cfg *config.Config, addr jid.JID) (net.Conn, error) {
	if cfg.Server != "" {
		hostPort := net.JoinHostPort(cfg.Server, fmt.Sprintf("%d", cfg.Port))
		var d net.Dialer
		return d.DialContext(ctx, "tcp", hostPort)
	}
	d := &dial.Dialer{
		TLSConfig: CreateTLSConfig(cfg, addr.Domain().String()),
	}
	return d.Dial(ctx, "tcp", addr)
}
