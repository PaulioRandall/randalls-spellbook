<script>
	import MediaButton from './MediaButton.svelte'

	const title = 'Submit observation'

	let { mediaSvox, textareaSvox, media } = $props()
	let disabled = $derived(
		!mediaSvox.loaded || textareaSvox.empty //
	)

	function onclick() {
		const ob = {
			mediaId: media.entityId,
			startTime: mediaSvox.currentTime,
			duration: 0,
			description: textareaSvox.text,
		}

		CastSpell('AddObservation', JSON.stringify(ob))
			.then(() => textareaSvox.setText(''))
			.catch(console.error)
	}
</script>

<MediaButton {disabled} {onclick} {title}>New Observation</MediaButton>
