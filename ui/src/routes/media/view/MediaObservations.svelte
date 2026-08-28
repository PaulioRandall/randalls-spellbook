<script>
	import { onMount } from 'svelte'
	import MediaButton from './MediaButton.svelte'

	let { mediaSvox, media } = $props()
	let observationList = $state([])

	function jumpToStartTime() {
		mediaSvox.seekTo(this.startTime)
	}

	onMount(async () => {
		observationList = await window.CastSpell(
			'ListObservationsByMediaId', //
			media.entityId
		)
	})
</script>

{#each observationList as ob (ob.entityId)}
	<MediaButton onclick={jumpToStartTime.bind(ob)} title={ob.description}>
		{ob.startTime} to {ob.startTime + ob.duration}
	</MediaButton>
{/each}
