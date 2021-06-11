const path = require("path");
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve:{
    alias: {
      // 键必须以斜线开始和结束
      "/@/": path.resolve(__dirname, "./src"),
    }
  },
  build:{
    outDir: "../public",
    assetsDir:"static",
  },
  server:{
    port: 8000,
    // 是否自动在浏览器打开
    open: true,
    // 是否开启 https
    https: false,
    // 服务端渲染
    ssr: false,
    proxy: {
    // 如果是 /api 打头，则访问地址如下
    "/api": {
      target: "http://localhost:8044/api/",
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, ""),
    },
  },
  },

})
