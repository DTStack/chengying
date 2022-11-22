const path = require('path');
const webpack = require("webpack");
const CopyWebpackPlugin = require('copy-webpack-plugin');
const TerserPlugin = require('terser-webpack-plugin');
const MonacoWebpackPlugin = require('monaco-editor-webpack-plugin');
const ReactRefreshPlugin = require('@pmmmwh/react-refresh-webpack-plugin');
const FigletConsoleWebpackPlugin = require('figlet-console-webpack-plugin');
const MonacoConfig = require('./monacoConfig');
const packageInfo = require("./package.json"); // app version.

const isProd = process.env.NODE_ENV === "production";

module.exports = {

  entry: {
    app: { import: './src/app.tsx', dependOn: 'default-vendors' },
    'default-vendors': ['react', 'react-dom', 'react-redux', 'react-router-dom', 'moment', 'js-cookie'],
  },

  module: {
    rules: [
      {
        test: /\.(js|ts[x]?)$/,
        exclude: /node_modules/,
        use: [
          {
            loader: "babel-loader",
            options: {
              cacheDirectory: true,
              plugins: !isProd ? ["react-refresh/babel"] : []
            },
          },
          {
            loader: 'thread-loader',
            options: {
              workers: 3,
            },
          }
        ]
      }
    ]
  },

  plugins: [

    new MonacoWebpackPlugin({
      features: MonacoConfig.features,
      languages: MonacoConfig.languages,
    }),

    new CopyWebpackPlugin([{
      from: path.resolve(__dirname, './public')
    }], {
      ignore: [
        'json/*.ts',
        'imgs/*',
        'styles/*',
        'fonts/emfont/*',
        'sw.js',
      ],
    }),

    new ReactRefreshPlugin(),

    new webpack.DefinePlugin({
      APP: {
        VERSION: JSON.stringify(packageInfo.version)
      }
    }),

    new FigletConsoleWebpackPlugin({
      name: "ChengYing",
      content: `App current version: ${packageInfo.version}`,
      options: {
        markMaxLength: 86
      }
    }),
  ],

  optimization: {
    usedExports: true,
    minimize: isProd,
    minimizer: [
      new TerserPlugin({
        parallel: true,
        cache: true,
        terserOptions: {
          compress: {
            drop_console: true,
            drop_debugger: true
          },
          format: {
            comments: false
          }
        },
        extractComments: false
      }),
    ],
  },

  resolve: {
    fallback: {
      fs: false,
      module: "empty",
      path: false,
      events: false,
    },
  },

  externals: {
    APPCONFIG: 'APPCONFIG',
  },

  devServer: {
    host: '0.0.0.0',
    port: 8099,
    open: false,
    client: {
      logging: "info",
      overlay: {
        errors: false,
        warnings: false,
      },
      progress: true,
    },
    proxy: [
      {
        path: '/api/v2/cluster/showShellLog',
        target: 'ws://172.16.82.176',
        ws: true,
      },
      {
        path: [
          '/api/**',
          '/gate/**',
        ],
        // target: 'http://172.16.10.191', // 测试环境
        target: 'http://172.16.82.176', // 开发环境          
      }
    ]
  },
};