import adapter from 'svelte-adapter-deno';
import { vitePreprocess } from '@sveltejs/kit/vite';

const config = {
  preprocess: [vitePreprocess()],
  kit: {
    adapter: adapter()
  }
};

export default config;
