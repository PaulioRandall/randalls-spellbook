<script>
	import { onMount } from 'svelte'
	import { page } from '$app/state'

	import Page from '#lib/Page.svelte'
	import Button from '#lib/Button.svelte'
	import MediaVideo from './MediaVideo.svelte'

	let entityId = $state('')

	onMount(() => {
		entityId = page.url.searchParams.get('entity_id')
	})
</script>

{#snippet header()}
	<a href="/media">Back</a>
{/snippet}

<Page {header}>
	{#if entityId}
		<MediaVideo {entityId} />
	{/if}
</Page>
