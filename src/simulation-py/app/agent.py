import asyncio
import random
import time
import logging
from app.models import TelemetryData
from app.client import GatewayClient

class Agent:
	def __init__(self, agent_id: str, client: GatewayClient):
		self.agent_id = agent_id
		self.client = client
		self.is_running = False
	
	async def run(self, interval_sec: float = 1.0):
		self.is_running = True
		logging.info(f"Agent [{self.agent_id}] started.")

		while self.is_running:
			telemetry = TelemetryData(
					agent_id=self.agent_id,
					timestamp=time.time(),
					status="ACTIVE",
					metrics={
						"cpu_usage": round(random.uniform(10.0, 90.0), 2),
						"memory_mb": random.randint(128, 1024),
					}
			)

			await self.client.send_telemetry(telemetry.to_dict())
			jitter = random.uniform(-0.1, 0.1)
			await asyncio.sleep(max(0.1, interval_sec + jitter))

	def stop(self):
		self.is_running = False



