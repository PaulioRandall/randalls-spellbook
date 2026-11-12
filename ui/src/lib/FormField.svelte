<script>
	let { label = '', name = '', error = '', children } = $props()
</script>

<div class="form-field">
	{#if label}
		<label for={name}>{label}</label>
	{/if}

	{@render children?.()}

	{#if error}
		<p class="error">
			{error}
		</p>
	{/if}
</div>

<style>
	.form-field {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	label {
		font-weight: bold;
	}

	.error {
		color: darkred;
		margin: 1px;
	}
</style>
