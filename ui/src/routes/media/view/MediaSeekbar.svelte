<script>
	import { onMount } from 'svelte'
	import eventUtil from '$lib/eventUtil.js'

	let { mediaSvox, ...attrs } = $props()

	let disabled = $derived(!mediaSvox.seekable)
	let max = $derived(mediaSvox.duration)
	let value = $state(0)

	onMount(() => {
		value = mediaSvox.currentTime
		return mediaSvox.on({ timeupdate })
	})

	function onpointerdown(event) {
		if (eventUtil.isPrimaryButton(event)) {
			mediaSvox.off('timeupdate', timeupdate)
		}
	}

	function onpointerup(event) {
		if (eventUtil.isPrimaryButton(event)) {
			mediaSvox.on('timeupdate', timeupdate)
		}
	}

	function timeupdate() {
		value = mediaSvox.currentTime
	}

	function oninput() {
		mediaSvox.seekTo(value)
	}
</script>

<svelte:window {onpointerup} />

<input
	{...attrs}
	bind:value
	class:media-seekbar={true}
	type="range"
	min="0"
	{max}
	{disabled}
	{onpointerdown}
	{oninput} />

<style>
	.media-seekbar {
		width: 99%;

		&:disabled {
			cursor: not-allowed;
		}
	}
</style>
