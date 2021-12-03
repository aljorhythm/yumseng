import HtmlWebPackPlugin from "html-webpack-plugin";
import webpack from "webpack";
import path, { dirname } from "path";
import { fileURLToPath } from "url";

import MiniCssExtractPlugin from "mini-css-extract-plugin";

import Dotenv from "dotenv-webpack";
import nodeExternals from "webpack-node-externals";

const __dirname = dirname(fileURLToPath(import.meta.url));

const htmlPlugin = new HtmlWebPackPlugin({
  template: "./src/index.html",
  filename: "./index.html",
});

const outPath = path.join(__dirname, "dist");
const outPathResource = path.join(outPath, "resources");

const configResource = {
  mode: "development",
  entry: ["react-hot-loader/patch", "./src/index.js"],
  plugins: [
    new Dotenv({ systemvars: true }),
    htmlPlugin,
    new webpack.HotModuleReplacementPlugin(),
    new MiniCssExtractPlugin(),
  ],
  output: {
    path: outPathResource,
    filename: "index.js",
    publicPath: "/",
    clean: true,
  },
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        exclude: /node_modules/,
        use: {
          loader: "babel-loader",
        },
      },
      {
        test: /\.css$/i,
        use: ["style-loader", "css-loader"],
      },
      {
        test: /\.png/,
        type: "asset/resource",
      },
    ],
  },
};

export default [configResource];
