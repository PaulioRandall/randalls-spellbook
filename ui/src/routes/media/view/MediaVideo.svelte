<script>
	import { onMount, onDestroy } from 'svelte'

	import backend from '$lib/backend.js'
	import SvelteMediaElement from '$lib/svelteMediaElement.svelte.js'

	import MediaSeekbar from './MediaSeekbar.svelte'
	import MediaButton from './MediaButton.svelte'
	import MediaButtonPlayPause from './MediaButtonPlayPause.svelte'
	import MediaButtonRestart from './MediaButtonRestart.svelte'

	const svelteMediaElement = new SvelteMediaElement()
	let { entityId } = $props()

	let media = $state(null)
	let mediaElement = $state(null)

	function onseekvalue(_, seekTime) {
		mediaElement.currentTime = seekTime
	}

	onMount(async () => {
		media = await backend.getMediaById(entityId)
		svelteMediaElement.setElement(mediaElement)
	})

	function logOnEvent(name) {
		svelteMediaElement.on(name, () => {
			console.log('EVENT', name)
		})
	}

	logOnEvent('elementset')
	logOnEvent('elementunset')
	logOnEvent('running')
	logOnEvent('flowing')
	logOnEvent('pausing')
	logOnEvent('buffering')

	function logReactor(name) {
		return () => {
			console.log('STATE', name, svelteMediaElement[name])
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
	<video
		bind:this={mediaElement}
		class="media-video"
		width="320"
		height="240"
		title={media?.name}
		alt={media?.description}
		onclick={() => svelteMediaElement.playPause()}>
		<source src="/media?entity_id={encodeURI(entityId)}" type="video/mp4" />
		HTML videos not supported by browser.
	</video>

	<div class="media-video-controls">
		<div class="media-video-seekbar">
			<MediaSeekbar {mediaElement} {onseekvalue} />
		</div>
		<div class="media-video-control-buttons">
			<MediaButtonPlayPause {svelteMediaElement} />
			<MediaButtonRestart {svelteMediaElement} />
			<MediaButton
				disabled={!svelteMediaElement.loaded}
				onclick={() => svelteMediaElement.reload()}>
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
