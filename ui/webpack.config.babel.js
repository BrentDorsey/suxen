const webpack = require('webpack');
const HTMLPlugin = require('html-webpack-plugin');
const path = require('path');

const NODE_ENV = process.env.NODE_ENV;
const isDev = NODE_ENV === 'development';

const prodEntry = {
  app: './src/index.jsx',
};
const devEntry = {
  'react-hot-loader/patch': 'react-hot-loader/patch',
  app: './src/index.jsx',
};

module.exports = {
  mode: NODE_ENV,
  entry: isDev ? devEntry : prodEntry,
  output: {
    path: path.resolve(__dirname, 'dist/'),
    filename: '[name].js',
  },
  resolve: {
    extensions: ['.jsx', '.js'],
  },
  module: {
    rules: [
      {
        test: /.jsx?$/,
        use: ['babel-loader'],
        exclude: /node_modules/,
      },
      {
        test: /.(ttf|otf|eot|svg|woff(2)?)(\?[a-z0-9]+)?$/,
        use: [{
          loader: 'file-loader',
          options: {
            name: '[name].[ext]',
            outputPath: 'fonts/',    // where the fonts will go
            publicPath: '../'       // override the default path
          }
        }]
      },
      {
        test: /\.s?css$/,
        use: [
          {
            loader: 'style-loader',
          },
          {
            loader: 'css-loader',
            options: {
              sourceMap: isDev,
              minimize: !isDev,
              localIdentName: isDev ? '[name]__[local]--[hash:base64:5]' : '[hash:base64:5]',
            },
          },
          {
            loader: 'sass-loader',
          },
        ],
      },
    ],
  },
  plugins: [
    new HTMLPlugin({
      chunks: ['react-hot-loader/patch', 'app'],
      chunksSortMode: 'manual',
      title: 'Suxen: Nexus Image Search'
    }),
    new webpack.DefinePlugin({
      'process.env.NODE_ENV': JSON.stringify(NODE_ENV),
    }),
    new webpack.NamedModulesPlugin(),
    new webpack.HotModuleReplacementPlugin(),
  ].filter(Boolean),
  devServer: {
    port: process.env.PORT || 8008,
    hotOnly: true,
    proxy: {
      '/query': {
        target: 'http://localhost:8080'
      }
    }
  },
  devtool: 'cheap-module-source-map',
};
