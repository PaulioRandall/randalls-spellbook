<script>
	import MediaButton from './MediaButton.svelte'

	let { mediaController } = $props()
</script>

<MediaButton
	disabled={!mediaController.loaded}
	onclick={() => mediaController.restart()}>
	Reset
</MediaButton>
