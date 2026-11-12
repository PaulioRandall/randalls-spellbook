<script>
	import { onMount } from 'svelte'
	import { goto } from '$app/navigation'

	import Page from '#lib/Page.svelte'
	import Button from '#lib/Button.svelte'
	import Dropdown from '#lib/Dropdown.svelte'
	import Form from '#lib/Form.svelte'
	import TextInput from '#lib/TextInput.svelte'
	import FileInput from '#lib/FileInput.svelte'

	let mediaType = $state('video')
	let mediaTypeError = $state('')

	let name = $state('')
	let nameError = $state('')

	let description = $state('')
	let descriptionError = $state('')

	let localPath = $state('')
	let localPathError = $state('')

	const mediaTypes = {
		video: 'Video',
	}

	function fileSelected(path) {
		localPath = path

		if (!name.trim()) {
			name = extractNameFromLocalPath(localPath)
			name = capitalise(name)
		}
	}

	function extractNameFromLocalPath(localPath) {
		return localPath.split(/[\\/]/).pop().split('.')[0]
	}

	function capitalise(s) {
		return s[0].toUpperCase() + s.slice(1)
	}

	function submit() {
		if (!validateFields()) {
			return
		}

		const media = {
			mediaType, //
			name,
			description,
			localPath,
		}

		Go('Project.AddMedia', media)
			.then((media) => {
				goto(`/media/view?entity_id=${media.EntityId}`)
			})
			.catch(console.error)
	}

	function validateFields() {
		mediaTypeError = ''
		nameError = ''
		descriptionError = ''
		localPathError = ''

		mediaType = mediaType.trim()
		name = name.trim()
		description = description.trim()
		localPath = localPath.trim()

		if (!mediaType) {
			mediaTypeError = 'Select a media type'
		}

		if (!Object.hasOwn(mediaTypes, mediaType)) {
			mediaTypeError = 'Select a media type. The one provided is invalid.'
		}

		if (!name) {
			nameError = 'Name the new media'
		}

		if (!localPath) {
			localPathError = 'Select a media file'
		}

		return (
			!mediaTypeError && //
			!nameError &&
			!descriptionError &&
			!localPathError
		)
	}
</script>

{#snippet header()}
	<a href="/media">Cancel</a>
{/snippet}

<Page {header}>
	<Form>
		<Dropdown
			disabled
			label="Media type"
			name="media-type"
			bind:value={mediaType}
			error={mediaTypeError}
			options={mediaTypes} />

		<TextInput label="Name" name="name" bind:value={name} error={nameError} />

		<TextInput
			multiline
			label="Description"
			name="description"
			bind:value={description}
			error={descriptionError} />

		<FileInput
			label="Local path"
			name="local-path"
			bind:value={localPath}
			error={localPathError} />

		<Button type="submit" onclick={submit}>Submit</Button>
	</Form>
</Page>
