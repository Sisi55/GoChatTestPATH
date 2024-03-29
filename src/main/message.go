package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mholt/binding"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"time"
)

const messageFetchSize = 10

type Message struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	RoomId    bson.ObjectId `bson:"room_ID" json:"room_id"`
	Content   string        `bson:"content" json:"content"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	User      *User         `bson:"user" json:"user"`
}

func (m *Message) FileMap(r *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&m.Content: "content", // 왜 이건 r 이 아닐까
	}
}

func (m *Message) create() error {
	// 몽고 DB 세션 생성
	session := mongoSession.Copy()
	// 몽고 DB 세션을 닫는 코드를 defer 로 등록
	defer session.Close()

	// 몽고 DB 아이디 생성
	m.ID = bson.NewObjectId()
	// 메시지 생성 시간 기록
	m.CreatedAt = time.Now()
	// message 정보 저장을 위한 몽고 DB 컬렉션 객체 생성
	c := session.DB("test").C("messages")

	// messages 컬렉션에 message 정보 저장
	if err := c.Insert(m); err != nil {
		return err
	}
	return nil
}

func retrieveMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// 몽고 DB 세션 생성
	session := mongoSession.Copy()
	// 콩고 DB 세션을 닫는 코드를 defer 로 등록
	defer session.Close()

	// 쿼리 파라미터로 전달된 limit 값 확인
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		// 정상적인 limit 값이 전달되지 않은 경우
		// limit를 messageFetchSize 로 셋팅
		limit = messageFetchSize
	}

	var messages []Message
	// _id 역순으로 정렬하여 limit 수만큼 message 조회
	err = session.DB("test").C("messages").
		Find(bson.M{"room_id": bson.ObjectIdHex(ps.ByName("id"))}).
		Sort("-_id").Limit(limit).All(&messages)
	if err != nil {
		// 오류 발생 시 500 에러 리턴
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	// message 조회 결과 리턴
	renderer.JSON(w, http.StatusOK, messages)
}
