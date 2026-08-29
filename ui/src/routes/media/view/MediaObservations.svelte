<script>
	import { onMount, untrack } from 'svelte'
	import MediaButton from './MediaButton.svelte'

	let { mediaSvox, media } = $props()
	let observationList = $state([])

	$effect(() => {
		if (!!media) {
			updateObservationList(media.entityId)
		}
	})

	function updateObservationList(entityId) {
		untrack(() => {
			CastSpell('ListObservationsByMediaId', entityId) //
				.then((result) => (observationList = result))
				.catch(console.error)
		})
	}

	function jumpToStartTime() {
		mediaSvox.seekTo(this.startTime)
	}
</script>

{#each observationList as ob (ob.entityId)}
	<MediaButton onclick={jumpToStartTime.bind(ob)} title={ob.description}>
		{ob.startTime} to {ob.startTime + ob.duration}
	</MediaButton>
{/each}
