package wisdom

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"wisdom/pkg/logger"

	"github.com/pkg/errors"
)

// Server is general app's entity that represents TCP
// server of Word Of Wisdom.
type Server struct {
	port int
	pow  PoW
	log  logger.Logger

	// ipConns serves to count connections from IPs
	// for calculation PoW complexity for new connection.
	ipConns map[string]int
	// mu is mutex to protect operations with addrConns.
	mu *sync.Mutex
}

// PoW is settings for Proof of Work algorithm of the server.
type PoW struct {
	// ComplexityFactor shows how fast does complexity grow.
	ComplexityFactor float64
	// MaxComplexity limits the maximum complexity level.
	MaxComplexity int
	// ComplexityDuration is value that shows how much time should pass
	// to restore complexity points after grow.
	ComplexityDuration time.Duration
}

// NewServer returns new Server.
func NewServer(port int, pow PoW, log logger.Logger) *Server {
	return &Server{
		port: port,
		pow:  pow,
		log:  log,

		ipConns: make(map[string]int),
		mu:      &sync.Mutex{},
	}
}

// Launch launches server.
// Returned error is nil if context closed.
func (s *Server) Launch(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrapf(err, "can't listen port %d", s.port)
	}

	connected := &sync.WaitGroup{}

	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					s.log.Error(errors.Wrap(err, "can't accept conn"))
				}

				continue
			}

			connected.Add(1)

			handleDone := make(chan struct{})

			go func() {
				err := s.handle(conn)
				if err != nil {
					s.log.Warn(errors.Wrap(err, "error handling connection"))
				}

				handleDone <- struct{}{}
			}()

			go func() {
				select {
				case <-ctx.Done():
					err := conn.Close()
					if err != nil {
						s.log.Warn(errors.Wrap(err, "can't close connection"))
					}
				case <-handleDone:
				}

				connected.Done()
			}()
		}
	}()

	<-ctx.Done()
	connected.Wait()

	err = lis.Close()
	if err != nil {
		return errors.Wrap(err, "can't close listener")
	}

	return nil
}
