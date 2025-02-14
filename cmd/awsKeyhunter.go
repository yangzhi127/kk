package main

import (
	"fmt"
	"log"
	"os"

	"github.com/iamlucif3r/aws-key-hunter/internal/pkg"
	"github.com/joho/godotenv"
)

const (
	Red    = "\033[31m"
	Yellow = "033[33m"
	Green  = "033[32m"
	Reset  = "\033[0m"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
func main() {
	fmt.Println(Red + "┏┓┓ ┏┏┓  ┓┏┓      ┓┏     		" + Reset)
	fmt.Println(Red + "┣┫┃┃┃┗┓━━┃┫ ┏┓┓┏━━┣┫┓┏┏┓╋┏┓┏┓	" + Reset)
	fmt.Println(Red + "┛┗┗┻┛┗┛  ┛┗┛┗ ┗┫  ┛┗┗┻┛┗┗┗ ┛ 	" + Reset)
	fmt.Println(Red + "               ┛   v1.0.0      	" + Reset)
	fmt.Println()
	log.Println(Yellow + "🚀 Starting AWS Key Scanner..." + Reset)

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("GITHUB_TOKEN is not set")
	}
	go pkg.WatchNewFilesatchNewFiles(githubToken)
	for {
		pkg.SearchGithub(githubToken)
		// time.Sleep(1 * time.Minute)
	}
}
