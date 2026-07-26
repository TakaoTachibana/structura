using System;
using Xunit;
using Structura.Core.Deconstruction;

namespace Structura.Tests;

public class CriticalityEvaluatorTests {
	[Fact]
	public void Evaluate_Should_TriggerAlert_On_HighVariance_And_Autocorrelation() {
		ReadOnlySpan<byte> traceData = stackalloc byte[] { 10, 12, 11, 200, 210, 205, 220, 215, 225 };
		var evaluator = new StreamingCriticalityEvaluator(varianceThreshold: 100.0, autocorrThreshold: 0.5);

		bool isCritical = evaluator.Evaluate(traceData, out double variance, out double autocorr);

		Assert.True(isCritical);
		Assert.True(variance > 100.0);
		Assert.True(autocorr > 0.5);
	}

	[Fact]
	public void Evaluate_Must_Be_Zero_Allocation() {
		ReadOnlySpan<byte> traceData = stackalloc byte[] { 50, 52, 48, 51, 49, 53, 50, 52 };
		var evaluator = new StreamingCriticalityEvaluator();

		evaluator.Evaluate(traceData, out _, out _);

		long allocatedBefore = GC.GetAllocatedBytesForCurrentThread();

		for (int i = 0; i < 100_000; i++) {
			evaluator.Evaluate(traceData, out _, out _);
		}

		long allocatedAfter = GC.GetAllocatedBytesForCurrentThread();
		long totalAllocated = allocatedAfter - allocatedBefore;

		Assert.Equal(0, totalAllocated);
	}
}

