package main

import (
	"app/internal/appstatus"
	"app/internal/utils"

	"github.com/google/uuid"
	//"encoding/json"
	//ap "github.com/go-ap/activitypub"
	//"github.com/go-ap/jsonld"
)

func main() {
    //person := ap.PersonNew("test")
    //person.Context = ap.IRIs([]ap.IRI{
    //    ap.ActivityBaseURI,
    //})
    //person.Name = ap.NaturalLanguageValuesNew(ap.LangRefValueNew("en", "user"), ap.LangRefValueNew("th", "กดหกหด"))
    //person.Icon = ap.ItemCollection{ap.Image{
    //    Type: ap.ImageType,
    //    URL: ap.IRI("dsfdf"),
    //}}
    //body,err := jsonld.WithContext(
    //jsonld.IRI(ap.ActivityBaseURI),).Marshal(person)
    //j, err := json.Marshal(person)
    //if err != nil {
    //    panic(err)
    //}
    //println(string(j))
    //println(string(body))

    key, err := utils.GenerateRandomKey(45)
    if err != nil {
        panic(err)
    }
    println(key)
    println(uuid.New().String())
    println(appstatus.BadLogin("bad value", "dsfdf").Error())
    println(appstatus.BadLogin("bad value").Code())
    println(appstatus.BadLogin("bad value").Status())
}
