<!-- @format -->

# 🚀 SynergySphere464 – Team Collaboration Platform

SynergySphere464 is a fast, responsive team collaboration MVP built with **Go (Fiber)**, **PostgreSQL**, **Tailwind CSS**, and **server-side HTML templates**. Designed for hackathons, it supports modern, mobile-friendly views with modular templates and easy deployment using Docker.

---

## 📸 Screenshots

> 🖼️ Includes:
>
> - Clean signup and login forms
> - Reusable header/footer components
> - Responsive Tailwind-based layout

---

## 📦 Tech Stack

- **Backend**: Go (Fiber framework)
- **Frontend**: Tailwind CSS + Go `html/template`
- **Database**: PostgreSQL
- **Templating**: Server-side partials for layout, header, footer
- **Runtime**: Docker + Docker Compose
- **Helpers**: Alpine.js or jQuery support if needed

---

## ✨ Features

- ✅ User Signup + Login (HTML forms)
- ✅ Reusable header & footer partials
- ✅ Static asset support via `/static`
- ✅ Tailwind CSS via CDN
- ✅ PostgreSQL connection using environment variables
- ✅ Fully Dockerized setup

---

## 🗂 Folder Structure

.
├── main.go
├── go.mod / go.sum
├── templates/
│ ├── layout.html
│ ├── signup.html
│ ├── login.html
│ └── partials/
│ ├── header.html
│ └── footer.html
├── static/
│ └── images/logo.png
├── Dockerfile
├── docker-compose.yml
└── db/init.sql

```bash
git clone https://github.com/yourname/synergysphere464.git
cd synergysphere464

docker compose up --build
```
