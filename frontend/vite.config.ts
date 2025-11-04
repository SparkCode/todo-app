import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig((configEnv) => {
  const isSsrBuild = configEnv.isSsrBuild ?? false;

  return {
  plugins: [react()],
  root: '.',
  build: isSsrBuild
    ? {
        outDir: 'build/server',
        target: 'node18',
        rollupOptions: {
          input: './node.ts',
        },
      }
    : {
        outDir: 'dist',
        rollupOptions: {
          input: {
            main: './index.html',
          },
        },
        manifest: true,
      },
  ssr: {
    noExternal: ['react', 'react-dom'],
  },
  };
});

