package repository

const createFile = `
	INSERT INTO video (url, type, filesize, width, height, md5, name, board, thread, created_at) VALUES
	(:v.url, :v.type, :v.filesize, :v.width, :v.height, :v.md5, :v.name, :v.board, :v.thread, :v.created_at)
`

const selectFiles = `
	select v.* from video as v order by created_at asc
`

const deleteFileByID = `
	delete from video where id = ?
`
