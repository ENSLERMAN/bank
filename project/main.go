/*
 Банковская система на обычном REST ( монолит ).
 */
package main

import (
	"github.com/ENSLERMAN/soft-eng/project/internal/app/apiserver"
	"log"
)

func main() {
	config := apiserver.NewConfig()
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
