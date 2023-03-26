const path = require("path");
const BundleAnalyzerPlugin =
  require("webpack-bundle-analyzer").BundleAnalyzerPlugin;

const isDevelopment = process.env.NODE_ENV === "development";

const plugins = [];

if (process.env.ANALYZE) {
  plugins.push(new BundleAnalyzerPlugin());
}

/** @type {import('@rspack/cli').Configuration} */
const config = {
  entry: "./src/app.client.tsx",
  output: {
    filename: "[name].client.js",
    path: path.resolve(__dirname, "..", ".artifacts"),
  },
  optimization: {
    minimize: isDevelopment ? false : true,
  },
  builtins: {
    treeShaking: true,
  },
  plugins: plugins,
};

module.exports = config;
