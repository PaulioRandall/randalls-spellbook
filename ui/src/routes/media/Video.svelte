<script>
	import { onMount, onDestroy } from 'svelte'

	import MediaSeekbar from './MediaSeekbar.svelte'
	import MediaButton from './MediaButton.svelte'

	let { src, type } = $props()
	let mediaElement = $state(null)
	let playing = $state(false)

	function playPause() {
		if (mediaElement.paused) {
			mediaElement.play()
		} else {
			mediaElement.pause()
		}
		syncVideoState()
	}

	function reset() {
		mediaElement.pause()
		mediaElement.load()
		syncVideoState()
	}

	function syncVideoState() {
		playing = !mediaElement.paused
	}

	function onseekvalue(_, seekTime) {
		mediaElement.currentTime = seekTime
		syncVideoState()
	}
</script>

<div class="video-container">
	<!-- TODO: Probably should be its own component -->
	<video
		bind:this={mediaElement}
		class="video"
		width="320"
		height="240"
		onclick={playPause}>
		<source src="/media?token=123&path={encodeURI(src)}" {type} />
		HTML videos not supported by browser.
	</video>

	<div class="video-controls">
		<div class="video-seekbar">
			<MediaSeekbar {mediaElement} {onseekvalue} />
		</div>
		<div class="video-control-buttons">
			<MediaButton onclick={playPause}>
				{#if playing}Pause{:else}Play{/if}
			</MediaButton>
			<MediaButton onclick={reset}>Reset</MediaButton>
		</div>
	</div>
</div>

<style>
	.video-container {
		width: 100%;
		height: 100%;

		display: flex;
		flex-direction: column;
	}

	.video {
		width: 100%;
		flex-grow: 1;
	}

	.video-controls {
		flex-basis: 60px;
		flex-grow: 0;
		flex-shrink: 0;
		width: 100%;
	}

	.video-seekbar {
		width: 100%;
	}

	.video-control-buttons {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
	}
</style>
