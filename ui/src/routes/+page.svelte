<script>
	import Video from '$lib/Video.svelte'
</script>

<main>
	<Video src="https://www.w3schools.com/html/movie.mp4" type="video/mp4" />
</main>

<style>
	main {
		width: 100%;
		height: 100%;

		display: flex;
		flex-direction: column;
	}
</style>
