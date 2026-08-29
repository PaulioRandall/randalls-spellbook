<script>
	import MediaButton from './MediaButton.svelte'

	const title = 'Submit observation'

	let { mediaSvox, textareaSvox, media } = $props()
	let disabled = $derived(
		!mediaSvox.loaded || textareaSvox.empty //
	)

	function onclick() {
		let ob = {
			mediaId: media.entityId,
			startTime: mediaSvox.currentTime,
			duration: 0,
			description: textareaSvox.text,
		}

		ob = window.CastSpell(
			'AddObservation', //
			JSON.stringify(ob)
		)

		if (ob) {
			textareaSvox.setText('')
		}
	}
</script>

<MediaButton {disabled} {onclick} {title}>New Observation</MediaButton>
