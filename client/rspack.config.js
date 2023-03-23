const path = require("path");

const isDevelopment = process.env.NODE_ENV === "development";

/** @type {import('@rspack/cli').Configuration} */
const config = {
  entry: {
    main: "./src/index.ts",
  },
  output: {
    filename: "main.js",
    path: path.resolve(__dirname, "..", ".artifacts"),
  },
  optimization: {
    minimize: isDevelopment ? false : true,
  },
  builtins: { emotion: true },
};

module.exports = config;
