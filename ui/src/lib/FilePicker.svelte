<script>
	import Modal from '#lib/Modal.svelte'

	let { folder = false, onPick } = $props()

	let modal = null
	let dir = $state(null)

	export function open() {
		queryDir('.')
		modal.open()
	}

	async function queryDir(path) {
		dir = await Go('FilePicker.ListFilesInDir', path) //
			.catch(console.error)
	}

	function navToParent() {
		queryDir(dir.AbsParent)
	}

	function enterDir() {
		queryDir(this.AbsPath)
	}

	function selectFile() {
		modal.close()
		onPick(this.AbsPath)
	}

	function selectDir() {
		modal.close()
		onPick(this.AbsPath)
	}

	function cancel() {
		modal.close()
	}
</script>

{#snippet forFolderPicking(entries)}
	{#each entries as entry}
		{#if entry.IsDir}
			<div class="dir-entry">
				<span>Dir:&nbsp;</span>
				<button onclick={selectDir.bind(entry)}>{entry.Name}</button>
				<button onclick={enterDir.bind(entry)}><b>Open</b></button>
			</div>
		{/if}
	{/each}
{/snippet}

{#snippet forFilePicking(entries)}
	{#each entries as entry}
		<div class="dir-entry">
			{#if entry.IsDir}
				<span>Dir:&nbsp;</span>
				<button onclick={enterDir.bind(entry)}>{entry.Name}</button>
			{:else}
				<span>File:</span>
				<button onclick={selectFile.bind(entry)}>{entry.Name}</button>
			{/if}
		</div>
	{/each}
{/snippet}

<Modal bind:this={modal}>
	<div class="file-picker">
		<b class="title">Select a File</b>
		<div class="dir-nav">
			<button onclick={navToParent}>Parent</button>
			<button onclick={cancel}>Cancel</button>
		</div>
		<div class="dir-entries">
			{#if !dir}
				<!-- Show nothing! -->
			{:else if folder}
				{@render forFolderPicking(dir.Entries)}
			{:else}
				{@render forFilePicking(dir.Entries)}
			{/if}
		</div>
	</div>
</Modal>

<style>
	button {
		cursor: pointer;
	}

	.title {
		color: white;
	}

	.file-picker {
		display: flex;
		flex-direction: column;

		padding: 1rem;

		max-height: 100%;
	}

	.dir-nav {
		display: flex;
		gap: 1rem;
		color: white;

		margin: 1rem;
	}

	.dir-entries {
		flex-grow: 1;

		display: flex;
		row-gap: 0.5rem;
		flex-direction: column;

		overflow-y: auto;
	}

	.dir-entry {
		display: flex;
		column-gap: 0.5rem;
		color: white;
	}
</style>
