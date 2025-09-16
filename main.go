package yora

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)



func InitialiseProject(projectName string) error {
	if projectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	subdirs := []string{"apps", "db", "middlewares"}
	for _, dir := range subdirs {
		path := filepath.Join(".", dir)
		if err := os.Mkdir(path, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %w", dir, err)
		}
	}

	// db file path
	dbPath := filepath.Join(".", "db")

	// Make sure project db folder exists (redundant but harmless since already created above)
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return fmt.Errorf("failed to create db folder: %w", err)
	}

	// Now create the sqlite file
	dbFile := filepath.Join(dbPath, fmt.Sprintf("%s.sqlite", projectName))
	dbContent := fmt.Sprintf("-- Placeholder schema for %s\n", projectName)
	if err := os.WriteFile(dbFile, []byte(dbContent), 0644); err != nil {
		return fmt.Errorf("unable to create db: %w", err)
	}

	// config 
	configContent := `package main;

	package yora

var HOST = "localhost"
var PORT = "2300"
var dbDriver = "sqlite"
var dbPath = "db/"
var auth = "session"
var middlewares = []string{}

var logfile = false
`
	configFile := filepath.Join(".", "config.go")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil { // 0644 for source files
		return fmt.Errorf("unable to create config: %w", err)
	}
	// Add routes.go with valid Go content
	routesContent := `package main

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	// TODO: Add your routes here
	// Example: mux.HandleFunc("/example", exampleHandler)
}
`
	routeFile := filepath.Join(".", "routes.go")

	if err := os.WriteFile(routeFile, []byte(routesContent), 0644); err != nil { // 0644 for source files
		return fmt.Errorf("unable to create routes: %w", err)
	}

	fmt.Println("Initialised project successfully")
	return nil
}

func RunServer() error {
	host := "localhost"
	port := "2300"

	if _, err := os.Stat("./config.go"); err == nil {
		// Assuming config.go defines/imports HOST and PORT; otherwise, this block won't execute
		host = HOST
		port = PORT
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>YORA Project</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background: #1f2a38;
					color: #e0e0e0;
					text-align: center;
					padding: 100px;
				}
				.container {
					background: #2a3a4f;
					border-radius: 16px;
					padding: 40px;
					display: inline-block;
					box-shadow: 0 8px 20px rgba(0,0,0,0.4);
				}
				h1 {
					font-size: 2.5rem;
					margin-bottom: 20px;
					color: #ffd166; /* soft gold accent */
				}
				p {
					font-size: 1.2rem;
					margin-bottom: 30px;
					color: #f0f0f0;
				}
				.button {
					background: #4dabf7; /* soft blue */
					color: #1f2a38;
					padding: 12px 24px;
					font-size: 1rem;
					font-weight: bold;
					border: none;
					border-radius: 8px;
					cursor: pointer;
					text-decoration: none;
					transition: background 0.3s;
				}
				.button:hover {
					background: #3a92f2;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>🚀 YORA Project</h1>
				<p>Your server is up and running successfully.</p>
				<a href="/documentation" class="button">Check documentation here!</a>
			</div>
		</body>
		</html>
		`)
	})

	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Server running at http://%s\n", addr)

	return http.ListenAndServe(addr, nil)
}

func main() {
	initFlag := flag.String("init", "", "Initialize the project with basic config")
	// createApp := flag.String("create-app", "", "Create an app with basic config")
	runserver := flag.Bool("runserver", false, "Run dev server")

	flag.Parse()

	switch {
	case *initFlag != "":
		if err := InitialiseProject(*initFlag); err != nil {
			log.Fatalf("Unable to initialize project: %v", err)
		}
		fmt.Println("Initialised project!!")

	// case *createApp != "":
	// 	if err := createAppFunc(*createApp); err != nil {
	// 		log.Fatalf("Unable to start app: %v", err)
	// 	}
	// 	fmt.Println("Created app!!")

	case *runserver:
		if err := RunServer(); err != nil {
			log.Fatalf("Unable to run server: %v", err)
		}
	}
}