package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func AESEncrypt(data []byte, key []byte) ([]byte, error) {
	if len(key) != 2*aes.BlockSize {
		return nil, fmt.Errorf("key length must be %d", 2*aes.BlockSize)
	}
	iv := make([]byte, aes.BlockSize)
	copy(iv, key[:aes.BlockSize])

	var dataBlock []byte
	length := len(data)
	if length%aes.BlockSize != 0 {
		extendBlock := aes.BlockSize - (length % aes.BlockSize)
		dataBlock = make([]byte, length+extendBlock)
	} else {
		dataBlock = make([]byte, length)
	}
	copy(dataBlock, data)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, len(dataBlock))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, dataBlock)
	return cipherText, nil
}

func AESDecrypt(data []byte, key []byte) ([]byte, error) {
	if len(key) != 2*aes.BlockSize {
		return nil, fmt.Errorf("key length must be %d", 2*aes.BlockSize)
	}
	iv := make([]byte, aes.BlockSize)
	copy(iv, key[:aes.BlockSize])

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plainText := make([]byte, len(data))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plainText, data)
	return plainText, nil
}
