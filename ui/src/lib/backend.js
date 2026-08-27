const testMedia = [
	{
		entityId: 'abcdefghijklmnopqrstuvwxyz-1111',
		mediaType: 'video',
		name: 'Alice',
		description: 'Alice in chains.',
		localPath: 'C://data/alice-in-chains.mp4',
	},
	{
		entityId: 'abcdefghijklmnopqrstuvwxyz-2222',
		mediaType: 'video',
		name: 'Bob',
		description: "Bob's your uncle.",
		localPath: 'C://data/bobs-your-uncle.mp4',
	},
	{
		entityId: 'abcdefghijklmnopqrstuvwxyz-3333',
		mediaType: 'video',
		name: 'Charlie',
		description: 'Charlie loves cheese.',
		localPath: 'C://data/charlie-loves-cheese.mp4',
	},
]

async function SelectLocalFile() {
	if (!!window?.CastSpell) {
		return window.CastSpell('SelectLocalFile', null)
	}
	return '/path/to/media file.mp4'
}

async function ListMedia() {
	if (!!window?.CastSpell) {
		return window.CastSpell('ListMedia', null)
	}
	return structuredClone(testMedia)
}

async function GetMediaById(entityId) {
	if (!!window?.CastSpell) {
		return window.CastSpell('GetMediaById', entityId)
	}

	const media = testMedia.find((m) => m.entityId === entityId)
	return media ? structuredClone(media) : null
}

async function AddMedia(
	mediaType, //
	name,
	description,
	localPath
) {
	if (!!window?.CastSpell) {
		const data = JSON.stringify({
			mediaType, //
			name,
			description,
			localPath,
		})

		return window.CastSpell('AddMedia', data)
	}

	const entityId = crypto.randomUUID()

	testMedia.push({
		entityId, //
		mediaType,
		name,
		description,
		localPath,
	})

	return entityId
}

export default {
	SelectLocalFile,
	ListMedia,
	GetMediaById,
	AddMedia,
}
