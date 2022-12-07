package mideaac

import (
	"log"
	"net"
	"os"
	"time"
)

type Lan struct {
	DeviceIP   string
	DeviceID   uint64
	DevicePort string
	Security   security
	retries    int
	socket     net.Conn
	token      []byte
	key        []byte
	timestamp  time.Time
	tcpKey     []byte
	local      string
	remote     string
	buffer     []byte
}

func newLan(deviceIP string, deviceID uint64, devicePort string) Lan {
	return Lan{
		DeviceIP:   deviceIP,
		DeviceID:   deviceID,
		DevicePort: devicePort,
		Security:   NewSecurity(),
		retries:    0,
		socket:     nil,
		token:      []byte{},
		key:        []byte{},
		timestamp:  time.Now(),
		tcpKey:     nil,
		local:      "",
		remote:     deviceIP + ":" + devicePort,
	}
}

func (lan *Lan) connect() (err error) {
	if lan.socket == nil {
		lan.disconnect()
		// TODO log connecting
		lan.buffer = nil
		d := net.Dialer{Timeout: time.Second * 2}
		if lan.socket, err = d.Dial("tcp", lan.remote); err != nil {
			return err
		}
	}
	return nil
}

func (lan *Lan) disconnect() {
	if lan.socket != nil {
		lan.socket.Close()
		lan.socket = nil
		lan.tcpKey = nil
	}
}

// TODO
func (lan Lan) GetSocketInfo() string {
	return ""
}

func (lan *Lan) Request(message []byte) ([]byte, bool) {
	if err := lan.connect(); err != nil {
		return nil, false
	}

	if _, err := lan.socket.Write(message); err != nil {
		lan.disconnect()
		lan.retries += 1
		return nil, true
	}

	response := make([]byte, 1024)
	if n, err := lan.socket.Read(response); err != nil || n == 0 {
		if os.IsTimeout(err) {
			lan.retries += 1
			return nil, true
		}
		lan.disconnect()
		lan.retries += 1
		return nil, true
	}
	return response, true
}

func (lan *Lan) Authenticate(token []byte, key []byte) bool {
	lan.token = token
	lan.key = key
	request := lan.Security.Encode8370(token, MSGTYPE_HANDSHAKE_REQUEST)

	response, _ := lan.Request(request)
	response = response[8:72]

	tcpKey, success := lan.Security.TcpKey(response, key)
	if success {
		lan.tcpKey = tcpKey
		log.Println("got TCP key")
		time.Sleep(1 * time.Second)
	} else {
		log.Println("auth fail")
	}
	return success
}

func (lan *Lan) ApplianceTransparentSend8370(data []byte, msgType byte) [][]byte {
	if lan.socket == nil || lan.tcpKey == nil {
		lan.disconnect()
		if lan.Authenticate(lan.token, lan.key) {
			return nil
		}
	}
	originalData := make([]byte, len(data))
	copy(originalData, data)
	data = lan.Security.Encode8370(data, msgType)

	time.Sleep(time.Second * time.Duration(lan.retries))
	response, b := lan.Request(data)
	if string(response[8:13]) == "ERROR" {
		lan.disconnect()
		return [][]byte{[]byte("ERROR")}
	}
	if response == nil && lan.retries < 2 && b {
		packets := lan.ApplianceTransparentSend8370(originalData, msgType)
		lan.retries = 0
		return packets
	}
	var responses [][]byte
	responses, lan.buffer = lan.Security.Decode8370(append(lan.buffer, response...))
	packets := [][]byte{}

	for _, response := range responses {
		if len(response) > 40+16 {
			response = lan.Security.AesDecrypt(response[40 : len(response)-16])
		}
		if len(response) > 10 {
			packets = append(packets, response)
		}
	}
	return packets
}

func (lan *Lan) ApplianceTransparentSend(data []byte) [][]byte {
	time.Sleep(time.Second * time.Duration(lan.retries))
	response, b := lan.Request(data)
	if response == nil && lan.retries < 2 && b {
		packets := lan.ApplianceTransparentSend(data)
		lan.retries = 0
		return packets
	}
	packets := [][]byte{}
	if response == nil {
		return packets
	}

	dlen := len(response)
	if response[0] == 0x5A && response[1] == 0x5A && dlen > 5 {
		i := 0
		for i < dlen {
			size := response[i+4]
			data = lan.Security.AesDecrypt(response[i : i+int(size)])
		}
	} else if response[0] == 0xAA && dlen > 2 {
		i := 0
		for i < dlen {
			size := response[i+1]
			data = response[i : i+int(size)+1]

			if len(data) > 10 {
				packets = append(packets, data)
			}
			i += int(size) + 1
		}
	} else {
		log.Println("unknown response")
	}

	return packets
}
