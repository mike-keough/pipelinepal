PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS leads (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  phone TEXT,
  email TEXT,
  source TEXT,
  kind TEXT NOT NULL DEFAULT 'lead', -- lead|buyer|seller
  status TEXT NOT NULL DEFAULT 'new', -- new|contacted|nurture|hot|cold|closed|dead
  notes TEXT,
  created_at TEXT NOT NULL,
  last_contact_at TEXT,
  next_follow_up_at TEXT
);

CREATE INDEX IF NOT EXISTS idx_leads_kind ON leads(kind);
CREATE INDEX IF NOT EXISTS idx_leads_next_follow_up ON leads(next_follow_up_at);
