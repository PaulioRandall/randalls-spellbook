// E.g. left mouse button
const PRIMARY_BUTTON = 0

// E.g. middle mouse button
const AUXILIARY_BUTTON = 1

// E.g. right mouse button
const SECONDARY_BUTTON = 2

function isPrimaryButton(event) {
	return event.button === PRIMARY_BUTTON
}

function isAuxiliaryButton(event) {
	return event.button === AUXILIARY_BUTTON
}

function isSecondaryButton(event) {
	return event.button === SECONDARY_BUTTON
}

export default {
	PRIMARY_BUTTON,
	AUXILIARY_BUTTON,
	SECONDARY_BUTTON,
	isPrimaryButton,
	isAuxiliaryButton,
	isSecondaryButton,
}
