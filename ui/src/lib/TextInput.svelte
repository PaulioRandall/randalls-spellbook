<script>
	import FormField from './FormField.svelte'
	import InputText from './InputText.svelte'

	let {
		readonly = false,
		multiline = false,
		rows = 3,
		label = '',
		name = '',
		value = $bindable(''),
		error = '',
	} = $props()
</script>

<FormField {label} {name} {error}>
	<InputText {readonly} {multiline} {name} {rows} bind:value />
</FormField>
