<script>
	let { mediaSvox, textareaSvox } = $props()

	function setMediaElement(node) {
		textareaSvox.setElement(node)
		return () => textareaSvox.unsetElement()
	}
</script>

<textarea use:setMediaElement class="textarea"></textarea>

<style>
	.textarea {
	}
</style>
