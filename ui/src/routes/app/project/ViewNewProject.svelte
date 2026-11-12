<script>
	import { getContext } from 'svelte'

	import Form from '#lib/Form.svelte'
	import TextInput from '#lib/TextInput.svelte'
	import FileInput from '#lib/FileInput.svelte'
	import Button from '#lib/Button.svelte'

	const enableTab = getContext('enableTab')
	const gotoTab = getContext('gotoTab')

	let filename = $state('')
	let filenameError = $state('')

	let folder = $state('')
	let folderError = $state('')

	let path = $state('')
	let isValidPath = $state(false)

	function submit() {
		validateFields()
			.then((isValid) => {
				if (isValid) {
					createProject()
				}
			})
			.catch(console.error)
	}

	function createProject() {
		if (!isValidPath) {
			return
		}

		Go('Project.Open', path)
			.then(() => enableTab('mediaList', true))
			.then(() => gotoTab('mediaList'))
			.catch(console.error)
	}

	async function validateFields() {
		filenameError = ''
		folderError = ''

		filename = filename.trim()
		folder = folder.trim()

		if (!filename) {
			filenameError = 'Enter a filename'
		}

		if (!folder) {
			folderError = 'Select a folder'
		}

		if (filenameError || folderError) {
			return false
		}

		if (!isValidPath) {
			filenameError = 'File already exists'
			return false
		}

		return true
	}

	function isPathValid() {
		if (filename.trim() && folder.trim()) {
			return Go('FilePicker.DoesFileExist', path)
		}
		return false
	}

	$effect(() => {
		path = folder.trim() + '/' + filename.trim() + '.sqlite'
		isValidPath = isPathValid()
	})
</script>

<Form>
	<TextInput
		label="Filename"
		name="Filename"
		bind:value={filename}
		error={filenameError} />

	<FileInput
		folder
		label="Folder"
		buttonLabel="Select Folder"
		name="Folder"
		error={folderError}
		bind:value={folder} />

	<p>
		<b>File:&nbsp;</b>
		{#if isValidPath}✅{:else}❌{/if}
		{path}
	</p>

	<Button disabled={!isValidPath} type="submit" onclick={submit}>Create</Button>
</Form>
