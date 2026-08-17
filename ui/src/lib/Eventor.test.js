import Eventor from './Eventor.js'

function mockEventTarget() {
	let addCount = 0
	let removeCount = 0

	return {
		addCount: () => addCount,
		removeCount: () => removeCount,
		addEventListener: () => addCount++,
		removeEventListener: () => removeCount++,
	}
}

describe('Eventor.js', () => {
	const listenerA = () => {}
	const listenerB = () => {}

	test('on(eventType, listener, options)', () => {
		const ev = new Eventor()
		const exp = Eventor.newListenerEntry(
			'testevent', //
			listenerA
		)

		const off = ev.on('testevent', listenerA)

		expect(ev._entries.length).toEqual(1)
		expect(ev._entries[0]).toEqual(exp)
		expect(typeof off).toEqual('function')
	})

	test('on(...) ignores duplicate requests', () => {
		const ev = new Eventor()

		ev.on('testevent', listenerA)
		ev.on('testevent', listenerA)

		expect(ev._entries.length).toEqual(1)
	})

	test('on(listenerObject)', () => {
		const ev = new Eventor()

		const expA = Eventor.newListenerEntry(
			'testeventA', //
			listenerA
		)

		const expB = Eventor.newListenerEntry(
			'testeventB', //
			listenerB
		)

		const off = ev.on({
			testeventA: listenerA,
			testeventB: listenerB,
		})

		expect(ev._entries.length).toEqual(2)
		expect(ev._entries[0]).toEqual(expA)
		expect(ev._entries[1]).toEqual(expB)
		expect(typeof off).toEqual('function')
	})

	test('on(...) with options', () => {
		const ev = new Eventor()

		const entryA = Eventor.newListenerEntry(
			'testevent', //
			listenerA,
			{ capture: true }
		)

		const entryB = Eventor.newListenerEntry(
			'testevent', //
			listenerA,
			{ capture: false }
		)

		const entryC = Eventor.newListenerEntry(
			'testevent', //
			listenerA,
			true
		)

		const entryD = Eventor.newListenerEntry(
			'testevent', //
			listenerA,
			false
		)

		ev.on('testevent', listenerA, { capture: true })
		ev.on('testevent', listenerA, { capture: false })
		ev.on('testevent', listenerA, true)
		ev.on('testevent', listenerA, false)

		expect(ev._entries.length).toEqual(2)
		expect(ev._entries[0]).toEqual(entryA)
		expect(ev._entries[1]).toEqual(entryB)
	})

	test('isOn(...)', () => {
		const ev = new Eventor()

		const entryA = Eventor.newListenerEntry(
			'testevent', //
			listenerA
		)

		const entryB = Eventor.newListenerEntry(
			'testevent', //
			listenerB
		)

		const entryC = Eventor.newListenerEntry(
			'testeventA', //
			listenerA
		)

		ev._entries.push(entryA)
		let act = null

		act = ev.isOn('testevent', listenerA)
		expect(act).toEqual(true)

		act = ev.isOn('testevent', listenerB)
		expect(act).toEqual(false)

		act = ev.isOn('testeventA', listenerA)
		expect(act).toEqual(false)
	})

	test('off(eventType, listener, options)', () => {
		const ev = new Eventor()
		const entry = Eventor.newListenerEntry(
			'testevent', //
			listenerA
		)

		ev._entries.push(entry)
		ev.off('testevent', listenerA)

		expect(ev._entries.length).toEqual(0)
	})

	test('off(listenerObject)', () => {
		const ev = new Eventor()

		const entryA = Eventor.newListenerEntry(
			'testeventA', //
			listenerA
		)

		const entryB = Eventor.newListenerEntry(
			'testeventB', //
			listenerB
		)

		ev._entries.push(entryA)
		ev._entries.push(entryB)

		ev.off({
			testeventA: listenerA,
			testeventB: listenerB,
		})

		expect(ev._entries.length).toEqual(0)
	})

	test('addTo(...) removeFrom(...)', () => {
		const ev = new Eventor()
		const target = mockEventTarget()

		const entryA = Eventor.newListenerEntry(
			'testeventA', //
			listenerA
		)

		const entryB = Eventor.newListenerEntry(
			'testeventB', //
			listenerB
		)

		ev._entries.push(entryA)
		ev._entries.push(entryB)

		expect(target.addCount()).toEqual(0)
		expect(target.removeCount()).toEqual(0)

		ev.addTo(target)
		expect(target.addCount()).toEqual(2)
		expect(target.removeCount()).toEqual(0)

		ev.removeFrom(target)
		expect(target.addCount()).toEqual(2)
		expect(target.removeCount()).toEqual(2)
	})
})
