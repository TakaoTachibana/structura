using System.Runtime.InteropServices;

namespace Structura.Core;

[StructLayout(LayoutKind.Sequential, Pack = 1)]
public readonly struct TelemetryPacket {
	public const uint ValidMagic = 0x53545255;

	public readonly uint Magic;
	public readonly uint AgentId;
	public readonly long Timestamp;
	public readonly float CpuUsage;
	public readonly float MemoryUsage;

	public bool IsValid => Magic == ValidMagic;
}

