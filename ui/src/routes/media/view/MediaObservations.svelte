<script>
	import { onMount, untrack } from 'svelte'
	import MediaButton from './MediaButton.svelte'

	let { mediaSvox, media } = $props()
	let observationList = $state([])

	$effect(() => updateObservationList(media))

	function updateObservationList(media) {
		untrack(() => {
			if (!!media) {
				Go('ListObservationsByMediaId', media.EntityId) //
					.then((result) => (observationList = result))
					.catch(console.error)
			}
		})
	}

	function jumpToStartTime() {
		mediaSvox.seekTo(this.StartTime)
	}
</script>

{#each observationList as ob (ob.EntityId)}
	<MediaButton onclick={jumpToStartTime.bind(ob)} title={ob.Description}>
		{ob.StartTime} ({ob.Duration})
	</MediaButton>
{/each}
