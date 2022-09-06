package pic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type InfoPic struct {
	pid       int
	uid       int
	series    int
	stat      int
	views     int
	bookmark  int
	count     int
	timestamp int
	title     string
	maps      map[string]string
}

func getter(sid string, client *http.Client) func(s string) io.ReadCloser {
	header := DefaultHeader()
	header.Add("referer", "https://www.pixiv.net/member_illust.php?mode=medium&illust_id="+sid)
	header.Add("accept-language", "en")
	return func(s string) io.ReadCloser {
		ur := &url.URL{}
		ur, _ = ur.Parse(s)
		request := &http.Request{
			Method: "GET",
			URL:    ur,
			Header: *header,
		}
		res, err := client.Do(request)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return res.Body
	}
}

func filer(dir string) func(string, io.Reader) error {
	return func(fname string, reader io.Reader) (e error) {
		file, e := os.OpenFile(filepath.Join(dir, fname), os.O_WRONLY|os.O_CREATE, 0644)
		if e != nil {
			return
		}
		defer file.Close()
		_, e = io.Copy(file, reader)
		if e != nil {
			return
		}
		return nil
	}
}

func Crawl(pid int, dir string, client *http.Client) (ans InfoPic, e error) {
	ans.maps = make(map[string]string)
	sid := fmt.Sprintf("%d", pid)

	f := getter(sid, client)
	g := filer(filepath.Join(dir, sid))

	infoBody := f("https://www.pixiv.net/touch/ajax/illust/details?illust_id=" + sid)
	if infoBody == nil {
		e = errors.New("infoBody == nil")
		return
	}

	defer infoBody.Close()

	con, e := io.ReadAll(infoBody)
	if e != nil {
		return
	}

	var x interface{}
	json.Unmarshal(con, &x)

	j := x.(map[string]interface{})
	body := j["body"].(map[string]interface{})
	detail := body["illust_details"].(map[string]interface{})

	// pid INTEGER PRIMARY KEY NOT NULL
	ans.pid, e = strconv.Atoi(detail["id"].(string))
	if e != nil {
		return
	}
	if ans.pid != pid {
		e = errors.New("ans.pid != pid")
		return
	}

	// uid INTEGER NOT NULL
	ans.uid, e = strconv.Atoi(detail["user_id"].(string))
	if e != nil {
		return
	}

	// series INTEGER NOT NULL
	if detail["series"] == nil {
		types, err := strconv.Atoi(detail["type"].(string))
		e = err
		if e != nil {
			return
		}
		ans.series = -types
	} else {
		series := detail["series"].(map[string]interface{})
		ans.series, e = strconv.Atoi(series["id"].(string))
		if e != nil {
			return
		}
	}

	// stat INTEGER NOT NULL
	ans.stat = 1

	// views INTEGER NOT NULL
	ans.views, e = strconv.Atoi(detail["rating_view"].(string))
	if e != nil {
		return
	}

	// bookmark INTEGER NOT NULL
	ans.bookmark = int(detail["bookmark_user_total"].(float64))

	// count INTEGER NOT NULL
	ans.count, e = strconv.Atoi(detail["rating_count"].(string))
	if e != nil {
		return
	}

	// timestamp INTEGER NOT NULL
	ans.timestamp = int(detail["upload_timestamp"].(float64))

	for _, j := range detail["display_tags"].([]interface{}) {
		tags := j.(map[string]interface{})
		if tags["translation"] == nil {
			ans.maps[tags["tag"].(string)] = ""
		} else {
			ans.maps[tags["tag"].(string)] = tags["translation"].(string)
		}
	}

	// title VARCHAR(64) NOT NULL
	comment := ""
	if detail["title"] == nil {
		ans.title = ""
	} else {
		ans.title = detail["title"].(string)
		comment = "#### " + ans.title + "\n"
	}

	dir = filepath.Join(dir, sid)
	os.MkdirAll(dir, 0644)

	if detail["comment"] != nil {
		comment += "\n" + detail["comment"].(string) + "\n"
	}
	fp := sid + ".md"
	e = g(fp, strings.NewReader(comment))
	if e != nil {
		return
	}

	if detail["manga_a"] == nil {
		s := detail["url_big"].(string)
		imgBody := f(s)
		if imgBody == nil {
			e = errors.New(s + " imgBody == nil")
			return
		}
		defer imgBody.Close()

		split := strings.Split(s, "/")
		fp := split[len(split)-1]
		e = g(fp, imgBody)
		return
	}

	for i, j := range detail["manga_a"].([]interface{}) {
		manga := j.(map[string]interface{})
		s := manga["url_big"].(string)
		imgBody := f(s)
		if imgBody == nil {
			e = errors.New(s + " imgBody == nil")
			ans.stat = -i * 2
			continue
		}
		defer imgBody.Close()

		split := strings.Split(s, "/")
		fp := split[len(split)-1]
		e = g(fp, imgBody)
		if e != nil {
			ans.stat = -i*2 - 1
		}
	}
	return
}
