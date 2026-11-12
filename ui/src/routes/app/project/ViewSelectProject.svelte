<script>
	import { getContext } from 'svelte'

	import Button from '#lib/Button.svelte'
	import FilePicker from '#lib/FilePicker.svelte'
	import Tab from '../Tab.svelte'

	const enableTab = getContext('enableTab')
	const gotoTab = getContext('gotoTab')
	const gotoNewProjectView = getContext('gotoNewProjectView')

	let filePicker = null

	function newProject() {
		gotoNewProjectView()
	}

	function openProject() {
		filePicker.open()
	}

	function onPick(path) {
		Go('Project.Open', path)
			.then(() => enableTab('mediaList', true))
			.then(() => gotoTab('mediaList'))
			.catch(console.error)
	}

	function openTestProject() {
		// TEMP
		Go('Project.Open', './testproject/data.sqlite')
			.then(() => enableTab('mediaList', true))
			.then(() => gotoTab('mediaList'))
			.catch(console.error)
	}
</script>

<div class="menu">
	<Button onclick={newProject}>
		<div class="expander">New Project</div>
	</Button>
	<Button onclick={openProject}>
		<div class="expander">Open Project</div>
	</Button>
	<Button onclick={openTestProject}>
		<div class="expander">Open Test Project</div>
	</Button>
</div>

<FilePicker bind:this={filePicker} {onPick} />

<style>
	.menu {
		height: 100%;

		display: flex;
		justify-content: space-around;
		align-items: center;
		gap: 2rem;
	}

	.expander {
		height: 10rem;
		width: 6rem;

		display: flex;
		justify-content: center;
		align-items: center;

		font-weight: bold;
	}
</style>
