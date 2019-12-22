package main

import (
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"strings"
)

const (
	// 세션에 저장되는 next page 의 키
	nextPageKey = "next_page"
	authSecurityKey = "auth_security_key"
)

func init(){
	// gomniauth 정보 셋팅
	gomniauth.SetSecurityKey(authSecurityKey)
	gomniauth.WithProviders(
		google.New("624827322356-eso1pumc9o3oau63t96ar5gbr0728tnp.apps.googleusercontent.com",
			"LIJzbAhN9zoexM9F-ozvuPBU",
			"http://localhost:3000/auth/callback/google"))
}

func LoginRequired(ignore ...string) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// ignore url 인 경우 다음 핸들러 실행
		for _, s := range ignore {
			if strings.HasPrefix(r.URL.Path, s){
				next(w, r)
				return
			}
		}
		// CurrentUser 정보 가져옴
		u := GetCurrentUser(r)

		// CurrentUser 정보가 유효한 경우 만료시간을 갱신하고
		// 다음 핸들러 실행
		if u != nil && u.Valid(){
			SetCurrentUser(r, u)
			next(w, r)
			return
		}

		// CurrentUser 정보가 유효하지 않은 경우
		// CurrentUser 를 nil 로 셋팅
		SetCurrentUser(r, nil)

		// 로그인 후 이동할 url 을 세션에 저장 r
		// 현재 사용자가 들어온 경로를 다음 페이지로 저장해서
		// 로그인 후 리다이렉트하는구나
		sessions.GetSession(r).Set(nextPageKey, r.URL.RequestURI())

		// login 페이지로 redirect
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	action := ps.ByName("action")
	provider := ps.ByName("provider")
	s := sessions.GetSession(r)

	switch action {
	case "login":
		// gomniauth Provider 의 login 페이지로 이동
		p, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln(err)
		}
		loginUrl, err := p.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln(err)
		}
		http.Redirect(w, r, loginUrl, http.StatusFound)
	case "callback":
		// gomniauth 콜백 처리
		p, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln(err)
		}

		creds, err := p.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln(err)
		}

		// 콜백 결과로부터 사용자 정보 확인
		user, err := p.GetUser(creds)
		if err != nil {
			log.Fatalln(err)
		}

		u := &User{
			Uid: user.Data().Get("id").MustStr(),
			Name: user.Name(),
			Email: user.Email(),
			AvatarUrl: user.AvatarURL(),
		}

		// 사용자 정보를 세션에 저장
		SetCurrentUser(r, u)

		http.Redirect(w, r, s.Get(nextPageKey).(string), http.StatusFound)
	default:
		http.Error(w, "Auth action '"+action+"' is not supported", http.StatusNotFound)
	}

}
