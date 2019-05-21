package youtube

const (
	dechipherFunc = `function\(a\){a.split\(""\);.+}`
)

type Dechiper interface {
	decryptSignature(sig string) (string, error)
}

type DechipherAction func(s []rune) []rune

type SimpleDechipher struct{}
