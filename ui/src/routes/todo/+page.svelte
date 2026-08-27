<script>
	import { onMount } from 'svelte'

	const todos = [
		`Remove EntityId type because it's causing dev mistakes`,
		`Experiment with typed Incantations that involve creating recursively wrapped thunks`,
		`Fix issue with play button and seekbar not updating to match MediaElement`,
		`Look into how to improve error messaging with spells, e.g. could we return the Effect and let the frontend determine check the error?`,
		`Try making Realm embed Spellbook again as it looks like the previous issues we resolved were not related`,
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
