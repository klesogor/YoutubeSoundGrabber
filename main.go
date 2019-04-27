package main

import (
	"fmt"
	"strings"
)

const rightReq = "https://r3---sn-3vh25opt5n8t51p-gv8e.googlevideo.com/videoplayback?id=o-ALFwVrf0d_uHFytxnIO_xMSdMSIMhV8gb_tjMa_wpwTe&itag=251&source=youtube&requiressl=yes&pl=24&ei=HVnEXKrXBoW67QTW4LPoAQ&mime=audio%2Fwebm&gir=yes&clen=9835106&dur=551.541&lmt=1540159339728523&fvip=16&keepalive=yes&c=WEB&txp=5411222&ip=193.23.58.204&ipbits=0&expire=1556393341&signature=1CA11B3E61FC37096689F4224B59F5FA7B26DD6D.5E3D1D16326C9383AC60FA7870E11C0F7C24E72C&key=yt8&alr=yes&cpn=dbNWCVj6f-X2DYv1&cver=2.20190423&mm=31%2C29&mn=sn-3vh25opt5n8t51p-gv8e%2Csn-n8v7znly&ms=au%2Crdu&mt=1556371609&mv=m&nh=%2CIgpwcjAzLnN2bzA2KgkxMjcuMC4wLjE&range=0-193127&rn=13&rbuf=78331&initcwndbps=1396250&sparams=ip%2Cipbits%2Cexpire%2Cid%2Citag%2Csource%2Crequiressl%2Cmm%2Cmn%2Cms%2Cmv%2Cpl%2Cnh%2Cei%2Cinitcwndbps%2Cmime%2Cgir%2Cclen%2Cdur%2Clmt"
const initReq = "https://r3---sn-3vh25opt5n8t51p-gv8e.googlevideo.com/videoplayback?id=o-AG83OgWATf7vvccGJAHyKCwgo3oHYNqnov-NZ3wKok6m&itag=250&source=youtube&requiressl=yes&mm=31,29&mn=sn-3vh25opt5n8t51p-gv8e,sn-n8v7kn7d&ms=au,rdu&mv=m&pl=24&nh=,IgpwcjAzLnN2bzA2KgkxMjcuMC4wLjE&ei=elrEXJu4JcLf7gSL8qmgAQ&initcwndbps=1467500&mime=audio/webm&gir=yes&clen=5091722&dur=551.541&lmt=1540159340834561&mt=1556371992&fvip=16&keepalive=yes&c=WEB&txp=5411222&ip=193.23.58.204&ipbits=0&expire=1556393690&sparams=ip,ipbits,expire,id,itag,source,requiressl,mm,mn,ms,mv,pl,nh,ei,initcwndbps,mime,gir,clen,dur,lmt&signature=B86F506911E7EEF657D9097AC526DDC5C553231D.2E8F8F1385E802C0F3648573EB763973E8E2C9F8&key=yt8&itag=250&init=0-258&type=audio/webm"

func main() {
	printDiffMap(parseRequest(rightReq), parseRequest(initReq))
}

func parseRequest(s string) map[string]string {
	splitted, res := strings.Split(s, "&"), make(map[string]string)

	for _, v := range splitted {
		if pair := strings.Split(v, "="); len(pair) == 2 {
			res[pair[0]] = res[pair[1]]
		}
	}

	return res
}

func printDiffMap(a, b map[string]string) {
	for k, v := range a {
		val := b[k]
		fmt.Printf("%v: %v - %v \n", k, v, val)
	}
}
