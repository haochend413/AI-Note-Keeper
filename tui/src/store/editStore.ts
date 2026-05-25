import type { EditEntry, EditType, NoteEditEntry, Superlink } from '../models/types.js';

export class EditStore {
  editMap = new Map<string, EditEntry>(); // key: "type:id"
  noteEditStack: NoteEditEntry[] = []; // recent edits history

  private key(entityType: 'note' | 'branch' | 'thread', id: number): string {
    return `${entityType}:${id}`;
  }

  addEdit(entityType: 'note' | 'branch' | 'thread', id: number, editType: EditType, link?: Superlink): void {
    const k = this.key(entityType, id);
    const existing = this.editMap.get(k);

    if (!existing) {
      this.editMap.set(k, { entityType, id, editType });
      if (link && (editType.startsWith('create') || editType.startsWith('update'))) {
        this.appendNoteEdit(link);
      }
      return;
    }

    const prev = existing.editType;

    // State machine transitions
    if (editType === `update-${entityType}` as EditType) {
      if (prev === `create-${entityType}`) {
        // Stay as create - not yet in DB
        if (link) this.appendNoteEdit(link);
        return;
      }
      if (prev === `delete-${entityType}`) return; // invalid, ignore
      // prev === update: idempotent
      existing.editType = editType;
      if (link) this.appendNoteEdit(link);
    } else if (editType === `delete-${entityType}` as EditType) {
      if (prev === `create-${entityType}`) {
        existing.editType = 'none'; // discard - never synced
      } else if (prev === `update-${entityType}`) {
        existing.editType = editType; // was updated, now deleted
      }
      // already delete: ignore
    }
    // create-X when prev exists: error state, ignore gracefully
  }

  removeEdit(entityType: 'note' | 'branch' | 'thread', id: number): void {
    this.editMap.delete(this.key(entityType, id));
  }

  getEdit(entityType: 'note' | 'branch' | 'thread', id: number): EditEntry | undefined {
    return this.editMap.get(this.key(entityType, id));
  }

  private appendNoteEdit(link: Superlink): void {
    const idx = this.noteEditStack.findIndex(
      e => e.link.noteId === link.noteId && e.link.branchId === link.branchId && e.link.threadId === link.threadId
    );
    if (idx !== -1) this.noteEditStack.splice(idx, 1);
    this.noteEditStack.push({ link });
    // Keep last 50 entries
    if (this.noteEditStack.length > 50) this.noteEditStack.shift();
  }

  clearOnSync(): void {
    this.editMap.clear();
    // Preserve noteEditStack (history survives sync)
  }

  clear(): void {
    this.editMap.clear();
    this.noteEditStack = [];
  }
}
