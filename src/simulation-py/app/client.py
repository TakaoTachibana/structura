import socket
import json
import logging

class GatewayClient:
	def __init__(self, host: str = "127.0.0.1", port: int = 8080):
		self.host = host
		self.port = port
		self.sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

	async def send_telemetry(self, data: dict) -> bool:
		try:
			payload = json.dumps(data).encode('utf-8')
			self.sock.sendto(payload, (self.host, self.port))
			return True
		except Exception as e:
			logging.error(f"Failed to send UDP packet: {e}")
			return False
	
	async def close(self):
		self.sock.close()
