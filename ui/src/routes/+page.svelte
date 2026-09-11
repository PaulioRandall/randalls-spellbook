<script>
	import { onMount } from 'svelte'
	onMount(() => {
		Go('abc', 123, [1, 2, 3], {
			One: 1,
			Two: 2,
			Three: 3,
		}).catch(console.error)
	})
</script>

<main>
	<div class="menu">
		<a class="option" href="/media">Media</a>
		<a class="option" href="/todo">Todo</a>
	</div>
</main>

<style>
	main {
		width: 100%;
		height: 100%;

		padding: 16px;
	}

	.menu {
		display: flex;
		flex-wrap: wrap;
		gap: 16px;
	}

	.option {
		width: 320px;
		height: 240px;

		display: flex;
		justify-content: center;
		align-items: center;

		border: 1px solid black;
		border-radius: 8px;

		cursor: pointer;
	}
</style>
