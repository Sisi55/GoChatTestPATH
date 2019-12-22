package main

import (
	"github.com/mholt/binding"
	"net/http"
)

type Room struct {
	Id   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
}

func (r *Room) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{&r.Name: "name"}
}
