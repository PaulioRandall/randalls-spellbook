import fs from 'node:fs'
import path from 'node:path'

const SOURCE_DIR = path.resolve('./src')
const BUILD_DIR = path.resolve('./build')

function name(f) {
	const n = path.basename(f)
	return n.split(n, '.')[0]
}

function basename(f) {
	return path.basename(f)
}

function extension(f) {
	return path.extname(f).replace('.', '')
}

function parent(f) {
	return path.dirname(f)
}

function absSrc(rel) {
	return path.join(SOURCE_DIR, rel)
}

function absBuild(rel) {
	return path.join(BUILD_DIR, rel)
}

function replaceExt(f, newExt) {
	const currExt = path.extname(f)
	f = f.slice(0, -currExt.length)
	return `${f}.${newExt}`
}

function exists(f) {
	return fs.existsSync(f)
}

function isDir(f) {
	if (!exists(f)) {
		return false
	}

	const lstat = fs.lstatSync(f)
	return lstat.isDirectory()
}


function read(f) {
	return fs.readFileSync(f, { encoding: 'utf-8' })
}

function listSrc() {
	return fs.readdirSync(SOURCE_DIR, {
		recursive: true,
	})
}

function listSrcForExt(ext) {
	return listSrc().filter(f => {
		return extension(f) === ext
	})
}

function write(f, data) {
	return fs.writeFileSync(f, data, {
		encoding: 'utf-8',
		flush: true,
	})
}

function writeMakePath(f, data) {
	makePath(parent(f))
	write(f, data)
}

function make(f) {
	fs.mkdirSync(f)
}

function makePath(f) {
	fs.mkdirSync(f, {
		recursive: true,
	})
}

function remake(f) {
	remove(f)
	make(f)
}

function remove(f) {
	if (isDir(f)) {
		fs.rmSync(f, {
			force: true, //
			recursive: true,
		})
	} else {
		fs.rmSync(f, { 
			force: true // 
		})
	}
}

export default {
	BUILD_DIR,
	SOURCE_DIR,
	name,
	basename,
	parent,
	absSrc,
	absBuild,
	replaceExt,
	exists,
	isDir,
	read,
	listSrc,
	listSrcForExt,
	write,
	writeMakePath,
	make,
	remake,
	makePath,
	remove,
}
