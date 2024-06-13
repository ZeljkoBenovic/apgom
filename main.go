package main

import (
	"log"

	"github.com/ZeljkoBenovic/apgom/internal/app"
	"github.com/ZeljkoBenovic/apgom/internal/config"
)

func main() {
	conf := config.NewConfig()
	log.Fatalln(app.NewApp(conf).Run())
}
