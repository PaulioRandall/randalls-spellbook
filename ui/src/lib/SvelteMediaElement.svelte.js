import { untrack } from 'svelte'

// TODO: Create test suite!! Good luck :)git status

// HTML_EVENT_TYPES is an array of all
// HTMLMediaELement event types that users can register
// listeners for.
const HTML_EVENT_TYPES = [
	'abort', //
	'canplay',
	'canplaythrough',
	'durationchange',
	'emptied',
	'ended',
	'error',
	'loadeddata',
	'loadedmetadata',
	'loadstart',
	'pause',
	'play',
	'playing',
	'progress',
	'ratechange',
	'seeked',
	'seeking',
	'stalled',
	'suspend',
	'timeupdate',
	'volumechange',
	'waiting',
	'waitingforkey',
]

// CUSTOM_EVENT_TYPES is an array of all custom event
// types fired by SvelteMediaElement, not by the
// underlying HTMLMediaELement.
const CUSTOM_EVENT_TYPES = [
	'elementset', //
	'elementunset',
	'running',
	'flowing',
	'buffering',
	'pausing',
]

const EVENT_TYPES = [
	...HTML_EVENT_TYPES, //
	...CUSTOM_EVENT_TYPES,
]

// SvelteMediaElement is a Reactive Adapter Box (READOX)
// for the standard HTMLMediaElement class.
//
// The primary purpose was to decouple Svelte components
// from the underlying HTMLMediaElement. Because the
// HTMLMediaElement is set and unset separate to
// instance construction, the instance passed to components
// can always be non-null. This moves most existence
// checking from the components into the class, where it
// can be better abstracted.
//
// The class stores listeners added by components and other
// code so when a HTMLMediaElement is set and unset the
// listeners will automatically be added and removed
// respectively. Components using this class can register
// 'elementset' and 'elementunset' listeners so they know
// when the underlying HTMLMediaElement is set and unset.
//
// State is updated during the capturing phase of events
// before user registered events so all user listeners
// have access to fresh state through this class. Any
// capturing listeners added to the HTMLMediaElement before
// them will fire before state is updated. Furthermore, if
// a listener i registered through this class is
// unregistered manually it will not prevent the listener
// from being added next time a HTMLMediaElement is set.
// Unless you have a particualar niche requirement, it is
// recommended that all listeners are registered and
// unregistered through this class to avoid nasty, hard to
// debug, errors.
//
// All properties are exposed as reactive state so they can
// be used directly in components. However, untracked
// values can sometimes be useful within the $effect rune
// so a full set of untracked 'get' and 'is' functions are
// provided. All public functions are untracked.
//
// This adapter has a number of additional fields and
// functions for common use cases. For instance,
// currentRemaining (duration - currentTime) is updated
// along with currentTime. Either the flowing or buffering
// property will be true when the media is not paused.
// The media is flowing when there some data is loaded,
// i.e. the user perceives the media to be playing. Where
// as the buffering property is true when waiting for data.
//
// 1. Straight after the metadata has loaded:
//    + loaded, paused, seekable (unless live stream)
//    - seeking, playable, running, flowing, buffering
//
// 2. Once data has loaded, the user has pressed play and
//    the media playing without needing to wait for data:
//    + loaded, seekable, playable, running, flowing
//    - paused, seeking, buffering
//
// 3. When the user is seeking far into the future during
//    the middle of playback:
//    + loaded, seekable, playable, running, seeking
//    - paused, buffering, flowing
//
// Documentation for each property and function will end
// with 'Tracked' if its use is tracked by Svelte, else it
// will end with 'Untracked'. This allows user to quickly
// see if use of a particular property or function will
// trigger execution of $effect runes.
// Properties and their getters should be tracked while
// public functions should be untracked. This makes it even
// easier for readers to quickly understand the reactivity
// of some code.
export default class SvelteMediaElement {
	// _stateListeners are capturing event listeners that
	// mutate the adpaters state. They are added first when
	// an element is set so they fire first.
	// Format: [{ eventType, listener, options }]
	_stateListeners = generateStateListeners(this)

	// _userListeners are the user defined listeners. They
	// include HTMLMediaElement event types along with
	// special custom types that are fired by this class,
	// e.g. 'elementset' and 'elementunset'.
	_userListeners = [
		/* {
			eventType,
			listener,
			options,

			// capture is true if options is true or
			// { capture: true }. This value is only used within
			// the class for ease of listener management.
			capture,
		} */
	]

	// element is the underlying HTMLMediaElement or null
	// when no element is set.
	//
	// Tracked.
	_element = $state(null)
	get element() {
		return this._element
	}

	// hasElement returns true if the element property is
	// set, i.e. not null.
	//
	// Untracked.
	hasElement() {
		return untrack(() => this._element === null)
	}

	// getElement returns the untracked value of the element
	// property.
	//
	// Untracked.
	getElement() {
		return untrack(() => this._element)
	}

	// setElement sets the underlying HTMLMediaElement. If
	// one is already set then unsetElement is called first
	// causing relevant events to fire.
	//
	// Untracked.
	setElement(element) {
		untrack(() => {
			this.unsetElement()

			const validType = element instanceof HTMLMediaElement
			if (!element || !validType) {
				throw new Error('Not a HTMLMediaElement')
			}

			this._element = element
			this._syncStatesInit()

			addListenersToElement(this._element, this._stateListeners)
			addListenersToElement(this._element, this._userListeners)

			dispatchEvent(this._element, 'elementset')
		})
	}

	// unsetElement sets the underlying HTMLMediaElement
	// to null. It will reset all state then fire the
	// 'elementunset' event.
	//
	// Untracked.
	unsetElement() {
		untrack(() => {
			if (this._element === null) {
				return
			}

			removeListenersToElement(this._element, this._userListeners)
			removeListenersToElement(this._element, this._stateListeners)

			this._element = null
			this._resetStates()

			dispatchEvent(this._element, 'elementunset')
		})
	}

	// loaded is true when the media metadata has loaded.
	// It tells you nothing about data loading.
	//
	// Tracked.
	_loaded = $state(false)
	get loaded() {
		return this._loaded
	}

	// isLoaded returns the untracked value of the loaded
	// property.
	//
	// Untracked.
	isLoaded() {
		return untrack(() => this._loaded)
	}

	//_playable is true when the media has data loaded for
	// some playback to occur.
	//
	// Tracked.
	_playable = $state(false)
	get playable() {
		return this._playable
	}

	// isPlayable returns the untracked value of the playable
	// property.
	//
	// Untracked.
	isPlayable() {
		return untrack(() => this._playable)
	}

	// running is true when paused is false.
	//
	// A media can be both running and not flowing which
	// almost always means the media doesn't have enough data
	// to continue playback right now but will continue when
	// data arrives.
	//
	// Tracked.
	_running = $state(false)
	get running() {
		return this._running
	}

	// isRunning returns the untracked value of the running
	// property.
	//
	// Untracked.
	isRunning() {
		return untrack(() => this._running)
	}

	// flowing is true when the media is not paused, not
	// ended, not seeking, and has enough data to play some
	// frames. When flowing is true then the user is actively
	// consuming the media, i.e. the video or audio are
	// playing without buffering. When the media freezes to
	// load more data the flowing state becomes false until
	// data is available.
	//
	// One can assume that the user's attention is on the
	// media when flowing is true (at least for video).
	// However, the flowing state is not enougth to determine
	// if a buffering icon should be shown because it will
	// also be false when the media is paused or not loaded.
	// Instead, the buffering state can be used.
	//
	// Tracked.
	_flowing = $state(false)
	get flowing() {
		return this._flowing
	}

	// isFlowing returns the untracked value of the flowing
	// property.
	//
	// Untracked.
	isFlowing() {
		return untrack(() => this._flowing)
	}

	// buffering is true when the media is not paused, not
	// ended, not seeking, but does not have enough data to
	// continue playing frames. The buffering state cannot be
	// true whilst the flowing state is true but both can be
	// false at once, such as when the media is paused.
	//
	// When buffering the media will appear frozen from the
	// user's perspective. The freeze is likely to have
	// broken the user's attention on the media (at least for
	// video) so the buffering state can be used
	// independently in a Svelte {#if} block to control the
	// display of a buffering icon.
	//
	// Tracked.
	_buffering = $state(false)
	get buffering() {
		return this._buffering
	}

	// isBuffering returns the untracked value of the
	// buffering property.
	//
	// Untracked.
	isBuffering() {
		return untrack(() => this._buffering)
	}

	// paused is true prior to loading, straight after
	// loading, before autoplay begins (if enabled), and when
	// the user has paused the media. It will still be false
	// when the media has stopped playing due to needing more
	// data, unlike the flowing field..
	//
	// Tracked.
	_paused = $state(true)
	get paused() {
		return this._paused
	}

	// isPaused returns the untracked value of the paused
	// property.
	//
	// Untracked.
	isPaused() {
		return untrack(() => this._paused)
	}

	// seekable is true when the underlying HTMLMediaElement
	// is seekable.
	//
	// Tracked.
	_seekable = $state(false)
	get seekable() {
		return this._seekable
	}

	// isSeekable returns the untracked value of the seekable
	// property.
	//
	// Untracked.
	isSeekable() {
		return untrack(() => this._seekable)
	}

	// seeking is true when the underlying HTMLMediaElement
	// is in seeking mode.
	//
	// Tracked.
	_seeking = $state(false)
	get seeking() {
		return this._seeking
	}

	// isSeeking returns the untracked value of the seeking
	// property.
	//
	// Untracked.
	isSeeking() {
		return untrack(() => this._seeking)
	}

	// duration is the total duration of the current media.
	// If element is not set or loaded then this will be 0.
	//
	// Tracked.
	_duration = $state(0)
	get duration() {
		return this._duration
	}

	// getDuration returns the untracked value of the
	// duration property.
	//
	// Untracked.
	getDuration() {
		return untrack(() => this._duration)
	}

	// currentTime is the current playback time. This
	// maps directly to currentTime within the underlying
	// HTMLMediaElement.
	//
	// Tracked.
	_currentTime = $state(0)
	get currentTime() {
		return this._currentTime
	}

	// getCurrentTime returns the untracked value of the
	// currentTime property.
	//
	// Untracked.
	getCurrentTime() {
		return untrack(() => this._currentTime)
	}

	// currentRemaining is the amount of time remaining in
	// playback based on the currentTime. It updates inline
	// with the currentTime field.
	//
	// Tracked.
	_currentRemaining = $state(0)
	get currentRemaining() {
		return this._currentRemaining
	}

	// getCurrentRemaining returns the untracked value of the
	// currentRemaining property.
	//
	// Untracked.
	getCurrentRemaining() {
		return untrack(() => this._currentRemaining)
	}

	// playtime is the current playback time, however it is
	// only updated during playback. This differs from
	// currentTime which is also updated during seeking.
	//
	// Tracked.
	_playtime = $state(0)
	get playtime() {
		return this._playtime
	}

	// getPlaytime returns the untracked value of the
	// playtime property.
	//
	// Untracked.
	getPlaytime() {
		return untrack(() => this._playtime)
	}

	// seektime is the current seek time which is only
	// updated when seeking is true.
	//
	// Tracked.
	_seektime = $state(0)
	get seektime() {
		return this._seektime
	}

	// getSeektime returns the untracked value of the
	// seektime property.
	//
	// Untracked.
	getSeektime() {
		return untrack(() => this._seektime)
	}

	// remainingTime is the amount of time remaining in
	// playback based on the playtime, thus updates inline
	// with the playtime field.
	//
	// Tracked.
	_remainingTime = $state(0)
	get remainingTime() {
		return this._remainingTime
	}

	// getRemainingTime returns the untracked value of the
	// getRemainingTime property.
	//
	// Untracked.
	getRemainingTime() {
		return untrack(() => this._remainingTime)
	}

	// reload reloads the media. Load related functions will
	// fire as if the load function was called with the
	// current source.
	//
	// Untracked.
	reload() {
		untrack(() => this._element?.load())
	}

	// play begins playing the media if it has loaded and is
	// not currently playing, i.e. media enters the running
	// state.
	//
	// Untracked.
	play() {
		untrack(() => this._element?.play())
	}

	// pause stops playing the media if is currently playing,
	// i.e. media enters the paused state.
	//
	// Untracked.
	pause() {
		untrack(() => this._element?.pause())
	}

	// playPause begins playing the media if it is paused or
	// pauses if it is running.
	//
	// Untracked.
	playPause() {
		untrack(() => {
			this._paused ? this.play() : this.pause()
		})
	}

	// restart sets the playback time to the start of the
	// video. It's a seek operation to time 0. Unlike the
	// reload function, restart does not reload the media so
	// load related events are not fired, but seek operations
	// will.
	//
	// Untracked.
	restart() {
		untrack(() => {
			if (this._element) {
				this._element.currentTime = 0
			}
		})
	}

	// on registers an event listener that is added to the
	// underlying HTMLMediaElement when it is set and removed
	// when it is unest. The listener is invoked when an
	// event with the given eventType occurs. Note that most
	// events fired on HTMLMediaElement are non-bubbling,
	// non-cancelable, and non-composed. Custom events
	// dispatched will always have those options.
	//
	// All standard HTMLMediaElement event types are
	// supported along with the custom ones: 'elementset'
	// and 'elementunset'.
	//
	// Returns true if the unique combination of eventType,
	// listener, and capturing phase were added, false if the
	// combination already existed.
	//
	// Untracked.
	on(eventType, listener, options = undefined) {
		return untrack(() => {
			const entry = createListenerEntry(
				eventType, //
				listener,
				options
			)

			const alreadyRegistered = containsListenerEntry(
				this._userListeners, //
				entry
			)

			if (alreadyRegistered) {
				return false
			}

			this._userListeners.push(entry)

			if (this._element) {
				addListenerToElement(this._element, entry)
			}

			return true
		})
	}

	// off unregisters an eventType, listener, and options
	// combination registered through the on function. True
	// is returned if the combination of eventType, listener,
	// and capturing phase (determined by the options
	// argument) was not registered.
	//
	// Untracked.
	off(eventType, listener, options = undefined) {
		return untrack(() => {
			let entry = createListenerEntry(
				eventType, //
				listener,
				options
			)

			const index = indexOfListenerEntry(
				this._userListener, //
				entry
			)

			if (index < 0) {
				return false
			}

			// Use the original entry for sanity's sake.
			entry = this._userListeners[index]
			this._userListeners.splice(index, 1)

			if (this._element) {
				removeListenerFromElement(this._element, entry)
			}

			return true
		})
	}

	// isOn returns true if the combination of eventType,
	// listener, and capturing phase (determined by the
	// options argument) is currently registered.
	//
	// Untracked.
	isOn(eventType, listener, options = undefined) {
		untrack(() => {
			const entry = createListenerEntry(
				eventType, //
				listener,
				options
			)

			return containsListenerEntry(
				this._userListeners, //
				entry
			)
		})
	}

	// _syncStatesInit performs an initial state syncing
	// with a newly set HTMLMediaElement. This includes
	// syncing metadata.
	_syncStatesInit() {
		const elem = this._element

		if (elem.readyState < HTMLMediaElement.HAVE_METADATA) {
			return
		}

		this._updateMetadata()
		this._syncStates()
	}

	// _updateMetadata updates general metadata state
	// including any state that may be available straight
	// after metadata has loaded.
	_updateMetadata() {
		this._duration = this._element.duration
		this._seekable = this._element.seekable.length > 0
	}

	// _syncStates performs a state syncing with the current
	// set HTMLMediaElement. It assumes metadata is already
	// synced.
	_syncStates() {
		const wasPaused = this._paused
		const wasFlowing = this._flowing
		const wasBuffering = this._buffering

		this._updateSeekStates()
		this._updateRunningPausedStates()
		this._updatePlayableState()
		this._updateFlowingBufferingStates()

		this._dispatchRunningPausingEvent(
			wasPaused //
		)

		this._dispatchFlowingBufferingEvent(
			wasFlowing, //
			wasBuffering
		)
	}

	// _updateSeekStates syncs the seekable and seeking
	// states to match the current HTMLMediaElement's state.
	_updateSeekStates() {
		this._seekable = this._element.seekable.length > 0
		this._seeking = this._element.seeking
	}

	// _updateRunningPausedStates sets the state of the
	// playable, running, flowing, buffering, and paused
	// states together. When unpausing, the 'running' event
	// will be fired before the 'flowing' or 'buffering'
	// events.
	_updateRunningPausedStates() {
		this._paused = this._element.paused
		this._running = !this._paused
	}

	// _updatePlayableState sets the playable set dependent
	// on not currently seeking and some future data being
	// loaded to play.
	_updatePlayableState() {
		this._playable =
			!this._element.seeking && //
			this._hasFutureData()
	}

	// _updateFlowingBufferingStates updates the playable,
	// flowing, and buffering states based upon the
	// underlying HTMLMediaElement's paused, ended, seeking,
	// and readyState values.
	_updateFlowingBufferingStates() {
		if (this._couldBeFlowingOrBuffering()) {
			this._flowing = this._playable
			this._buffering = !this._playable
		} else {
			this._flowing = false
			this._buffering = false
		}
	}

	// _couldBeFlowingOrBuffering returns true if the media
	// can be placed into the flowing or buffering states.
	_couldBeFlowingOrBuffering() {
		const media = this._element
		return !media.paused && !media.seeking && !media.ended
	}

	// _dispatchRunningPausingEvent will dispatch the
	// 'running' or 'pausing' events if either satisfy
	// their dispatch conditions.
	_dispatchRunningPausingEvent(wasPaused) {
		if (wasPaused && !this._paused) {
			dispatchEvent(this._element, 'running')
		} else if (!wasPaused && this._paused) {
			dispatchEvent(this._element, 'pausing')
		}
	}

	// _dispatchFlowingBufferingEvent will dispatch the
	// 'flowing' or 'buffering' events if either satisfy
	// their dispatch conditions.
	_dispatchFlowingBufferingEvent(wasFlowing, wasBuffering) {
		if (!wasFlowing && this._flowing) {
			dispatchEvent(this._element, 'flowing')
		} else if (!wasBuffering && this._buffering) {
			dispatchEvent(this._element, 'buffering')
		}
	}

	// _hasFutureData returns true if enough data has loaded
	// to allow playing without buffering for some time into
	// the future.
	_hasFutureData() {
		const HAS_DATA = HTMLMediaElement.HAVE_FUTURE_DATA
		return this._element.readyState >= HAS_DATA
	}

	// _updateCurrentTime sets the currentTime and
	// currentRemaining states based on the element's
	// currentTime. Depending on the current seeking state it
	// will update either the seekTime or playtime and
	// remainingTime states.
	_updateCurrentTime() {
		const time = this._element.currentTime

		this._currentTime = time
		this._currentRemaining = this._duration - time

		if (this.seeking) {
			this._seektime = time
		} else {
			this._playtime = time
			this._remainingTime = this._duration - time
		}
	}

	// _resetStates resets all properties to their preload
	// defaults. This should only be called when media load
	// starts, after unloading the media, or unsetting the
	// HTMLMediaElement.
	_resetStates() {
		this._duration = 0
		this._currentTime = 0
		this._currentRemaining = 0
		this._playtime = 0
		this._remainingTime = 0

		this._paused = true
		this._flowing = false
		this._buffering = false
		this._running = false
		this._seeking = false
		this._playable = false
		this._seekable = false

		this._loaded = false
	}
}

// generateStateListeners creates the set of
// HTMLMediaElement listeners that manage the state of a
// SvelteMediaElement. It returns an array of listener
// entries.
function generateStateListeners(sme) {
	function abort() {
		sme._syncStates()
	}

	function canplay() {
		sme._syncStates()
	}

	function canplaythrough() {
		sme._syncStates()
	}

	function durationchange() {
		sme._updateMetadata()
		sme._syncStates()
	}

	function emptied() {
		sme._resetStates()
	}

	function ended() {
		sme._syncStates()
	}

	function error() {
		sme._syncStates()
	}

	function loadeddata() {
		sme._updateMetadata()
		sme._syncStates()
	}

	function loadedmetadata() {
		sme._loaded = true

		sme._updateMetadata()
		sme._syncStates()
	}

	function loadstart() {
		sme._resetStates()
		sme._syncStates()
	}

	function pause() {
		sme._syncStates()
	}

	function play() {
		sme._syncStates()
	}

	function playing() {
		sme._syncStates()
	}

	function progress() {
		sme._syncStates()
	}

	function ratechange() {
		// Do nothing.
	}

	function seeked() {
		sme._syncStates()
	}

	function seeking() {
		sme._syncStates()
	}

	function stalled() {
		sme._syncStates()
	}

	function suspend() {
		sme._syncStates()
	}

	function timeupdate() {
		sme._updateCurrentTime()
	}

	function volumechange() {
		// Do nothing.
	}

	function waiting() {
		sme._syncStates()
	}

	function waitingforkey() {
		// Do nothing.
	}

	const listeners = {
		abort, //
		canplay,
		canplaythrough,
		durationchange,
		emptied,
		ended,
		error,
		loadeddata,
		loadedmetadata,
		loadstart,
		pause,
		play,
		playing,
		progress,
		ratechange,
		seeked,
		seeking,
		stalled,
		suspend,
		timeupdate,
		volumechange,
		waiting,
		waitingforkey,
	}

	return Object.keys(listeners).map((eventType) => {
		return createListenerEntry(
			eventType, //
			listeners[eventType],
			true
		)
	})
}

// isObject returns true if v is a non-null plain object.
function isObject(v) {
	return (
		typeof v === 'object' && //
		!Array.isArray(v) &&
		v !== null
	)
}

// isCapturing returns true if the listener options denote
// the listener will be called during the capturing phase
// of an event.
function isCapturing(options) {
	if (options === true) {
		return true
	}

	if (isObject(options)) {
		return options.capture === true
	}

	return false
}

// createListenerEntry creates an object holding the
// eventType, listener, options, and capture phase. The
// object is easier to use and pass around than the
// event listener parameters themselves.
function createListenerEntry(eventType, listener, options) {
	return {
		eventType, //
		listener,
		options,
		capture: isCapturing(options),
	}
}

// containsListenerEntry returns true if the listener
// entry is present within listeners.
//
// Untracked.
function containsListenerEntry(listeners, entry) {
	const index = indexOfListenerEntry(
		listeners, //
		entry
	)

	return index > -1
}

// indexOfListenerEntry returns the index of the listener
// entry containing a unique combination of eventType,
// listener, and capturing phase (determined by the options
// argument) within listeners. -1 is returned if the
// combination does not exist.
//
// Untracked.
function indexOfListenerEntry(listeners, entry) {
	for (let i = 0; i < listeners.length; i++) {
		const currEntry = listeners[i]

		const match =
			currEntry.eventType === entry.eventType && //
			currEntry.listener === entry.listener &&
			currEntry.capture === entry.capture

		if (match) {
			return i
		}
	}

	return -1
}

// addListenersToElement adds all listeners to the
// element.
//
// Tracked.
function addListenersToElement(element, listeners) {
	for (const entry of listeners) {
		addListenerToElement(element, entry)
	}
}

// removeListenersFromElement removes all listeners from
// the element.
//
// Tracked.
function removeListenersFromElement(element, listeners) {
	for (const entry of listeners) {
		removeListenerFromElement(element, entry)
	}
}

// addListenerToElement adds a listener to the element.
//
// Tracked.
function addListenerToElement(element, entry) {
	element.addEventListener(
		entry.eventType, //
		entry.listener,
		entry.options
	)
}

// removeListenerFromElement removes a listener from the
// element.
//
// Tracked.
function removeListenerFromElement(element, entry) {
	element.removeEventListener(
		entry.eventType, //
		entry.listener,
		entry.options
	)
}

// dispatchEvent dispatches an event on the element
// with bubbling, cancelable, composed options all false.
//
// Tracked.
function dispatchEvent(element, eventType) {
	const event = new Event(eventType, {
		bubbles: false, //
		cancelable: false,
		composed: false,
	})

	element.dispatchEvent(event)
}
