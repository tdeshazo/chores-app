# Chore Chart (Go + SQLite)

Lightweight Go 1.25.1 + SQLite app to track kids' daily chores with swipeable cards.

## Features
- Daily status tracking (pending/done/skipped) with animated swipe gestures.
- All-chores view with kid tabs, plus per-kid pages at `/kid/<name>`.
- Long-press completed/skipped cards to restore to pending.
- Mobile-friendly UI (touch + mouse).

## Getting started
```bash
# Requires Go 1.25.1+
go run ./cmd/server
```
Then open http://localhost:5000.

The database is created/seeded automatically on first run. Delete `chores.db` for a clean slate.

## Project layout
- `cmd/server/main.go` — entrypoint wiring DB, templates, handlers, and static assets.
- `internal/db/` — DB open/init/seed and queries.
- `internal/handlers/` — HTML and JSON handlers.
- `internal/models/` — task + view models.
- `internal/web/` — template parsing with helpers.
- `templates/index.html` — main template shared by all views.
- `static/static.css` — styling.
- `static/app.js` — swipe, restore, and kid-tab logic.
- `chores.db` — generated SQLite database (ignored by git).

## API
- `POST /api/update_status` with JSON `{ "task_id": <int>, "status": "pending"|"done"|"skipped" }` updates today's status for that task.

## Notes
- The app binds to `0.0.0.0:5000` by default (see `cmd/server/main.go`).

## Roadmap

- [ ] Replace hard-coded seed data with configurable chore/kid management.
- [ ] Support external Postgres deployment.
- [ ] Add responsive desktop layout w/ wider cards, multi-column view, mouse-friendly interactions.
- [ ] Create history view.
