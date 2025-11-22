from flask import Flask, render_template, request, jsonify
import sqlite3
import datetime

app = Flask(__name__)
DB_PATH = "chores.db"


def get_db():
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row
    return conn


def init_db():
    conn = get_db()
    cur = conn.cursor()

    # Tasks (static-ish)
    cur.execute("""
        CREATE TABLE IF NOT EXISTS tasks (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            kid TEXT NOT NULL,
            title TEXT NOT NULL,
            sort_order INTEGER NOT NULL DEFAULT 0
        )
    """)

    # Log of status per day
    cur.execute("""
        CREATE TABLE IF NOT EXISTS task_log (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            task_id INTEGER NOT NULL,
            date TEXT NOT NULL,
            status TEXT NOT NULL,
            UNIQUE(task_id, date),
            FOREIGN KEY(task_id) REFERENCES tasks(id)
        )
    """)

    # Seed some example chores if empty
    cur.execute("SELECT COUNT(*) FROM tasks")
    if cur.fetchone()[0] == 0:
        sample_tasks = [
            ("Griffin", "Make bed", 1),
            ("Griffin", "Brush teeth", 2),
            ("Griffin", "Feed the cat", 3),
            ("Garreth", "Put toys away", 1),
            ("Garreth", "Set the table", 2),
        ]
        cur.executemany(
            "INSERT INTO tasks (kid, title, sort_order) VALUES (?, ?, ?)",
            sample_tasks,
        )

    conn.commit()
    conn.close()


def fetch_tasks(for_kid=None):
    """Return task rows for today, optionally filtered to a specific kid."""
    conn = get_db()
    cur = conn.cursor()
    today = datetime.date.today().isoformat()

    query = """
        SELECT
            t.id,
            t.kid,
            t.title,
            COALESCE(l.status, 'pending') AS status
        FROM tasks t
        LEFT JOIN task_log l
            ON l.task_id = t.id AND l.date = ?
    """
    params = [today]

    if for_kid:
        query += " WHERE t.kid = ?"
        params.append(for_kid)

    query += " ORDER BY t.kid, t.sort_order"

    cur.execute(query, params)
    tasks = [dict(row) for row in cur.fetchall()]
    conn.close()
    return tasks, today


@app.route("/")
def index():
    tasks, today = fetch_tasks()
    kid_list = sorted({task["kid"] for task in tasks})
    return render_template(
        "index.html",
        tasks=tasks,
        today=today,
        kid_name=None,
        kid_list=kid_list,
    )


@app.route("/kid/<kid_name>")
def kid_page(kid_name):
    tasks, today = fetch_tasks(for_kid=kid_name)
    return render_template(
        "index.html",
        tasks=tasks,
        today=today,
        kid_name=kid_name,
        kid_list=[],
    )


@app.route("/api/update_status", methods=["POST"])
def update_status():
    data = request.get_json() or {}
    task_id = data.get("task_id")
    status = data.get("status")

    if status not in ("pending", "done", "skipped"):
        return jsonify({"error": "Invalid status"}), 400

    if task_id is None:
        return jsonify({"error": "Missing task_id"}), 400

    today = datetime.date.today().isoformat()
    conn = get_db()
    cur = conn.cursor()

    cur.execute("""
        INSERT INTO task_log (task_id, date, status)
        VALUES (?, ?, ?)
        ON CONFLICT(task_id, date)
        DO UPDATE SET status = excluded.status
    """, (task_id, today, status))
    conn.commit()
    conn.close()

    return jsonify({"ok": True})


if __name__ == "__main__":
    init_db()
    app.run(debug=True, host="0.0.0.0", port=5000)
