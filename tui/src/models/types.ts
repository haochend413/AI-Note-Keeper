export interface Thread {
  id: number;
  name: string;
  summary: string;
  lastEdit: Date;
  highlight: boolean;
  private: boolean;
  frequency: number;
  createdAt: Date;
  updatedAt: Date;
  branches: Branch[];
}

export interface Branch {
  id: number;
  threadId: number;
  name: string;
  summary: string;
  lastEdit: Date;
  highlight: boolean;
  private: boolean;
  frequency: number;
  createdAt: Date;
  updatedAt: Date;
  notes: Note[];
}

export interface Note {
  id: number;
  threadId: number;
  content: string;
  diff: string;
  lastEdit: Date;
  highlight: boolean;
  private: boolean;
  frequency: number;
  createdAt: Date;
  updatedAt: Date;
  branchIds: number[]; // many-to-many
}

export interface Superlink {
  threadId: number;
  branchId: number;
  noteId: number;
}

export type EditType =
  | 'create-note' | 'update-note' | 'delete-note'
  | 'create-thread' | 'update-thread' | 'delete-thread'
  | 'create-branch' | 'update-branch' | 'delete-branch'
  | 'none';

export interface EditEntry {
  entityType: 'note' | 'branch' | 'thread';
  id: number;
  editType: EditType;
}

export interface NoteEditEntry {
  link: Superlink;
}

export type FocusState =
  | 'threads' | 'branches' | 'notes'
  | 'edit' | 'changelog' | 'recent' | 'diff';

export type ViewMode = 'application' | 'quit-confirm' | 'help';

export interface AppState {
  threadCursor: number; // saved cursor positions
  branchCursor: number;
  noteCursor: number;
}
