const { defineConfig } = require('@vue/cli-service');

module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    proxy: 'http://localhost:8080' // Adjust the port if your Express server is on a different port
  }
});
