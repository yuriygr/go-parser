package main

import "log"

// Actualizer - Проверяем доступность ссылок
type Actualizer struct {
	client  *ClientInstance
	storage *Storage
}

// Run - метод, вызываемый кроном в Job
func (a *Actualizer) Run() {
	log.Print("===> 40 minutes has passed, let's truncate some links?")

	files, err := a.storage.GetFiles()
	if err != nil {
		log.Print(err)
	}

	for _, file := range files {
		if err := a.client.IsExist(file.URL); err != nil {
			if err := a.storage.DeleteFile(file.ID); err != nil {
				log.Print(err)
			}
		}
	}
}
