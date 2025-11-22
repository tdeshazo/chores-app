# Chore Chart

Small Flask + SQLite app to track kids' daily chores with swipeable cards.

## Features
- Daily status tracking (pending/done/skipped) with animated swipe gestures.
- All-chores view with kid tabs, plus per-kid pages at `/kid/<name>`.
- Long-press completed/skipped cards to restore to pending.
- Mobile-friendly UI (touch + mouse).

## Getting started
```bash
python -m venv .venv
source .venv/bin/activate  # On Windows: .venv\Scripts\activate
pip install flask
python app.py
```
Then open http://localhost:5000.

The database is presently created/seeded automatically on first run via `init_db()` in `app.py`. If you want a fresh start, delete `chores.db`.

## Project layout
- `app.py` — Flask app, routes, and DB init.
- `templates/index.html` — main template shared by all views.
- `static/static.css` — styling.
- `static/app.js` — swipe, restore, and kid-tab logic.
- `chores.db` — generated SQLite database (ignored by git).

## API
- `POST /api/update_status` with JSON `{ "task_id": <int>, "status": "pending"|"done"|"skipped" }` updates today's status for that task.

## Notes
- The app binds to `0.0.0.0:5000` in debug mode by default (see `app.py`).
- For production, add a proper requirements file and run via a WSGI server (gunicorn, etc.).***

## Roadmap

- [ ] Replace hard-coded seed data with configurable chore/kid management.
- [ ] Support external Postgres deployment.
- [ ] Add responsive desktop layout w/ wider cards, multi-column view, mouse-friendly interactions.
- [ ] Create history view.
