import type { Thread, Branch, Note, EditEntry } from '../models/types.js';
import { openDb, loadAllThreads } from './database.js';

export function syncData(
  threads: Thread[],
  editMap: Map<string, EditEntry>
): Thread[] {
  if (editMap.size === 0) return loadAllThreads();

  const db = openDb();
  const now = new Date().toISOString();

  // Build lookup maps
  const threadsById = new Map<number, Thread>();
  const branchesById = new Map<number, Branch>();
  const notesById = new Map<number, Note>();

  for (const t of threads) {
    threadsById.set(t.id, t);
    for (const b of t.branches) {
      branchesById.set(b.id, b);
      for (const n of b.notes) {
        notesById.set(n.id, n);
      }
    }
  }

  const creates = { threads: [] as number[], branches: [] as number[], notes: [] as number[] };
  const updates = { threads: [] as number[], branches: [] as number[], notes: [] as number[] };
  const deletes = { threads: [] as number[], branches: [] as number[], notes: [] as number[] };

  for (const [, entry] of editMap) {
    if (entry.editType === 'none') continue;
    const { entityType, id, editType } = entry;
    if (editType === `create-${entityType}`) {
      if (entityType === 'thread') creates.threads.push(id);
      else if (entityType === 'branch') creates.branches.push(id);
      else creates.notes.push(id);
    } else if (editType === `update-${entityType}`) {
      if (entityType === 'thread') updates.threads.push(id);
      else if (entityType === 'branch') updates.branches.push(id);
      else updates.notes.push(id);
    } else if (editType === `delete-${entityType}`) {
      if (entityType === 'thread') deletes.threads.push(id);
      else if (entityType === 'branch') deletes.branches.push(id);
      else deletes.notes.push(id);
    }
  }

  const runInTransaction = db.transaction(() => {
    // 1. Create threads
    const insertThread = db.prepare(`
      INSERT OR REPLACE INTO threads (id, name, summary, last_edit, highlight, private, frequency, created_at, updated_at)
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `);
    for (const id of creates.threads) {
      const t = threadsById.get(id);
      if (!t) continue;
      insertThread.run(t.id, t.name, t.summary, t.lastEdit.toISOString(), t.highlight ? 1 : 0, t.private ? 1 : 0, t.frequency, t.createdAt.toISOString(), now);
    }

    // 2. Create branches
    const insertBranch = db.prepare(`
      INSERT OR REPLACE INTO branches (id, thread_id, name, summary, last_edit, highlight, private, frequency, created_at, updated_at)
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `);
    for (const id of creates.branches) {
      const b = branchesById.get(id);
      if (!b) continue;
      insertBranch.run(b.id, b.threadId, b.name, b.summary, b.lastEdit.toISOString(), b.highlight ? 1 : 0, b.private ? 1 : 0, b.frequency, b.createdAt.toISOString(), now);
    }

    // 3. Update threads
    const updateThread = db.prepare(`UPDATE threads SET name=?, summary=?, last_edit=?, highlight=?, private=?, frequency=?, updated_at=? WHERE id=?`);
    for (const id of updates.threads) {
      const t = threadsById.get(id);
      if (!t) continue;
      updateThread.run(t.name, t.summary, t.lastEdit.toISOString(), t.highlight ? 1 : 0, t.private ? 1 : 0, t.frequency, now, t.id);
    }

    // 4. Create notes
    const insertNote = db.prepare(`
      INSERT OR REPLACE INTO notes (id, thread_id, content, diff, last_edit, highlight, private, frequency, created_at, updated_at)
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `);
    const insertBranchNote = db.prepare(`INSERT OR IGNORE INTO branch_notes (branch_id, note_id) VALUES (?, ?)`);
    for (const id of creates.notes) {
      const n = notesById.get(id);
      if (!n) continue;
      insertNote.run(n.id, n.threadId, n.content.trim(), n.diff, n.lastEdit.toISOString(), n.highlight ? 1 : 0, n.private ? 1 : 0, n.frequency, n.createdAt.toISOString(), now);
      for (const bId of n.branchIds) {
        insertBranchNote.run(bId, n.id);
      }
    }

    // 5. Update notes
    const updateNote = db.prepare(`UPDATE notes SET content=?, diff=?, last_edit=?, highlight=?, private=?, frequency=?, updated_at=? WHERE id=?`);
    for (const id of updates.notes) {
      const n = notesById.get(id);
      if (!n) continue;
      updateNote.run(n.content.trim(), n.diff, n.lastEdit.toISOString(), n.highlight ? 1 : 0, n.private ? 1 : 0, n.frequency, now, n.id);
    }

    // 6. Update branches
    const updateBranch = db.prepare(`UPDATE branches SET name=?, summary=?, last_edit=?, highlight=?, private=?, frequency=?, updated_at=? WHERE id=?`);
    for (const id of updates.branches) {
      const b = branchesById.get(id);
      if (!b) continue;
      updateBranch.run(b.name, b.summary, b.lastEdit.toISOString(), b.highlight ? 1 : 0, b.private ? 1 : 0, b.frequency, now, b.id);
    }

    // 7. Delete notes
    const softDeleteNote = db.prepare(`UPDATE notes SET deleted_at=? WHERE id=?`);
    for (const id of deletes.notes) {
      softDeleteNote.run(now, id);
    }

    // 8. Delete branches
    const softDeleteBranch = db.prepare(`UPDATE branches SET deleted_at=? WHERE id=?`);
    for (const id of deletes.branches) {
      softDeleteBranch.run(now, id);
    }

    // 9. Delete threads
    const softDeleteThread = db.prepare(`UPDATE threads SET deleted_at=? WHERE id=?`);
    for (const id of deletes.threads) {
      softDeleteThread.run(now, id);
    }
  });

  runInTransaction();
  return loadAllThreads();
}
