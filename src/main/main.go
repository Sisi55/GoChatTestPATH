package main

import (
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"gopkg.in/mgo.v2"
	"net/http"
)

const (
	// 어플리케이션에서 사용할 세션의 키 정보
	sessionKey    = "simple_cat_session"
	sessionSecret = "simple_chat_session_secret"

	socketBufferSize = 1024
)

var (
	mongoSession *mgo.Session
	renderer     = render.New()
	//upgrader     = &websocket.Upgrader{
	//	ReadBufferSize:  socketBufferSize,
	//	WriteBufferSize: socketBufferSize,
	//}
)

func init() {
	// 렌더러 생성
	renderer = render.New()

	s, err := mgo.Dial("")
	if err != nil {
		panic(err)
	}

	//// 세션 스위치. 단조로운 행동으로 ?
	//s.SetMode(mgo.Monotonic, true)
	mongoSession = s
}

func main() {
	// 라우터 생성
	router := httprouter.New()

	// 핸들러 정의
	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// 렌더러를 사용하여 템플릿 렌더링
		renderer.HTML(w, http.StatusOK, "index", map[string]interface{}{"title": "simple chat!"})
	})

	//router.GET("/rooms", retrieveRooms)
	//router.GET("/rooms/:id", retrieveRoom)
	//router.POST("/rooms", createRoom)
	//router.DELETE("/rooms/:id", deleteRoom)
	//
	//router.GET("/rooms/:id/message", retrieveMessages)
	////router.POST("/rooms/:id/messages", createMessage)
	//
	//router.GET("/info", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	//	u := GetCurrentUser(r)
	//	info := map[string]interface{}{"current_user": u, "clients":clients}
	//	renderer.JSON(w, http.StatusOK, info)
	//})
	//
	router.GET("/login", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
		// login 페이지 렌더링
		renderer.HTML(w, http.StatusOK, "login", nil)
	})
	router.GET("/logout", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
		// 세션에서 사용자 정보 제거 후 로그인 페이지로 이동
		sessions.GetSession(r).Delete(currentUserKey)
		http.Redirect(w, r, "/login", http.StatusFound)
	})
	router.GET("auth/:action/:provider", loginHandler)
	//
	//router.GET("/ws/:room_id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	//	socket, err := upfrader.Upfrade(w, r, nil)
	//	if err != nil {
	//		log.Fatal("ServeHTTP: ", err)
	//		return
	//	}
	//	newClient(socket, ps.ByName("room_id"), GetCurrentUser(r))
	//})

	// negroni 미들웨어 생성
	n := negroni.Classic()

	store := cookiestore.New([]byte(sessionSecret))
	n.Use(sessions.Sessions(sessionKey, store))

	n.Use(LoginRequired("/login", "/auth"))

	// negroni 에 router 를 핸들러로 등록
	n.UseHandler(router)

	// 웹 서버 실행
	n.Run(":3000")
}
