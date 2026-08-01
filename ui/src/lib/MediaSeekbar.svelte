<script>
	import { onMount } from 'svelte'
	import eventUtil from '$lib/eventUtil.js'

	let {
		mediaElement, //
		name,
		initValue,
		onseekstart,
		onseekvalue,
		onseekend,
		...attrs
	} = $props()

	let input = null
	let trackingMediaTime = false
	let value = $state(initValue)

	function onpointerdown(event) {
		if (eventUtil.isPrimaryButton(event)) {
			stopTrackingMediaTime()
			onseekstart?.(name)
		}
	}

	function oninput(event) {
		if (!trackingMediaTime) {
			onseekvalue?.(name, value)
		}
	}

	function onpointerup(event) {
		if (!trackingMediaTime && eventUtil.isPrimaryButton(event)) {
			startTrackingMediaTime()
			onseekend?.(name)
		}
	}

	onMount(() => {
		input.max = mediaElement.duration
		input.disabled = !mediaElement.seekable
		startTrackingMediaTime()
		return stopTrackingMediaTime
	})

	function trackMediaTime() {
		value = mediaElement.currentTime
	}

	function startTrackingMediaTime() {
		if (mediaElement) {
			mediaElement.addEventListener('timeupdate', trackMediaTime)
			trackingMediaTime = true
		}
	}

	function stopTrackingMediaTime() {
		if (mediaElement) {
			mediaElement.removeEventListener('timeupdate', trackMediaTime)
			trackingMediaTime = false
		}
	}
</script>

<svelte:window {onpointerup} />

<input
	{...attrs}
	bind:this={input}
	bind:value
	class:media-seekbar={true}
	{name}
	disabled={!mediaElement}
	type="range"
	min="0"
	step="1"
	{onpointerdown}
	{oninput} />

<style>
	.media-seekbar {
		width: 99%;
	}
</style>
