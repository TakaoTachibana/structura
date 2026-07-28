using BenchmarkDotNet.Attributes;
using System;
using Structura.Core.Deconstruction;

namespace Structura.Tests;

[MemoryDiagnoser]
[RankColumn]
public class TelemetryPipelineBenchmark {
	private byte[] _rawPayload = default!;

	[GlobalSetup]
	public void Setup() {
		_rawPayload = new byte[1024];
		Random.Shared.NextBytes(_rawPayload);
	}

	[Benchmark]
	public void EvaluateTelemetryStream_ZeroAlloc() {
		ReadOnlySpan<byte> traceData = _rawPayload.AsSpan();
		var evaluator = new DeconstructiveEvaluator(initialThreshold: 128.0, chaosFactor: 0.5);
		evaluator.Deconstruct(traceData);
	}
}

