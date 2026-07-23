package ebpf

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"log"
	"strconv"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"gateway-go/internal/listener"
)

type socketEvent struct {
	Pid uint32
	Comm [16]byte
	TimestampNs uint64
}

type EBPFListener struct {
	objs socketTraceObjects
	kp link.Link
}

func NewEBPFListener() *EBPFListener {
	return &EBPFListener{}
}

func (l *EBPFListener) Name() string {
	return "eBPF-Kernel-Listener"
}

func (l *EBPFListener) Start(ctx context.Context, dataCh chan<- listener.TelemetryData) error {
	log.Printf("[%s] Loading BPF objects into Linux kernel...", l.Name())

	if err := loadSocketTraceObjects(&l.objs, nil); err != nil {
		return err 
	}

	kp, err := link.Kprobe("tcp_v4_connect", l.objs.KprobeTcpV4Connect, nil)
	if err != nil {
		l.objs.Close()
		return err
	}
	l.kp = kp

	rd, err := ringbuf.NewReader(l.objs.Events)
	if err != nil {
		l.kp.Close()
		l.objs.Close()
		return err
	}

	log.Printf("[%s] Successfully attached kprobe to tcp_v4_connect! Listening for real socket event...", l.Name())


	go func() {
		defer func() {
			rd.Close()
			l.kp.Close()
			l.objs.Close()
		}()

		go func() {
			<-ctx.Done()
			rd.Close()
		}()

		for {
			record, err := rd.Read()
			if err != nil {
				if errors.Is(err, ringbuf.ErrClosed) {
					log.Printf("[%s] RingBuffer closed, exiting reader...", l.Name())
					return
				}
				log.Printf("[%s] RingBuffer read error: %v", l.Name(), err)
				continue
			}

			var event socketEvent
			if err := binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &event); err != nil {
				log.Printf("[%s] Failed to parse ringbuf record: %v", l.Name(), err)
				continue
			}

			commStr := string(bytes.TrimRight(event.Comm[:], "\x00"))


			dataCh <- listener.TelemetryData {
				Timestamp: int64(event.TimestampNs),
				Source: "ebpf",
				Metric: "kernel_tcp_connect_event",
				Value: 1.0,
				Tags: map[string]string {
					"pid": strconv.Itoa(int(event.Pid)),
					"comm": commStr,
				},
			}
		}
	}()

	return nil
}

func (l *EBPFListener) Stop() error {
	return nil
}

