import { untrack } from 'svelte'
import ElementSvox from './ElementSvox.svelte.js'

// TextareaSvox is an ElementSvox extended for the standard
// HTMLTextAreaElement class.
export default class TextareaSvox extends ElementSvox {
	// isValidElement overides, HTMLTextAreaElement only.
	//
	// Untracked.
	isValidElement(element) {
		return element instanceof HTMLTextAreaElement
	}

	// text is the user's text input.
	//
	// Tracked.
	_text = $state('')
	get text() {
		return this._text
	}

	// getText returns the untracked value of the text
	// property.
	//
	// Untracked.
	getText() {
		return untrack(() => this._text)
	}

	// setText sets the text value of the
	// HTMLTextAreaElement.
	setText(value) {
		untrack(() => {
			if (this.hasElement()) {
				this.getElement().value = value
				this.syncStates()
			}
		})
	}

	// empty is true when the textarea contains no text or
	// whitespace.
	//
	// Tracked.
	_empty = $state(true)
	get empty() {
		return this._empty
	}

	// isEmpty returns the untracked value of the empty
	// property.
	//
	// Untracked.
	isEmpty() {
		return untrack(() => this._empty)
	}

	// syncStates overides.
	//
	// Untracked.
	syncStates() {
		untrack(() => {
			const elem = this.getElement()

			if (!elem) {
				this._text = ''
				this._empty = true
				return
			}

			this._text = elem.value
			this._empty = !this._text.trim()
		})
	}

	// generateStateListeners overides.
	generateStateListeners() {
		const svox = this

		function input() {
			svox.syncStates()
		}

		return {
			input, //
		}
	}
}
