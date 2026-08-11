<script>
	import MediaButton from './MediaButton.svelte'

	let { svelteMediaElement } = $props()

	let label = $derived(
		svelteMediaElement.paused ? 'Play' : 'Pause' //
	)
</script>

<MediaButton
	disabled={!svelteMediaElement.playable}
	onclick={() => svelteMediaElement.playPause()}>
	{label}
</MediaButton>
