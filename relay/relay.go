package relay

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/netip"
	"os/signal"
	"syscall"
	"time"

	"github.com/edup2p/common/types/key"
	"github.com/edup2p/common/types/relay"
	"github.com/edup2p/common/types/relay/relayhttp"
	"github.com/edup2p/common/types/stun"
	"github.com/edup2p/relay-server/config"
)

func RunRelay(cfg config.Config, privKey key.NodePrivate) error {
	slog.Info(
		"starting relay...",
		"bind", cfg.Bind,
		"port", cfg.Port,
		"stun_port", cfg.STUNPort,
		"public_key", privKey.Public().Debug(),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	netip.AddrPortFrom(cfg.Bind, cfg.Port)

	httpAp := netip.AddrPortFrom(cfg.Bind, cfg.Port)
	stunAp := netip.AddrPortFrom(cfg.Bind, cfg.STUNPort)

	stunServer := stun.NewServer(ctx)

	if err := stunServer.Listen(stunAp); err != nil {
		return fmt.Errorf("stun.Listen: %w", err)
	}

	go func() {
		defer cancel()

		if err := stunServer.Serve(); err != nil {
			slog.Error("stun server serve error", "err", err)
		}
	}()

	relayServer := relay.NewServer(privKey)

	mux := http.NewServeMux()

	mux.Handle("/relay", relayhttp.ServerHandler(relayServer))

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		browserHeaders(w)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		w.WriteHeader(http.StatusOK)

		if _, err := io.WriteString(w, ToverSokRelayDefaultHTML); err != nil {
			slog.Error("failed to write default HTML response", "err", err)
		}
	}))

	mux.Handle("/robots.txt", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		browserHeaders(w)
		if _, err := io.WriteString(w, "User-agent: *\nDisallow: /\n"); err != nil {
			slog.Error("failed to write robots.txt", "err", err)
		}
	}))
	mux.Handle("/generate_204", http.HandlerFunc(serverCaptivePortalBuster))

	httpServer := &http.Server{
		Addr:    httpAp.String(),
		Handler: mux,

		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	context.AfterFunc(ctx, func() {
		if err := httpServer.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdown server", "err", err)
		}
	})

	// TODO setup TLS with autocert: https://github.com/eduP2P/relay-server/issues/2

	slog.Info("relay started")
	err := httpServer.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("relay server errored", "err", err)
	} else {
		slog.Info("relay stopped")
	}

	return nil
}
