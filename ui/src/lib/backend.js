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

async function getAllMedia() {
	if (!!window?.getAllMedia) {
		return window.getAllMedia()
	}
	return structuredClone(testMedia)
}

async function getMediaById(entityId) {
	if (!!window?.getMediaById) {
		return window.getMediaById(entityId)
	}

	const media = testMedia.find((m) => m.entityId === entityId)
	return media ? structuredClone(media) : null
}

async function selectLocalMediaFile() {
	if (!!window?.selectLocalMediaFile) {
		return window.selectLocalMediaFile()
	}
	return '/path/to/media file.mp4'
}

async function addMedia(
	mediaType, //
	name,
	description,
	localPath
) {
	if (!!window?.addMediaToProject) {
		return window.addMediaToProject(
			mediaType, //
			name,
			description,
			localPath
		)
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
	getAllMedia,
	getMediaById,
	addMedia,
	selectLocalMediaFile,
}
