package network

import (
	"TimeSeriesData/crypto"
	"crypto/rsa"
	"crypto/sha256"
	"log"
	"net"
)

/*

Client message structure (RSA encoded, maximum request data length ~ 462 bytes for RSA 4096):
|AES key - 32 bytes|AES gcm nonce - 12 bytes|Request data|

Server message structure:
|Response + sha256 of response data encrypted with AES-GCM|

*/

const (
	ERROR uint8 = 1
	OK    uint8 = 0
)

type TcpServer[T any] struct {
	port     int
	key      *rsa.PrivateKey
	label    []byte
	handler  func([]byte, *T) ([]byte, error)
	userData *T
}

func NewTcpServer[T any](port int, keyFileName string, label string, userData *T,
	handler func([]byte, *T) ([]byte, error)) (*TcpServer[T], error) {
	key, err := crypto.LoadRSAPrivateKey(keyFileName)
	if err != nil {
		return nil, err
	}
	return &TcpServer[T]{
		port:     port,
		key:      key,
		label:    []byte(label),
		handler:  handler,
		userData: userData,
	}, nil
}

func (s *TcpServer[T]) Start() error {
	addr := net.TCPAddr{Port: s.port}
	l, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		return err
	}
	defer func() { _ = l.Close() }()
	log.Printf("TCP server started on port %d\n", s.port)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			return err
		}
		// Handle connections in a new goroutine.
		go s.handleTcp(conn, l)
	}
}

func logTcpRequest(addr net.Addr, prefix string) {
	log.Printf("%v TCP request from address %s\n", prefix, addr.String())
}

func (s *TcpServer[T]) handleTcp(conn net.Conn, l *net.TCPListener) {
	defer func() { _ = conn.Close() }()
	logTcpRequest(conn.RemoteAddr(), "[Start]")
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	if err != nil {
		log.Printf("conn.Read error %v\n", err.Error())
		return
	}
	decrypted, err := rsa.DecryptOAEP(sha256.New(), nil, s.key, buf[:reqLen], s.label)
	if err != nil {
		log.Printf("rsa.DecryptOAEP error %v\n", err.Error())
		return
	}
	if len(decrypted) <= 44 {
		log.Printf("wrong decoded data length %v\n", err.Error())
		_ = l.Close()
		return
	}
	response, err := s.handler(decrypted[44:], s.userData)
	aesKey := decrypted[:32]
	aesNonce := decrypted[32:44]
	if err != nil {
		log.Printf("handler error %v\n", err.Error())
		sendResponse(conn, aesKey, aesNonce, ERROR, []byte(err.Error()))
	} else if response != nil {
		sendResponse(conn, aesKey, aesNonce, OK, response)
	}
	logTcpRequest(conn.RemoteAddr(), "[Done]")
}

func sendResponse(conn net.Conn, key []byte, nonce []byte, responseType uint8, responseData []byte) {
	aes, err := crypto.NewAesGcm(key)
	if err != nil {
		log.Printf("crypto.NewAesGcm error %v\n", err.Error())
		return
	}
	response := append([]byte{responseType}, responseData...)
	sha := sha256.New()
	hash := sha.Sum(response)
	encrypted := aes.EncryptWithNonce(append(response, hash...), nonce)
	_, err = conn.Write(encrypted)
	if err != nil {
		log.Printf("conn.Write error %v\n", err.Error())
	}
}
