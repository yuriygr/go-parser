package main

import (
	"log"
	"os"
	"os/signal"

	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron"
)

func main() {
	client := NewClient()
	storage := NewStorage()

	parser := &BoardParser{client, storage, "2ch.hk", []string{"b", "po", "mov", "tv", "a", "v", "cg", "vg", "kpop"}}
	actualizer := &Actualizer{client, storage}

	log.Println("===> Cron created")

	c := cron.New()
	c.AddJob("@every 30m", parser)
	c.AddJob("@every 40m", actualizer)
	c.Start()

	// Debug
	// parser.Run()
	// actualizer.Run()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
