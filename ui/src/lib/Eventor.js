// ListenerEntry encapsulates the values to be passed to
// EventTarget.addEventListener and
// EventTarget.removeEventListener.
class ListenerEntry {
	// listenerObjectToEntries converts an object with
	// eventTypes as property names and either a listener, or
	// object containing a listener and options, as values
	// into an array of ListenerEntry.
	static listenerObjectToEntries(listenerObject) {
		function createEntry(eventType, listenerOrObject) {
			if (isObject(listenerOrObject)) {
				return new ListenerEntry(
					eventType, //
					listenerOrObject.listener,
					listenerOrObject.options
				)
			}

			return new ListenerEntry(
				eventType, //
				listenerOrObject
			)
		}

		const result = []
		const eventTypes = Object.getOwnPropertyNames(
			listenerObject //
		)

		for (const eventType of eventTypes) {
			const listenerOrObject = listenerObject[eventType]
			const entry = createEntry(eventType, listenerOrObject)
			result.push(entry)
		}

		return result
	}

	// isCapturing returns true if capture is true based upon
	// the passed listener options. It returns true if
	// options is true or is an object containing a 'capture'
	// property which is true.
	static isCapturing(options) {
		if (options === true) {
			return true
		}

		if (isObject(options)) {
			return options.capture === true
		}

		return false
	}

	// normaliseOptions accepts listener options, which may
	// be undefined, null, a bool, or a plain object, and
	// returns a plain object with the option properties.
	// Properties not set in the passed options will not be
	// set in the returned options, and any unrelated
	// properties will be discarded.
	//
	// If options is undefined or null then the result will
	// be an empty object. If a bool is passed then it will
	// have a capture property. If a plain object is passed
	// then the following properties will be set in the
	// result object only if they are present in the passed
	// options: capture, once, passive, and signal. If any
	// other type is passed then an error is thrown.
	static normaliseOptions(options) {
		const result = {}

		if (options === undefined || options === null) {
			return result
		}

		if (typeof options === 'boolean') {
			result.capture = options
			return result
		}

		if (!isObject(options)) {
			throw new Error(
				'Listener options must be undefined, null, a bool, or a plain object'
			)
		}

		if (Object.hasOwn(options, 'capture')) {
			result.capture = options.capture
		}

		if (Object.hasOwn(options, 'once')) {
			result.once = options.once
		}

		if (Object.hasOwn(options, 'passive')) {
			result.passive = options.passive
		}

		if (Object.hasOwn(options, 'signal')) {
			result.signal = options.signal
		}

		return result
	}

	_eventType = ''
	_listener = null
	_options = undefined
	_capture = false

	// constructor assumes the arguments have been validated.
	// The capture property is determined based upon the
	// options argument using ListenerEntry.isCapturing and
	// the options are normalised via
	// ListenerEntry.normaliseOptions. All arguments are
	// validated except for the values of some options. An
	// error is thrown if a validation fails.
	constructor(eventType, listener, options) {
		if (!eventType || typeof eventType !== 'string') {
			throw new Error('eventType must be a string')
		}

		if (!eventType.trim()) {
			throw new Error('eventType must be a non-empty string')
		}

		if (!listener || typeof listener !== 'function') {
			throw new Error('listener must be a function')
		}

		this._eventType = eventType
		this._listener = listener
		this._options = ListenerEntry.normaliseOptions(options)
		this._capture = ListenerEntry.isCapturing(options)
	}

	// eventType is the event type to pass to
	// EventTarget functions. It is always a non-null and
	// non-empty string.
	get eventType() {
		return this._eventType
	}

	// listener is the event listener to pass to
	// EventTarget functions. It is always a non-null
	// function.
	get listener() {
		return this._listener
	}

	// options is the event options to pass to
	// EventTarget functions. It may be undefined, a
	// bool, or a plain object.
	get options() {
		return this._options
	}

	// capture is true if options is true or has property
	// a 'capture' property that is true. This value is used
	// for convenience and not passed to EventTarget
	// functions.
	get capture() {
		return this._capture
	}

	// equals returns true if this ListenerEntry and the
	// argument ListenerEntry are equal from the perspective
	// of EventTarget, i.e. If you try adding both to an
	// EventTarget the second will not be added. It returns
	// true if the argument is an instance of Eventor and
	// are strictly equal for eventType, listener, and
	// capture properties. Note that some options can differ
	// and are not compared unless checkOptions is true,
	// which will also use strict equality on each option.
	equals(other, checkOptions = false) {
		if (!(other instanceof ListenerEntry)) {
			return false
		}

		if (this._eventType !== other._eventType) {
			return false
		}

		if (this._listener !== other._listener) {
			return false
		}

		const equals = this._capture === other._capture

		if (!equals || !checkOptions) {
			return equals
		}

		const thisOpts = this._options
		const otherOpts = other._options

		if (thisOpts.capture !== otherOpts.capture) {
			return false
		}

		if (thisOpts.once !== otherOpts.once) {
			return false
		}

		if (thisOpts.passive !== otherOpts.passive) {
			return false
		}

		if (thisOpts.signal !== otherOpts.signal) {
			return false
		}

		return true
	}

	// addTo adds the listener to the passed eventTarget.
	addTo(eventTarget) {
		eventTarget.addEventListener(
			this._eventType, //
			this._listener,
			this._options
		)
	}

	// removeFrom removes the listener from the passed
	// eventTarget.
	removeFrom(eventTarget) {
		eventTarget.removeEventListener(
			this._eventType, //
			this._listener,
			this._options
		)
	}
}

// Eventor is a store of event lsiteners that can be added
// and removed from instances of EventTarget.
export default class Eventor {
	// newListenerEntry creates and returns a new listener
	// entry. This is primarily for testing.
	static newListenerEntry(eventType, listener, options) {
		return new ListenerEntry(eventType, listener, options)
	}

	// _entries is a unique set of ListenerEntry instances.
	_entries = []

	// on registers an event listener that is added and
	// removed from elements. Listeners are added to
	// EventTargets in the order they are registered here.
	// Either an eventType, listener, and optional options
	// may be passed or a plain object with property names as
	// event types and values as either a listener or a plain
	// object containing a listener and options. The latter
	// approach enables bulk adding and removing of
	// listeners.
	//
	// For each listener, if the unique combination of
	// eventType, listener, and capturing phase were added,
	// then the combination already exists and will not be
	// registered.
	//
	// Returns a function that unregisters the listener or
	// listeners when called.
	on(eventTypeOrObject, listener, options = undefined) {
		if (isObject(eventTypeOrObject)) {
			return this._onAll(eventTypeOrObject)
		}

		const entry = new ListenerEntry(
			eventTypeOrObject, //
			listener,
			options
		)

		return this._on(entry)
	}

	// _onAll registers all event listeners that are own
	// properties of a plain object.
	_onAll(listenerObject) {
		const listeners = ListenerEntry.listenerObjectToEntries(
			listenerObject //
		)
		const unregFuncs = []

		for (const entry of listeners) {
			unregFuncs.push(this._on(entry))
		}

		return () => unregFuncs.forEach((f) => f())
	}

	// _on registers a single event listener entry.
	_on(entry) {
		if (this._indexOfEntry(entry) === -1) {
			this._entries.push(entry)
		}

		return () => this._off(entry)
	}

	// off unregisters a listeners registered through the on
	// function.
	off(eventTypeOrObject, listener, options = undefined) {
		if (isObject(eventTypeOrObject)) {
			this._offAll(eventTypeOrObject)
			return
		}

		const entry = new ListenerEntry(
			eventTypeOrObject, //
			listener,
			options
		)

		this._off(entry)
	}

	// off unregisters the listeners registered through the
	// _onAll function.
	_offAll(listenerObject) {
		const listeners = ListenerEntry.listenerObjectToEntries(
			listenerObject //
		)

		for (const entry of listeners) {
			this._off(entry)
		}
	}

	// _off unregisters a listener entry registered through
	// the _on function.
	_off(entry) {
		const i = this._indexOfEntry(entry)

		if (i > -1) {
			this._entries.splice(i, 1)
		}
	}

	// isOn returns true if the combination of eventType,
	// listener, and capturing phase (determined by the
	// options argument) is currently registered.
	isOn(eventType, listener, options = undefined) {
		const entry = new ListenerEntry(
			eventType, //
			listener,
			options
		)

		return this._indexOfEntry(entry) > -1
	}

	// addTo adds the registered listeners to the passed
	// eventTarget.
	addTo(eventTarget) {
		for (const entry of this._entries) {
			entry.addTo(eventTarget)
		}
	}

	// removeFrom removes the registered listeners from the
	// passed eventTarget.
	removeFrom(eventTarget) {
		for (const entry of this._entries) {
			entry.removeFrom(eventTarget)
		}
	}

	// _indexOfEntry returns the index of the passed entry,
	// or -1 if no entry is found.
	_indexOfEntry(entry) {
		const entries = this._entries

		for (let i = 0; i < entries.length; i++) {
			if (entries[i].equals(entry)) {
				return i
			}
		}

		return -1
	}
}

// isObject returns true if v is a non-null plain object.
function isObject(v) {
	return (
		typeof v === 'object' && //
		!Array.isArray(v) &&
		v !== null
	)
}
