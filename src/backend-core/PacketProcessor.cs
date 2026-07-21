using System;
using System.Diagnostics;
using System.Runtime.InteropServices;

namespace Structura.Core;

public static class PacketProcessor {
	private const float CpuCriticalThreshold = 85.0f;
	private const float MemoryCriticalThreshold = 90.0f;

	private static long _packetCount = 0;
	private static long _alertCount = 0;
	private static readonly Stopwatch _stopwatch = Stopwatch.StartNew();
	private static long _lastReportTimeMs = 0;

	public static void ProcessBuffer(ReadOnlySpan<byte> buffer) {
		int packetSize = Marshal.SizeOf<TelemetryPacket>();

		if (buffer.Length < packetSize) {
			return;
		}

		ref readonly var packet = ref MemoryMarshal.AsRef<TelemetryPacket>(buffer[..packetSize]);

		if (!packet.IsValid) {
			return;
		}

		_packetCount++;
		bool isCpuAlert = packet.CpuUsage >= CpuCriticalThreshold;
		bool isMemAlert = packet.MemoryUsage >= MemoryCriticalThreshold;

		if (isCpuAlert || isMemAlert) {
			_alertCount++;
		}

		long currentMs = _stopwatch.ElapsedMilliseconds;
		if (currentMs - _lastReportTimeMs >= 1000) {
			double elapsedSeconds = (currentMs - _lastReportTimeMs) / 1000.0;
			double tps = _packetCount / elapsedSeconds;

			Console.WriteLine($"[PREF] Speed: {tps:N0} TPS (pkts/sec) | Alerts Detected: {_alertCount:N0} / {_packetCount:N0}");

			_packetCount = 0;
			_alertCount = 0;
			_lastReportTimeMs = currentMs;
		}
	}
}




