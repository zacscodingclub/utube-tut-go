package main

import (
	"fmt"
	"net/http"

	"github.com/zacscodingclub/utube-tut/models"
	"github.com/zacscodingclub/utube-tut/routes"
	"github.com/zacscodingclub/utube-tut/utils"
)

func main() {
	models.Init()
	utils.LoadTemplates("templates/*.html")
	r := routes.NewRouter()
	port := ":8080"

	http.Handle("/", r)
	fmt.Println("Server up and running on localhost" + port)
	http.ListenAndServe(port, nil)
}
