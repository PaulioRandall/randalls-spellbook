<script>
	import MediaButton from './MediaButton.svelte'

	const title = 'Quick jump back to the beginning'

	let { mediaSvox } = $props()
	let disabled = $derived(!mediaSvox.loaded)

	function onclick() {
		mediaSvox.seekTo(0)
	}
</script>

<MediaButton {disabled} {onclick} {title}>Reset</MediaButton>
