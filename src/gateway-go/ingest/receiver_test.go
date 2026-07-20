package ingest

import (
	"testing"
)

func BenchmarkPacketReceiver_PoolAllocation(b *testing.B) {
	receiver := NewPacketReceiver("127.0.0.1")
	testData := []byte("LOG_LEVEL=WARN msg='Logos threshold shifting initiated'")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		bufPtr := receiver.pool.Get().(*[]byte)
		buf := *bufPtr

		copy(buf, testData)

		receiver.pool.Put(bufPtr)
	}
}

