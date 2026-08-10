<script>
	import MediaButton from './MediaButton.svelte'

	let { svelteMediaElement } = $props()
</script>

<MediaButton
	disabled={!svelteMediaElement.loaded}
	onclick={() => svelteMediaElement.restart()}>
	Reset
</MediaButton>
