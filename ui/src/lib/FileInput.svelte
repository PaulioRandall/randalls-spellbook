<script>
	import Button from './Button.svelte'
	import TextInput from './TextInput.svelte'
	import FormField from './FormField.svelte'
	import InputText from './InputText.svelte'
	import FilePicker from './FilePicker.svelte'

	let {
		folder = false,
		label = '',
		buttonLabel = 'Select File',
		name = '',
		value = $bindable(''),
		error = '',
	} = $props()

	let filePicker = null

	function openFilePicker() {
		filePicker.open()
	}

	function onPick(path) {
		value = path
	}
</script>

<FormField {label} {name} {error}>
	<Button onclick={openFilePicker}>
		{buttonLabel}
	</Button>

	<InputText readonly {name} bind:value />

	<FilePicker bind:this={filePicker} {folder} {onPick} />
</FormField>
