import { untrack } from 'svelte'
import ArrayUtil from '$lib/ArrayUtil.js'

// HTML_MEDIA_ELEMENT_EVENT_TYPES is an array of all
// HTMLMediaELement event types that users can register
// listeners for.
const HTML_MEDIA_ELEMENT_EVENT_TYPES = [
	'abort',
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

// SPECIAL_EVENT_TYPES is an array of all custom event
// types fired by SvelteMediaElement, not by the
// underlying HTMLMediaELement.
const SPECIAL_EVENT_TYPES = [
	'elementset',
	'elementunset',
	'loaded',
	'running',
	'paused',
]

const EVENT_TYPES = [
	...HTML_MEDIA_ELEMENT_EVENT_TYPES, //
	...SPECIAL_EVENT_TYPES,
]

// SvelteMediaElement is a Svelte adapter for the built-in
// HTMLMediaElement class that is allows some functionality
// even when the underlying HTMLMediaElement is not set.
//
// The original purpose was to decouple Svelte components
// from the underlying HTMLMediaElement. Because the
// HTMLMediaElement is set and unset separate to
// instance construction, the SvelteMediaElement instance
// passed to components can always be non-null. This
// moves most existence checking from the components into
// this class.
//
// The adapter stores listeners added by components so when
// a HTMLMediaElement is set and unset the listeners will
// automatically be added or remove from it. Most media
// control and display related components don't need to do
// anything when the underlying HTMLMediaElement changes.
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
// trigger execution of $effect and $derived runes.
// Properties and their getters should be tracked while
// public functions should be untracked. This makes it even
// easier for readers to quickly understand the reactivity
// of some code.
//
// TODO: replace all $derived with $state so users can
//       access the updated values within the same tick,
//       e.g. in events fired after root $state is updated.
export default class SvelteMediaElement {
	// _muxCallbacks are middleware event listeners. Instead
	// of passing the users event listeners directly to the
	// underlying HTMLMediaElement, they are called through
	// by these mux callbacks.
	_muxCallbacks = generateMuxCallbacks(this)

	// _userCallbacks are the user defined listeners. They
	// include HTMLMediaElement event types and special
	// custom types that are fired only by this class, e.g.
	// 'elementset' and 'elementunset' event types.
	_userCallbacks = {/* eventType: [callback] */}

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
			this._addMuxCallbacks()
			this._fireEvent('elementset')
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

			this._removeMuxCallbacks()
			this._element = null
			this._resetStates()

			this._fireEvent('elementunset')
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

	// playing is true when the media is playing. E.g.
	// when the user is actively watching a video, but will
	// change to false when playback stops to buffer the
	// video. Thus, this is not quite the opposite of paused.
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
	_running = $derived(!this._paused)
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
	_currentRemaining = $derived(
		this._duration - this._currentTime //
	)
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
	_remainingTime = $derived(this._duration - this._playtime)
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
		untrack(() => this._element.load())
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
	// underlying HTMLMediaElement when set and removed when
	// unest. The listener is called when the eventType
	// occurs. All standard HTMLMediaElement event types are
	// supported along with additional ones: 'elementset',
	// 'elementunset', 'loaded', 'running', and 'paused'.
	//
	// Returns true if the eventType and callback pair were
	// added, false if the pair already existed. Tracking is
	// prevented.
	//
	// Untracked.
	//
	// TODO: Support options starting with { once: bool }
	on(eventType, callback) {
		return untrack(() => {
			if (!EVENT_TYPES.includes(eventType)) {
				throw new Error(`Unknown event type: '${eventType}'`)
			}

			const callbackType = typeof callback
			if (!callback || callbackType !== 'function') {
				throw new Error(
					`Callback must be a non-null function, not '${callbackType}'`
				)
			}

			return this._registerUserCallback(
				eventType, //
				callback
			)
		})
	}

	// off unregisters an event handler registered through
	// the on function.
	//
	// Untracked.
	off(eventType, callback) {
		return untrack(() => {
			return this._unregisterUserCallback(
				eventType, //
				callback
			)
		})
	}

	// _resetStates resets all fields to their preload
	// defaults. This should only be called when media load
	// starts or after unloading the media.
	//
	// Tracked.
	_resetStates() {
		this._playing = false
		this._paused = true
		this._seeking = false

		this._duration = 0
		this._currentTime = 0
		this._playtime = 0

		this._playable = false
		this._seekable = false
		this._loaded = false
	}

	// _addMuxCallbacks adds all middleware listeners to the
	// current HTMLMediaElement.
	//
	// Tracked.
	_addMuxCallbacks() {
		for (const eventType in this._muxCallbacks) {
			const callback = this._muxCallbacks[eventType]
			this._element.addEventListener(eventType, callback)
		}
	}

	// _removeMuxCallbacks removes all middleware listeners
	// from the current HTMLMediaElement.
	//
	// Tracked.
	_removeMuxCallbacks() {
		for (const eventType in this._muxCallbacks) {
			const callback = this._muxCallbacks[eventType]
			this._element.removeEventListener(eventType, callback)
		}
	}

	// _registerUserCallback registers the callback to a
	// specific event type. True is returned if the callback
	// was registered and false returned if the unique pair
	// of eventType and callback already exists.
	//
	// Untracked.
	_registerUserCallback(eventType, callback) {
		let callbackSet = this._userCallbacks[eventType]

		if (!callbackSet) {
			callbackSet = []
			this._userCallbacks[eventType] = callbackSet
		}

		if (callbackSet.includes(callback)) {
			return false
		}

		callbackSet.push(callback)
		return true
	}

	// _unregisterUserCallback unregisters a callback
	// registered through _registerUserCallback. True is
	// returned if the callback was unregistered and false
	// returned if the unique pair of eventType and callback
	// didn't exist.
	//
	// Untracked.
	_unregisterUserCallback(eventType, callback) {
		const callbackSet = this._userCallbacks[eventType]

		if (!callbackSet) {
			return false
		}

		if (callbackSet.includes(callback)) {
			return false
		}

		ArrayUtil.remove(callbackSet, callback)

		if (callbackSet.length === 0) {
			delete this._userCallbacks[eventType]
		}

		return true
	}

	_fireEvent(eventType, event = null) {
		const callbackSet = this._userCallbacks[eventType]

		if (!callbackSet) {
			return
		}

		for (const callback of callbackSet) {
			callback(this, event)
		}
	}

	_fireCallback(callback, event = null) {
		callback(this, event)
	}

	_firePremisedCallback(callback, event, ...premises) {
		if (premises.find((p) => p === true)) {
			this._fireCallback(callback, event)
		}
	}
}

function generateMuxCallbacks(mc) {
	function updateMetadata() {
		mc._seekable = mc._element.seekable
		mc._duration = mc._element.duration
		mc._currentTime = mc._element.currentTime
		updateLoadStates()
	}

	function updateLoadStates(event) {
		// TODO: Clean up and maybe split up concerns of
		//       data from metadata.
		const wasLoaded = mc._loaded

		const METADATA_READY = HTMLMediaElement.HAVE_METADATA
		mc._loaded = mc._element.readyState >= METADATA_READY

		const DATA_READY = HTMLMediaElement.HAVE_FUTURE_DATA
		mc._playable = mc._element.readyState >= DATA_READY

		if (!wasLoaded && mc._loaded) {
			mc._fireEvent('loaded', event)
		}
	}

	function abort(event) {
		updateLoadStates(event)
		mc._fireEvent('abort', event)
	}

	function canplay(event) {
		updateLoadStates(event)
		mc._fireEvent('canplay', event)
	}

	function canplaythrough(event) {
		updateLoadStates(event)
		mc._fireEvent('canplaythrough', event)
	}

	function durationchange(event) {
		mc._duration = mc._element.duration
		mc._fireEvent('durationchange', event)
	}

	function emptied(event) {
		mc._resetStates()
		mc._fireEvent('emptied', event)
	}

	function ended(event) {
		mc._playing = false
		mc._paused = true

		mc._fireEvent('paused', event)
		mc._fireEvent('ended', event)
	}

	function error(event) {
		updateLoadStates()
		mc._fireEvent('error', event)
	}

	function loadeddata(event) {
		updateLoadStates()
		mc._fireEvent('loadeddata', event)
	}

	function loadedmetadata(event) {
		updateLoadStates()
		mc._fireEvent('loadedmetadata', event)
	}

	function loadstart(event) {
		mc._resetStates()
		updateLoadStates()
		mc._fireEvent('loadstart', event)
	}

	function pause(event) {
		mc._playing = false
		mc._paused = true
		mc._fireEvent('pause', event)
		mc._fireEvent('paused', event)
	}

	function play(event) {
		mc._paused = false
		mc._fireEvent('running', event)
		mc._fireEvent('play', event)
	}

	function playing(event) {
		mc._playing = true
		mc._fireEvent('playing', event)
	}

	function progress(event) {
		updateLoadStates()
		mc._fireEvent('progress', event)
	}

	function ratechange(event) {
		mc._fireEvent('ratechange', event)
	}

	function seeked(event) {
		mc._seeking = false
		mc._fireEvent('seeked', event)
	}

	function seeking(event) {
		mc._seeking = true
		mc._fireEvent('seeking', event)
	}

	function stalled(event) {
		mc._fireEvent('stalled', event)
	}

	function suspend(event) {
		mc._fireEvent('suspend', event)
	}

	function timeupdate(event) {
		mc._currentTime = mc._element.currentTime

		if (mc.seeking) {
			mc._seektime = mc._currentTime
		} else {
			mc._playtime = mc._currentTime
		}

		mc._fireEvent('timeupdate', event)
	}

	function volumechange(event) {
		mc._fireEvent('volumechange', event)
	}

	function waiting(event) {
		mc._playing = false
		mc._fireEvent('waiting', event)
	}

	function waitingforkey(event) {
		mc._fireEvent('waitingforkey', event)
	}

	return {
		abort,
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
}
