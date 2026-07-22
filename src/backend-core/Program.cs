using System;
using System.IO;
using System.Net.Sockets;
using System.Threading.Tasks;

namespace Structura.Core;

internal class Program {
	private const string SocketPath = "/tmp/structura.sock";

	static async Task Main(string[] args) {
		await PacketProcessor.Broadcaster.StartAsync(8080);
		if (File.Exists(SocketPath)) {
			File.Delete(SocketPath);
		}

		Console.WriteLine("=== Structura Backend Core (.NET 10) ===");
		using var listenSocket = new Socket(AddressFamily.Unix, SocketType.Stream, ProtocolType.Unspecified);
		listenSocket.Bind(new UnixDomainSocketEndPoint(SocketPath));
		listenSocket.Listen(10);

		Console.WriteLine($"[C# Core] Waiting for Go Gateway connection on {SocketPath}...");

		using var clientSocket = await listenSocket.AcceptAsync();
		Console.WriteLine("[C# Core] Gateway connected. Zero-Alloc evaluation loop running...");

		var buffer = new byte[1024];

		while (true) {
			int readBytes = await clientSocket.ReceiveAsync(buffer, SocketFlags.None);
			if (readBytes == 0) {
				break;
			}
			PacketProcessor.ProcessBuffer(buffer.AsSpan(0, readBytes));
		}
	}
}

