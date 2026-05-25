import Database from 'better-sqlite3';
import path from 'path';
import os from 'os';
import fs from 'fs';
import type { Thread, Branch, Note } from '../models/types.js';

function getBasePath(): string {
  if (process.platform === 'darwin') {
    return path.join(os.homedir(), 'Library', 'Application Support', 'ntkpr');
  } else if (process.platform === 'win32') {
    return path.join(process.env['APPDATA'] || os.homedir(), 'ntkpr');
  } else {
    return path.join(os.homedir(), '.local', 'state', 'ntkpr');
  }
}

export function getDbPath(): string {
  return path.join(getBasePath(), 'notes_dev.db');
}

export function getStatePath(): string {
  return path.join(getBasePath(), 'state.json');
}

let db: Database.Database | null = null;

export function openDb(): Database.Database {
  if (db) return db;
  const dbPath = getDbPath();
  fs.mkdirSync(path.dirname(dbPath), { recursive: true });
  db = new Database(dbPath);
  db.pragma('journal_mode = WAL');
  db.pragma('busy_timeout = 5000');
  initSchema(db);
  return db;
}

function tableHasColumn(db: Database.Database, table: string, column: string): boolean {
  const rows = db.prepare(`PRAGMA table_info(${table})`).all() as Array<{ name: string }>;
  return rows.some(r => r.name === column);
}

function addColumnIfMissing(db: Database.Database, table: string, column: string, type: string, defaultVal: string) {
  if (!tableHasColumn(db, table, column)) {
    db.exec(`ALTER TABLE ${table} ADD COLUMN ${column} ${type} DEFAULT ${defaultVal}`);
  }
}

function initSchema(db: Database.Database) {
  // Create tables if they don't exist (compatible with Go-generated schema)
  db.exec(`
    CREATE TABLE IF NOT EXISTS threads (
      id INTEGER PRIMARY KEY,
      name TEXT,
      summary TEXT,
      highlight NUMERIC DEFAULT false,
      private NUMERIC DEFAULT false,
      frequency INTEGER NOT NULL DEFAULT 0,
      created_at DATETIME,
      updated_at DATETIME,
      deleted_at DATETIME
    );
    CREATE TABLE IF NOT EXISTS branches (
      id INTEGER PRIMARY KEY,
      thread_id INTEGER,
      name TEXT,
      summary TEXT,
      highlight NUMERIC DEFAULT false,
      private NUMERIC DEFAULT false,
      frequency INTEGER NOT NULL DEFAULT 0,
      created_at DATETIME,
      updated_at DATETIME,
      deleted_at DATETIME
    );
    CREATE TABLE IF NOT EXISTS notes (
      id INTEGER PRIMARY KEY,
      thread_id INTEGER,
      content TEXT,
      highlight NUMERIC DEFAULT false,
      private NUMERIC DEFAULT false,
      frequency INTEGER NOT NULL DEFAULT 0,
      created_at DATETIME,
      updated_at DATETIME,
      deleted_at DATETIME
    );
    CREATE TABLE IF NOT EXISTS branch_notes (
      branch_id INTEGER NOT NULL,
      note_id INTEGER NOT NULL,
      PRIMARY KEY (branch_id, note_id)
    );
  `);

  // Migrate: add columns that may be missing in older Go-created databases
  addColumnIfMissing(db, 'threads', 'last_edit', 'DATETIME', 'NULL');
  addColumnIfMissing(db, 'branches', 'last_edit', 'DATETIME', 'NULL');
  addColumnIfMissing(db, 'notes', 'last_edit', 'DATETIME', 'NULL');
  addColumnIfMissing(db, 'notes', 'diff', 'TEXT', "''");
}

export function closeDb() {
  if (db) { db.close(); db = null; }
}

// --- Read helpers ---
export function loadAllThreads(): Thread[] {
  const d = openDb();
  const threads = d.prepare(`SELECT * FROM threads WHERE deleted_at IS NULL ORDER BY created_at ASC`).all() as Record<string, unknown>[];
  return threads.map(t => {
    const branches = loadBranchesForThread(t['id'] as number);
    return dbRowToThread(t, branches);
  });
}

function loadBranchesForThread(threadId: number): Branch[] {
  const d = openDb();
  const rows = d.prepare(`SELECT * FROM branches WHERE thread_id = ? AND deleted_at IS NULL ORDER BY created_at ASC`).all(threadId) as Record<string, unknown>[];
  return rows.map(b => {
    const notes = loadNotesForBranch(b['id'] as number);
    return dbRowToBranch(b, notes);
  });
}

function loadNotesForBranch(branchId: number): Note[] {
  const d = openDb();
  const rows = d.prepare(`
    SELECT n.* FROM notes n
    JOIN branch_notes bn ON bn.note_id = n.id
    WHERE bn.branch_id = ? AND n.deleted_at IS NULL
    ORDER BY n.created_at ASC
  `).all(branchId) as Record<string, unknown>[];
  return rows.map(n => dbRowToNote(n));
}

function dbRowToThread(row: Record<string, unknown>, branches: Branch[]): Thread {
  return {
    id: row['id'] as number,
    name: row['name'] as string,
    summary: row['summary'] as string,
    lastEdit: row['last_edit'] ? new Date(row['last_edit'] as string) : new Date(row['created_at'] as string),
    highlight: !!(row['highlight'] as number),
    private: !!(row['private'] as number),
    frequency: row['frequency'] as number,
    createdAt: new Date(row['created_at'] as string),
    updatedAt: new Date(row['updated_at'] as string),
    branches,
  };
}

function dbRowToBranch(row: Record<string, unknown>, notes: Note[]): Branch {
  return {
    id: row['id'] as number,
    threadId: row['thread_id'] as number,
    name: row['name'] as string,
    summary: row['summary'] as string,
    lastEdit: row['last_edit'] ? new Date(row['last_edit'] as string) : new Date(row['created_at'] as string),
    highlight: !!(row['highlight'] as number),
    private: !!(row['private'] as number),
    frequency: row['frequency'] as number,
    createdAt: new Date(row['created_at'] as string),
    updatedAt: new Date(row['updated_at'] as string),
    notes,
  };
}

function dbRowToNote(row: Record<string, unknown>): Note {
  const d = openDb();
  const branchIds = (d.prepare(`SELECT branch_id FROM branch_notes WHERE note_id = ?`).all(row['id'] as number) as Record<string, unknown>[]).map(r => r['branch_id'] as number);
  return {
    id: row['id'] as number,
    threadId: row['thread_id'] as number,
    content: row['content'] as string,
    diff: row['diff'] as string,
    lastEdit: row['last_edit'] ? new Date(row['last_edit'] as string) : new Date(row['created_at'] as string),
    highlight: !!(row['highlight'] as number),
    private: !!(row['private'] as number),
    frequency: row['frequency'] as number,
    createdAt: new Date(row['created_at'] as string),
    updatedAt: new Date(row['updated_at'] as string),
    branchIds,
  };
}

// --- ID generation ---
export interface NextIds {
  noteId: number;
  branchId: number;
  threadId: number;
}

export function getNextIds(): NextIds {
  const d = openDb();
  const noteMax = (d.prepare(`SELECT COALESCE(MAX(id), 0) as m FROM notes`).get() as Record<string, number>)['m'];
  const branchMax = (d.prepare(`SELECT COALESCE(MAX(id), 0) as m FROM branches`).get() as Record<string, number>)['m'];
  const threadMax = (d.prepare(`SELECT COALESCE(MAX(id), 0) as m FROM threads`).get() as Record<string, number>)['m'];
  return { noteId: noteMax + 1, branchId: branchMax + 1, threadId: threadMax + 1 };
}
