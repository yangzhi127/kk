# AWS-Key-Hunter

AWS Key Hunter is a powerful and automated tool that scans GitHub repositories for exposed AWS keys. It continuously monitors commits, detects AWS secrets in both base64 and plaintext formats, and alerts users about potential security risks on Discord.

## Features 🚀

- Real-time Monitoring: Watches for new commits in GitHub repositories.
- AWS Key Detection: Identifies AWS keys in both plaintext and base64-encoded formats.
- Automated Scanning: Runs periodic searches for exposed AWS credentials.
- Efficient & Secure: Optimized for minimal resource usage and packaged in a secure Docker container.
- Discord integration to get alerts for any valid findings. 

## Installation 📥

Create a `.env` file and add your **Github** token and your **Discord** Server's web hook in the file. 

### Using Docker

Build the Docker image
```bash
docker build -t aws-key-scanner .
```
Run the container
```bash
docker run --rm -d --name aws-scanner aws-key-scanner
```

## Usage 🛠

Running Locally
```bash
go run main.go
```

## Disclaimer

This tool was created for educational and experimental purposes only. They are not intended to be used for malicious activities or to harm others in any way. I do not endorse or encourage the use of this tool or information for illegal, unethical, or harmful actions.

By using this tool or reading the article, you agree to accept full responsibility for any consequences that may arise from its use. I will not be held accountable for any damages, losses, or legal repercussions resulting from the misuse of this tool or the information provided.

Use at your own risk.

## Contributing 🤝

Contributions are welcome! Feel free to open an issue or submit a PR.