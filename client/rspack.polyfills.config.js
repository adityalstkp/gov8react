const path = require("path");

const isDevelopment = process.env.NODE_ENV === "development";

/** @type {import('@rspack/cli').Configuration} */
const config = {
  entry: {
    text_encoder: "./polyfills/text_encoder.js",
    buffer: "./polyfills/buffer.js",
  },
  output: {
    filename: "polyfills.[name].js",
    path: path.resolve(__dirname, "..", ".artifacts"),
  },
  optimization: { minimize: isDevelopment ? false : true },
};

module.exports = config;
