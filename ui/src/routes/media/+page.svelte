<script>
	import { onMount } from 'svelte'
	import MediaVideo from './MediaVideo.svelte'

	let entityId = $state('')

	onMount(async () => {
		const localPath = await window.selectLocalMediaFile()

		if (!localPath) {
			alert('Refresh page and select a video this time!')
			return
		}

		const name = extractNameFromLocalPath(localPath)
		const description = ''

		entityId = await addVideoToProject(
			name, //
			description,
			localPath
		)
	})

	function extractNameFromLocalPath(localPath) {
		return localPath.split(/[\\/]/).pop().split('.')[0]
	}
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
