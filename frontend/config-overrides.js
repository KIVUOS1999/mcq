const webpack = require('webpack')
module.exports = function override(config) {
    // Initialize fallback object, defaulting to an empty object if not present
    const fallback = config.resolve.fallback || {};

    // Assign the polyfill configuration to the fallback object
    Object.assign(fallback, {
        stream: require.resolve("stream-browserify"),
    });

    // Set the updated fallback object back to config.resolve
    config.resolve.fallback = fallback;
    config.plugins = (config.plugins || []).concat([
        new webpack.ProvidePlugin({
            process: "process/browser"
        })
    ])

    return config;
};
