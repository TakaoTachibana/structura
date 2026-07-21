import asyncio
import random
import socket
import struct
import time

GATEWAY_ADDR = ("127.0.0.1", 8080)
MAGIC = 0x53545255

def create_telemetry_packet(agent_id: int) -> bytes:
	timestamp_ms = int(time.time() * 1000)
	cpu_usage = (
			random.uniform(86.0, 98.0)
			if random.random() < 0.2
			else random.uniform(20.0, 70.0)
	)
	mem_usage = (
			random.uniform(91.0, 99.0)
			if random.random() < 0.1
			else random.uniform(30.0, 60.0)
	)

	return struct.pack(
			"<IIqff", MAGIC, agent_id, timestamp_ms, cpu_usage, mem_usage
	)

async def run_simulator():
	sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
	print(
			f"[Python Simulator] Sending binary packets to Go Gateway ({GATEWAY_ADDR[0]}:{GATEWAY_ADDR[1]})..."
	)

	agent_id = 101
	packet_count = 0

	try:
		while True:
			packet = create_telemetry_packet(agent_id)
			sock.sendto(packet, GATEWAY_ADDR)

			packet_count += 1
			if packet_count % 100 == 0:
				print (
						f"[Python Simulator] Sent {packet_count} binary packets."
				)

			await asyncio.sleep(0)
	finally:
		sock.close()

if __name__ == "__main__":
	asyncio.run(run_simulator())





