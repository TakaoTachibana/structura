import numpy as np
import pandas as pd
from scipy.stats import gumbel_r

class CelebrityFanaticismSim:
	def __init__(self, num_agents=10000, alpha=0.3, gamma=1.5, mu_bias=0.25):
		self.N = num_agents
		self.alpha = alpha
		self.gamma = gamma
		self.mu_bias = mu_bias

		self.states = np.random.choice([-1, 1], size=self.N, p=[0.9, 0.1])

		self.A_i = 0.75
	
	def compute_celeb_influence(self, current_followers):
		L_i = gumbel_r.rvs(loc=0, scale=1.0)
		k_ratio = current_followers / self.N
		effective_quality = self.alpha * self.A_i + (1 - self.alpha) * max(0, L_i)
		h_i = effective_quality * (k_ratio**self.gamma)
		return h_i

	def run(self, steps=1000, J=1.2, initial_temp=5.0, cooling_rate=0.97):
		history = []
		T = initial_temp

		for step in range(steps):
			m = np.mean(self.states)
			followers = np.sum(self.states == 1)
			h_i = self.compute_celeb_influence(followers)
			eff_field = J * m + h_i + self.mu_bias
			p_fanatic = 1.0 / (1.0 + np.exp(-2.0 * eff_field / T))
			self.states = np.where(np.random.rand(self.N) < p_fanatic, 1, -1)
			history.append(
					{
						"step": step,
						"temperature_T": T,
						"order_parameter_m": m,
						"fanatic_ratio": followers / self.N,
						"celeb_influence_h": h_i,
					}
			)
			T *= cooling_rate
		return pd.DataFrame(history)

if __name__ == "__main__":
	sim = CelebrityFanaticismSim(num_agents=10000)
	df_result = sim.run(steps=1000, J=1.2, initial_temp=5.0)
	df_result.to_parquet("fanaticism_simulation_log.parquet")
	print("Simulation completed. Saved to 'fanaticism_simulation_log.parquet'.")
	print(df_result.tail(10))


			
