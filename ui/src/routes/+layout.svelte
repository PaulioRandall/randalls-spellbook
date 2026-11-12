<script>
	import '../app.css'
	import { goto } from '$app/navigation'

	import Overlay from '#lib/Overlay.svelte'
	import OverlayItem from '#lib/OverlayItem.svelte'
	import Button from '#lib/Button.svelte'

	let { children } = $props()
	let overlay = null

	async function openTestProject() {
		// Auto closes current project.
		Go('Project.Open', './testproject/data.sqlite')
			.then(() => {
				overlay.hide()
				goto('/media')
			})
			.catch(console.error)
	}
</script>

{@render children()}

<Overlay bind:this={overlay}>
	<OverlayItem bottom right>
		<Button onclick={openTestProject}>
			<div class="padder">Open test project</div>
		</Button>
	</OverlayItem>
</Overlay>

<style>
	.padder {
		padding: 1rem;
		font-weight: bold;
	}
</style>
