using System.Runtime.InterpoServices;

namespace Structura.Core;

[StructLayout(Layout.Sequential, Pack = 1)]
public readonly struct SocialTelemetryPacket {
	public const uint ValidMagic = 0x534F4349;

	public readonly uint Magic;
	public readonly long TimestampMs;
	public readonly float SentimentIndex;
	public readonly float PostVolumeRate;

	public SocialTelemetryPacket(long timestampMs, float sentimentIndex, float postVolumeRate) {
		Magic = ValidMagic;
		TimestampMs = timestampMs;
		SentimentIndex = sentimentIndex;
		PostVolumeRate = postVolumeRate;
	}
	public bool isValid => Magic == ValidMagic;
}

