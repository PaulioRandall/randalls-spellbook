<script>
	import { onMount, setContext } from 'svelte'

	import Tab from '../Tab.svelte'
	import ViewSpinner from './ViewSpinner.svelte'
	import ViewSelectProject from './ViewSelectProject.svelte'
	import ViewDisplayProject from './ViewDisplayProject.svelte'
	import ViewNewProject from './ViewNewProject.svelte'

	import ComponentSwitch from './ComponentSwitch.svelte'

	let view = null

	setContext('loadProject', loadProject)
	setContext('gotoProjectView', (name) => view.setValue(name))

	function loadProject() {
		Go('Project.Path') //
			.then(onProjectLoaded)
			.catch(console.error)
	}

	function onProjectLoaded(projectPath) {
		view.setValue('working')

		if (projectPath != '') {
			view.setValue('project')
		} else {
			view.setValue('menu')
		}
	}

	onMount(loadProject)
</script>

<Tab>
	<ComponentSwitch
		bind:this={view}
		defaultValue="menu"
		components={{
			menu: ViewSelectProject, // Default, no project opened
			working: ViewSpinner,
			project: ViewDisplayProject,
			new: ViewNewProject,
		}} />
</Tab>
