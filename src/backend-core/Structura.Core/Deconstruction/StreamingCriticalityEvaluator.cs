using System;
using System.Runtime.CompilerServices;

namespace Structura.Core.Deconstruction;

public ref struct StreamingCriticalityEvaluator {
	private readonly double _varianceThreshold;
	private readonly double _autocorrThreshold;

	public StreamingCriticalityEvaluator(double varianceThreshold = 1.5, double autocorrThreshold = 0.7) {
		_varianceThreshold = varianceThreshold;
		_autocorrThreshold = autocorrThreshold;
	}

	[MethodImpl(MethodImplOptions.AggressiveInlining)]
	public readonly bool Evaluate(ReadOnlySpan<byte> buffer, out double variance, out double autocorr) {
		if (buffer.Length < 2) {
			variance = 0;
			autocorr = 0;
			
			return false;
		}

		double sum = 0;
		for (int i = 0; i < buffer.Length; i++) {
			sum += buffer[i];
		}
		double mean = sum / buffer.Length;

		double varianceSum = 0;
		double covarianceSum = 0;

		double prevDev = buffer[0] - mean;
		varianceSum += prevDev * prevDev;

		for (int i = 1; i < buffer.Length; i++) {
			double currDev = buffer[i] - mean;
			varianceSum += currDev * currDev;
			covarianceSum += prevDev * currDev;
			prevDev = currDev;
		}

		variance = varianceSum / buffer.Length;
		autocorr = variance > 1e-9 ? (covarianceSum / buffer.Length) / variance : 0;

		return variance > _varianceThreshold && autocorr > _autocorrThreshold;
	}
}

