const path = require('path');

module.exports = {
  mode: "development",
  entry: {
    index: './web/src/index.js',
    print: './web/src/print.js',
  },
  output: {
    filename: '[name].bundle.js',
    path: path.resolve(__dirname, 'web/dist'),
  },
};