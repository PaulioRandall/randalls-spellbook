<script>
	import { onMount } from 'svelte'
	import Video from './Video.svelte'

	let src = $state(null)
	let type = $state('video/mp4')

	onMount(async () => {
		const filepath = await window.selectVideoFile()

		if (!filepath) {
			alert('Refresh page and select a video this time!')
			return
		}

		const data = await window.readVideoFile(filepath)

		if (!data) {
			alert('An error occurred loading the video :(')
			return
		}

		const unit8Data = Uint8Array.fromBase64(data, (c) => c.charCodeAt(0))
		const blob = new Blob([unit8Data], { type })
		src = URL.createObjectURL(blob)
	})

	// TODO: test out Vidstack Svelte video player
	// TODO: Check state of video player before trying to
	//       interact with it.
</script>

<main>
	{#if src}
		<Video {src} {type} />
	{/if}
</main>

<style>
	main {
		width: 100%;
		height: 100%;

		display: flex;
		flex-direction: column;
	}
</style>
