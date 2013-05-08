package channelExample

import (
	"appengine"
	"appengine/memcache"
	"time"
)

type Users struct {
	Name    string
	context appengine.Context
	member  map[string]int32
}

func (u *Users) loadUserList() {
	_, err := memcache.JSON.Get(u.context, u.Name, &u.member)
	if err != nil && err != memcache.ErrCacheMiss {
		u.context.Infof("warning")
		return
	}
}

func (u *Users) saveUserList() {
	_ = memcache.JSON.Set(u.context, &memcache.Item{
		Key: u.Name, Object: u.member,
	})
}

func (u *Users) RegistUser(key string) {
	u.loadUserList()
	u.member[key] = int32(time.Now().Unix())
	u.saveUserList()
}

func (u *Users) DeleteUser(key string) {
	u.loadUserList()
	delete(u.member, key)
	u.saveUserList()
}

func (u *Users) checkExpire() {
	now := int32(time.Now().Unix())
	for i, j := range u.member {
		if now-j > 60*60*2 {
			delete(u.member, i)
		}
	}
}

func getUsers(s string, c appengine.Context) *Users {
	data := &Users{
		Name:    s,
		context: c,
		member:  make(map[string]int32)}
	data.loadUserList()
	return data
}
