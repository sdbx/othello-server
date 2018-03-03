package ws

const jsonErrorMsg = `
{
	"type":"error",
	"from":"none",
	"msg":"json error"
}
`
const typeErrorMsg = `
{
	"type":"error",
	"from":"none",
	"msg":"undefined type"
}
`

const pongMsg = `
{
	"type":"pong"
}
`

const userNoMsg = `
{
	"type":"error",
	"from":"%s",
	"msg":"user doesn't exist"
}
`

const gameNoMsg = `
{
	"type":"error",
	"from":"enter",
	"msg":"game doesn't exist"
}
`

const onceMsg = `
{
	"type":"error",
	"from":"enter",
	"msg":"login should be occured once in a session"
}
`

const disconnectMsg = `
{
	"type":"disconnect",
	"username":"%s"
}
`

const connectMsg = `
{
	"type":"connect",
	"username":"%s"
}
`
