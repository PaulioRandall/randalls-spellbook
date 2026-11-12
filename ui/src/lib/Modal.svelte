<script>
	import { untrack } from 'svelte'

	const id = crypto.randomUUID()

	let { children } = $props()
	let active = $state(false)

	export function isOpen() {
		return untrack(() => active)
	}

	export function open() {
		active = true
	}

	export function close() {
		active = false
	}

	export function toggle() {
		active = !active
	}

	function closeModal(event) {
		if (event.target.id === id) {
			active = false
		}
	}
</script>

{#if active}
	<div {id} class="modal" onclick={closeModal}>
		<div class="modal-content">
			{@render children?.()}
		</div>
	</div>
{/if}

<style>
	.modal {
		position: fixed;
		top: 0;
		left: 0;

		width: 100vw;
		height: 100vh;
		background: rgba(0, 0, 0, 0.5);
	}

	.modal-content {
		--margin: 2rem;

		margin: var(--margin);
		width: calc(100% - (var(--margin) + var(--margin)));
		height: calc(100% - (var(--margin) + var(--margin)));

		background: var(--color-dark_navy);
		border-radius: 6px;

		overflow: hidden;
	}
</style>
