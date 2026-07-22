import pandas as pd
import matplotlib.pyplot as plt

df = pd.read_parquet("fanaticism_simulation_log.parquet")

fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(14, 5))

color_m = 'tab:red'
ax1.set_xlabel('Step', fontsize=11)
ax1.set_ylabel('Order Parameter (m)', color=color_m, fontsize=11)
ax1.plot(df['step'], df['order_parameter_m'], color=color_m, linewidth=2.5, label='Order Parameter (m)')
ax1.tick_params(axis='y', labelcolor=color_m)
ax1.set_ylim(-1.1, 1.1)
ax1_twin = ax1.twinx()
color_t = 'tab:blue'
ax1_twin.set_ylabel('Tempertature / Stress (T)', color=color_t, fontsize=11)
ax1_twin.plot(df['step'], df['temperature_T'], color=color_t, linestyle='--', linewidth=1.5, label='Temperature (T)')
ax1.set_title('Time Series: Jump to Fanaticism', fontsize=12, fontweight='bold')
ax1.grid(True, alpha=0.3)

ax2.plot(df['temperature_T'], df['order_parameter_m'], marker='o', markersize=4, color='purple', linestyle='-')
ax2.invert_xaxis()
ax2.set_xlabel('Social Stress / Temperature (T) [Decaying]', fontsize=11)
ax2.set_ylabel('Order Parameter (m)', fontsize=11)
ax2.set_title('Phase Diagram: Criticality Boundary', fontsize=12, fontweight='bold')
ax2.set_ylim(-1.1, 1.1)
ax2.grid(True, alpha=0.3)

plt.tight_layout()
plt.savefig("phase_transition_plot.png", dpi=300)
print("Graph saved successfully as 'phase_transition_plot.png'.")

