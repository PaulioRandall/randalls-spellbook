import ArrayUtil from '$lib/ArrayUtil.js'

// TODO: Create functions for each public field so $effect
//       runes won't trigger if devs don't want them to.

// MediaController is an adapter for the built-in
// HTMLMediaElement class that is usable when the
// HTMLMediaElement is not set. This class serves multiple
// purposes:
// 1. Decouples Svelte components from the underlying
//    HTMLMediaElement and removes the need to pass the
//    HTMLMediaElement around when it becomes available.
// 2. Minimise the need to handle null HTMLMediaElement
//    within each Svelte class and other code.
// 3. Provides a simplified bespoke interface for playing
//    and managing media within this app.
// 4. Manages event adding and removal as the underlying
//    HTMLMediaElement changes, so Svelte components don't
//    need to monitor its existance.
export default class MediaController {
	// loaded is a readonly state that is true when the
	// media has loaded.
	loaded = $state(false)

	// playable is a readonly state that is true when the
	// media is in a playable state.
	playable = $state(false)

	// playing is a readonly state that is true when the
	// media is playing.
	playing = $state(false)

	// paused is a readonly state that is true when the
	// playing state is false, that is, it's always opposite
	// of the playing state.
	paused = $derived(!this.playing)

	// seekable is a readonly state that is true when the
	// underlying HTMLMediaElement is seekable.
	seekable = $state(false)

	// seeking is a readonly state that is true when the
	// underlying HTMLMediaElement is in seeking mode.
	seeking = $state(false)

	// duration is a readonly state for the current media
	// duration.
	duration = $state(0)

	// currentTime is a readonly state for the current
	// playback time. This differs from playtime which is not
	// updated during seeking.
	currentTime = $state(0)

	// playtime is a readonly state for the current playback
	// time, however it is only updated during playback.
	// This differs from currentTime which is also updated
	// during seeking.
	playtime = $state(0)

	// remaining is a readonly state for the amount of time
	// remaining in playback, i.e. duration - playtime. It
	// updates inline with the playtime field.
	remaining = $derived(this.duration - this.playtime)

	// _element is an updatable state for the underlying
	// HTMLMediaElement.
	_element = $state(null)

	// element is a readonly state for the underlying
	// HTMLMediaElement.
	element = $derived(this._element)

	_muxListeners = generateMuxListeners(this)
	_userListeners = {/* eventType: [listener] */}

	// isLoaded returns the raw unreactive value of the
	// loaded field.
	isLoaded() {
		const v = $state.raw(this.loaded)
		// svelte-ignore state_referenced_locally
		return v
	}

	// isPlayable returns the raw unreactive value of the
	// playable field.
	isPlayable() {
		const v = $state.raw(this.playable)
		// svelte-ignore state_referenced_locally
		return v
	}

	// isPlaying returns the raw unreactive value of the
	// playing field.
	isPlaying() {
		const v = $state.raw(this.playing)
		// svelte-ignore state_referenced_locally
		return v
	}

	// isPaused returns the raw unreactive value of the
	// paused field.
	isPaused() {
		const v = $state.raw(this.paused)
		// svelte-ignore state_referenced_locally
		return v
	}

	// isSeekable returns the raw unreactive value of the
	// seekable field.
	isSeekable() {
		const v = $state.raw(this.seekable)
		// svelte-ignore state_referenced_locally
		return v
	}

	// isSeeking returns the raw unreactive value of the
	// seeking field.
	isSeeking() {
		const v = $state.raw(this.seeking)
		// svelte-ignore state_referenced_locally
		return v
	}

	// getDuration returns the raw unreactive value of the
	// duration field.
	getDuration() {
		const v = $state.raw(this.duration)
		// svelte-ignore state_referenced_locally
		return v
	}

	// getCurrentTime returns the raw unreactive value of the
	// currentTime field.
	getCurrentTime() {
		const v = $state.raw(this.currentTime)
		// svelte-ignore state_referenced_locally
		return v
	}

	// getPlaytime returns the raw unreactive value of the
	// playtime field.
	getPlaytime() {
		const v = $state.raw(this.playtime)
		// svelte-ignore state_referenced_locally
		return v
	}

	// getRemaining returns the raw unreactive value of the
	// remaining field.
	getRemaining() {
		const v = $state.raw(this.remaining)
		// svelte-ignore state_referenced_locally
		return v
	}

	// hasElement returns true if a HTMLMediaElement is set,
	// else returns false.
	hasElement() {
		const v = $state.raw(this._element)
		// svelte-ignore state_referenced_locally
		return !!v
	}

	// getElement returns the current HTMLMediaElement or
	// null if it is not set.
	getElement() {
		const v = $state.raw(this._element)
		// svelte-ignore state_referenced_locally
		return v
	}

	// setElement sets the underlying HTMLMediaElement. If
	// one is already set then it will be unset first,
	// invoking the relevant events to fire. If the
	// mediaElement argument is falsey it is interpreted as
	// an unset only with the element field bein set to null.
	setElement(mediaElement) {
		this._unsetElement()

		if (!mediaElement) {
			return
		}

		const validType = mediaElement instanceof HTMLMediaElement
		if (!validType) {
			throw new Error('Not a HTMLMediaElement')
		}

		this._element = mediaElement
		this._addMuxListeners()
		this._callUserListeners('elementset')
	}

	// _unsetElement sets the underlying HTMLMediaElement
	// to null. It will reset all state then fire the
	// 'elementunset' event.
	_unsetElement(mediaElement) {
		if (this._element === null) {
			return
		}

		this._removeMuxListeners()

		this.hasElement = false
		this._element = null
		this._resetStates()

		this._callUserListeners('elementunset')
	}

	// reload reloads the media. Load related functions will
	// fire as if the load function was called with the
	// current source.
	reload() {
		this._resetStates()
		this._element.load()
	}

	/*
	// load first unloads the current media and removes any
	// existing sources, then adds the specified source and
	// begins loading it.
	loadSource(src, type) {
		if (this._element) {
			return
		}

		// 1. Unload current media & remove source
		// 2. Add new source
	}
	*/

	// play begins playing the media if it has loaded and is
	// not currently playing.
	play() {
		this._element?.play()
	}

	// pause stops playing the media if is currently playing.
	pause() {
		this._element?.pause()
	}

	// playPause begins playing the media if it is not
	// playing and stops playing if it is.
	playPause() {
		this.playing ? this.pause() : this.play()
	}

	// restart sets the playback time to 0. Unlike the reload
	// function, restart does not reload the media so load
	// related events are not fired.
	restart() {
		if (this.hasElement()) {
			this._element.currentTime = 0
		}
	}

	_resetStates() {
		this.loaded = false
		this.playable = false
		this.playing = false
		this.seekable = false
		this.seeking = false

		this.duration = 0
		this.currentTime = 0
		this.playtime = 0
	}

	// _addMuxListeners adds all middleware listeners to the
	// current MediaElement.
	_addMuxListeners() {
		for (const eventType in this._muxListeners) {
			const listener = this._muxListeners[eventType]
			this._element.addEventListener(eventType, listener)
		}
	}

	// _removeMuxListeners removes all middleware listeners
	// from the current MediaElement.
	_removeMuxListeners() {
		for (const eventType in this._muxListeners) {
			const listener = this._muxListeners[eventType]
			this._element.removeEventListener(eventType, listener)
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
		mc.seekable = mc._element.seekable
		mc.duration = mc._element.duration
		mc.currentTime = 0
		updateLoadStates()
	}

	function updateLoadStates() {
		const METADATA_READY = HTMLMediaElement.HAVE_METADATA
		mc.loaded = mc._element.readyState >= METADATA_READY

		const DATA_READY = HTMLMediaElement.HAVE_FUTURE_DATA
		mc.playable = mc._element.readyState >= DATA_READY
	}

	function loadedmetadata(event) {
		updateLoadStates()
		mc._callUserListeners('loadedmetadata', event)
	}

	function loadstart(event) {
		updateLoadStates()
		mc._callUserListeners('loadstart', event)
	}

	function loadeddata(event) {
		updateLoadStates()
		mc._callUserListeners('loadeddata', event)
	}

	function progress(event) {
		updateLoadStates()
		mc._callUserListeners('progress', event)
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
		mc.duration = mc._element.duration
		mc._callUserListeners('durationchange', event)
	}

	function timeupdate(event) {
		mc.currentTime = mc._element.currentTime
		if (!mc.seeking) {
			mc._callUserListeners('playback', event)
		}
	}

	function play(event) {
		mc.playing = true
		mc._callUserListeners('play', event)
	}

	function pause(event) {
		mc.playing = false
		mc._callUserListeners('pause', event)
	}

	function seeking(event) {
		mc.seeking = true
		mc._callUserListeners('seeking', event)
	}

	function seeked(event) {
		mc.seeking = false
		mc._callUserListeners('seeked', event)
	}

	return {
		loadedmetadata,
		loadstart,
		loadeddata,
		progress,
		canplay,
		canplaythrough,
		durationchange,
		timeupdate,
		play,
		pause,
		seeking,
		seeked,
	}
}
