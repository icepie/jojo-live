package mideaac

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"log"
	"strings"

	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
)

const appKey = "434a209a5ce141c3b726de067835d7f0"
const signKey = "xhdiwjnchekd4d512chdjx5d8e4c394D2D7S"
const loginKey = "3742e9e5842d4ad59c2db887e12449f9"

type security struct {
	AppKey    []byte
	SignKey   []byte
	BlockSize int
	Iv        []byte
	// EncKey        []byte
	// DynamicKey    string
	tcpKey        []byte
	requestCount  int
	responseCount int
}

func NewSecurity() security {
	return security{
		AppKey:    []byte(appKey),
		SignKey:   []byte(signKey),
		BlockSize: 16,
		Iv:        []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		// EncKey:        encKey(),
		// DynamicKey:    dynamicKey,
		tcpKey:        nil,
		requestCount:  0,
		responseCount: 0,
	}
}

func (s security) AesDecrypt(raw []byte) []byte {
	k := s.EncKey()
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil
	}
	mode := ecb.NewECBDecrypter(block)
	pt := make([]byte, len(raw))
	mode.CryptBlocks(pt, raw)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	pt, err = padder.Unpad(pt)
	return pt
}

func (s security) AesEncrypt(raw []byte) []byte {
	k := s.EncKey()
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil
	}
	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	raw, err = padder.Pad(raw)

	ct := make([]byte, len(raw))
	mode.CryptBlocks(ct, raw)
	return ct
}

func (s security) AesCBCEncrypt(plaintext, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	// ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	mode := cipher.NewCBCEncrypter(block, s.Iv)
	mode.CryptBlocks(plaintext, plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	return plaintext
}

func (s security) AesCBCDecrypt(ciphertext, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := s.Iv
	// ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	return ciphertext
}

func (s security) EncKey() [16]byte {
	return md5.Sum(s.SignKey)
}

func (s security) DynamicKey() []byte {
	b := md5.Sum(s.AppKey)
	return b[:8]
}

func (s security) Encode32Data(raw []byte) [16]byte {
	return md5.Sum(append(raw, s.SignKey...))
}

func (s security) LocalKey(macS string, ssid string, pw string) [32]byte {
	mac, _ := hex.DecodeString(strings.Replace(macS, ":", "", -1))
	if len(mac) != 6 {
		panic("bad MAC address")
	}
	return sha256.Sum256(append([]byte(ssid+pw), mac...))
}

func (s security) TokenKeyPair(mac string, ssid string, pw string) ([]byte, []byte) {
	localKey := s.LocalKey(mac, ssid, pw)

	rnd := make([]byte, 32)
	rand.Read(rnd)

	key := strxor(rnd, localKey[:])
	token := s.AesCBCEncrypt(key, localKey[:])
	sign := sha256.Sum256(key)
	return append(token, sign[:]...), key
}

func (s *security) TcpKey(response, key []byte) ([]byte, bool) {
	if string(response) == "ERROR" {
		return []byte{}, false
	}
	if len(response) != 64 {
		return []byte{}, false
	}
	payload := response[:32]
	sign := response[32:]
	plain := s.AesCBCDecrypt(payload, key)
	b := sha256.Sum256(plain)
	if !bytes.Equal(b[:], sign) {
		return []byte{}, false
	}

	s.tcpKey = strxor(plain, key)
	s.requestCount = 0
	s.responseCount = 0
	return s.tcpKey, true
}

func (s *security) Encode8370(data []byte, msgtype byte) []byte {
	header := []byte{0x83, 0x70}
	size, padding := len(data), 0
	if msgtype == MSGTYPE_ENCRYPTED_RESPONSE || msgtype == MSGTYPE_ENCRYPTED_REQUEST {
		if (size+2)%16 != 0 {
			padding = 16 - (size + 2&0xf)
			size += padding + 32
			log.Println("padding", padding, "size", size)
			var b []byte
			rand.Read(b)
			data = append(data, b...)
		}
	}
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(size))
	header = append(header, b...)
	header = append(header, 0x20, byte(padding)<<4|byte(msgtype))
	if s.requestCount >= 0xFFF {
		s.requestCount = 0
	}
	d := make([]byte, 2)
	binary.BigEndian.PutUint16(d, uint16(s.requestCount))
	data = append(d, data...)
	s.requestCount += 1
	if msgtype == MSGTYPE_ENCRYPTED_RESPONSE || msgtype == MSGTYPE_ENCRYPTED_REQUEST {
		sign := sha256.Sum256(append(header, data...))
		data = append(s.AesCBCEncrypt(data, s.tcpKey), sign[:]...)
	}
	return append(header, data...)
}

func (s *security) Decode8370(data []byte) ([][]byte, []byte) {
	if len(data) < 6 {
		return [][]byte{}, data
	}
	header := data[:6]
	if header[0] != 0x83 || header[1] != 0x70 {
		// TODO error
		return nil, nil
	}

	size := binary.BigEndian.Uint16(header[2:4]) + 8
	var leftover []byte
	if len(data) < int(size) {
		return [][]byte{}, []byte{}
	} else if len(data) > int(size) {
		leftover = data[size:]
		data = data[:size]
	}
	if header[4] != 0x20 {
		log.Panic("missing byte 4")
	}
	padding := header[5] >> 4
	msgtype := header[5] & 0xF
	data = data[6:]
	if msgtype == MSGTYPE_ENCRYPTED_RESPONSE || msgtype == MSGTYPE_ENCRYPTED_REQUEST {
		sign := data[len(data)-32:]
		data = data[:len(data)-32]
		data = s.AesCBCDecrypt(data, s.tcpKey)
		b := sha256.Sum256(append(header, data...))
		if !bytes.Equal(b[:], sign) {
			log.Panic("sign does not match")
		}
		if padding > 0 {
			data = data[:len(data)-int(padding)]
		}
	}
	s.responseCount = int(binary.BigEndian.Uint16(data[:2]))
	data = data[2:]
	if leftover != nil {
		packets, incomplete := s.Decode8370(leftover)
		a := append([][]byte{}, data)
		return append(a, packets...), incomplete
	}
	return [][]byte{data}, []byte{}
}

// TODO
func (s *security) Sign(urli string, payload []byte) []byte {
	// u, _ := url.Parse(urli)
	// path := u.Path
	// query := sort.
	return []byte{}
}

func strxor(str1, str2 []byte) []byte {
	a := make([]byte, len(str1))
	for i, b := range str1 {
		a[i] = b ^ str2[i]
	}
	return a
}

func GetUdpID(data []byte) []byte {
	b := sha256.Sum256(data)
	b1, b2 := b[:16], b[16:]
	b3 := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < len(b1); i++ {
		b3[i] = b1[i] ^ b2[i]
	}
	return b3
}
