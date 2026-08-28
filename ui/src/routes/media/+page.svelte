<script context="module">
	// This page lists all the media added to the current
	// project.
	//
	// Users can:
	// + remove a media from the project
	//
	// Users can access pages to:
	// + add a new media
	// + edit an existing media
	// + view a media
</script>

<script>
	import { onMount } from 'svelte'
	import backend from '$lib/backend.js'

	let mediaList = $state([])
	let selectedMedia = $state(null)

	function selectMedia(event) {
		const entityId = event.target.id

		selectedMedia = mediaList.find((m) => {
			return m.entityId === entityId
		})
	}

	onMount(async () => {
		mediaList = await backend.ListMedia()

		console.log(mediaList[0])
	})
</script>

<main>
	<div class="media-menu-bar">
		<a href="/media/add">Add New Media</a>
	</div>
	<div class="media-list" role="list">
		{#each mediaList as media (media.entityId)}
			<div class="media-item" role="listitem">
				<div
					id={media.entityId}
					class="selectable-media"
					role="button"
					onclick={selectMedia}>
					<span class="media-name">{media.name}</span>
					<a
						class="media-view-button"
						href="/media/view?entity_id={media.entityId}">
						View
					</a>
				</div>
			</div>
		{/each}
	</div>
	<p class="selected-media-description">
		{#if selectedMedia}
			<span class="media-name">{selectedMedia.name}:</span>
			{selectedMedia.description}
		{/if}
	</p>
</main>

<style>
	main {
		flex-grow: 1;
		width: 100%;

		display: flex;
		flex-direction: column;
		gap: 8px;

		height: 100%;
		max-height: 100%;
	}

	.media-menu-bar {
		flex: 0 0 40px;

		padding: 8px;
		border-bottom: 2px solid grey;
	}

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
	}

	.media-item {
		width: 100%;
		height: 100%;
	}

	.selectable-media {
		cursor: pointer;

		&:hover {
			background: #cccccc;
		}

		display: flex;
		gap: 16px;
		justify-content: space-between;

		border: 1px solid grey;
		border-radius: 4px;

		padding: 8px;
	}

	.media-name {
		font-weight: bold;
	}

	.media-view-button {
		cursor: pointer;

		&:hover {
			background: #cccccc;
		}
	}

	.selected-media-description {
		flex: 0 0 100px;
		width: 100%;

		margin: 0;
		padding: 16px;

		border-top: 2px solid grey;
	}
</style>
