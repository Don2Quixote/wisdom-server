package wisdom

import (
	"io"
	"math/big"
)

func readUint32(r io.Reader) (uint32, error) {
	var buf [4]byte
	var value big.Int

	_, err := io.ReadFull(r, buf[:])
	if err != nil {
		return 0, err
	}

	value.SetBytes(buf[:])

	return uint32(value.Uint64()), nil
}

func readBigInt(r io.Reader, length int) (*big.Int, error) {
	buf := make([]byte, length)
	value := &big.Int{}

	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	value.SetBytes(buf)

	return value, nil
}
