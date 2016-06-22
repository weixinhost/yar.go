package yar

type YarOpt int

const (
	//YarOptMagicNumber 代表协议中规定的相关信息，用于进行请求验证使用
	YarOptMagicNumber YarOpt = 1
	//YarOptTimeout 处理超时
	YarOptTimeout = 2
	//YarOptConnectTimeout 连接超时，只有在非http类server中有意义
	YarOptConnectTimeout = 3
	//YarOptPackager 数据打包协议，目前只支持json
	YarOptPackager = 4
	//YarOptEncrypt 是否启用加密
	YarOptEncrypt = 5
	//YarOptEncryptPrivateKey 用于加密的aes key
	YarOptEncryptPrivateKey = 6
)

type Opt struct {
	MagicNumber       uint32
	Timeout           uint32
	ConnectTimeout    uint32
	Packager          string
	Encrypt           bool
	EncryptPrivateKey string
	DynamicParam      bool
}

func NewOpt() *Opt {
	opt := new(Opt)
	opt.MagicNumber = MagicNumber
	opt.Encrypt = false
	opt.EncryptPrivateKey = ""
	opt.Packager = "json"
	opt.ConnectTimeout = 1000 * 5
	opt.Timeout = 30 * 1000
	opt.DynamicParam = false
	return opt
}
