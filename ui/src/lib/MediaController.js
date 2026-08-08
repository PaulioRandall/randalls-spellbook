export default class MediaController {
	_mediaElement = null
	_listeners = {
		//[eventType]: [...listeners]
	}

	// on adds an event listener called when state changes
	// happen on the underlying MediaElement.
	//
	// The passed listener should have the following
	// signiture: (MediaController, Event) => void
	//
	// Event names are the same as MediaElement so the
	// documentation for MediaElement listeners is still
	// valid.
	//
	// However, there are a few additonal event types:
	// - 'set' is called when the MediaElement is set.
	// - 'unset' is called when the MediaElement is unset.
	// These are also called when the media is reset, all
	// 'unset' listeners then all 'set' listeners even though
	// the MediaElement hasn't actually changed.
	on(eventType, listener, once = false) {}

	// off removes an event listener added via the on
	// function.
	off(eventType, listener) {}

	// play ...
	play() {}

	// pause ...
	pause() {}
}
