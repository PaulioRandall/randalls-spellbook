<script>
	import MediaButton from './MediaButton.svelte'

	const title = 'Submit observation'

	let { mediaSvox, textareaSvox, media } = $props()
	let disabled = $derived(
		!mediaSvox.loaded || textareaSvox.empty //
	)

	function onclick() {
		const ob = {
			MediaId: media.EntityId,
			StartTime: mediaSvox.currentTime,
			Duration: 0,
			Description: textareaSvox.text,
		}

		Go('AddObservation', ob)
			.then(() => textareaSvox.setText(''))
			.catch(console.error)
	}
</script>

<MediaButton {disabled} {onclick} {title}>New Observation</MediaButton>
