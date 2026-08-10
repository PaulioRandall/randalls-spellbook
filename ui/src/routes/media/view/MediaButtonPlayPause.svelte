<script>
	import MediaButton from './MediaButton.svelte'

	let { mediaController } = $props()
	let label = $derived(
		mediaController.playing ? 'Pause' : 'Play' //
	)
</script>

<MediaButton
	disabled={!mediaController.playable}
	onclick={() => mediaController.playPause()}>
	{label}
</MediaButton>
