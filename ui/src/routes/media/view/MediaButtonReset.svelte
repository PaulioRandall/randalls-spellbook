<script>
	import { onMount, onDestroy } from 'svelte'
	import MediaButton from './MediaButton.svelte'

	let { mediaElement } = $props()
	let mediaLoaded = $state(false)

	function reset() {
		mediaElement.pause()
		mediaLoaded = false
		mediaElement.load()
	}

	onMount(startTracking)
	onDestroy(stopTracking)

	function onloadeddata() {
		mediaLoaded = true
	}

	function startTracking() {
		mediaElement.addEventListener('loadeddata', onloadeddata)
	}

	function stopTracking() {
		mediaElement?.removeEventListener('loadeddata', onloadeddata)
	}
</script>

<MediaButton disabled={!mediaLoaded} onclick={reset}>Reset</MediaButton>
