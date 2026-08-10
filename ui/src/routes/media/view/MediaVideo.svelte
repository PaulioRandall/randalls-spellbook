<script>
	import { onMount, onDestroy } from 'svelte'

	import backend from '$lib/backend.js'
	import MediaController from '$lib/MediaController.svelte.js'

	import MediaSeekbar from './MediaSeekbar.svelte'
	import MediaButton from './MediaButton.svelte'
	import MediaButtonPlayPause from './MediaButtonPlayPause.svelte'
	import MediaButtonRestart from './MediaButtonRestart.svelte'

	const mediaController = new MediaController()
	let { entityId } = $props()

	let media = $state(null)
	let mediaElement = $state(null)

	function onseekvalue(_, seekTime) {
		mediaElement.currentTime = seekTime
	}

	onMount(async () => {
		media = await backend.getMediaById(entityId)
		mediaController.setElement(mediaElement)
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
		onclick={() => mediaController.playPause()}>
		<source src="/media?entity_id={encodeURI(entityId)}" type="video/mp4" />
		HTML videos not supported by browser.
	</video>

	<div class="media-video-controls">
		<div class="media-video-seekbar">
			<MediaSeekbar {mediaElement} {onseekvalue} />
		</div>
		<div class="media-video-control-buttons">
			<MediaButtonPlayPause {mediaController} />
			<MediaButtonRestart {mediaController} />
			<MediaButton
				disabled={!mediaController.loaded}
				onclick={() => mediaController.reload()}>
				Reload
			</MediaButton>
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
