import { untrack } from 'svelte'
import ElementalSvox from './ElementalSvox.svelte.js'
import Eventor from './Eventor.js'

// SvelteMediaElement is a Reactive Element Adapter Box
// for the standard HTMLMediaElement class.
//
// READOX classes follow a set of rules aimed at providing
// clarity for programmers and minimising reactivity issues
// that can be hard to debug:
// 1. All public properties are readonly, i.e. they can
//    only be set internally or through public functions.
// 2. All public properties are reactive and tracked,
//    i.e. they are created through runes such as $state
//    and $derived, and trigger $effect and $derived runes.
// 3. All public functions are non-reactive and untracked,
//    i.e. they do not trigger $effect or $derived runes.
//    Use common and iconic verbs for functions that get or
//    set properties: e.g. get, is, has, set, put, update.
// 4. Instances of the class must never be null. The
//    purpose of being a box is to minimise the need for
//    existence checking.
//
// The primary purpose was to decouple Svelte components
// from the underlying HTMLMediaElement. Because the
// HTMLMediaElement is set and unset separate to
// instance construction, the instance passed to components
// can always be non-null.
//
// The class uses an Eventor to manage state and user
// listeners. When a HTMLMediaElement is set and unset the
// listeners will automatically be added and removed
// respectively. Components using this class can use the
// onElement function to register callbacks that are called
// when the underlying HTMLMediaElement is set and unset
// (these are callbacks, not events, so accept an instance
// of this class instead of an event object).
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
//    + loaded, paused
//    - seeking, playable, running, flowing, buffering
//
// 2. Once data has loaded, the user has pressed play and
//    the media playing without needing to wait for data:
//    + loaded, playable, running, flowing
//    - paused, seeking, buffering
//
// 3. When the user is seeking far into the future during
//    the middle of playback:
//    + loaded, playable, running, seeking
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
export default class SvelteMediaElement extends ElementalSvox {
	// _stateEventor are capturing event listeners that
	// mutate the adpaters state. They are added first when
	// an element is set so they fire first.
	_stateEventor = generateStateEventor(this)

	// isValidElement override to constrain elements to
	// HTMLMediaElements only.
	//
	// Untracked.
	isValidElement(element) {
		return element instanceof HTMLMediaElement
	}

	// afterSet override to add state listeners.
	//
	// Untracked.
	afterSet() {
		this._stateEventor.addTo(this.getElement())
	}

	// beforeUnset override to remove state listeners.
	//
	// Untracked.
	beforeUnset() {
		this._stateEventor.removeFrom(this.getElement())
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
		untrack(() => this.element?.load())
	}

	// play begins playing the media if it has loaded and is
	// not currently playing, i.e. media enters the running
	// state.
	//
	// Untracked.
	play() {
		untrack(() => this.element?.play())
	}

	// pause stops playing the media if is currently playing,
	// i.e. media enters the paused state.
	//
	// Untracked.
	pause() {
		untrack(() => this.element?.pause())
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

	// seekTo seeks to the specified time. If the time is
	// greater than the duration then the duration is used
	// instead, i.e. seeks to the end of the media.
	//
	// Untracked.
	seekTo(time) {
		untrack(() => {
			if (!this.element) {
				return
			}

			if (time > this._duration) {
				time = this._duration
			}

			this.element.currentTime = time
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
		this.seekTo(0)
	}

	// syncStates performs a state syncing with the current
	// set HTMLMediaElement.
	syncStates() {
		const media = this.getElement()

		if (!media) {
			this._resetStates()
			return
		}

		if (!this._loaded) {
			if (media.readyState < HTMLMediaElement.HAVE_METADATA) {
				return
			}

			this._loaded = true
			this._syncMetadataStates()
		}

		this._syncMediaStates()
	}

	// _syncMetadataStates updates general metadata state
	// including any state that may be available straight
	// after metadata has loaded.
	_syncMetadataStates() {
		this._duration = this.getElement().duration
		this._seekable = this.getElement().seekable.length > 0
	}

	// _syncMediaStates updates general states such as
	// playable, paused, flowing, etc.
	_syncMediaStates() {
		const wasPaused = this._paused
		const wasFlowing = this._flowing
		const wasBuffering = this._buffering

		this._updateStates()

		this._dispatchStateChanges(
			wasPaused, //
			wasFlowing,
			wasBuffering
		)
	}

	// _updateStates updates states to match the current
	// HTMLMediaElement's state.
	_updateStates() {
		this._seekable = this.getElement().seekable.length > 0

		this._seeking = this._element.seeking
		this._paused = this._element.paused
		this._running = !this._paused

		this._playable = !this._seeking && this._hasFutureData()

		if (this._couldBeFlowingOrBuffering()) {
			this._flowing = this._playable
			this._buffering = !this._playable
		} else {
			this._flowing = false
			this._buffering = false
		}
	}

	// _hasFutureData returns true if enough data has loaded
	// to allow playing without buffering for some time into
	// the future.
	_hasFutureData() {
		const HAS_DATA = HTMLMediaElement.HAVE_FUTURE_DATA
		return this.getElement().readyState >= HAS_DATA
	}

	// _couldBeFlowingOrBuffering returns true if the media
	// can be placed into the flowing or buffering states.
	_couldBeFlowingOrBuffering() {
		const media = this.getElement()
		return !media.paused && !media.seeking && !media.ended
	}

	// _dispatchStateChanges dispatches custom events after
	// state syncing.
	_dispatchStateChanges(wasPaused, wasFlowing, wasBuffering) {
		if (wasPaused && !this._paused) {
			this.dispatchEvent('running')
		} else if (!wasPaused && this._paused) {
			this.dispatchEvent('pausing')
		}

		if (!wasFlowing && this._flowing) {
			this.dispatchEvent('flowing')
		} else if (!wasBuffering && this._buffering) {
			this.dispatchEvent('buffering')
		}
	}

	// _updateCurrentTime sets the currentTime and
	// currentRemaining states based on the element's
	// currentTime. Depending on the current seeking state it
	// will update either the seekTime or playtime and
	// remainingTime states.
	_updateCurrentTime() {
		const time = this.getElement().currentTime

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
function generateStateEventor(sme) {
	function abort() {
		sme._syncMediaStates()
	}

	function canplay() {
		sme._syncMediaStates()
	}

	function canplaythrough() {
		sme._syncMediaStates()
	}

	function durationchange() {
		sme._syncMetadataStates()
		sme._syncMediaStates()
	}

	function emptied() {
		sme._resetStates()
	}

	function ended() {
		sme._syncMediaStates()
	}

	function error() {
		sme._syncMediaStates()
	}

	function loadeddata() {
		sme._syncMetadataStates()
		sme._syncMediaStates()
	}

	function loadedmetadata() {
		sme.syncStates()
	}

	function loadstart() {
		sme._resetStates()
		sme.syncStates()
	}

	function pause() {
		sme._syncMediaStates()
	}

	function play() {
		sme._syncMediaStates()
	}

	function playing() {
		sme._syncMediaStates()
	}

	function progress() {
		sme._syncMediaStates()
	}

	function ratechange() {
		// Do nothing.
	}

	function seeked() {
		sme._syncMediaStates()
	}

	function seeking() {
		sme._syncMediaStates()
	}

	function stalled() {
		sme._syncMediaStates()
	}

	function suspend() {
		sme._syncMediaStates()
	}

	function timeupdate() {
		sme._updateCurrentTime()
	}

	function volumechange() {
		// Do nothing.
	}

	function waiting() {
		sme._syncMediaStates()
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

	// Transform to enable capture.
	for (const eventType in listeners) {
		listeners[eventType] = {
			listener: listeners[eventType], //
			options: { capture: true },
		}
	}

	const eventor = new Eventor()
	eventor.on(listeners)
	console.log(eventor._entries)
	return eventor
}
