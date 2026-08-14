import { untrack } from 'svelte'
import ArrayUtil from '$lib/ArrayUtil.js'

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
	'continuing',
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
// The original purpose was to decouple Svelte components
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
// respectively. Most media control and display related
// components don't need to do anything when the underlying
// HTMLMediaElement changes. They just register listeners
// for the events they're interested in, including when the
// element is set and unset.
//
// State is updated during the capturing phase of events
// before user registered events so all user listeners
// have access to fresh state through this class's
// instances. Any capturing listeners added to the
// HTMLMediaElement before it is set here will fire before
// state is updated. Furthermore, if a listener is
// registered through this class is unregistered manually
// it will not prevent the listener from being added to the
// next HTMLMediaElement. Unless you have a particualar
// niche requirement, it is recommended that all listeners
// are registered and unregistered through this class to
// avoid nasty, hard to debug, errors.
//
// SvelteMediaElement exposes all fields as reactive state
// so they can be used directly in components. However,
// untracked values can sometimes be useful within the
// $effect rune so functions are provided that return
// untracked raw values. In fact, none of the public class
// functions will be tracked.
//
// This adapter has a number of additional fields and
// functions for common use cases. For instance, playtime
// and remaining fields provide the currentTime and
// remaining time (duration - playtime) that are only
// updated when the video is actually playing, i.e. not
// during seeking.
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
			this._syncStates()

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

	// playing is true when the media is not paused, not
	// ended, not seeking, and has enough data to play some
	// frames. When playing is true then the user is actively
	// consuming the media, i.e. the video or audio are
	// actively playing without buffering. When the media
	// freezes to load more data the playing state will
	// become false until data is available again.
	//
	// One can assume that the user's attention is on the
	// media when playing is true (at least for video).
	// However, the playing state is not enougth to determine
	// if a buffering icon should be shown because it will
	// also be false when the media is paused. Instead, the
	// buffering state can be used.
	//
	// Tracked.
	_playing = $state(false)
	get playing() {
		return this._playing
	}

	// isPlaying returns the untracked value of the playing
	// property.
	//
	// Untracked.
	isPlaying() {
		return untrack(() => this._playing)
	}

	// buffering is true when the media is not paused, not
	// ended, not seeking, but does not have enough data to
	// continue playing frames. The media will appear frozen
	// from the user's perspective.
	//
	// The freeze is likely to have broken the user's
	// attention on the media (at least for video) so
	// the buffering state can be used independently in a
	// Svelte IF block to control the display of a buffering
	// icon, e.g:
	// {#if svelteMediaElement.buffering}
	//   <img src="buffering.gif" />
	// {/if}
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
	// data, unlike the playing field..
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

	// running is true when paused is false. A media can be
	// both running and not playing, usually indicating the
	// media doesn't have enough data to continue playback
	// right now but will continue when data arrives.
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
	// not currently playing.
	//
	// Untracked.
	play() {
		untrack(() => this._element?.play())
	}

	// pause stops playing the media if is currently playing.
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

	// restart sets the playback time to 0. Unlike the reload
	// function, restart does not reload the media so load
	// related events are not fired. Because it's a seek,
	// seek related events will fire.
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
	// supported along with additional ones: 'elementset',
	// 'elementunset', 'running', 'playing', 'buffering',
	// and 'pausing'.
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

	// _syncStates sets all states based on the underlying
	// HTMLMediaElement's state. If metadata has yet to be
	// loaded then no state changes occur.
	_syncStates() {
		const elem = this._element

		if (elem.readyState < HTMLMediaElement.HAVE_METADATA) {
			return
		}

		this._updateMetadata()
		this._updatePlayStatesWithoutDispatch()
	}

	// _updateMetadata updates general metadata state such as
	// duration and seekable states.
	_updateMetadata() {
		// TODO: When does seekable need updating?
		this._seekable = this._element.seekable.length > 0
		this._duration = this._element.duration
	}

	// _updatePlayStates calls
	// _updatePlayStatesWithoutDispatch then
	// _dispatchContinueBufferingEvent since this is
	// a common pattern. The operations were separated to
	// allow the 'running' event to be fired before 'playing'
	// or 'buffering' when the media is unpaused.
	_updatePlayStates() {
		const wasPlaying = this._playing
		const wasBuffering = this._buffering
		this._updatePlayStatesWithoutDispatch()
		this._dispatchContinueBufferingEvent(wasPlaying, wasBuffering)
	}

	// _updatePlayStatesWithoutDispatch updates the playable,
	// playing, and buffering states based upon the
	// underlying HTMLMediaElement's paused, ended, seeking,
	// and readyState values.
	_updatePlayStatesWithoutDispatch() {
		this._playable = this._hasFutureData()

		if (this._couldBePlayingOrBuffering()) {
			this._playing = this._playable
			this._buffering = !this._playable
		} else {
			this._playing = false
			this._buffering = false
		}
	}

	// _couldBePlayingOrBuffering returns true if the media
	// can be placed into the playing or buffering states.
	_couldBePlayingOrBuffering() {
		const media = this._element
		return !media.paused && !media.ended && !media.seeking
	}

	// _dispatchContinueBufferingEvent will dispatch the
	// 'playing' or 'buffering' events if either satisfy
	// their dispatch conditions.
	_dispatchContinueBufferingEvent(wasPlaying, wasBuffering) {
		if (!wasPlaying && this._playing) {
			dispatchEvent(this._element, 'continuing')
			return
		}

		if (!wasBuffering && this._buffering) {
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

	// _updatePausedState sets the state of the pause and
	// running states together. It also updates the play
	// states. When unpausing, the 'running' event will be
	// fired before the 'playing' or 'buffering' states.
	_updatePausedState() {
		const wasPaused = this._paused
		const paused = this._element.paused
		const wasPlaying = this._playing
		const wasBuffering = this._buffering

		this._paused = paused
		this._running = !paused

		this._updatePlayStatesWithoutDispatch()

		if (wasPaused && !paused) {
			dispatchEvent(this._element, 'running')
		} else if (!wasPaused && paused) {
			dispatchEvent(this._element, 'pausing')
		}

		this._dispatchContinueBufferingEvent(wasPlaying, wasBuffering)
	}

	// _updateSeekingState sets the state of seeking and
	// aligns the playing states accordingly.
	_updateSeekingState() {
		this._seeking = this._element.seeking
		this._updatePlayStates()
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
	// starts or after unloading the media.
	_resetStates() {
		this._updatePausedState()

		this._seekable = false
		this._seeking = false

		this._duration = 0
		this._currentTime = 0
		this._currentRemaining = 0
		this._playtime = 0
		this._remainingTime = 0

		this._playable = false
		this._loaded = false
	}
}

// generateStateListeners creates the set of
// HTMLMediaElement listeners that manage the state of a
// SvelteMediaElement. It returns an array of listener
// entries.
function generateStateListeners(sme) {
	function abort() {
		sme._updatePlayStates()
	}

	function canplay() {
		sme._updatePlayStates()
	}

	function canplaythrough() {
		sme._updatePlayStates()
	}

	function durationchange() {
		sme._updateMetadata()
	}

	function emptied() {
		sme._resetStates()
	}

	function ended() {
		sme._updatePausedState()
	}

	function error() {
		sme._updatePlayStates()
	}

	function loadeddata() {
		sme._updateMetadata()
		sme._updatePlayStates()
	}

	function loadedmetadata() {
		sme._loaded = true

		sme._updateMetadata()
		sme._updatePlayStates()
	}

	function loadstart() {
		sme._resetStates()
		sme._updatePlayStates()
	}

	function pause() {
		sme._updatePausedState()
	}

	function play() {
		sme._updatePausedState()
	}

	function playing() {
		sme._updatePlayStates()
	}

	function progress() {
		sme._updatePlayStates()
	}

	function ratechange() {
		// Do nothing.
	}

	function seeked() {
		sme._updateSeekingState()
	}

	function seeking() {
		sme._updateSeekingState()
	}

	function stalled() {
		// Do nothing.
	}

	function suspend() {
		// Do nothing.
	}

	function timeupdate() {
		sme._updateCurrentTime()
	}

	function volumechange() {
		// Do nothing.
	}

	function waiting() {
		sme._updatePlayStates()
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
