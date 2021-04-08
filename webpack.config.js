const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
  mode: "development",
  plugins: [
    new HtmlWebpackPlugin({
      title: "Nuts Demo registry admin",
      template: "./web/src/index.html"
    }),
  ],
  entry: {
    index: './web/src/index.js',
    print: './web/src/print.js',
  },
  output: {
    filename: '[name].bundle.js',
    path: path.resolve(__dirname, 'web/dist'),
    clean: true,
  },
};