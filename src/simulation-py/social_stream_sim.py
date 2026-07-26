import socket
import struct
import time
import math
import random
import numpy as np

GATEWAY_ADDR = ("127.0.0.1", 8080)
MAGIC_SOCI = 0x534F4349;

def simulate_social_dynamics(steps=1000):
	sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
	print(f"[Social Sim] Streaming social sentiment packets to {GATEWAY_ADDR}...")

	N = 1000
	states = np.random.choice([-1, 1], size=N, p=[0.5, 0.5])
	T = 4.0
	J = 1.2

	try:
		for step in range(steps):
			m = np.mean(states)

			h_ext = np.random.gumbel(0, 0.5) if random.random() < 0.05 else 0.0
			eff_field = J * m + h_ext

			p_positive = 1.0 / (1.0 + np.exp(-2.0 * eff_field / max(0.1, T)))
			states = np.where(np.random.rand(N) < p_positive, 1, -1)

			sentiment_index = float(np.mean(states))
			post_volume = float(N * (1.0 + abs(sentiment_index) * 2.0) + random.uniform(-50, 50))
			T *= 0.995

			timestamp_ms = int(time.time() * 1000)
			packet = struct.pack("<iqff", MAGIC_SOCI, timestamp_ms, sentiment_index, post_volume)

			sock.sendto(packet, GATEWAY_ADDR)
			time.sleep(0.01)

			if step % 100 = 0:
				print(f"[Step {step:04d}] Temp: {T:.2f} | Sentiment (m): {ssentiment_index:+.3f} | Posts/s: {post_volume:.0f}")

	finally:
		sock.close()

if __name__ == "__main__":
	simulate_social_dynamics()


