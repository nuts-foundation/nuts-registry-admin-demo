const path = require('path');
const glob = require('glob')
const HtmlWebpackPlugin = require('html-webpack-plugin');
const webpack = require("webpack");
const {VueLoaderPlugin} = require('vue-loader')

module.exports = {
  plugins: [
    new VueLoaderPlugin(),
    new HtmlWebpackPlugin({
      title: "Nuts Demo registry admin",
      template: "./web/src/index.html"
    }),
    new webpack.DefinePlugin({
      __VUE_OPTIONS_API__: true,
      __VUE_PROD_DEVTOOLS__: false,
    })
  ],
  optimization: {
    splitChunks: {
      cacheGroups: {
        styles: {
          name: 'styles',
          test: /\.css$/,
          chunks: 'all',
          enforce: true
        }
      }
    }
  },
  // devtool: 'inline-source-map',
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
      // Use the runtime here since it is smaller and we use precompiled .vue components
      'vue': 'vue/dist/vue.runtime.esm-bundler.js'
      // 'vue': 'vue/dist/vue.esm-bundler.js'
    }
  },
  module: {
    rules: [
      {
        test: /\.vue$/,
        loader: 'vue-loader'
      },
      {
        test: /\.(woff|woff2|eot|ttf|otf)$/i,
        type: 'asset/resource',
      },
      {
        test: /\.css$/,
        use: [
          'vue-style-loader',
          'style-loader',
          'css-loader',
          'postcss-loader'
        ]
      }
    ]
  }
};