<script>
	let { mediaSvox, media } = $props()

	function setMediaElement(node) {
		mediaSvox.setElement(node)
		return () => mediaSvox.unsetElement()
	}
</script>

<video
	use:setMediaElement
	class="media-video-player"
	class:media-video-loaded={mediaSvox.loaded}
	width="320"
	height="240"
	title={media?.name}
	alt={media?.description}
	onclick={() => mediaSvox.playPause()}>
	{#if !media}
		Waiting for media...
	{:else if media?.entityId}
		<source
			src="/media?entity_id={encodeURI(media.entityId)}"
			type="video/mp4" />
	{:else}
		HTML videos not supported by the WebView.
	{/if}
</video>

<style>
	.media-video-player {
	}

	.media-video-loaded {
		object-fit: contain;

		width: 100%;
		height: auto;

		max-width: 100%;
		max-height: 100%;
	}
</style>
