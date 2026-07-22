const PACKET_SIZE = 24;
const BUFFER_CAPACITY = 10000;
const TOTAL_BUFFER_SIZE = 4 + (BUFFER_CAPACITY * PACKET_SIZE);

async function init() {
	console.log('[Main] Initializing Structura Zero-Copy Engine...');

	if (typeof SharedArrayBuffer === 'undefined') {
		alert('SharedArrayBuffer is not supported. Check Cross-Origin Isolation headers.');
		return;
	}

	const sharedBuffer = new SharedArrayBuffer(TOTAL_BUFFER_SIZE);
	const worker = new Worker(new URL('./worker.ts', import.meta.url), { type: 'module' });
	worker.postMessage({ type: 'INIT', buffer: sharedBuffer });
	const atomicIndex = new Int32Array(sharedBuffer, 0, 1);
	let lastIndex = 0;

	setInterval(() => {
		const  currentIndex = Atomics.load(atomicIndex, 0);
		const tps = currentIndex - lastIndex;
		lastIndex = currentIndex;

		const statusEl = document.getElementById('tps-counter');
		if (statusEl) {
			statusEl.innerText = `Current Ingestion Rate: ${tps.toLocaleString()} TPS (Total Received: ${currentIndex.toLocaleString()})`;
		}
	}, 1000);
}

window.addEventListener('DOMContentLoaded', init);

