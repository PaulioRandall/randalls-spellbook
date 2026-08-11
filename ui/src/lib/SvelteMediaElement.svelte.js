import { untrack } from 'svelte'
import ArrayUtil from '$lib/ArrayUtil.js'

// SvelteMediaElement is a Svelte adapter for the built-in
// HTMLMediaElement class that is allows some functionality
// even when the underlying HTMLMediaElement is not set.
//
// The original purpose was to decouple Svelte components
// from the underlying HTMLMediaElement. Because the
// HTMLMediaElement is set and unset separate to
// instance construction, the SvelteMediaElement instance
// passed to components can always be non-null. This
// removes a lot of existence checking that I had before.
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
export default class SvelteMediaElement {
	_muxListeners = generateMuxListeners(this)
	_userListeners = {/* eventType: [listener] */}

	// element is a Svelte state for the underlying
	// HTMLMediaElement.
	element = $state(null)

	// loaded is a Svelte state that is true when the
	// media has loaded.
	loaded = $state(false)

	// playable is a Svelte state that is true when the
	// media is in a playable state.
	playable = $state(false)

	// playing is a Svelte state that is true when the
	// media is playing. E.g. playing is true when the user
	// actively watching a video but will change to false
	// when playback stops to buffer the video. It's
	// useful for toggling independent spinners and messaging
	// related to data loading.
	playing = $state(false)

	// paused is a Svelte state that is true prior to
	// loading, straight after loading, before autoplay
	// begins (if enabled), and when the user has paused the
	// video. It will be false when the media has stopped
	// playing because it's waiting for data, unlike the
	// playing field which is only true during user observed
	// playback.
	paused = $state(true)

	// running is a Svelte state that is true when paused
	// is false. A media can be both running and not
	// playing, usually indicating the media doesn't have
	// enough data to continue play right now but will
	// continue playing when data arrives.
	running = $derived(!this.paused)

	// seekable is a Svelte state that is true when the
	// underlying HTMLMediaElement is seekable.
	seekable = $state(false)

	// seeking is a Svelte state that is true when the
	// underlying HTMLMediaElement is in seeking mode.
	seeking = $state(false)

	// duration is a Svelte state for the current media
	// duration.
	duration = $state(0)

	// currentTime is a Svelte state for the current
	// playback time. This differs from playtime which is not
	// updated during seeking.
	currentTime = $state(0)

	// currentRemaining is a Svelte state for the
	// amount of time remaining in playback based on the
	// currentTime, i.e. duration - currentTime. It
	// updates inline with the currentTime field.
	currentRemaining = $derived(
		this.duration - this._currentTime //
	)

	// playtime is a Svelte state for the current playback
	// time, however it is only updated during playback.
	// This differs from currentTime which is also updated
	// during seeking. When currentTime is updated then
	// either playtime or seektime will update, but never
	// both.
	playtime = $state(0)

	// seektime is a Svelte state for the currentTime that
	// is only updated when seeking is true. When currentTime
	// is updated then either playtime or seektime will
	// update, but never both.
	seektime = $state(0)

	// remainingTime is a Svelte state for the amount of
	// time remaining in playback, i.e. duration - playtime.
	// It updates inline with the playtime field.
	remainingTime = $derived(this.duration - this.playtime)

	// hasElement returns true if a HTMLMediaElement is set,
	// else returns false. Tracking is prevented.
	hasElement() {
		return untrack(() => !!this.element)
	}

	// getElement returns the current HTMLMediaElement or
	// null if it is not set. Tracking is prevented.
	getElement() {
		return untrack(() => this.element)
	}

	// setElement sets the underlying HTMLMediaElement. If
	// one is already set then it will be unset first,
	// invoking the relevant events to fire. If the
	// mediaElement argument is falsey it is interpreted as
	// an unset only with the element field bein set to null.
	// Tracking is prevented.
	setElement(mediaElement) {
		untrack(() => {
			this._unsetElement()

			if (!mediaElement) {
				return
			}

			const validType = mediaElement instanceof HTMLMediaElement
			if (!validType) {
				throw new Error('Not a HTMLMediaElement')
			}

			this.element = mediaElement
			this._addMuxListeners()
			this._callUserListeners('elementset')
		})
	}

	// _unsetElement sets the underlying HTMLMediaElement
	// to null. It will reset all state then fire the
	// 'elementunset' event.
	_unsetElement(mediaElement) {
		if (this.element === null) {
			return
		}

		this._removeMuxListeners()

		this.hasElement = false
		this.element = null
		this._resetStates()

		this._callUserListeners('elementunset')
	}

	// isLoaded returns the loaded field without tracking.
	isLoaded() {
		return untrack(() => this.loaded)
	}

	// isPlayable returns the playable field without
	// tracking.
	isPlayable() {
		return untrack(() => this.playable)
	}

	// isPlaying returns the playing field without tracking.
	isPlaying() {
		return untrack(() => this.playing)
	}

	// isPaused returns the paused field without tracking.
	isPaused() {
		return untrack(() => this.paused)
	}

	// isRunning returns the running field without tracking.
	isRunning() {
		return untrack(() => this.running)
	}

	// isSeekable returns the seekable field without
	// tracking.
	isSeekable() {
		return untrack(() => this.seekable)
	}

	// isSeeking returns the seeking field without tracking.
	isSeeking() {
		return untrack(() => this.seeking)
	}

	// getDuration returns the duration field without
	// tracking.
	getDuration() {
		return untrack(() => this.duration)
	}

	// getCurrentTime returns the currentTime field without
	// tracking.
	getCurrentTime() {
		return untrack(() => this.currentTime)
	}

	// getCurrentRemaining returns the currentRemaining field
	// without tracking.
	getCurrentRemaining() {
		return untrack(() => this.currentRemaining)
	}

	// getPlaytime returns the playtime field without
	// tracking.
	getPlaytime() {
		return untrack(() => this.playtime)
	}

	// getSeektime returns the seektime field without
	// tracking.
	getSeektime() {
		return untrack(() => this.seektime)
	}

	// getRemainingTime returns the remaining field without
	// tracking.
	getRemainingTime() {
		return untrack(() => this.remainingTime)
	}

	// reload reloads the media. Load related functions will
	// fire as if the load function was called with the
	// current source. Tracking is prevented.
	reload() {
		untrack(() => this.element.load())
	}

	// play begins playing the media if it has loaded and is
	// not currently playing. Tracking is prevented.
	play() {
		untrack(() => this.element?.play())
	}

	// pause stops playing the media if is currently playing.
	// Tracking is prevented.
	pause() {
		untrack(() => this.element?.pause())
	}

	// playPause begins playing the media if it is not
	// playing and stops playing if it is.
	playPause() {
		untrack(() => {
			return this.paused ? this.play() : this.pause()
		})
	}

	// restart sets the playback time to 0. Unlike the reload
	// function, restart does not reload the media so load
	// related events are not fired.
	restart() {
		untrack(() => {
			if (this.element) {
				this.element.currentTime = 0
			}
		})
	}

	// _resetStates resets all fields to their preload
	// defaults. This should only be called when media load
	// starts or after unloading the media.
	_resetStates() {
		this.playing = false
		this.paused = true
		this.seeking = false

		this.duration = 0
		this.currentTime = 0
		this.playtime = 0

		this.playable = false
		this.seekable = false
		this.loaded = false
	}

	// _addMuxListeners adds all middleware listeners to the
	// current MediaElement.
	_addMuxListeners() {
		for (const eventType in this._muxListeners) {
			const listener = this._muxListeners[eventType]
			this.element.addEventListener(eventType, listener)
		}
	}

	// _removeMuxListeners removes all middleware listeners
	// from the current MediaElement.
	_removeMuxListeners() {
		for (const eventType in this._muxListeners) {
			const listener = this._muxListeners[eventType]
			this.element.removeEventListener(eventType, listener)
		}
	}

	/*

	_registerUserListener(eventType, callback) {
		let listeners = this._userListeners[eventType]

		if (!listeners) {
			listeners = []
			this._userListeners[eventType] = 	listeners
		}

		if (listeners.includes(callback)) {
			return
		}

		listeners.push(callback)
	}

	_unregisterUserListener(eventType, callback) {
		const listeners = this._userListeners[eventType]

		if (!listeners) {
			return
		}

		removeFromArray(listeners, callback)

		if (listeners.length === 0) {
			delete this._userListeners[eventType]
		}
	}
	*/

	_callUserListeners(eventType, event = null) {
		const listeners = this._userListeners[eventType]

		if (!listeners) {
			return
		}

		for (const func of listeners) {
			func(this, event)
		}
	}
}

function generateMuxListeners(mc) {
	function updateMetadata() {
		mc.seekable = mc.element.seekable
		mc.duration = mc.element.duration
		mc.currentTime = mc.currentTime
		updateLoadStates()
	}

	function updateLoadStates() {
		const METADATA_READY = HTMLMediaElement.HAVE_METADATA
		mc.loaded = mc.element.readyState >= METADATA_READY

		const DATA_READY = HTMLMediaElement.HAVE_FUTURE_DATA
		mc.playable = mc.element.readyState >= DATA_READY
	}

	function abort() {
		updateLoadStates()
		mc._callUserListeners('abort', event)
	}

	function canplay(event) {
		updateLoadStates()
		mc._callUserListeners('canplay', event)
	}

	function canplaythrough(event) {
		updateLoadStates()
		mc._callUserListeners('canplaythrough', event)
	}

	function durationchange(event) {
		mc.duration = mc.element.duration
		mc._callUserListeners('durationchange', event)
	}

	function emptied(event) {
		mc._resetStates()
		mc._callUserListeners('emptied', event)
	}

	function ended(event) {
		mc.playing = false
		mc.paused = mc.element.paused
		mc._callUserListeners('ended', event)
	}

	function error(event) {
		updateLoadStates()
		mc._callUserListeners('error', event)
	}

	function loadeddata(event) {
		updateLoadStates()
		mc._callUserListeners('loadeddata', event)
	}

	function loadedmetadata(event) {
		updateLoadStates()
		mc._callUserListeners('loadedmetadata', event)
	}

	function loadstart(event) {
		mc._resetStates()
		updateLoadStates()
		mc._callUserListeners('loadstart', event)
	}

	function pause(event) {
		mc.paused = true
		mc.playing = false
		mc._callUserListeners('pause', event)
	}

	function play(event) {
		mc.paused = false
		mc._callUserListeners('play', event)
	}

	function playing(event) {
		mc.playing = true
		mc._callUserListeners('playing', event)
	}

	function progress(event) {
		updateLoadStates()
		mc._callUserListeners('progress', event)
	}

	function ratechange(event) {
		mc._callUserListeners('ratechange', event)
	}

	function seeked(event) {
		mc.seeking = false
		mc._callUserListeners('seeked', event)
	}

	function seeking(event) {
		mc.seeking = true
		mc._callUserListeners('seeking', event)
	}

	function stalled(event) {
		mc._callUserListeners('stalled', event)
	}

	function suspend(event) {
		mc._callUserListeners('suspend', event)
	}

	function timeupdate(event) {
		mc.currentTime = mc.element.currentTime
		mc._callUserListeners('timeupdate', event)

		if (mc.seeking) {
			mc.seektime = mc.currentTime
			mc._callUserListeners('seektimeupdate', event)
		} else {
			mc.playtime = mc.currentTime
			mc._callUserListeners('playtimeupdate', event)
		}
	}

	function volumechange(event) {
		mc._callUserListeners('volumechange', event)
	}

	function waiting(event) {
		mc.playing = false
		mc._callUserListeners('waiting', event)
	}

	function waitingforkey(event) {
		mc._callUserListeners('waitingforkey', event)
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
