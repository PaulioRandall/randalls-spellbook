import path from 'node:path'

import files from './files.js'
import { compile } from 'svelte/compiler'

function compileSrcFile(relSrcFile) {
	const inFile = files.absSrc(relSrcFile)
	const src = files.read(inFile)
	const result = compileSvelteFile(inFile, src)
  const json = JSON.stringify(result, null, 2)

  let outFile = files.absBuild(relSrcFile)
	outFile = files.replaceExt(outFile, 'json')

	files.writeMakePath(outFile, json)
}

function compileSvelteFile(f) {
  return compile(
  	f, //
  	{
    	filename: files.basename(f),
    	css: 'injected',
    	generate: 'client',
    	runes: true,
    	modernAst: true,
  	},
  )
}

files.remake(files.BUILD_DIR)

const svelteFiles = files.listSrcForExt('svelte')

for (const f of svelteFiles) {
	compileSrcFile(f)
}
