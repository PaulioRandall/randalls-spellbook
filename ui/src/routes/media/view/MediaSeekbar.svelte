<script>
	import { onMount, onDestroy } from 'svelte'
	import eventUtil from '$lib/eventUtil.js'

	let { svelteMediaElement, ...attrs } = $props()
	let value = $state(0)

	onMount(() => {
		value = svelteMediaElement.currentTime
		svelteMediaElement.on('timeupdate', timeupdate)
	})

	onDestroy(() => {
		svelteMediaElement.off('timeupdate', timeupdate)
	})

	function onpointerdown(event) {
		if (eventUtil.isPrimaryButton(event)) {
			svelteMediaElement.off('timeupdate', timeupdate)
		}
	}

	function onpointerup(event) {
		if (eventUtil.isPrimaryButton(event)) {
			svelteMediaElement.on('timeupdate', timeupdate)
		}
	}

	function timeupdate() {
		value = svelteMediaElement.currentTime
	}

	function oninput() {
		svelteMediaElement.seekTo(value)
	}
</script>

<svelte:window {onpointerup} />

<input
	{...attrs}
	bind:value
	class:media-seekbar={true}
	disabled={!svelteMediaElement.seekable}
	type="range"
	min="0"
	max={svelteMediaElement.duration}
	{onpointerdown}
	{oninput} />

<style>
	.media-seekbar {
		width: 99%;
	}
</style>
