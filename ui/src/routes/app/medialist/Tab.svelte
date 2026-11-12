<script>
	import { onMount, getContext } from 'svelte'

	import Button from '#lib/Button.svelte'
	import Tab from '../Tab.svelte'

	const enableTab = getContext('enableTab')
	const gotoTab = getContext('gotoTab')

	let mediaList = $state([])

	function deleteMedia(event) {
		Go('Project.DeleteMedia', this.EntityId) //
			.then(updateMediaList)
			.catch(console.error)
	}

	function updateMediaList() {
		Go('Project.ListMedia') //
			.then((result) => (mediaList = result))
			.catch(console.error)
	}

	function openMedia() {
		const media = this

		Go('Project.SetCurrentMedia', media.EntityId)
			.then(() => {
				enableTab('media', true)
				gotoTab('media')
			})
			.catch(console.error)
	}

	onMount(updateMediaList)
</script>

<Tab>
	<div class="media-list" role="list">
		{#each mediaList as media (media.EntityId)}
			<div class="media-item" role="listitem">
				<Button onclick={openMedia.bind(media)}>{media.Name}</Button>
				<Button onclick={deleteMedia.bind(media)}>Delete</Button>
			</div>
		{/each}
	</div>
</Tab>

<style>
	.media-list {
		flex: 1 0 100px;

		display: flex;
		flex-direction: column;
		gap: 16px;

		width: 100%;
		padding: 16px;

		overflow: scroll;
	}

	.media-item {
		flex: 0 0 40px;
		width: 100%;
		height: 100%;

		border: 1px solid grey;
		border-radius: 4px;

		gap: 1rem;

		display: flex;
	}
</style>
