package secure

import (
	"math/big"
	"net"
)

func Recieve3Pass(conn net.Conn) ([]byte, error) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	prime := new(big.Int).SetBytes(buffer[:n])
	key, err := KeyFromPrime(prime)
	if err != nil {
		return nil, err
	}

	n, err = conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	ciphertext := Encrypt(new(big.Int).SetBytes(buffer[:n]), key)

	_, err = conn.Write(ciphertext.Bytes())
	if err != nil {
		return nil, err
	}

	n, err = conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	message := Decrypt(new(big.Int).SetBytes(buffer[:n]), key)
	return message.Bytes(), nil
}

func Send3Pass(conn net.Conn, message []byte, prime *big.Int) error {
	_, err := conn.Write(prime.Bytes())
	if err != nil {
		return err
	}
	key, err := KeyFromPrime(prime)
	if err != nil {
		return err
	}

	ciphertext := Encrypt(new(big.Int).SetBytes(message), key)
	_, err = conn.Write(ciphertext.Bytes())
	if err != nil {
		return err
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return err
	}

	message = Decrypt(new(big.Int).SetBytes(buffer[:n]), key).Bytes()
	_, err = conn.Write(message)
	return err
}
