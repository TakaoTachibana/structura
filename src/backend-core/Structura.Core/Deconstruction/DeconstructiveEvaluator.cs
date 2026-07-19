using System;

namespace Structura.Core.Deconstruction;

public ref struct DeconstructiveEvaluator {
	private double _logosThreshold;
	private readonly double _chaosFactor;

	public DeconstructiveEvaluator(double initialThreshold, double chaosFactor) {
		_logosThreshold = initialThreshold;
		_chaosFactor = chaosFactor;
	}

	public readonly double CurrentThreshold => _logosThreshold;

	public void Deconstruct(ReadOnlySpan<byte> traceData) {
		if (traceData.IsEmpty) {
			return;
		}
		double signalSum = 0;
		for (int i = 0; i < traceData.Length; i++) {
			signalSum += traceData[i];
		}
		double meanSignal = signalSum / traceData.Length;

		if (meanSignal > _logosThreshold) {
			_logosThreshold -= (meanSignal - _logosThreshold) * _chaosFactor * 0.01;
		}
	}
}

