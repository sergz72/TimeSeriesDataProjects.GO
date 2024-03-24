package network

import (
	"TimeSeriesData/crypto"
	"crypto/rsa"
	"crypto/sha256"
	"log"
	"net"
	"sync"
)

/*

Client message structure (RSA encoded, maximum request data length ~ 462 bytes for RSA 4096):
|AES key - 32 bytes|AES gcm nonce - 12 bytes|Request data|

Server message structure:
|Response + sha256 of response data encrypted with AES-GCM|

*/

const (
	OK       uint8 = 0
	OK_BZIP2 uint8 = 1
	ERROR    uint8 = 0x7F
)

type TcpServer[T any] struct {
	port     int
	key      *rsa.PrivateKey
	label    []byte
	handler  func([]byte, *T) ([]byte, error, bool)
	userData *T
	listener *net.TCPListener
}

func NewTcpServer[T any](port int, keyFileName string, label string, userData *T,
	handler func([]byte, *T) ([]byte, error, bool)) (*TcpServer[T], error) {
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

func (s *TcpServer[T]) Terminate() {
	l := s.listener
	s.listener = nil
	if l != nil {
		_ = l.Close()
	}
}

func (s *TcpServer[T]) Start() error {
	addr := net.TCPAddr{Port: s.port}
	var err error
	s.listener, err = net.ListenTCP("tcp", &addr)
	if err != nil {
		return err
	}
	log.Printf("TCP server started on port %d\n", s.port)
	var wg sync.WaitGroup
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Error accepting: %v\nWaiting for all goroutines to finish...\n", err.Error())
			wg.Wait()
			s.Terminate()
			log.Println("TCP server terminated")
			return err
		}
		// Handle connections in a new goroutine.
		go func() {
			wg.Add(1)
			s.handleTcp(conn)
			wg.Done()
		}()
	}
}

func logTcpRequest(addr net.Addr, prefix string) {
	log.Printf("%v TCP request from address %s\n", prefix, addr.String())
}

func (s *TcpServer[T]) handleTcp(conn net.Conn) {
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
		s.Terminate()
		return
	}
	response, err, terminate := s.handler(decrypted[44:], s.userData)
	if err != nil {
		log.Printf("handler error %v\n", err.Error())
	}
	if terminate {
		log.Println("fatal error, terminating tcp server")
		s.Terminate()
		return
	}
	aesKey := decrypted[:32]
	aesNonce := decrypted[32:44]
	if err != nil {
		sendResponse(conn, aesKey, aesNonce, ERROR, []byte(err.Error()))
	} else if response != nil {
		compressed, err := bzipData(response)
		if err != nil || len(compressed) >= len(response) {
			sendResponse(conn, aesKey, aesNonce, OK, response)
		} else {
			sendResponse(conn, aesKey, aesNonce, OK_BZIP2, compressed)
		}
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
	sha.Write(response)
	hash := sha.Sum(nil)
	encrypted := aes.EncryptWithNonce(append(response, hash...), nonce)
	_, err = conn.Write(encrypted)
	if err != nil {
		log.Printf("conn.Write error %v\n", err.Error())
	}
}
