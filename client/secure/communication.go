package secure

import (
	"encoding/base64"
	"math/big"
	"net"
)

func Receive3Pass(conn net.Conn) ([]byte, error) {
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

	ciphertext := ShamirEncrypt(new(big.Int).SetBytes(buffer[:n]), key)

	_, err = conn.Write(ciphertext.Bytes())
	if err != nil {
		return nil, err
	}

	n, err = conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	message := ShamirDecrypt(new(big.Int).SetBytes(buffer[:n]), key)
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

	ciphertext := ShamirEncrypt(new(big.Int).SetBytes(message), key)
	_, err = conn.Write(ciphertext.Bytes())
	if err != nil {
		return err
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return err
	}

	message = ShamirDecrypt(new(big.Int).SetBytes(buffer[:n]), key).Bytes()
	_, err = conn.Write(message)
	return err
}

func ReceiveAES(conn net.Conn, key string) ([]byte, error) {
	keyData, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	return AESDecrypt(buffer[:n], keyData)
}

func SendAES(conn net.Conn, data []byte, key string) error {
	keyData, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return err
	}
	ciphertext, err := AESEncrypt(data, keyData)
	if err != nil {
		return err
	}
	_, err = conn.Write(ciphertext)
	return err
}
