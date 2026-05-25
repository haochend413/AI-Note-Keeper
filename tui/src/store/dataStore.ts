import type { Thread, Branch, Note } from '../models/types.js';

export class DataStore {
  threads: Thread[] = [];
  branches: Branch[] = []; // active thread's branches
  notes: Note[] = [];      // active branch's notes

  activeThreadIdx = 0;
  activeBranchIdx = 0;
  activeNoteIdx = 0;

  activeThreadId = 0;
  activeBranchId = 0;
  activeNoteId = 0;

  // O(1) maps
  private threadById = new Map<number, Thread>();
  private branchById = new Map<number, Branch>();
  private noteById = new Map<number, Note>();
  private threadIdxById = new Map<number, number>();
  private branchIdxById = new Map<number, number>();
  private noteIdxById = new Map<number, number>();

  load(threads: Thread[]) {
    this.threads = threads;
    this.rebuildAll();
    if (threads.length > 0) {
      this.activeThreadId = threads[0].id;
      this.activeThreadIdx = 0;
      this.branches = threads[0].branches;
      this.rebuildBranchIdx();
      if (this.branches.length > 0) {
        this.activeBranchId = this.branches[0].id;
        this.activeBranchIdx = 0;
        this.notes = this.branches[0].notes;
        this.rebuildNoteIdx();
        if (this.notes.length > 0) this.activeNoteId = this.notes[0].id;
      } else {
        this.activeBranchId = 0; this.activeBranchIdx = 0;
        this.notes = []; this.activeNoteId = 0; this.activeNoteIdx = 0;
      }
    } else {
      this.activeThreadId = 0; this.activeThreadIdx = 0;
      this.branches = []; this.activeBranchId = 0; this.activeBranchIdx = 0;
      this.notes = []; this.activeNoteId = 0; this.activeNoteIdx = 0;
    }
  }

  private rebuildAll() {
    this.threadById.clear(); this.branchById.clear(); this.noteById.clear();
    this.threadIdxById.clear();
    for (let i = 0; i < this.threads.length; i++) {
      const t = this.threads[i];
      this.threadById.set(t.id, t);
      this.threadIdxById.set(t.id, i);
      for (const b of t.branches) {
        this.branchById.set(b.id, b);
        for (const n of b.notes) this.noteById.set(n.id, n);
      }
    }
  }

  private rebuildBranchIdx() {
    this.branchIdxById.clear();
    for (let i = 0; i < this.branches.length; i++) this.branchIdxById.set(this.branches[i].id, i);
  }

  private rebuildNoteIdx() {
    this.noteIdxById.clear();
    for (let i = 0; i < this.notes.length; i++) this.noteIdxById.set(this.notes[i].id, i);
  }

  getActiveThread(): Thread | undefined { return this.threadById.get(this.activeThreadId); }
  getActiveBranch(): Branch | undefined { return this.branchById.get(this.activeBranchId); }
  getActiveNote(): Note | undefined { return this.noteById.get(this.activeNoteId); }

  findThreadById(id: number): Thread | undefined { return this.threadById.get(id); }
  findBranchById(id: number): Branch | undefined { return this.branchById.get(id); }
  findNoteById(id: number): Note | undefined { return this.noteById.get(id); }

  switchThread(idx: number) {
    if (idx < 0 || idx >= this.threads.length) return;
    const t = this.threads[idx];
    this.activeThreadIdx = idx;
    this.activeThreadId = t.id;
    this.branches = t.branches;
    this.rebuildBranchIdx();
    this.activeBranchIdx = 0;
    if (this.branches.length > 0) {
      this.activeBranchId = this.branches[0].id;
      this.notes = this.branches[0].notes;
    } else {
      this.activeBranchId = 0; this.notes = [];
    }
    this.rebuildNoteIdx();
    this.activeNoteIdx = 0;
    this.activeNoteId = this.notes.length > 0 ? this.notes[0].id : 0;
  }

  switchBranch(idx: number) {
    if (idx < 0 || idx >= this.branches.length) return;
    const b = this.branches[idx];
    this.activeBranchIdx = idx;
    this.activeBranchId = b.id;
    this.notes = b.notes;
    this.rebuildNoteIdx();
    this.activeNoteIdx = 0;
    this.activeNoteId = this.notes.length > 0 ? this.notes[0].id : 0;
  }

  switchNote(idx: number) {
    if (idx < 0 || idx >= this.notes.length) return;
    this.activeNoteIdx = idx;
    this.activeNoteId = this.notes[idx].id;
  }

  switchThreadById(id: number): boolean {
    const idx = this.threadIdxById.get(id);
    if (idx === undefined) return false;
    this.switchThread(idx); return true;
  }

  switchBranchById(id: number): boolean {
    const idx = this.branchIdxById.get(id);
    if (idx === undefined) return false;
    this.switchBranch(idx); return true;
  }

  switchNoteById(id: number): boolean {
    const idx = this.noteIdxById.get(id);
    if (idx === undefined) return false;
    this.switchNote(idx); return true;
  }

  addThread(t: Thread) {
    this.threads.push(t);
    this.threadById.set(t.id, t);
    this.threadIdxById.set(t.id, this.threads.length - 1);
  }

  addBranch(b: Branch) {
    const t = this.getActiveThread();
    if (!t) return;
    t.branches.push(b);
    this.branches = t.branches;
    this.branchById.set(b.id, b);
    this.rebuildBranchIdx();
  }

  addNote(n: Note) {
    const b = this.getActiveBranch();
    if (!b) return;
    b.notes.push(n);
    this.notes = b.notes;
    this.noteById.set(n.id, n);
    this.rebuildNoteIdx();
  }

  removeThread(idx: number) {
    if (idx < 0 || idx >= this.threads.length) return;
    this.threads.splice(idx, 1);
    this.rebuildAll();
    if (this.threads.length === 0) {
      this.activeThreadId = 0; this.activeThreadIdx = 0;
      this.branches = []; this.activeBranchId = 0;
      this.notes = []; this.activeNoteId = 0;
      return;
    }
    const newIdx = Math.min(idx, this.threads.length - 1);
    this.switchThread(newIdx);
  }

  removeBranch(idx: number) {
    const t = this.getActiveThread();
    if (!t || idx < 0 || idx >= this.branches.length) return;
    this.branches.splice(idx, 1);
    t.branches = this.branches;
    this.rebuildBranchIdx();
    if (this.branches.length === 0) {
      this.activeBranchId = 0; this.activeBranchIdx = 0;
      this.notes = []; this.activeNoteId = 0; return;
    }
    const newIdx = Math.min(idx, this.branches.length - 1);
    this.switchBranch(newIdx);
  }

  removeNote(idx: number) {
    const b = this.getActiveBranch();
    if (!b || idx < 0 || idx >= this.notes.length) return;
    this.notes.splice(idx, 1);
    b.notes = this.notes;
    this.rebuildNoteIdx();
    if (this.notes.length === 0) { this.activeNoteId = 0; this.activeNoteIdx = 0; return; }
    const newIdx = Math.min(idx, this.notes.length - 1);
    this.switchNote(newIdx);
  }

  refreshByIds(threads: Thread[], threadId: number, branchId: number, noteId: number) {
    this.threads = threads;
    this.rebuildAll();
    if (!this.switchThreadById(threadId)) this.switchThread(0);
    if (branchId) this.switchBranchById(branchId);
    if (noteId) this.switchNoteById(noteId);
  }
}
