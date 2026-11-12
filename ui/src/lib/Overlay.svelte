<script>
	import { untrack } from 'svelte'

	let { children } = $props()
	let active = $state(true)

	export function isShowing() {
		return untrack(() => active)
	}

	export function show() {
		active = true
	}

	export function hide() {
		active = false
	}

	export function toggle() {
		active = !active
	}
</script>

{#if active}
	<div class="overlay">
		<div class="content">
			{@render children?.()}
		</div>
	</div>
{/if}

<style>
	.overlay {
		position: fixed;
		top: 0;
		left: 0;

		width: 100vw;
		height: 100vh;

		pointer-events: none;
	}

	.content {
		position: relative;
		width: 100%;
		height: 100%;
		overflow: hidden;
	}
</style>
