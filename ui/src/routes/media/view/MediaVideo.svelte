<script>
	import { onMount, onDestroy } from 'svelte'

	import backend from '$lib/backend.js'
	import MediaSeekbar from './MediaSeekbar.svelte'
	import MediaButton from './MediaButton.svelte'
	import MediaButtonPlayPause from './MediaButtonPlayPause.svelte'
	import MediaButtonReset from './MediaButtonReset.svelte'

	let { entityId } = $props()
	let media = $state(null)
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

	function syncVideoState() {
		playing = !mediaElement.paused
	}

	function onpause() {
		syncVideoState()
	}

	function onseekvalue(_, seekTime) {
		mediaElement.currentTime = seekTime
		syncVideoState()
	}

	onMount(async () => {
		media = await backend.getMediaById(entityId)
	})
</script>

<div class="media-video-container">
	<video
		bind:this={mediaElement}
		class="media-video"
		width="320"
		height="240"
		title={media?.name}
		alt={media?.description}
		onclick={playPause}
		{onpause}>
		<source src="/media?entity_id={encodeURI(entityId)}" type="video/mp4" />
		HTML videos not supported by browser.
	</video>

	<div class="media-video-controls">
		<div class="media-video-seekbar">
			<MediaSeekbar {mediaElement} {onseekvalue} />
		</div>
		<div class="media-video-control-buttons">
			<MediaButtonPlayPause {mediaElement} />
			<MediaButtonReset {mediaElement} />
		</div>
	</div>
</div>

<style>
	.media-video-container {
		width: 100%;
		height: 100%;

		display: flex;
		flex-direction: column;
	}

	.media-video {
		width: 100%;
		flex-grow: 1;
	}

	.media-video-controls {
		flex-basis: 60px;
		flex-grow: 0;
		flex-shrink: 0;
		width: 100%;
	}

	.media-video-seekbar {
		width: 100%;
	}
</style>
