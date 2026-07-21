using System;
using System.Threading;
using System.Threading.Tasks;
using Structura.Core.Buffers;
using Structura.Core.Ipc;

namespace Structura.Core;

internal class Program {
	private const string SocketPath = "/tmp/structura.sock";

	static async Task Main(string[] args) {
		Console.WriteLine("=== Structura Backend Core (.NET 10) ===");
		using var cts = new CancellationTokenSource();
		Console.CancelKeyPress += (s, e) => {
			e.Cancel = true;
			cts.Cancel();
		};

		var buffer = new AllocationFreeBuffer(1024 * 1024);
		var ipcListener = new UnixDomainSocketListener(SocketPath, buffer);
		Task listenerTask = ipcListener.StartAsync(cts.Token);
		Console.WriteLine("Press Ctrl+C to exit.");

		try {
			await listenerTask;
		} catch (OperationCanceledException) {
			Console.WriteLine("Shutting down gracefully...");
		}
	}
}

