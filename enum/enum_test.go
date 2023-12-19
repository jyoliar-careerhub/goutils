package enum_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jae2274/goutils/enum"
)

type SiteValues struct{}

type Site = enum.Enum[SiteValues]

const (
	jumpit = Site("jumpit")
	wanted = Site("wanted")
)

func (SiteValues) Values() []string {
	return []string{string(jumpit), string(wanted)}
}

func Test_MarshalJumpint(t *testing.T) {
	site, err := jumpit.MarshalText()

	checkErr(err, t)

	checkEqual[string](string(site), "jumpit", t)
}

func Test_MarshalWanted(t *testing.T) {
	site, err := wanted.MarshalText()

	checkErr(err, t)

	checkEqual[string](string(site), "wanted", t)
}

func Test_UnmarshalJumpit(t *testing.T) {
	var site Site

	err := site.UnmarshalText([]byte("jumpit"))
	checkErr(err, t)

	checkEqual[string](string(site), "jumpit", t)
}

func Test_UnmarshalWanted(t *testing.T) {
	var site Site

	err := site.UnmarshalText([]byte("wanted"))
	checkErr(err, t)

	checkEqual[string](string(site), "wanted", t)
}

func Test_MarshalJobPosting(t *testing.T) {
	posintg := jobPosting{
		Title: "Hiring Google",
		Site:  jumpit,
	}

	postingStr, err := json.Marshal(posintg)
	checkErr(err, t)

	checkEqual[string](string(postingStr), "{\"title\":\"Hiring Google\",\"site\":\"jumpit\"}", t)
}

func Test_UnmarshalJobPosting(t *testing.T) {
	jsonStr := "{\"title\":\"Hiring Google\", \"site\":\"jumpit\"}"

	var posting jobPosting
	err := json.Unmarshal([]byte(jsonStr), &posting)

	checkErr(err, t)
	checkEqual[string](posting.Title, "Hiring Google", t)

	checkEqual[Site](posting.Site, jumpit, t)
	fmt.Println(posting)

	_, err = json.Marshal(posting)

	checkErr(err, t)
}

type jobPosting struct {
	Title string `json:"title"`
	Site  Site   `json:"site"`
}

func checkErr(err error, t *testing.T) {
	if err != nil {
		t.Error("error caused: ", err)
	}
}

func checkEqual[T comparable](actual, expected T, t *testing.T) {
	if expected != actual {

		t.Errorf("Expected: %v, Actual: %v\n\n", expected, actual)
	}
}

// func assertNil(i interface{}, t *testing.T) {
// 	if i != nil {
// 		t.Errorf("% is not nil", i)
// 	}
// }
