import { defineConfig, loadEnv } from 'vite'
// @ts-ignore
import path from 'path'
// @ts-ignore
import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'

function resolveAllowedHosts(raw: string): string[] | true {
  const hosts = raw
    .split(',')
    .map((host) => host.trim())
    .filter(Boolean);

  return hosts.length > 0 ? hosts : true;
}

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');
  const proxyTarget = (env.VITE_API_PROXY_TARGET || 'http://localhost:28080').trim();
  const allowedHosts = resolveAllowedHosts(env.VITE_DEV_ALLOWED_HOSTS || '');

  return {
    plugins: [
      react(),
      tailwindcss(),
    ],

    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },

    assetsInclude: ['**/*.svg', '**/*.csv'],

    server: {
      host: '0.0.0.0',
      allowedHosts,
      proxy: {
        '/api': {
          target: proxyTarget,
          changeOrigin: true,
        },
      },
    },

    preview: {
      host: '0.0.0.0',
      allowedHosts,
    },
  };
})
