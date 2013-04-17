package channelExample

import (
	"html/template"
	"net/http"
	"time"

	"appengine"
	"appengine/channel"
	"appengine/memcache"
	"appengine/user"
)

var m map[string]int32

var mainTemplate = template.Must(template.ParseFiles("channelDart/main.html"))

func init() {
	http.HandleFunc("/_ah/channel/connected/", connected)
	http.HandleFunc("/_ah/channel/disconnected/", disconnected)
	http.HandleFunc("/", main)
	http.HandleFunc("/receive", receive)

}

func main(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c) // assumes 'login: required' set in app.yaml
	key := r.FormValue("gamekey")
	tok, err := channel.Create(c, u.ID+key)
	if err != nil {
		http.Error(w, "Couldn't create Channel", http.StatusInternalServerError)
		c.Errorf("channel.Create: %v", err)
		return
	}

	err = mainTemplate.Execute(w, map[string]string{
		"token":    tok,
		"me":       u.ID,
		"game_key": key,
	})
	if err != nil {
		c.Errorf("mainTemplate: %v", err)
	}

	//dataStore(c, u.ID+key)
}

func dataStore(c appengine.Context, id string) {

	m = make(map[string]int32)
	_, err := memcache.JSON.Get(c, "users", &m)
	if err != nil && err != memcache.ErrCacheMiss {
		return
	}
}

func getUserList(c appengine.Context) {
	m = make(map[string]int32)
	_, err := memcache.JSON.Get(c, "users", &m)
	if err != nil && err != memcache.ErrCacheMiss {
		return
	}
}

func setUserList(c appengine.Context) {
	_ = memcache.JSON.Set(c, &memcache.Item{
		Key: "users", Object: m,
	})
}

func receive(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	//key := r.FormValue("g")
	getUserList(c)
	for i, _ := range m {
		// c.Infof("%v", i)
		// c.Infof("%v", v)
		channel.Send(c, i, "go receive!"+time.Now().String())
	}

}
func connected(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	key := r.FormValue("from")

	getUserList(c)
	m[key] = int32(time.Now().Unix())
	setUserList(c)

	c.Infof("connected")
	c.Infof("%s", key)
}

func disconnected(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	key := r.FormValue("from")

	getUserList(c)
	delete(m, key)
	setUserList(c)

	c.Infof("disconnected")
	c.Infof("%s", key)
}
