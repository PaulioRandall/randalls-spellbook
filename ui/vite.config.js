import { sveltekit } from '@sveltejs/kit/vite'
import adapter from '@sveltejs/adapter-static'

export default {
	plugins: [
		sveltekit({
			adapter: adapter(),
			preprocess: [
				// Used often.
			],
		}),
	],
}
