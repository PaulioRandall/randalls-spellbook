import { untrack } from 'svelte'
import ElementSvox from './ElementSvox.svelte.js'

// MediaSvox is a ElementSvox specific for the
// standard HTMLMediaElement class.
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
export default class MediaSvox extends ElementSvox {
	// isValidElement override to constrain elements to
	// HTMLMediaElements only.
	//
	// Untracked.
	isValidElement(element) {
		return element instanceof HTMLMediaElement
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

	// syncStates performs a state syncing with the current
	// set HTMLMediaElement.
	//
	// Untracked.
	syncStates() {
		untrack(() => {
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
		})
	}

	// generateStateListeners overide to monitor media state.
	generateStateListeners() {
		const svox = this

		function abort() {
			svox._syncMediaStates()
		}

		function canplay() {
			svox._syncMediaStates()
		}

		function canplaythrough() {
			svox._syncMediaStates()
		}

		function durationchange() {
			svox._syncMetadataStates()
			svox._syncMediaStates()
		}

		function emptied() {
			svox._resetStates()
		}

		function ended() {
			svox._syncMediaStates()
		}

		function error() {
			svox._syncMediaStates()
		}

		function loadeddata() {
			svox._syncMetadataStates()
			svox._syncMediaStates()
		}

		function loadedmetadata() {
			svox.syncStates()
		}

		function loadstart() {
			svox._resetStates()
			svox.syncStates()
		}

		function pause() {
			svox._syncMediaStates()
		}

		function play() {
			svox._syncMediaStates()
		}

		function playing() {
			svox._syncMediaStates()
		}

		function progress() {
			svox._syncMediaStates()
		}

		function ratechange() {
			// Do nothing.
		}

		function seeked() {
			svox._syncMediaStates()
		}

		function seeking() {
			svox._syncMediaStates()
		}

		function stalled() {
			svox._syncMediaStates()
		}

		function suspend() {
			svox._syncMediaStates()
		}

		function timeupdate() {
			svox._updateCurrentTime()
		}

		function volumechange() {
			// Do nothing.
		}

		function waiting() {
			svox._syncMediaStates()
		}

		function waitingforkey() {
			// Do nothing.
		}

		return {
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
