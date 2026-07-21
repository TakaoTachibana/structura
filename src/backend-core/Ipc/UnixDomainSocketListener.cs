using System;
using System.Buffers;
using System.IO;
using System.Net.Sockets;
using System.Threading;
using System.Threading.Tasks;
using Structura.Core.Buffers;

namespace Structura.Core.Ipc;

public sealed class UnixDomainSocketListener {
	private readonly string _socketPath;
	private readonly AllocationFreeBuffer _buffer;
	private Socket? _listenerSocket;

	public UnixDomainSocketListener(string socketPath, AllocationFreeBuffer buffer) {
		_socketPath = socketPath;
		_buffer = buffer;
	}

	public async Task StartAsync(CancellationToken cancellationToken) {
		if (File.Exists(_socketPath)) {
			File.Delete(_socketPath);
		}

		var endPoint = new UnixDomainSocketEndPoint(_socketPath);
		_listenerSocket = new Socket(AddressFamily.Unix, SocketType.Stream, ProtocolType.Unspecified);

		_listenerSocket.Bind(endPoint);
		_listenerSocket.Listen(128);

		Console.WriteLine($"[C# Core] IPC Listener started on {_socketPath}");

		try {
			while (!cancellationToken.IsCancellationRequested) {
				Socket clientSocket = await _listenerSocket.AcceptAsync(cancellationToken);
				_ = HandleClientAsync(clientSocket, cancellationToken);
			}
		} catch (OperationCanceledException) {
		} finally {
			Cleanup();
		}
	}

	private async Task HandleClientAsync(Socket clientSocket, CancellationToken cancellationToken) {
		using (clientSocket) {
			byte[] receiveBuffer = ArrayPool<byte>.Shared.Rent(65535);
			try {
				while (!cancellationToken.IsCancellationRequested) {
					int bytesRead = await clientSocket.ReceiveAsync(
							receiveBuffer.AsMemory(),
							SocketFlags.None,
							cancellationToken
					);
					
					if (bytesRead == 0) {
						break;
					}

					_buffer.Write(receiveBuffer.AsSpan(0, bytesRead));
					Console.WriteLine($"[C# Core] IPC Received: {bytesRead} byte from Go Gateway");
				}
			} catch (Exception ex) when (ex is not OperationCanceledException) {
				Console.WriteLine($"[C# Core] IPC Connection error: {ex.Message}");
			} finally {
				ArrayPool<byte>.Shared.Return(receiveBuffer);
			}
		}
	}

	private void Cleanup() {
		_listenerSocket?.Close();
		if (File.Exists(_socketPath)) {
			File.Delete(_socketPath);
		}
		Console.WriteLine("[C# Core] IPC Listener stopped & socket cleaned up.");
	}
}



