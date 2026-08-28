import { untrack } from 'svelte'
import ElementSvox from './ElementSvox.svelte.js'
import Eventor from './Eventor.js'

// TextareaSvox is a ElementSvox specific for the standard
// HTMLTextAreaElement class.
export default class TextareaSvox extends ElementSvox {
	// isValidElement override to constrain elements to
	// HTMLTextAreaElement only.
	//
	// Untracked.
	isValidElement(element) {
		return element instanceof HTMLTextAreaElement
	}

	// text is the users text input.
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

	// syncStates performs a state syncing with the currently
	// set HTMLTextAreaElement.
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

	// reset removes all text from the textarea.
	//
	// Untracked.
	reset() {
		untrack(() => {
			if (this.hasElement()) {
				this.getElement().value = ''
				this.syncStates()
			}
		})
	}

	// generateStateListeners overide to return monitor value
	// state.
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
