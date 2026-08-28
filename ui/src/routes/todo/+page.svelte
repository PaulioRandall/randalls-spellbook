<script>
	import { onMount } from 'svelte'

	const todos = [
		`Experiment with typed Incantations that involve creating recursively wrapped thunks`,
		`Fix issue with play button and seekbar not updating to match MediaElement`,
		`Look into how to improve error messaging with spells, e.g. could we return the Effect and let the frontend determine check the error?`,
	]
</script>

<main>
	<div>
		<a href="/">&larr; Home</a>
	</div>
	<h2>TODO:</h2>
	<ul>
		{#each todos as todo}
			<li>{todo}</li>
		{/each}
	</ul>
</main>

<style>
	main {
		width: 100%;
		height: 100%;

		display: flex;
		flex-direction: column;

		padding: 1rem;
	}
</style>
