<script>
	import MediaButton from './MediaButton.svelte'
	import backend from '$lib/backend.js'

	const title = 'Submit observation'

	let { mediaSvox, textareaSvox } = $props()
	let disabled = $derived(
		!mediaSvox.loaded || textareaSvox.empty //
	)

	$effect(() => {
		console.log('empty:', textareaSvox.empty)
	})

	function onclick() {
		const data = {
			time: mediaSvox.currentTime,
			description: textareaSvox.text,
		}

		console.table(data)

		/*
		backend.CastSpell('AddObservation', {
			time: mediaSvox.currentTime,
			description: textareaSvox.text,
		})
		*/
	}
</script>

<MediaButton {disabled} {onclick} {title}>New Observation</MediaButton>
