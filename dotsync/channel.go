package channelExample

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/channel"
	"appengine/user"
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
	getUsers("users", c)
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
}

func receive(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	//key := r.FormValue("g")
	//param := r.FormValue("p")
	msg := r.FormValue("m")

	users := getUsers("users", c)
	for i, _ := range users.member {
		//channel.Send(c, i, "go receive!"+time.Now().String())
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
