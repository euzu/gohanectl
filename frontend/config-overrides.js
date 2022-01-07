module.exports = function override(config, env) {
    const loaders = config.resolve;
    loaders.fallback = {
        "fs": false,
        "tls": false,
        "net": false,
        "http": require.resolve("stream-http"),
        "https": false,
        "zlib": require.resolve("browserify-zlib"),
        "path": require.resolve("path-browserify"),
        "stream": require.resolve("stream-browserify"),
        //"util": require.resolve("util/"),
        "crypto": require.resolve("crypto-browserify")
    };

    config.module.rules.push(...[
        {
            test: /\.(js|jsx)$/,
            exclude: /node_modules/,
            use: {
                loader: 'babel-loader'
            }
        },
        {
            test: /\.js$/,
            enforce: "pre",
            use: ["source-map-loader"],
        }
    ]);

    config.ignoreWarnings = [/Failed to parse source map/];

    return config;
}