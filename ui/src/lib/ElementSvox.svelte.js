// SVOX (Svelte Adapter Box) classes and their usage follow
// a strict set of rules aimed at providing clarity for
// programmers, minimising reactivity issues that can be
// hard to debug, and minimising existence checking of the
// underlying value:
// 1. Identifiers holding SVOX instances should never be
//    null.
// 2. All public properties must be readonly, i.e. they can
//    only be set internally or through public functions
//    and internal processes. Keep them non-public and use
//    getters to minimise falling foul. Reactivity,
//    particularly Svelte's binding feature, is powerful
//    and highly convenient but can cause difficult to
//    diagnose errors.
// 3. All public properties must be reactive and tracked,
//    i.e. they are created through runes such as $state
//    and $derived, and trigger $effect and $derived runes.
// 2. All public functions must be non-reactive and
//    untracked, i.e. they do not trigger $effect or
//    $derived runes. Use Svelte's untrack function if
//    there's a need to prevent triggering runes when
//    accessing state. Together with the prior rule, this
//    makes it easy for readers to quickly understand the
//    reactivity of some code.
// 4. Users may be allowed to perform certain actions even
//    when the underlying value is not set, e.g. the
//    ElementSvox allows event listeners to be registered
//    and unregistered at anytime. These listeners will be
//    added and removed from the element as it's set and
//    unset.
// 5. Prioritise readability and changability, but use
//    common and iconic verbs for functions that get or
//    set properties as much as possible: e.g. get, is,
//    has, set, put, update, etc.

import { untrack } from 'svelte'
import Eventor from './Eventor.js'

// ElementSvox is a SVOX for the standard Element class.
// It is intended for extension, but may be used as is.
export default class ElementSvox {
	_stateEventor = this._generateStateEventor()
	_eventor = new Eventor()

	// _element is the boxed HTMLElement, which may be null.
	_element = $state(null)

	// _onElement is a set of functions called each time an
	// element is set, including when the elemnet is set to
	// null (unset).
	_onElement = new Set()

	// element is the underlying HTMLElement or null when no
	// element is set.
	//
	// Tracked.
	get element() {
		return this._element
	}

	// hasElement returns true if the element is set.
	//
	// Untracked.
	hasElement() {
		return untrack(() => this._element === null)
	}

	// getElement returns the element if set, else returns
	// null.
	//
	// Untracked.
	getElement() {
		return untrack(() => this._element)
	}

	// onElement adds a listener to be called when an element
	// is set or unset. If the listener is not a function an
	// error is thrown. A function is returned that
	// unregisters the listener when called.
	//
	// Untracked.
	onElement(listener) {
		return untrack(() => {
			const t = typeof listener

			if (t !== 'function') {
				throw new Error(`listener must be a function, not '${t}'`)
			}

			this._onElement.add(listener)
			return () => this.offElement(listener)
		})
	}

	// offElement removes a listener registered via
	// onElement.
	//
	// Untracked.
	offElement(listener) {
		untrack(() => {
			this._onElement.delete(listener)
		})
	}

	// setElement sets the element. If one is already set
	// then unsetElement is called first. If null is passed
	// the function returns after unsetting the prior element.
	// After setting it will call syncState, then call
	// afterSet, then call the element listeners, and finally
	// add the event listeners to the element.
	//
	// Untracked.
	setElement(element) {
		untrack(() => {
			this.unsetElement()

			if (element === null) {
				return
			}

			if (!this.isValidElement(element)) {
				throw new Error(`Invalid element type '${typeof element}'`)
			}

			this._element = element
			this.syncStates()

			this._stateEventor.addTo(this.getElement())
			this.afterSet()

			this._callElementListeners()
			this._eventor.addTo(this._element)
		})
	}

	// unsetElement sets the element to null. It will remove
	// the event listeners first, then set as null, then
	// call syncState, finally call the element listeners.
	//
	// Untracked.
	unsetElement() {
		untrack(() => {
			if (this._element === null) {
				return
			}

			this.beforeUnset()
			this._stateEventor.removeFrom(this.getElement())

			this._eventor.removeFrom(this._element)
			this._element = null

			this.syncStates()
			this._callElementListeners()
		})
	}

	// _callbackListeners calls all element listeners set.
	// All errors are caught and printed to error console.
	_callElementListeners() {
		for (const f of this._onElement) {
			try {
				f(this)
			} catch (err) {
				console.error(err)
			}
		}
	}

	// isValidElement returns true if the element type is a
	// valid HTMLElement. If extending this class to a
	// specific element type, the method should be overridden
	// to return true for that specific type.
	//
	// Untracked.
	isValidElement(element) {
		return element instanceof Element
	}

	// afterSet is called after setting an element and
	// syncing, but not when setting to null (unset). This
	// allows additional custom setup to occur before
	// onElement listeners are called.
	//
	// Untracked.
	afterSet() {
		// Do nothing.
	}

	// beforeUnset is called before unsetting the current
	// element. This allows clean up of additional custom
	// setup.
	//
	// Untracked.
	beforeUnset() {
		// Do nothing.
	}

	// generateStateListeners is called during construction
	// to create the listeners that are added to the element
	// before user listeners. These listeners modify the
	// SVOX's internal state.
	generateStateListeners() {
		return {}
	}

	// syncStates synchronises the state of the class
	// instance with the element. Any overidden
	// implementation must handle when the element is set to
	// null, which usually means resetting state to initial
	// values. syncStates may be called at any time to
	// synchronise state; it's often called by event
	// listeners that were added to the element.
	//
	// Untracked.
	syncStates() {
		// Do nothing.
	}

	// on registers event listeners that are added to the
	// element when set and removed when unest.
	//
	// See Eventor.on docs for more info.
	//
	// Untracked.
	on(...args) {
		return this._eventor.on(...args)
	}

	// off unregisters event listener registered through
	// the on function.
	//
	// See Eventor.off docs for more info.
	//
	// Untracked.
	off(...args) {
		this._eventor.off(...args)
	}

	// isOn returns true if the event listener is currently
	// registered.
	//
	// See Eventor.isOn docs for more info.
	//
	// Untracked.
	isOn(...args) {
		this._eventor.isOn(...args)
	}

	// dispatchEvent dispatches an event on the element,
	// but does nothing if the elemnt is not set.
	//
	// Untracked.
	dispatchEvent(eventType, options) {
		untrack(() => {
			if (this._element) {
				return
			}

			if (!options) {
				options = {
					bubbles: false, //
					cancelable: false,
					composed: false,
				}
			}

			const event = new Event(eventType, options)
			this._element.dispatchEvent(event)
		})
	}

	_generateStateEventor() {
		const listeners = this.generateStateListeners()
		const captureListeners = {}

		for (const eventType in listeners) {
			captureListeners[eventType] = {
				listener: listeners[eventType], //
				options: { capture: true },
			}
		}

		const eventor = new Eventor()
		eventor.on(listeners)
		return eventor
	}
}
