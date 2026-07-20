import time
from dataclasses import dataclass, asdict

@dataclass
class TelemetryData:
	agent_id: str
	timestamp: float
	status: str
	metrics: dict

	def to_dict(self) -> dict:
		return asdict(self)
