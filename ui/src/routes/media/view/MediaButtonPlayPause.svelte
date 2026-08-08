<script>
	import { onMount, onDestroy } from 'svelte'
	import MediaButton from './MediaButton.svelte'

	let { mediaElement } = $props()
	let mediaLoaded = $state(false)
	let playing = $state(false)

	onMount(startTracking)
	onDestroy(stopTracking)

	function onloadstart() {
		mediaLoaded = false
		playing = false
	}

	function onloadeddata() {
		mediaLoaded = true
	}

	function onplaypause() {
		playing = !mediaElement.paused
	}

	function startTracking() {
		mediaElement.addEventListener('loadstart', onloadstart)
		mediaElement.addEventListener('loadeddata', onloadeddata)
		mediaElement.addEventListener('play', onplaypause)
		mediaElement.addEventListener('pause', onplaypause)
	}

	function stopTracking() {
		if (mediaElement) {
			mediaElement.removeEventListener('pause', onplaypause)
			mediaElement.removeEventListener('play', onplaypause)
			mediaElement.removeEventListener('loadeddata', onloadeddata)
			mediaElement.removeEventListener('loadstart', onloadstart)
		}
	}

	function playPause() {
		if (mediaElement.paused) {
			mediaElement.play()
		} else {
			mediaElement.pause()
		}
	}
</script>

<MediaButton disabled={!mediaLoaded} onclick={playPause}>
	{#if playing}
		Pause
	{:else}
		Play
	{/if}
</MediaButton>
