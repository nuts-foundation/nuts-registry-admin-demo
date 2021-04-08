const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const { VueLoaderPlugin } = require('vue-loader')


module.exports = {
  mode: "development",
  plugins: [
    new VueLoaderPlugin(),
    new HtmlWebpackPlugin({
      title: "Nuts Demo registry admin",
      template: "./web/src/index.html"
    }),
  ],
  devtool: 'inline-source-map',
  entry: {
    index: './web/src/index.js',
  },
  output: {
    filename: '[name].bundle.js',
    path: path.resolve(__dirname, 'web/dist'),
    clean: true,
  },
  resolve: {
    alias: {
      'vue': 'vue/dist/vue.esm-bundler.js'
    }
  },
  module: {
    rules: [
      {
        test: /\.vue$/,
        loader: 'vue-loader'
      },
      {
        test: /\.css$/,
        use: [
          'vue-style-loader',
          'css-loader'
        ]
      }
    ]
  }
};