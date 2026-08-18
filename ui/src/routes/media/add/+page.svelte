<script>
	import { onMount } from 'svelte'
	import { goto } from '$app/navigation'
	import backend from '$lib/backend.js'

	let mediaType = $state('video')
	let mediaTypeError = $state('')

	let name = $state('')
	let nameError = $state('')

	let description = $state('')
	let descriptionError = $state('')

	let localPath = $state('')
	let localPathError = $state('')

	async function selectFile() {
		localPath = await backend.selectLocalMediaFile()

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

	async function submit() {
		if (!validateFields()) {
			return
		}

		const entityId = await backend.addMedia(
			mediaType, //
			name,
			description,
			localPath
		)

		if (entityId) {
			goto('/media')
		}
	}

	function validateFields() {
		nameError = ''
		descriptionError = ''
		localPathError = ''

		name = name.trim()
		description = description.trim()
		localPath = localPath.trim()

		if (!name) {
			nameError = 'Provide a name for the media file'
		}

		if (!localPath) {
			localPathError = 'Select a media file'
		}

		return !nameError && !descriptionError && !localPathError
	}
</script>

<main>
	<div class="form">
		<div class="form-field">
			<label class="form-field-media-type" for="mediaType"> Media type </label>
			{#if mediaTypeError}
				<p class="form-field-error">
					{mediaTypeErrorError}
				</p>
			{/if}
			<!-- TODO: Convert to select -->
			<input readonly type="text" name="mediaType" bind:value={mediaType} />
		</div>
		<div class="form-field">
			<label class="form-field-label" for="name"> Name </label>
			{#if nameError}
				<p class="form-field-error">
					{nameError}
				</p>
			{/if}
			<input type="text" name="name" bind:value={name} />
		</div>
		<div class="form-field">
			<label class="form-field-label" for="description"> Description </label>
			{#if descriptionError}
				<p class="form-field-error">
					{descriptionError}
				</p>
			{/if}
			<textarea name="description" rows="3" bind:value={description} />
		</div>
		<div class="form-field">
			<label class="form-field-label" for="local-path"> Local path </label>
			{#if localPathError}
				<p class="form-field-error">
					{localPathError}
				</p>
			{/if}
			<button onclick={selectFile}> Select Media File </button>
			<input type="text" name="local-path" bind:value={localPath} />
		</div>
		<button type="submit" onclick={submit}>Submit</button>
	</div>
</main>

<style>
	main {
		width: 100%;
		height: 100%;

		display: flex;
		flex-direction: column;
	}

	.form {
		display: flex;
		flex-direction: column;
		gap: 16px;

		padding: 16px;
	}

	.form-field {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.form-field-label {
		/* TODO */
	}

	.form-field-error {
		color: darkred;
	}

	label {
		font-weight: bold;
	}

	p {
		margin: 0;
	}

	button,
	input,
	textarea {
		border-radius: 4px;
		padding: 8px;
	}

	textarea {
		min-height: 60px;
		resize: vertical;
	}

	input,
	textarea {
		border: 1px solid grey;
	}

	button {
		cursor: pointer;
	}
</style>
