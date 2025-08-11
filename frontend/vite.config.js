import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

export default defineConfig({
  plugins: [
    vue(),
    AutoImport({
      resolvers: [ElementPlusResolver()],
      imports: ['vue', 'vue-router', 'pinia']
    }),
    Components({
      resolvers: [ElementPlusResolver()]
    })
  ],
  server: {
    host: "0.0.0.0",
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://192.168.6.22:8080',
        changeOrigin: true
      }
    }
  },
  build: {
    outDir: '../static',
    rollupOptions: {
      output: {
        entryFileNames: (chunkInfo) => {
          return chunkInfo.name === 'index' ? 'assets/index.js' : 'assets/[name]-[hash].js'
        },
        chunkFileNames: 'assets/[name]-[hash].js',
        assetFileNames: (assetInfo) => {
          if (assetInfo.name === 'index.css') {
            return 'assets/index.css'
          }
          return 'assets/[name]-[hash].[ext]'
        }
      }
    }
  }
})
