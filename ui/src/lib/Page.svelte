<script>
	let { header, children, ...rest } = $props()
</script>

<div class="page">
	{#if header}
		<header>
			<nav>
				{@render header()}
			</nav>
		</header>
	{/if}

	<main {...rest}>
		{@render children?.()}
	</main>
</div>

<style>
	.page {
		display: flex;
		flex-direction: column;

		width: 100vw;
		min-width: 100vw;
		max-width: 100vw;

		height: 100vh;
		min-height: 100vh;
		max-height: 100vh;

		& > * {
			width: 100%;
		}
	}

	header {
		flex: 0 0 2.5rem;
		padding: 0.5rem;
		border-bottom: 2px solid grey;
	}

	nav {
		display: flex;
		gap: 1rem;
	}

	main {
		flex-grow: 1;

		display: flex;
		flex-direction: column;
		gap: 8px;
	}
</style>
