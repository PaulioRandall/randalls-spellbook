<script>
	import { onMount } from 'svelte'
	import { getContext } from 'svelte'
	import Button from '#lib/Button.svelte'

	const enableTab = getContext('enableTab')
	const gotoProjectView = getContext('gotoProjectView')

	let path = $state('')

	onMount(() => {
		Go('Project.Path') //
			.then((p) => (path = p))
			.catch(console.error)
	})

	function closeProject() {
		path = ''
		Go('Project.Close') //
			.catch(console.error)
			.finally(() => enableTab('media', false))
			.finally(() => enableTab('mediaList', false))
			.finally(() => gotoProjectView('menu'))
	}
</script>

<div class="view-display-project">
	<dl>
		<dt>Path:</dt>
		<dd>{path}</dd>
	</dl>

	<Button onclick={closeProject}>Close Project</Button>
</div>

<style>
	.view-display-project {
		color: white;
	}
</style>
