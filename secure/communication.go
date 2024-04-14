package secure

import (
	"encoding/base64"
	"encoding/binary"
	"math/big"
	"net"
)

func Read(conn net.Conn) ([]byte, error) {
	var length int64
	err := binary.Read(conn, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, length)
	_, err = conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func Write(conn net.Conn, data []byte) error {
	length := int64(len(data))
	err := binary.Write(conn, binary.BigEndian, length)
	if err != nil {
		return err
	}
	_, err = conn.Write(data)
	return err
}

func Receive3Pass(conn net.Conn) ([]byte, error) {
	buffer, err := Read(conn)
	if err != nil {
		return nil, err
	}
	prime := new(big.Int).SetBytes(buffer)
	key, err := KeyFromPrime(prime)
	if err != nil {
		return nil, err
	}

	buffer, err = Read(conn)
	if err != nil {
		return nil, err
	}

	ciphertext := ShamirEncrypt(new(big.Int).SetBytes(buffer), key)

	err = Write(conn, ciphertext.Bytes())
	if err != nil {
		return nil, err
	}

	buffer, err = Read(conn)
	if err != nil {
		return nil, err
	}

	message := ShamirDecrypt(new(big.Int).SetBytes(buffer), key)
	return message.Bytes(), nil
}

func Send3Pass(conn net.Conn, message []byte, prime *big.Int) error {
	err := Write(conn, prime.Bytes())
	if err != nil {
		return err
	}
	key, err := KeyFromPrime(prime)
	if err != nil {
		return err
	}

	ciphertext := ShamirEncrypt(new(big.Int).SetBytes(message), key)
	err = Write(conn, ciphertext.Bytes())
	if err != nil {
		return err
	}

	buffer, err := Read(conn)
	if err != nil {
		return err
	}

	message = ShamirDecrypt(new(big.Int).SetBytes(buffer), key).Bytes()
	err = Write(conn, message)
	return err
}

func ReceiveAES(conn net.Conn, key string) ([]byte, error) {
	keyData, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	buffer, err := Read(conn)
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
	err = Write(conn, ciphertext)
	return err
}
