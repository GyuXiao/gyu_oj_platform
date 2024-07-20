const { defineConfig } = require("@vue/cli-service");
const MonacoWebpackPlugin = require("monaco-editor-webpack-plugin");

module.exports = defineConfig({
  transpileDependencies: true,
  chainWebpack(config) {
    config.plugin("manaco").use(MonacoWebpackPlugin, [
      {
        // Languages are loaded on demand at runtime
        languages: [
          "json",
          "go",
          "css",
          "html",
          "java",
          "javascript",
          "less",
          "markdown",
          "mysql",
          "php",
          "python",
          "scss",
          "shell",
          "redis",
          "sql",
          "typescript",
          "xml",
        ],

        features: [
          "format",
          "find",
          "contextmenu",
          "gotoError",
          "gotoLine",
          "gotoSymbol",
          "hover",
          "documentSymbols",
        ],
      },
    ]);
  },
});
