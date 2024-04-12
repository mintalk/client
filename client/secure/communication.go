package secure

import (
	"encoding/base64"
	"io"
	"math/big"
	"net"
)

func ReadData(conn net.Conn) ([]byte, error) {
	result := make([]byte, 0)
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		result = append(result, buffer[:n]...)
		if n < len(buffer) {
			break
		}
	}
	return result, nil
}

func BlockRead(conn net.Conn) ([]byte, error) {
	for {
		data, err := ReadData(conn)
		if err == nil {
			return data, nil
		}
	}
}

func Receive3Pass(conn net.Conn) ([]byte, error) {
	buffer, err := BlockRead(conn)
	if err != nil {
		return nil, err
	}
	prime := new(big.Int).SetBytes(buffer)
	key, err := KeyFromPrime(prime)
	if err != nil {
		return nil, err
	}

	buffer, err = BlockRead(conn)
	if err != nil {
		return nil, err
	}

	ciphertext := ShamirEncrypt(new(big.Int).SetBytes(buffer), key)

	_, err = conn.Write(ciphertext.Bytes())
	if err != nil {
		return nil, err
	}

	buffer, err = BlockRead(conn)
	if err != nil {
		return nil, err
	}

	message := ShamirDecrypt(new(big.Int).SetBytes(buffer), key)
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

	buffer, err := BlockRead(conn)
	if err != nil {
		return err
	}

	message = ShamirDecrypt(new(big.Int).SetBytes(buffer), key).Bytes()
	_, err = conn.Write(message)
	return err
}

func ReceiveAES(conn net.Conn, key string) ([]byte, error) {
	keyData, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	buffer, err := BlockRead(conn)
	if err != nil {
		return nil, err
	}
	return AESDecrypt(buffer, keyData)
}

func SendAES(conn net.Conn, data []byte, key string) error {
	keyData, err := base64.URLEncoding.DecodeString(key)
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
