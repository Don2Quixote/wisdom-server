package wisdom

import (
	"crypto/rand"
	"encoding/binary"
	"math/big"
	"net"
	"time"

	"github.com/pkg/errors"
)

// handle handles connection.
// Any blocking read/write operation with connection
// will be unblocked if connection will be closed outside
// this handle method.
func (s *Server) handle(conn net.Conn) error {
	defer conn.Close()

	// Get IP of connection.
	ip, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		return errors.Wrap(err, "can't split remote addr to host port")
	}

	// Get challenge according to IP.
	challenge, err := s.getChallenge(ip)
	if err != nil {
		return errors.Wrap(err, "can't get challenge")
	}

	// Send challenge to client.
	_, err = conn.Write(challenge)
	if err != nil {
		return errors.Wrap(err, "can't write to conn")
	}

	// Read response and check answer.
	err = s.checkAnswer(conn, challenge)
	if err != nil {
		return errors.Wrap(err, "can't check answer")
	}

	// Send quote.
	err = s.sendQuote(conn)
	if err != nil {
		return errors.Wrap(err, "can't send quote")
	}

	return nil
}

// getChallendge gets challenge number encoded as []byte for IP.
// Challenge's complexity depends on conns count from same IP.
func (s *Server) getChallenge(ip string) ([]byte, error) {
	// Complexity depends on conns count from same IP.
	complexity := s.increaseComplexity(ip)
	if complexity > s.pow.MaxComplexity {
		complexity = s.pow.MaxComplexity
	}

	// Buffer for challenge number.
	challenge := make([]byte, s.pow.MaxComplexity)

	// Put random number in buffer.
	// Number's byte-length depends on complexity.
	_, err := rand.Read(challenge[len(challenge)-complexity:])
	if err != nil {
		return nil, errors.Wrap(err, "can't get random")
	}

	// Prevent task with number < 2 (due to no solution).
	const minSolvableValue = 2
	if binary.BigEndian.Uint64(challenge[len(challenge)-8:]) < minSolvableValue {
		binary.BigEndian.PutUint64(challenge[len(challenge)-8:], minSolvableValue)
	}

	return challenge, nil
}

// checkAnswer reads answer and checks it.
// If it is incorrect - returned error is not nil.
func (s *Server) checkAnswer(conn net.Conn, challenge []byte) error {
	// Read count of factors in response.
	factorsCount, err := readUint32(conn)
	if err != nil {
		return errors.Wrap(err, "can't read uint32 from conn")
	}

	// factorsCountLimit should be configurable value but I have no time))
	if factorsCount > factorsCountLimit {
		return errors.Errorf("factors count (%d) > limit (%d)", factorsCount, factorsCountLimit)
	}

	// acc is accumulator of factors multiplication.
	acc := &big.Int{}

	for i := uint32(0); i < factorsCount; i++ {
		// Read count of bytes used to encode next factor.
		factorSize, err := readUint32(conn)
		if err != nil {
			return errors.Wrap(err, "can't read uint32 from conn")
		}

		// Any factor can't be bigger than initial challendge number.
		if factorSize > uint32(s.pow.MaxComplexity) {
			return errors.Errorf("factor size (%d) > max complexity (%d)", factorSize, s.pow.MaxComplexity)
		}

		// Read factor itself.
		factor, err := readBigInt(conn, int(factorSize))
		if err != nil {
			return errors.Wrap(err, "can't read big int from conn")
		}

		// If number is not prime - answer is not correct by defenition.
		const basesCount = 20
		if !factor.ProbablyPrime(basesCount) {
			return errors.Errorf("factor %v is not prime", factor)
		}

		// If acc hasn't been set yet - it is first number.
		if acc.BitLen() == 0 { // if acc == 0
			acc.Set(factor)
		} else {
			acc.Mul(acc, factor)
		}
	}

	challengeNumber := (&big.Int{}).SetBytes(challenge)
	if acc.Cmp(challengeNumber) != 0 {
		return errors.Errorf("wrong answer for challenge %v", challengeNumber)
	}

	return nil
}

// sendQuote gets random quote and writes it to conn.
func (s *Server) sendQuote(conn net.Conn) error {
	quote := randomQuote()

	s.log.Infof("sending wise quote %q", quote)

	// quoteLen encoded as uint32.
	var quoteLen [4]byte
	binary.BigEndian.PutUint32(quoteLen[:], uint32(len(quote)))

	_, err := conn.Write(quoteLen[:])
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(quote))
	if err != nil {
		return err
	}

	return nil
}

// increaseComplexity increments s.ipConns[ip],
// returns complexity for PoW and sets timer to restore complexity.
// Complexity >= 1.
func (s *Server) increaseComplexity(ip string) (complexity int) {
	s.mu.Lock()
	conns := s.ipConns[ip]
	s.ipConns[ip] = conns + 1
	s.mu.Unlock()

	go func() {
		time.Sleep(s.pow.ComplexityDuration)
		s.restoreComplexity(ip)
	}()

	// Complexity is >= 1.
	return int(float64(conns)*s.pow.ComplexityFactor) + 1
}

// restoreComplexity decrements s.addrConns[addr].
func (s *Server) restoreComplexity(addr string) {
	s.mu.Lock()
	s.ipConns[addr]--
	if s.ipConns[addr] == 0 {
		delete(s.ipConns, addr)
	}
	s.mu.Unlock()
}
