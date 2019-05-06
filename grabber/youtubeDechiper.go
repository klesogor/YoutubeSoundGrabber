package grabber;

type Decryptor interface{
	Decrypt(signature string) string
}

type YoutubeDechiper struct{}

func (d *YoutubeDechiper) Decrypt(s string) string{
	s = xe(s,20)
	s = ls(s, 58)
	s = vt(s, 2)
	s = xe(s, 65)
	s = vt(s, 3)
    s = xe(s, 32)
	s = vt(s, 1)
	s = ls(s, 70)
	s = xe(s, 38)
	return s
}

func xe(s string, a int) (result string){
	for _,v := range s {
		result = string(v) + result
	  }
	  return
}
func ls(s string, a int) string{
	temp := []rune(s)
	c := temp[0]
	temp[0] = temp[a%len(temp)]
	temp[a%len(temp)] = c
	return string(temp)
}
func vt(s string, a int) string{
	return s[a:]
}
