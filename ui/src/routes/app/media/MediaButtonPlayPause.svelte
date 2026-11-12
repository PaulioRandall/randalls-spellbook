<script>
	import MediaButton from './MediaButton.svelte'

	const textStates = {
		playing: { label: 'Pause', title: 'Pause the media' },
		paused: { label: 'Play', title: 'Play the media' },
	}

	let { mediaSvox } = $props()

	let textState = $derived.by(() => {
		if (mediaSvox.paused) {
			return textStates.playing
		}
		return textStates.paused
	})

	let disabled = $derived(
		!mediaSvox.running && //
			!mediaSvox.playable
	)

	function onclick() {
		mediaSvox.playPause()
	}
</script>

<MediaButton {disabled} {onclick} title={textState.title}>
	{textState.label}
</MediaButton>
