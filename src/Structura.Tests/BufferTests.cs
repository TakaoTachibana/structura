using System;
using Xunit;
using Structura.Core.Buffers;
using Structura.Core.Deconstruction;

namespace Structura.Tests;

public class BufferTests {
	[Fact]
	public void Pipeline_Shold_Be_Zero_Allocation() {
		using var buffer = new AllocationFreeBuffer(1024);
		ReadOnlySpan<byte> incomingChaos = stackalloc byte[] { 120, 150, 200, 180, 250 };
		long allocatedBefore = GC.GetAllocatedBytesForCurrentThread();

		for (int i = 0; i < 100000; i++) {
			buffer.Write(incomingChaos);
			ReadOnlySpan<byte> readTrace = buffer.Read();

			var evaluator = new DeconstructiveEvaluator(100.0, 0.5);
			evaluator.Deconstruct(readTrace);

			buffer.Clear();
		}

		long allocatedAfter = GC.GetAllocatedBytesForCurrentThread();
		long totalAllocatedBytes = allocatedAfter - allocatedBefore;

		Assert.Equal(0, totalAllocatedBytes);
	}
}


