using System;
using Xunit;
using Structura.Core.Deconstruction;

namespace Structura.Tests;

public class SocialCriticalityTests {
	[Fact]
	public void Evaluate_Should_Detect_CriticalSlowingDown_Before_SocialPhaseTransition() {
		ReadOnlySpan<byte> normalSentimentStream = stackalloc byte[] { 120, 130, 122, 128, 121, 125, 127, 123, 129, 124 };
		ReadOnlySpan<byte> criticalPrecursorStream = stackalloc byte[] { 10, 20, 30, 40, 50, 210, 220, 230, 240, 250 };

		var evaluator = new StreamingCriticalityEvaluator(varianceThreshold: 1000.0, autocorrThreshold: 0.4);

		bool isNormalCritical = evaluator.Evaluate(normalSentimentStream, out double normalVar, out double nomalAutocorr);
		Assert.False(isNormalCritical, "Normal sentiment data should not trigger a criticality alert.");

		bool isPrecursorCritical = evaluator.Evaluate(criticalPrecursorStream, out double precursorVar, out double precursorAutocorr);
		Assert.True(isPrecursorCritical, "Phase transition precursor (critical slowing down must trigger an alert.");
		Assert.True(precursorVar > 1000.0, $"Variance ({precursorVar:F2}) must exceed the defined threshold.");
		Assert.True(precursorAutocorr > 0.4, $"Autocorrelation ({precursorAutocorr:F2}) must exceed the defined threshold.");
	}

	[Fact]
	public void Evaluate_SocialSentiment_Must_Be_Zero_Allocation() {
		ReadOnlySpan<byte> sentimentStream = stackalloc byte[] { 100, 110, 105, 200, 210, 205, 220, 215 };
		var evaluator = new StreamingCriticalityEvaluator(varianceThreshold: 500.0, autocorrThreshold: 0.5);

		evaluator.Evaluate(sentimentStream, out _, out _);

		long allocatedBefore = GC.GetAllocatedBytesForCurrentThread();

		for (int i = 0; i < 100_000; i++) {
			evaluator.Evaluate(sentimentStream, out _, out _);
		}

		long allocatedAfter = GC.GetAllocatedBytesForCurrentThread();
		long totalAllocated = allocatedAfter - allocatedBefore;

		Assert.Equal(0, totalAllocated);
	}
}

