<!--
<script>
	import { onMount, onDestroy } from 'svelte'

	import backend from '$lib/backend.js'
	import MediaButton from './MediaButton.svelte'

	let { media, mediaElement } = $props()
	let playing = $state(false)

	function playPause() {
		if (mediaElement.paused) {
			mediaElement.play()
		} else {
			mediaElement.pause()
		}
		syncMediaState()
	}

	function reset() {
		mediaElement.pause()
		mediaElement.load()
		syncMediaState()
	}

	function onpause() {
		syncMediaState()
	}

	function syncMediaState() {
		playing = !mediaElement.paused
	}

	onMount(async () => {
		
	})
</script>

		<div class="media-control-bar">
			<MediaButton onclick={playPause}>
				{#if playing}Pause{:else}Play{/if}
			</MediaButton>
			<MediaButton onclick={reset}>Reset</MediaButton>
		</div>

<style>
		.media-control-bar {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
	}
</style>
-->
