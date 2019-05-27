package youtube

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	dechipherFunc       = `(?s)function\(a\){a=a\.split\(""\);.+?;return a\.join\(""\)};`
	chiperObjPattern    = `(?s)%s={.+?};`
	jsObjectFuncPattern = `(?s)\w+:function\(a(,b)?\){.+?}(,|})`
	basePath            = `http://www.youtube.com`
	argsRegexp          = `\w+\(a,(\d+)`
)

type Dechiper interface {
	decryptSignature(sig string) (string, error)
}

type ActionFactory func(int) DechipherAction

type DechipherAction func(s []rune) []rune

type chipherObject struct {
	swapFunc    string
	reverseFunc string
	spliceFunc  string
}

type SimpleDechipher struct {
	actions []DechipherAction
}

func (c *PlayerConfig) createDechipher() (SimpleDechipher, error) {
	res, err := http.Get(basePath + c.Assets.Js)
	if err != nil {
		return SimpleDechipher{}, err
	}
	body, err := parseBody(res)
	if err != nil {
		return SimpleDechipher{}, err
	}

	return createDechipherFromPlayreCode(string(body))
}

func createDechipherFromPlayreCode(code string) (SimpleDechipher, error) {
	function := parseDecryptionFunction(code)
	actions, err := getDechipherActions(function)
	if err != nil {
		return SimpleDechipher{}, err
	}
	chipherObject := getChipherObjectFromActions(code, actions)

	return createDechipherFromChipherAndActions(chipherObject, actions), nil
}

func parseDecryptionFunction(code string) string {
	reg := regexp.MustCompile(dechipherFunc)
	return reg.FindString(code)
}

func getDechipherActions(function string) ([]string, error) {
	code, err := getBetween(function, `a=a.split("");`, `;return a.join("")`)
	if err != nil {
		return nil, err
	}
	return strings.Split(code, ";"), nil
}

func getChipherObjectFromActions(code string, actions []string) chipherObject {
	objectName := strings.Split(actions[0], ".")[0]
	regex := regexp.MustCompile(fmt.Sprintf(chiperObjPattern, objectName))
	objectString := regex.FindString(code)
	regex = regexp.MustCompile(jsObjectFuncPattern)
	funcs := regex.FindAllString(objectString, -1)
	object := chipherObject{}
	for _, v := range funcs {
		split := strings.Split(v, ":")
		if strings.Contains(split[1], "splice") {
			object.spliceFunc = split[0]
		} else if strings.Contains(split[1], "reverse") {
			object.reverseFunc = split[0]
		} else {
			object.swapFunc = split[0]
		}
	}

	return object
}

func createDechipherFromChipherAndActions(c chipherObject, a []string) SimpleDechipher {
	actions := make([]DechipherAction, 0)
	for _, v := range a {
		action := strings.Split(v, ".")[1]
		actionName := action[:strings.Index(action, "(")]
		switch actionName {
		case c.reverseFunc:
			actions = append(actions, createReverseAction())
		case c.swapFunc:
			actions = append(actions, createActionFromArg(action, createSwapAction))
		case c.spliceFunc:
			actions = append(actions, createActionFromArg(action, createSpliceAction))
		}
	}
	return SimpleDechipher{actions: actions}
}

func createSwapAction(offset int) DechipherAction {
	return func(a []rune) []rune {
		c := a[0]
		a[0] = a[offset%len(a)]
		a[offset%len(a)] = c
		return a
	}
}

func createSpliceAction(offset int) DechipherAction {
	return func(a []rune) []rune {
		return a[offset:]
	}
}

func createReverseAction() DechipherAction {
	return func(a []rune) []rune {
		for i := len(a)/2 - 1; i >= 0; i-- {
			opp := len(a) - 1 - i
			a[i], a[opp] = a[opp], a[i]
		}

		return a
	}
}

func createActionFromArg(s string, f ActionFactory) DechipherAction {
	regex := regexp.MustCompile(argsRegexp)
	val, err := strconv.Atoi(regex.FindStringSubmatch(s)[1])
	if err != nil {
		panic(err)
	}
	return f(val)
}

func (d SimpleDechipher) decryptSignature(s string) string {
	runes := append([]rune(s))
	for _, callback := range d.actions {
		runes = callback(runes)
	}

	return string(runes)
}
