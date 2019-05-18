package youtube

import "testing"

func TestParseJsonFromDirtyStringSuccess(t *testing.T) {
	dirtyStr := `aaaa{"a":"b"}asdas`
	str, err := extractJsonFromString(dirtyStr)
	if err != nil {
		t.Error(err.Error())
	}
	if str != `{"a":"b"}` {
		t.Error("Strings are not equal")
	}
}

func TestParseJsonFromStringSuccess(t *testing.T) {
	dirtyStr := `{"a":"b"}`
	str, err := extractJsonFromString(dirtyStr)
	if err != nil {
		t.Error(err.Error())
	}
	if str != `{"a":"b"}` {
		t.Error("Strings are not equal")
	}
}

func TestParseJsonFromNoJsonString(t *testing.T) {
	dirtyStr := `asdasdagdfsdg`
	_, err := extractJsonFromString(dirtyStr)
	if err == nil {
		t.Error("err is nil!!!")
	}
}

func TestGetBetweenPositive(t *testing.T) {
	dirtyStr := `dirt;json.config={"a":"b"};more.config={"b":"a"}`
	res, err := getBetween(dirtyStr, "json.config=", ";more.config")
	if err != nil {
		t.Error(err.Error())
	}

	if res != `{"a":"b"}` {
		t.Errorf("Incorrect string! Expected: %s, got: %s", `{"a":"b"}`, res)
	}
}

func TestGetBetweenNoFirstMatch(t *testing.T) {
	dirtyStr := `{"a":"b"};more.config={"b":"a"}`
	_, err := getBetween(dirtyStr, "json.config=", ";more.config")
	if err == nil {
		t.Error("err is nil!!!")
	}
}

func TestGetBetweenNoLastMatch(t *testing.T) {
	dirtyStr := `json.config={"a":"b"}`
	_, err := getBetween(dirtyStr, "json.config=", ";more.config")
	if err == nil {
		t.Error("err is nil!!!")
	}
}
