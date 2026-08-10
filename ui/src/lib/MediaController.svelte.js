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
	_element = null
	hasElement = $state(false)

	loaded = $state(false)
	playable = $state(false)
	playing = $state(false)

	seekable = $state(false)
	seeking = $state(false)

	duration = $state(0)
	currentTime = $state(0)
	playtime = $state(0)
	remaining = $derived(this.duration - this.playtime)

	_muxListeners = generateMuxListeners(this)
	_userListeners = {/* eventType: [listener] */}

	// setElement sets the underlying HTMLMediaElement. If
	// one is already set then it will be unset first,
	// invoking the relevant events to fire.
	setElement(mediaElement) {
		this.unsetElement()

		const validType = mediaElement instanceof HTMLMediaElement
		if (!mediaElement || !validType) {
			throw new Error('Not a HTMLMediaElement')
		}

		this._element = mediaElement
		this.hasElement = true

		this._addMuxListeners()
		this._callUserListeners('elementset')
	}

	// unsetElement removes the underlying MediaElement if it
	// is set.
	unsetElement(mediaElement) {
		if (!this._element) {
			return
		}

		this._removeMuxListeners()

		this._element = null
		this.hasElement = false
		this._resetStates()

		this._callUserListeners('elementunset')
	}

	// get returns the current MediaElement or null if it is
	// not set.
	get() {
		return this._element
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
		if (this._element) {
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
		this.remaining = 0
	}

	_updateMetadata() {
		this.seekable = this._element.seekable
		this.duration = this._element.duration
		this.currentTime = 0
		this._updateLoadStates()
	}

	_updateLoadStates() {
		const METADATA_READY = HTMLMediaElement.HAVE_METADATA
		this.loaded = this._element.readyState >= METADATA_READY

		const DATA_READY = HTMLMediaElement.HAVE_FUTURE_DATA
		this.playable = this._element.readyState >= DATA_READY
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

function removeFromArray(array, item) {
	const i = array.indexOf(item)

	if (i > -1) {
		array.splice(i, 1)
	}

	return item
}

function generateMuxListeners(mc) {
	function loadedmetadata(event) {
		mc._updateMetadata()
		mc._callUserListeners('loadedmetadata', event)
	}

	function loadstart(event) {
		mc._updateLoadStates()
		mc._callUserListeners('loadstart', event)
	}

	function loadeddata(event) {
		mc._updateLoadStates()
		mc._callUserListeners('loadeddata', event)
	}

	function progress(event) {
		mc._updateLoadStates()
		mc._callUserListeners('progress', event)
	}

	function canplay(event) {
		mc._updateLoadStates()
		mc._callUserListeners('canplay', event)
	}

	function canplaythrough(event) {
		mc._updateLoadStates()
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
