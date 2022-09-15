package udp

import (
	"fmt"
	"net"
	"sync"

	"go.uber.org/zap"

	"github.com/goforbroke1006/net-conn-proxy/internal/common"
)

func NewProxyDataGram(downstreamConn *net.UDPConn, upstreamAddr string, bufferSizeArg uint64) *proxyUDP {
	upstreamConn, err := net.Dial("udp", upstreamAddr)
	if err != nil {
		zap.L().Error("dial upstream", zap.Error(err))
		return nil
	}

	return &proxyUDP{
		bufferSizeArg:  bufferSizeArg,
		upstreamConn:   upstreamConn,
		downstreamConn: downstreamConn,

		clientAddresses: make(map[string]*net.UDPAddr, 128),
	}
}

type proxyUDP struct {
	bufferSizeArg  uint64
	upstreamConn   net.Conn
	downstreamConn *net.UDPConn

	clientAddresses   map[string]*net.UDPAddr
	clientAddressesMx sync.RWMutex
}

func (p *proxyUDP) Run() {
	go p.clientsToUpstreamLoop()
	go p.upstreamToAllClientsLoop()
}

func (p *proxyUDP) clientsToUpstreamLoop() {
	var (
		buffer     = make([]byte, p.bufferSizeArg)
		readLen    int
		clientAddr *net.UDPAddr
		readErr    error
		writeLen   int
		writeErr   error
	)

	for {
		// downstream >>> buffer
		readLen, clientAddr, readErr = p.downstreamConn.ReadFromUDP(buffer[0:])

		if readLen > 0 {
			_, writeErr = p.upstreamConn.Write(buffer[:readLen])

			fmt.Printf("%s (%d) >>> %s (%d)\n%s\n%s\n",
				clientAddr.String(), readLen,
				p.upstreamConn.RemoteAddr().String(), writeLen,
				common.GetPrettyHexString(buffer[:readLen]),
				string(buffer[:readLen]),
			)
		}

		if readErr != nil {
			zap.L().Error("read from downstream", zap.Error(readErr))
		}

		if writeErr != nil {
			zap.L().Error("write to upstream", zap.Error(readErr))
		}
		_ = writeLen

		p.clientAddressesMx.Lock()
		p.clientAddresses[clientAddr.String()] = clientAddr
		p.clientAddressesMx.Unlock()
	}
}

func (p *proxyUDP) upstreamToAllClientsLoop() {
	var (
		buffer   = make([]byte, p.bufferSizeArg)
		readLen  int
		readErr  error
		writeLen int
		writeErr error
	)

	for {
		// buffer <<< upstream
		readLen, readErr = p.upstreamConn.Read(buffer[0:])

		if readLen > 0 {
			p.clientAddressesMx.RLock()
			for _, clientAddr := range p.clientAddresses {
				writeLen, writeErr = p.downstreamConn.WriteToUDP(buffer[:readLen], clientAddr)
			}

			fmt.Printf("%d client (%d) <<< %s (%d)\n%s\n%s\n",
				len(p.clientAddresses), readLen,
				p.upstreamConn.RemoteAddr().String(), writeLen,
				common.GetPrettyHexString(buffer[:readLen]),
				string(buffer[:readLen]),
			)
			p.clientAddressesMx.RUnlock()
		}

		if readErr != nil {
			zap.L().Error("read from upstream", zap.Error(readErr))
			break
		}

		if writeErr != nil {
			zap.L().Error("write to downstream", zap.Error(writeErr))
		}
	}
}
