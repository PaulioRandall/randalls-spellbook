<script>
	import FormField from './FormField.svelte'

	let {
		disabled = false,
		label = '',
		name = crypto.randomUUID(),
		value = $bindable(''),
		error = '',
		options = {
			// [value]: label
		},
	} = $props()
</script>

<FormField {label} {name} {error}>
	<select {disabled} {name} bind:value>
		{#each Object.entries(options) as [value, label] (value)}
			<option {value}>{label}</option>
		{/each}
	</select>
</FormField>

<style>
	select {
		border-radius: 4px;
	}
</style>
