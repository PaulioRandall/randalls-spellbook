<script>
	import { setContext } from 'svelte'
	import { fade, slide } from 'svelte/transition'

	import Button from '#lib/Button.svelte'

	import ProjectTab from './project/Tab.svelte'
	import MediaListTab from './medialist/Tab.svelte'
	import MediaTab from './media/Tab.svelte'
	import DataTab from './data/Tab.svelte'

	const tabs = $state([
		{
			id: 'project',
			label: 'Project',
			component: ProjectTab,
			disabled: false,
		},
		{
			id: 'mediaList',
			label: 'Media List',
			component: MediaListTab,
			disabled: true,
		},
		{
			id: 'media',
			label: 'Media',
			component: MediaTab,
			disabled: true,
		},
		{
			id: 'data',
			label: 'Data',
			component: DataTab,
			disabled: false,
		},
	])

	let currTab = $state(tabs[0])

	setContext('enableTab', (id, state = true) => {
		const tab = tabs.find((t) => t.id === id)

		if (tab) {
			tab.disabled = !state
		}
	})

	setContext('gotoTab', (id) => {
		const tab = tabs.find((t) => t.id === id)

		if (tab && !tab.disabled) {
			currTab = tab
		}
	})
</script>

<div class="page">
	<header>
		<nav>
			{#each tabs as tab}
				<Button
					disabled={tab.disabled || currTab.id === tab.id}
					onclick={() => (currTab = tab)}>
					{#if currTab.id === tab.id}
						<b>{tab.label}</b>
					{:else}
						{tab.label}
					{/if}
				</Button>
			{/each}
		</nav>
	</header>

	<main>
		{@render currTab.component()}
	</main>
</div>

<style>
	.page {
		display: flex;

		width: 100vw;
		min-width: 100vw;
		max-width: 100vw;

		height: 100vh;
		min-height: 100vh;
		max-height: 100vh;
	}

	header {
		flex: 0 0 8rem;
		height: 100%;
		padding: 0.5rem;
		background: #181010;
	}

	nav {
		height: 100%;

		display: flex;
		flex-direction: column;
		justify-content: center;
		gap: 1rem;
	}

	:global(nav button) {
		padding: 1rem 0;
	}

	main {
		position: relative;

		flex-grow: 1;
		height: 100%;

		display: flex;
		flex-direction: column;

		gap: 1rem;

		background: #303030;
	}
</style>
