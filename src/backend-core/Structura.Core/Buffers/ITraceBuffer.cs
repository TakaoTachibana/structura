using System;

namespace Structura.Core.Buffers;

public interface ITraceBuffer : IDisposable {
	void Write(ReadOnlySpan<byte> data);
	ReadOnlySpan<byte> Read();
	void Clear();
}

