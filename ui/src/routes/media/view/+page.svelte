<script>
	import { onMount } from 'svelte'
	import { page } from '$app/state'

	import MediaVideo from './MediaVideo.svelte'

	let entityId = $state('')

	onMount(() => {
		entityId = page.url.searchParams.get('entity_id')
	})
</script>

<main>
	{#if entityId}
		<MediaVideo {entityId} />
	{/if}
</main>

<style>
	main {
		width: 100%;
		height: 100%;

		display: flex;
		flex-direction: column;
	}
</style>
