using System;
using System.Net;
using System.Net.WebSockets;
using System.Threading;
using System.Threading.Tasks;

namespace Structura.Core;

public class WebBroadcaster {
	private HttpListener? _listener;
	private WebSocket? _activeSocket;
	private readonly CancellationTokenSource _cts = new();

	public async Task StartAsync(int port = 8080) {
		_listener = new HttpListener();
		_listener.Prefixes.Add($"http://localhost:{port}/ws/");
		_listener.Start();
		Console.WriteLine($"[WebBroadcaster] High-speed WebSocket listening on ws://localhost:{port}/ws/");

		_ = Task.Run(async () => {
				while (!_cts.Token.IsCancellationRequested) {
					try {
						var context = await _listener.GetContextAsync();
						if (context.Request.IsWebSocketRequest) {
							var wsContext = await context.AcceptWebSocketAsync(subProtocol: null);
							_activeSocket = wsContext.WebSocket;
							Console.WriteLine("[WebBroadcaster] Frontend Client Connection via Binary Stream!");
						}
					} catch (Exception) {
					}
				}
		}, _cts.Token);
	}

	public void BroadcastBuffer(ReadOnlySpan<byte> buffer) {
		if (_activeSocket == null || _activeSocket.State != WebSocketState.Open) {
			return;
		}
		var segment = new ArraySegment<byte>(buffer.ToArray());
		_ = _activeSocket.SendAsync(segment, WebSocketMessageType.Binary, true, CancellationToken.None);
	}

	public void Stop() {
		_cts.Cancel();
		_listener?.Stop();
	}
}

