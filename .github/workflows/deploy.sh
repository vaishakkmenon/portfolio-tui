name: Deploy Portfolio TUI

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4

      - name: Deploy to VPS via SSH
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.HOST }}
          username: ${{ secrets.USERNAME }}
          key: ${{ secrets.SSH_KEY }}
          script: |
            # 1. Navigate to the project directory
            cd ~/portfolio-tui

            # 2. Pull the latest code (using the credential store you set up)
            git pull origin main

            # 3. Rebuild the image with the refactored ui.go
            docker build -t portfolio-tui .

            # 4. Remove the old container
            docker rm -f portfolio-app || true

            # 5. Launch the hardened, color-forced container
            docker run -d \
              --name portfolio-app \
              -p 22:2222 \
              --restart unless-stopped \
              --read-only \
              --cap-drop=ALL \
              --security-opt no-new-privileges:true \
              -e TERM=xterm-256color \
              -e COLORTERM=truecolor \
              -v $(pwd)/.ssh:/home/vaishak/.ssh:ro \
              portfolio-tui