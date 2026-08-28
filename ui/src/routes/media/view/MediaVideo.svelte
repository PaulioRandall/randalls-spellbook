<script>
	import { onMount } from 'svelte'

	import backend from '$lib/backend.js'

	import MediaSvox from '$lib/MediaSvox.svelte.js'
	import TextareaSvox from '$lib/TextareaSvox.svelte.js'

	import MediaSeekbar from './MediaSeekbar.svelte'
	import MediaButtonPlayPause from './MediaButtonPlayPause.svelte'
	import MediaButtonRestart from './MediaButtonRestart.svelte'
	import MediaButtonReload from './MediaButtonReload.svelte'
	import MediaButtonSubmitObservation from './MediaButtonSubmitObservation.svelte'
	import MediaVideoPlayer from './MediaVideoPlayer.svelte'
	import MediaObservationInput from './MediaObservationInput.svelte'

	const mediaSvox = new MediaSvox()
	const textareaSvox = new TextareaSvox()

	let { entityId } = $props()
	let media = $state(null)

	onMount(async () => {
		media = await backend.GetMediaById(entityId)
	})

	mediaSvox.onElement(() => {
		console.log('EVENT', 'element', mediaSvox.getElement())
	})
</script>

<div class="media-video-container">
	<div class="media-video-player">
		<MediaVideoPlayer {mediaSvox} {media} />
	</div>

	<div class="media-video-controls">
		<div class="media-video-seekbar">
			<MediaSeekbar {mediaSvox} />
		</div>
		<div>
			<MediaButtonPlayPause {mediaSvox} />
			<MediaButtonRestart {mediaSvox} />
			<MediaButtonReload {mediaSvox} />
			<MediaButtonSubmitObservation {mediaSvox} {textareaSvox} />
		</div>
		<div>
			<MediaObservationInput {mediaSvox} {textareaSvox} />
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

	.media-video-player {
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
