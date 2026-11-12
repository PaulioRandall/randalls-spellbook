<script>
	let { components = {}, defaultValue = '', args = {} } = $props()

	let currentValue = $state(defaultValue)
	let currentArgs = $state(args)
	let current = $derived(components[currentValue])

	export function setValue(value, ...args) {
		currentValue = value
		currentArgs = args
		current = components[currentValue]

		if (!current) {
			current = components[defaultValue]
		}
	}

	export function getValue() {
		return currentValue
	}
</script>

<svelte:component this={current} {...currentArgs} />
