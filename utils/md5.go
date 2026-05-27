package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: MD5V
//@description: md5加密
//@param: str []byte
//@return: string

func MD5V(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}

func Hash256(buf []byte) string {
	h := md5.New()
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func MaskString(s string, n int) string {
	length := len(s)
	if length <= 2*n {

		return s
	}

	prefix := s[:n]

	suffix := s[length-n:]

	maskedPart := strings.Repeat("*", length-2*n)

	return prefix + maskedPart + suffix
}

// encrypt 函数：使用 AES-GCM 对明文进行加密
func GCMencrypt(plaintext []byte, key []byte) (string, error) {
	// 创建一个新的 AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 创建 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 生成随机的 nonce（必须与 GCM 的 NonceSize 相同）
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 使用 GCM 进行加密
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// 返回 base64 编码的密文（便于存储或传输）
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt 函数：使用 AES-GCM 对密文进行解密
func GCMdecrypt(encryptedString string, key []byte) ([]byte, error) {
	// 先对 base64 解码
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedString)
	if err != nil {
		return nil, err
	}

	// 创建 cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建 GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 获取 nonce 的长度
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// 分离 nonce 和实际密文
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// func main() {
//     // 示例明文
//     plaintext := []byte("Hello, this is a secret message!")

//     // 密钥（AES-256 需要 32 字节）
//     key := []byte("thisis32byteslongpassphrasewith!@#") // 32 bytes

//     // 加密
//     encrypted, err := encrypt(plaintext, key)
//     if err != nil {
//         panic(err)
//     }
//     fmt.Printf("Encrypted: %s\n", encrypted)

//     // 解密
//     decrypted, err := decrypt(encrypted, key)
//     if err != nil {
//         panic(err)
//     }
//     fmt.Printf("Decrypted: %s\n", string(decrypted))
// }
