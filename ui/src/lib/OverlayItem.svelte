<script>
	let {
		children,
		top = true,
		bottom = false,
		yCenter = false,
		left = true,
		right = false,
		xCenter = false,
	} = $props()

	let xClass = 'x-left'
	let yClass = 'y-top'

	if (right) {
		xClass = 'x-right'
	}

	if (xCenter) {
		xClass = 'x-center'
	}

	if (bottom) {
		yClass = 'y-bottom'
	}

	if (yCenter) {
		yClass = 'y-center'
	}
</script>

<div class="overlay-item {xClass} {yClass}">
	{@render children?.()}
</div>

<style>
	.overlay-item {
		position: absolute;
		--pad: 1rem;

		pointer-events: auto;
	}

	.y-top {
		top: var(--pad);
	}

	.y-bottom {
		bottom: var(--pad);
	}

	.y-center {
		top: 50%;
		transform: translateY(-50%);
	}

	.x-left {
		left: var(--pad);
	}

	.x-right {
		right: var(--pad);
	}

	.x-center {
		left: 50%;
		transform: translateX(-50%);
	}
</style>
