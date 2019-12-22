package main

import (
	"encoding/json"
	sessions "github.com/goincremental/negroni-sessions"
	"net/http"
	"time"
)

const (
	// 세션에 저장되는 CurrentUser 의 키
	currentUserKey = "oauth2_current_user"
	// 로그인 세션 유지 시간
	sessionDuration = time.Hour
)

type User struct {
	Uid string `json:"uid"`
	Name string `json:"name"`
	Email string `json:"user"`
	AvatarUrl string `json:"avatar_url"`
	Expired time.Time `json:"expired"`
}

func (u *User) Valid() bool {
	// 현재 시간 기준으로 만료 시간 확인
	return u.Expired.Sub(time.Now()) > 0
}

func (u *User) Refresh(){
	// 만료 시간 연장
	u.Expired = time.Now().Add(sessionDuration)
}

func GetCurrentUser(r *http.Request) *User {
	// 세션에서 CurrentUser 정보를 가져옴
	s := sessions.GetSession(r)
	// Request 에 있는 Session 을 반환한다

	if s.Get(currentUserKey) == nil {
		// 세션이 없으면 nil 을 반환한다 -- Get
		return nil
	}

	data := s.Get(currentUserKey).([]byte)
	var u User
	json.Unmarshal(data, &u)
	return &u
}

func SetCurrentUser(r *http.Request, u *User) {
	if u != nil {
		// CurrentUser 만료 시간 갱신
		u.Refresh()
	}

	// 세션에 CurrentUser 정보를 json 으로 저장
	s := sessions.GetSession(r)
	val, _ := json.Marshal(u)
	s.Set(currentUserKey, val)
}