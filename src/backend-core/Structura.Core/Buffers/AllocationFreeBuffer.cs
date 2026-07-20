using System;
using System.Buffers;

namespace Structura.Core.Buffers;

public sealed class AllocationFreeBuffer : ITraceBuffer {
	private readonly int _capacity;
	private byte[]? _rentedArray;
	private int _writePosition;
	private int _readPosition;
	private int _count;

	public AllocationFreeBuffer(int capacity) {
		_capacity = capacity;
		_rentedArray = ArrayPool<byte>.Shared.Rent(capacity);
	}

	public void Write(ReadOnlySpan<byte> data) {
		ObjectDisposedException.ThrowIf(_rentedArray is null, this);

		if (data.Length > _capacity) {
			throw new ArgumentOutOfRangeException(nameof(data), "datasize exceeded.");
		}

		int availableSpace = _capacity - _count;
		if (data.Length > availableSpace) {
			_readPosition = (_readPosition + (data.Length - availableSpace)) % _capacity;
			_count = _capacity - data.Length;
		}

		Span<byte> bufferSpan = _rentedArray.AsSpan();

		if (_writePosition + data.Length <= _capacity) {
			data.CopyTo(bufferSpan.Slice(_writePosition, data.Length));
			_writePosition = (_writePosition + data.Length) % _capacity;
		} else {
			int firstPartSize = _capacity - _writePosition;
			data[..firstPartSize].CopyTo(bufferSpan.Slice(_writePosition, firstPartSize));
			data[firstPartSize..].CopyTo(bufferSpan[.. (data.Length - firstPartSize)]);
			_writePosition = data.Length - firstPartSize;
		}

		_count += data.Length;
	}

	public ReadOnlySpan<byte> Read() {
		ObjectDisposedException.ThrowIf(_rentedArray is null, this);
		if (_count == 0) {
			return ReadOnlySpan<byte>.Empty;
		}

		Span<byte> bufferSpan = _rentedArray.AsSpan();
		if (_readPosition + _count <= _capacity) {
			return bufferSpan.Slice(_readPosition, _count);
		} else {
			return bufferSpan.Slice(_readPosition, _capacity - _readPosition);
		}
	}

	public void Clear() {
		_writePosition = 0;
		_readPosition = 0;
		_count = 0;
	}

	public void Dispose() {
		if (_rentedArray != null) {
			ArrayPool<byte>.Shared.Return(_rentedArray);
			_rentedArray = null;
		}
	}
}


