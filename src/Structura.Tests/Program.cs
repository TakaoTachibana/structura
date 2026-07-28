using BenchmarkDotNet.Running;

namespace Structura.Tests;

public class Program {
	public static void Main(string[] args) {
		var summay = BenchmarkRunner.Run(typeof(TelemetryPipelineBenchmark));
	}
}

