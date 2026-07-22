const path = require('path');
const rspack = require('@rspack/core');

/** @type {import('@rspack/cli').Configuration} */
module.exports = {
	entry: {
		main: './src/main.ts',
	},
	output: {
		path: path.resolve(__dirname, 'dist'),
		filename: '[name].js',
		publicPath: '/',
	},
	resolve: {
		extensions: ['.ts', '.js'],
	},
	module: {
		rules: [
			{
				test: /\.ts$/,
				exclude: [/node_modules/],
				loader: 'builtin:swc-loader',
				options: {
					jsc: {
						parser: {
							syntax: 'typescript',
						},
					},
				},
				type: 'javascript/auto',
			},
		],
	},
	plugins: [
		new rspack.HtmlRspackPlugin({
			template: './index.html',
		}),
	],
	devServer: {
		port: 3000,
		hot: true,
		headers: {
			'Cross-Origin-Opener-Policy': 'same-origin',
			'Cross-Origin-Embedder-Policy': 'require-corp',
			'Content-Security-Policy': "default-src 'self'; connect-src 'self' ws://localhost:* http://localhost:*; worker-src 'self' blob:; style-src 'self' 'unsafe-inline';",
		},
	},
};

