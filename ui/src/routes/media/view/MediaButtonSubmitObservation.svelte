<script>
	import MediaButton from './MediaButton.svelte'
	import backend from '$lib/backend.js'

	const title = 'Submit observation'

	let { mediaSvox, textareaSvox, media } = $props()
	let disabled = $derived(
		!mediaSvox.loaded || textareaSvox.empty //
	)

	$effect(() => {
		console.log('empty:', textareaSvox.empty)
	})

	function onclick() {
		const ob = {
			mediaId: media.entityId,
			startTime: mediaSvox.currentTime,
			duration: 0,
			description: textareaSvox.text,
		}

		window.CastSpell(
			'AddObservation', //
			JSON.stringify(ob)
		)
	}
</script>

<MediaButton {disabled} {onclick} {title}>New Observation</MediaButton>
