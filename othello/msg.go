package othello

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

const roomNoMsg = `
{
	"type":"error",
	"from":"%s",
	"msg":"room doesn't exist"
}
`

const onceMsg = `
{
	"type":"error",
	"from":"login",
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