<script>
	import { onMount } from 'svelte'

	import backend from '$lib/backend.js'

	import HTMLMediaElementSvox from '$lib/HTMLMediaElementSvox.svelte.js'

	import MediaSeekbar from './MediaSeekbar.svelte'
	import MediaButtonPlayPause from './MediaButtonPlayPause.svelte'
	import MediaButtonRestart from './MediaButtonRestart.svelte'
	import MediaButtonReload from './MediaButtonReload.svelte'
	import MediaVideoPlayer from './MediaVideoPlayer.svelte'
	import MediaObservationInput from './MediaObservationInput.svelte'

	const mediaSvox = new HTMLMediaElementSvox()

	let { entityId } = $props()
	let media = $state(null)

	onMount(async () => {
		media = await backend.getMediaById(entityId)
	})

	mediaSvox.onElement(() => {
		console.log('EVENT', 'element', mediaSvox.getElement())
	})

	function logOnEvent(name) {
		mediaSvox.on(name, () => {
			console.log('EVENT', name)
		})
	}

	logOnEvent('running')
	logOnEvent('flowing')
	logOnEvent('pausing')
	logOnEvent('buffering')

	function logReactor(name) {
		return () => {
			console.log('STATE', name, mediaSvox[name])
		}
	}

	$effect(logReactor('element'))
	$effect(logReactor('loaded'))
	$effect(logReactor('playable'))
	$effect(logReactor('running'))
	$effect(logReactor('flowing'))
	$effect(logReactor('buffering'))
	$effect(logReactor('paused'))
	$effect(logReactor('seekable'))
	$effect(logReactor('seeking'))
	$effect(logReactor('duration'))
	$effect(logReactor('currentTime'))
	$effect(logReactor('currentRemaining'))
	$effect(logReactor('playtime'))
	$effect(logReactor('seektime'))
	$effect(logReactor('remainingTime'))
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
		</div>
		<div>
			<MediaObservationInput {mediaSvox} {media} />
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
