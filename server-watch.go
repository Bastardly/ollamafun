package main

import (
	"log"
	"os"
	"strings"
	"time"
)

func checkForUpdatedFiles(lastModifiedMap map[string]time.Time, directory, extention string) {
	for {
		watcher, err := os.Open(directory)
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		dirEntries, err := watcher.Readdir(-1)
		if err != nil {
			log.Println("Error reading directory:", err)
			return
		}

		updated := false // Flag to track if an update has already occurred

		for _, entry := range dirEntries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), extention) {
				filePath := directory + entry.Name()
				lastModifiedTime, ok := lastModifiedMap[filePath]

				if !ok || entry.ModTime().After(lastModifiedTime) {
					if !updated {
						log.Println("File changed. Reloading server...")
						updated = true
					}
					lastModifiedMap[filePath] = entry.ModTime()
				}
			}
		}

		if updated {
			updateUITemplates()
		}

		time.Sleep(1 * time.Second) // Adjust the interval as needed
	}
}

func watchFiles() {
	println("Watching files...")
	lastModifiedMap := make(map[string]time.Time)

	// extensions := []string{".js", ".html"}

	go checkForUpdatedFiles(lastModifiedMap, "./templates/", ".html")
	go checkForUpdatedFiles(lastModifiedMap, "./static/", ".js")

	// for _, ext := range extensions {
	// 	go checkForUpdatedFiles(lastModifiedMap, "./webserver/clientTemplates/", ext)
	// }

	// Wait indefinitely
	select {}
}
