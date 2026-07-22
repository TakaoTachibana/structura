let sharedBuffer: SharedArrayBuffer;
let uint8Array: Uint8Array;
let atomicIndex: Int32Array;

const PACKET_SIZE = 24;
const BUFFER_CAPACITY = 10000;

self.onmessage = (event: MessageEvent) => {
	if (event.data.type === 'INIT') {
		sharedBuffer = event.data.buffer;
		atomicIndex = new Int32Array(sharedBuffer, 0, 1);
		uint8Array = new Uint8Array(sharedBuffer, 4);

		console.log('[Worker] SharedArrayBuffer initialized successfully.');
		connectWebSocket();
	}
};

function connectWebSocket() {
	const ws = new WebSocket('ws://localhost:8080/ws/');
	ws.binaryType = 'arraybuffer';
	ws.onopen = () => {
		console.log('[Worker] Connected to C# Binary WebSocket Stream!');
	};

	ws.onmessage = (event: MessageEvent) => {
		if (event.data instanceof ArrayBuffer) {
			const incomingBuffer = new Uint8Array(event.data);
			const totalBytes = incomingBuffer.byteLength;

			if (incomingBuffer.byteLength !== PACKET_SIZE) {
				return;
			}

			const packetCount = Math.floor(totalBytes / PACKET_SIZE);

			for (let i  = 0; i < packetCount; i++) {
				const offset = i * PACKET_SIZE;
				const packetData = incomingBuffer.subarray(offset, offset + PACKET_SIZE);
				const currentIndex = Atomics.load(atomicIndex, 0);
				const targetOffset = (currentIndex % BUFFER_CAPACITY) * PACKET_SIZE;
				uint8Array.set(incomingBuffer, targetOffset);
				Atomics.add(atomicIndex, 0, 1);
			}
		}
	};

	ws.onerror = (err) => console.error('[Worker] WS Error:', err);
	ws.onclose = () => setTimeout(connectWebSocket, 1000);
}

