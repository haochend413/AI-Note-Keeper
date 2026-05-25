import path from 'path';
import fs from 'fs';
import { createRequire } from 'module';
import { DataStore } from './dataStore.js';
import { EditStore } from './editStore.js';
import { syncData } from '../db/sync.js';
import { loadAllThreads, getNextIds, getStatePath } from '../db/database.js';
import type { Thread, Branch, Note, FocusState, ViewMode, Superlink } from '../models/types.js';

// Load diff-match-patch synchronously via CJS interop
const _require = createRequire(import.meta.url);
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const DiffMatchPatchLib: any = _require('diff-match-patch');
// The library exports the constructor as the module itself (CJS default export)
// eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/no-unsafe-assignment
const DmpCtor: new () => any =
  typeof DiffMatchPatchLib === 'function'
    ? DiffMatchPatchLib
    : (DiffMatchPatchLib['diff_match_patch'] ?? DiffMatchPatchLib['DiffMatchPatch'] ?? DiffMatchPatchLib);

function makeDiff(oldText: string, newText: string): string {
  try {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-call
    const dmp = new DmpCtor();
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-call, @typescript-eslint/no-unsafe-member-access
    const diffs = dmp.diff_main(oldText, newText);
    // eslint-disable-next-line @typescript-eslint/no-unsafe-call, @typescript-eslint/no-unsafe-member-access
    dmp.diff_cleanupSemantic(diffs);
    // eslint-disable-next-line @typescript-eslint/no-unsafe-return, @typescript-eslint/no-unsafe-call, @typescript-eslint/no-unsafe-member-access
    return dmp.diff_prettyHtml(diffs);
  } catch {
    return '';
  }
}

export class AppStore {
  data: DataStore;
  edits: EditStore;
  synced = true;
  focus: FocusState = 'threads';
  previousFocus: FocusState = 'threads';
  viewMode: ViewMode = 'application';
  actionStatus = '';

  nextThreadId = 1;
  nextBranchId = 1;
  nextNoteId = 1;

  private onChange: () => void;

  constructor(onChange: () => void) {
    this.onChange = onChange;
    this.data = new DataStore();
    this.edits = new EditStore();
  }

  private notify() { this.onChange(); }

  init() {
    const threads = loadAllThreads();
    this.data.load(threads);
    const ids = getNextIds();
    this.nextThreadId = ids.threadId;
    this.nextBranchId = ids.branchId;
    this.nextNoteId = ids.noteId;
  }

  setStatus(msg: string) { this.actionStatus = msg; this.notify(); }
  setFocus(f: FocusState) { this.focus = f; this.notify(); }
  setViewMode(m: ViewMode) { this.viewMode = m; this.notify(); }

  // --- Thread ops ---
  createThread() {
    const now = new Date();
    const t: Thread = {
      id: this.nextThreadId++,
      name: '', summary: '', lastEdit: now, highlight: false, private: false,
      frequency: 0, createdAt: now, updatedAt: now, branches: []
    };
    this.edits.addEdit('thread', t.id, 'create-thread');
    this.data.addThread(t);
    this.data.switchThread(this.data.threads.length - 1);
    this.synced = false;
    this.notify();
  }

  deleteThread() {
    const t = this.data.getActiveThread();
    if (!t) return;
    // clean up child edits
    for (const b of t.branches) {
      const be = this.edits.getEdit('branch', b.id);
      if (be?.editType === 'create-branch') this.edits.removeEdit('branch', b.id);
      for (const n of b.notes) {
        const ne = this.edits.getEdit('note', n.id);
        if (ne?.editType === 'create-note') this.edits.removeEdit('note', n.id);
      }
    }
    const te = this.edits.getEdit('thread', t.id);
    if (te?.editType === 'create-thread') {
      this.edits.removeEdit('thread', t.id);
    } else {
      this.edits.addEdit('thread', t.id, 'delete-thread');
    }
    this.data.removeThread(this.data.activeThreadIdx);
    this.synced = false;
    this.notify();
  }

  setThreadSummary(summary: string) {
    const t = this.data.getActiveThread();
    if (!t || t.summary === summary) return;
    t.summary = summary;
    t.name = summary.split('\n')[0] ?? '';
    t.frequency++;
    t.lastEdit = new Date();
    this.edits.addEdit('thread', t.id, 'update-thread');
    this.synced = false;
    this.notify();
  }

  toggleThreadHighlight() {
    const t = this.data.getActiveThread();
    if (!t) return;
    t.highlight = !t.highlight;
    this.edits.addEdit('thread', t.id, 'update-thread');
    this.synced = false;
    this.notify();
  }

  toggleThreadPrivate() {
    const t = this.data.getActiveThread();
    if (!t) return;
    t.private = !t.private;
    this.edits.addEdit('thread', t.id, 'update-thread');
    this.synced = false;
    this.notify();
  }

  // --- Branch ops ---
  createBranch() {
    const t = this.data.getActiveThread();
    if (!t) return;
    const now = new Date();
    const b: Branch = {
      id: this.nextBranchId++,
      threadId: t.id,
      name: '', summary: '', lastEdit: now, highlight: false, private: false,
      frequency: 0, createdAt: now, updatedAt: now, notes: []
    };
    this.edits.addEdit('branch', b.id, 'create-branch');
    this.data.addBranch(b);
    this.data.switchBranch(this.data.branches.length - 1);
    this.synced = false;
    this.notify();
  }

  deleteBranch() {
    const b = this.data.getActiveBranch();
    if (!b) return;
    for (const n of b.notes) {
      const ne = this.edits.getEdit('note', n.id);
      if (ne?.editType === 'create-note') this.edits.removeEdit('note', n.id);
    }
    const be = this.edits.getEdit('branch', b.id);
    if (be?.editType === 'create-branch') {
      this.edits.removeEdit('branch', b.id);
    } else {
      this.edits.addEdit('branch', b.id, 'delete-branch');
    }
    this.data.removeBranch(this.data.activeBranchIdx);
    this.synced = false;
    this.notify();
  }

  setBranchSummary(summary: string) {
    const b = this.data.getActiveBranch();
    if (!b || b.summary === summary) return;
    b.summary = summary;
    b.name = summary.split('\n')[0] ?? '';
    b.frequency++;
    b.lastEdit = new Date();
    this.edits.addEdit('branch', b.id, 'update-branch');
    this.synced = false;
    this.notify();
  }

  toggleBranchHighlight() {
    const b = this.data.getActiveBranch();
    if (!b) return;
    b.highlight = !b.highlight;
    this.edits.addEdit('branch', b.id, 'update-branch');
    this.synced = false;
    this.notify();
  }

  toggleBranchPrivate() {
    const b = this.data.getActiveBranch();
    if (!b) return;
    b.private = !b.private;
    this.edits.addEdit('branch', b.id, 'update-branch');
    this.synced = false;
    this.notify();
  }

  // --- Note ops ---
  createNote() {
    const t = this.data.getActiveThread();
    const b = this.data.getActiveBranch();
    if (!t || !b) return;
    const now = new Date();
    const n: Note = {
      id: this.nextNoteId++,
      threadId: t.id,
      content: '', diff: '', lastEdit: now, highlight: false, private: false,
      frequency: 0, createdAt: now, updatedAt: now,
      branchIds: [b.id],
    };
    this.edits.addEdit('note', n.id, 'create-note');
    // mark branch for update if it exists in DB
    const be = this.edits.getEdit('branch', b.id);
    if (!be || be.editType !== 'create-branch') {
      this.edits.addEdit('branch', b.id, 'update-branch');
    }
    this.data.addNote(n);
    this.data.switchNote(this.data.notes.length - 1);
    this.synced = false;
    this.notify();
  }

  deleteNote() {
    const n = this.data.getActiveNote();
    const b = this.data.getActiveBranch();
    if (!n) return;
    const ne = this.edits.getEdit('note', n.id);
    if (ne?.editType === 'create-note') {
      this.edits.removeEdit('note', n.id);
    } else {
      this.edits.addEdit('note', n.id, 'delete-note');
    }
    if (b) {
      const be = this.edits.getEdit('branch', b.id);
      if (!be || be.editType !== 'create-branch') {
        this.edits.addEdit('branch', b.id, 'update-branch');
      }
    }
    this.data.removeNote(this.data.activeNoteIdx);
    this.synced = false;
    this.notify();
  }

  setNoteContent(content: string) {
    const n = this.data.getActiveNote();
    const t = this.data.getActiveThread();
    const b = this.data.getActiveBranch();
    if (!n || n.content === content) return;

    n.diff = makeDiff(n.content, content);
    n.content = content;
    n.frequency++;
    n.lastEdit = new Date();

    const link: Superlink = { threadId: t?.id ?? 0, branchId: b?.id ?? 0, noteId: n.id };
    this.edits.addEdit('note', n.id, 'update-note', link);

    if (t) { t.frequency++; t.lastEdit = new Date(); }
    if (b) { b.frequency++; b.lastEdit = new Date(); }

    this.synced = false;
    this.notify();
  }

  toggleNoteHighlight() {
    const n = this.data.getActiveNote();
    if (!n) return;
    n.highlight = !n.highlight;
    this.edits.addEdit('note', n.id, 'update-note');
    this.synced = false;
    this.notify();
  }

  toggleNotePrivate() {
    const n = this.data.getActiveNote();
    if (!n) return;
    n.private = !n.private;
    this.edits.addEdit('note', n.id, 'update-note');
    this.synced = false;
    this.notify();
  }

  // --- Sync ---
  syncWithDb() {
    const threadId = this.data.activeThreadId;
    const branchId = this.data.activeBranchId;
    const noteId = this.data.activeNoteId;
    const updatedThreads = syncData(this.data.threads, this.edits.editMap);
    this.data.refreshByIds(updatedThreads, threadId, branchId, noteId);
    const ids = getNextIds();
    this.nextThreadId = ids.threadId;
    this.nextBranchId = ids.branchId;
    this.nextNoteId = ids.noteId;
    this.edits.clearOnSync();
    this.synced = true;
    this.notify();
  }

  saveState() {
    try {
      const statePath = getStatePath();
      const state = {
        threadCursor: this.data.activeThreadIdx,
        branchCursor: this.data.activeBranchIdx,
        noteCursor: this.data.activeNoteIdx,
      };
      fs.mkdirSync(path.dirname(statePath), { recursive: true });
      fs.writeFileSync(statePath + '.tmp', JSON.stringify(state, null, 2));
      fs.renameSync(statePath + '.tmp', statePath);
    } catch { /* ignore */ }
  }

  loadState(): { threadCursor: number; branchCursor: number; noteCursor: number } | null {
    try {
      const raw = fs.readFileSync(getStatePath(), 'utf-8');
      return JSON.parse(raw) as { threadCursor: number; branchCursor: number; noteCursor: number };
    } catch { return null; }
  }

  restoreState() {
    const s = this.loadState();
    if (!s) return;
    if (s.threadCursor < this.data.threads.length) {
      this.data.switchThread(s.threadCursor);
      if (s.branchCursor < this.data.branches.length) {
        this.data.switchBranch(s.branchCursor);
        if (s.noteCursor < this.data.notes.length) {
          this.data.switchNote(s.noteCursor);
        }
      }
    }
  }
}
