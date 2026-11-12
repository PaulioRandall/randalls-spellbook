<script>
	import { onMount } from 'svelte'
	import Page from '#lib/Page.svelte'
	import Button from '#lib/Button.svelte'

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

	onMount(updateMediaList)
</script>

{#snippet header()}
	<a href="/">Back</a>
	<a href="/media/add">Add New Media</a>
{/snippet}

<Page {header}>
	<div class="media-list" role="list">
		{#each mediaList as media (media.EntityId)}
			<div class="media-item" role="listitem">
				<a class="media-name" href="/media/view?entity_id={media.EntityId}">
					{media.Name}
				</a>
				<Button onclick={deleteMedia.bind(media)}>Delete</Button>
			</div>
		{/each}
	</div>
</Page>

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

		display: flex;
	}

	.media-name {
		flex-grow: 1;
		cursor: pointer;

		&:hover {
			background: #cccccc;
		}
	}
</style>
