package channelExample

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/channel"
	"time"
)

var mainTemplate = template.Must(template.ParseFiles("dotsync/canvas.html"))

func init() {
	http.HandleFunc("/_ah/channel/connected/", connected)
	http.HandleFunc("/_ah/channel/disconnected/", disconnected)
	http.HandleFunc("/", main)
	http.HandleFunc("/receive", receive)
}

func main(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	u := time.Now().Format("20060102150405")
	//users := getUsers("users", c)
	//users.RegistUser(u)

	key := r.FormValue("gamekey")
	tok, err := channel.Create(c, u+key)
	if err != nil {
		http.Error(w, "Couldn't create Channel", http.StatusInternalServerError)
		c.Errorf("channel.Create: %v", err)
		return
	}

	err = mainTemplate.Execute(w, map[string]string{
		"token":    tok,
		"me":       u,
		"game_key": key,
	})
	if err != nil {
		c.Errorf("mainTemplate: %v", err)
	}
}

func receive(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	c.Infof("receive")
	msg := r.FormValue("m")

	users := getUsers("users", c)
	for i, _ := range users.member {
		channel.Send(c, i, msg)
	}
}
func connected(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	key := r.FormValue("from")

	users := getUsers("users", c)
	users.RegistUser(key)

	c.Infof("connected")
	c.Infof("%s", key)
}

func disconnected(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	key := r.FormValue("from")

	users := getUsers("users", c)
	users.DeleteUser(key)

	c.Infof("disconnected")
	c.Infof("%s", key)
}
