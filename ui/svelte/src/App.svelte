<script>
	import P69 from "npm:@paulio/p69@next"

	import Button from "./Button.svelte";
	import TextInput from "./TextInput.svelte";

	let inputText = $state("Hello, world!");
	let displayedText = inputText;

	function updateText() {
		inputText = inputText.trim();
		if (inputText) {
			displayedText = inputText;
		}
	}
</script>

<TextInput bind:value={inputText} />

<Button onclick={updateText}>
	Update text
</Button>

<p>
	{displayedText}
</p>
