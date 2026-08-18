<script>
	import MediaButton from './MediaButton.svelte'

	const title = 'Reload from file'

	let { mediaSvox } = $props()
	let disabled = $derived(!mediaSvox.loaded)

	function onclick() {
		mediaSvox.reload()
	}
</script>

<MediaButton {disabled} {onclick} {title}>Reload</MediaButton>
