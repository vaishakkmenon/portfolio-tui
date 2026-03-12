# 🛰️ Vaishak Portfolio TUI

[![Deploy Status](https://github.com/vaishakkmenon/portfolio-tui/actions/workflows/deploy.yml/badge.svg)](https://github.com/vaishakkmenon/portfolio-tui/actions)
![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)
![Docker](https://img.shields.io/badge/Docker-Hardened-2496ED?style=flat&logo=docker)

A high-performance, terminal-based portfolio built in **Go** using the **Bubble Tea** framework. This project serves as a live demonstration of full-stack engineering, infrastructure automation, and secure system design.

---

## 🚀 Live Demo

You can access the live portfolio directly from your terminal — no installation required:

```bash
ssh 154.53.34.147
```

> Supports 24-bit TrueColor. Best viewed in **Windows Terminal**, **iTerm2**, or **Alacritty**.

---

## 🛠️ Technical Architecture

### Core Stack

| Component        | Technology              |
| ---------------- | ----------------------- |
| **Language**      | Go (Golang)            |
| **TUI Framework** | Bubble Tea & Lip Gloss |
| **SSH Server**    | Wish                   |
| **Infrastructure**| Docker & GitHub Actions|

### Key Features

- **Session-Aware Rendering** — Custom rendering logic that detects SSH terminal capabilities to provide the best visual experience, including TrueColor support.
- **Dynamic UI Themes** — Five built-in themes (`NEON`, `FORGE`, `TERMINAL`, and more) toggled via hotkeys.
- **Animated "Chip" Grid** — A procedurally animated circuit board visualization that simulates data flow.
- **Responsive Layout** — Standardized height locking to prevent layout jitter across different terminal sizes.

---

## 🛡️ DevOps & Security

As a **CKA**, I've designed the deployment pipeline to follow industry best practices for security and automation.

### Automated CI/CD

Pushing to `main` triggers a GitHub Action that builds and deploys the updated TUI to my VPS in Dillon, MT.

### Container Hardening

- **Non-Root User** — The application runs as a restricted user inside the container.
- **Read-Only Filesystem** — The container is mounted with `--read-only` to prevent unauthorized modifications.
- **Capability Dropping** — All Linux capabilities are dropped (`--cap-drop=ALL`) to minimize the attack surface.

### Network Segregation

Host management is rerouted to Port `2022`, while public SSH traffic on Port `22` is piped directly into the isolated Docker container.

---

## 📂 Project Structure

```
.
├── ui.go          # Main UI state machine, theme definitions, and rendering logic
├── chip.go        # Procedural animation engine for the circuit grid
├── ssh.go         # Wish SSH server configuration and middleware
├── Dockerfile     # Multi-stage build for a minimal and secure final image
└── .github/
    └── workflows/
        └── deploy.yml  # CI/CD pipeline definition
```

---

## 👤 About Me

I am a Software Engineer currently functioning as a Business Analyst within the Speridian Technologies "New Wave" program. My focus lies at the intersection of high-performance backend systems (Go/Rust) and Cloud Native infrastructure.

- **Certifications:** Certified Kubernetes Administrator (CKA), AWS Certified AI Practitioner
- **Interests:** Chess engine development ([Vantage](https://github.com/vaishakkmenon)), Machine Learning, and PC hardware

---

## 📄 License

This project is open source. See the [LICENSE](LICENSE) file for details.