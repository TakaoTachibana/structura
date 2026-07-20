import asyncio
import logging
from app.client import GatewayClient
from app.agent import Agent

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")

HOST = "127.0.0.1"
PORT = 8080
NUM_AGENTS = 5

async def main():
	client = GatewayClient(host=HOST, port=PORT)
	agents = [Agent(agent_id=f"agent-{i:03d}", client=client) for i in range(1, NUM_AGENTS + 1)]
	try:
		agent_tasks = [asyncio.create_task(agent.run(interval_sec=2.0)) for agent in agents]
		await asyncio.gather(*agent_tasks)
	except asyncio.CancelledError:
		logging.info("Shutting down agents...")
	finally:
		await client.close()

if __name__ == "__main__":
	try:
		asyncio.run(main())
	except KeyboardInterrupt:
		print("\nSimulation stopped by user.")

