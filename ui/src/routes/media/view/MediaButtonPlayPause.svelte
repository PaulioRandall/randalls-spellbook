<script>
	import MediaButton from './MediaButton.svelte'

	let { svelteMediaElement } = $props()

	let label = $derived(
		svelteMediaElement.paused ? 'Play' : 'Pause' //
	)

	let disabled = $derived(
		!svelteMediaElement.running && //
			!svelteMediaElement.playable
	)

	function onclick() {
		svelteMediaElement.playPause()
	}
</script>

<MediaButton {disabled} {onclick}>
	{label}
</MediaButton>
