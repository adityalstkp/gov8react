const path = require("path");

const isDevelopment = process.env.NODE_ENV === "development";

/** @type {import('@rspack/cli').Configuration} */
const config = {
  entry: {
    main: "./src/app.server.tsx",
  },
  output: {
    filename: "server.js",
    path: path.resolve(__dirname, "..", ".artifacts"),
  },
  optimization: {
    minimize: isDevelopment ? false : true,
  },
  builtins: {
    treeShaking: true,
  },
};

module.exports = config;
