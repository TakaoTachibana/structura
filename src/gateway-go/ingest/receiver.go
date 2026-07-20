package ingest

import (
	"context"
	"fmt"
	"net"
	"sync"
)

const MaxUDPBufferSize = 65535

type PacketReceiver struct {
	addr string
	pool *sync.Pool
}

func NewPacketReceiver(addr string) *PacketReceiver {
	return &PacketReceiver {
		addr : addr,
		pool: &sync.Pool {
			New: func() any {
				b := make([]byte, MaxUDPBufferSize)
				return &b
			},
		},
	}
}

func (r *PacketReceiver) Start(ctx context.Context, handler func([]byte)) error {
	pc, err := net.ListenPacket("udp", r.addr)
	if err != nil {
		return fmt.Errorf("failed to listen UDP: %w", err)
	}
	defer pc.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			bufPtr := r.pool.Get().(*[]byte)
			buf := *bufPtr

			n, _, err := pc.ReadFrom(buf)
			if err != nil {
				r.pool.Put(bufPtr)
				continue
			}
			handler(buf[:n])
			r.pool.Put(bufPtr)
		}
	}
}

