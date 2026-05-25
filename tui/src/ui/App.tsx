import React, { useState, useEffect, useRef } from 'react';
import { Box, Text, useInput, useApp, useStdout } from 'ink';
import { SelectTable } from './components/SelectTable.js';
import { Editor } from './components/Editor.js';
import { StatusBar } from './components/StatusBar.js';
import { RecentPanel } from './components/RecentPanel.js';
import { DiffView } from './components/DiffView.js';
import { AppStore } from '../store/appStore.js';
import type { FocusState } from '../models/types.js';

// Helper to format date for table display
function fmtDate(d: Date): string {
  return `${String(d.getFullYear() - 2000).padStart(2, '0')}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
}

function flags(h: boolean, p: boolean): string {
  return (h ? 'H' : '') + (p ? 'P' : '');
}

export function App() {
  const [, forceUpdate] = useState(0);
  const { exit } = useApp();
  const { stdout } = useStdout();

  const storeRef = useRef<AppStore | null>(null);
  if (!storeRef.current) {
    storeRef.current = new AppStore(() => forceUpdate(n => n + 1));
    storeRef.current.init();
    storeRef.current.restoreState();
  }
  const store = storeRef.current;

  const [recentCursor, setRecentCursor] = useState(0);
  const [editingContent, setEditingContent] = useState('');

  // Clock tick every second for status bar
  useEffect(() => {
    const t = setInterval(() => forceUpdate(n => n + 1), 1000);
    return () => clearInterval(t);
  }, []);

  const W = stdout?.columns ?? 120;
  const H = stdout?.rows ?? 40;

  const tableW = Math.floor(W * 0.33);
  const editorW = W - tableW - 4;
  const tableH = Math.floor((H - 4) / 3);
  const editorH = H - 3;

  const { data, edits, focus, viewMode, synced, actionStatus } = store;

  // Build table rows
  const threadRows = data.threads.map(t => ({
    id: String(t.id),
    time: fmtDate(t.createdAt),
    name: t.name || '(untitled)',
    branches: String(t.branches.length),
    flags: flags(t.highlight, t.private),
  }));

  const branchRows = data.branches.map(b => ({
    id: String(b.id),
    time: fmtDate(b.createdAt),
    name: b.name || '(untitled)',
    notes: String(b.notes.length),
    flags: flags(b.highlight, b.private),
  }));

  const noteRows = data.notes.map(n => ({
    id: String(n.id),
    time: fmtDate(n.createdAt),
    content: (n.content || '(empty)').split('\n')[0] ?? '(empty)',
    flags: flags(n.highlight, n.private),
  }));

  const threadCols = [
    { key: 'id', title: 'ID', width: 4 },
    { key: 'time', title: 'Time', width: 14 },
    { key: 'name', title: 'Name', width: Math.max(4, tableW - 35) },
    { key: 'branches', title: '#B', width: 3 },
    { key: 'flags', title: 'F', width: 4 },
  ];

  const branchCols = [
    { key: 'id', title: 'ID', width: 4 },
    { key: 'time', title: 'Time', width: 14 },
    { key: 'name', title: 'Name', width: Math.max(4, tableW - 35) },
    { key: 'notes', title: '#N', width: 3 },
    { key: 'flags', title: 'F', width: 4 },
  ];

  const noteCols = [
    { key: 'id', title: 'ID', width: 4 },
    { key: 'time', title: 'Time', width: 14 },
    { key: 'content', title: 'Content', width: Math.max(4, tableW - 28) },
    { key: 'flags', title: 'F', width: 4 },
  ];

  const activeThread = data.getActiveThread();
  const activeBranch = data.getActiveBranch();
  const activeNote = data.getActiveNote();

  // Get status bar info based on focus
  let statusId = 0, statusFreq = 0;
  let statusLastEdit: Date | null = null;
  if (focus === 'threads' && activeThread) {
    statusId = activeThread.id; statusFreq = activeThread.frequency; statusLastEdit = activeThread.lastEdit;
  } else if (focus === 'branches' && activeBranch) {
    statusId = activeBranch.id; statusFreq = activeBranch.frequency; statusLastEdit = activeBranch.lastEdit;
  } else if ((focus === 'notes' || focus === 'edit') && activeNote) {
    statusId = activeNote.id; statusFreq = activeNote.frequency; statusLastEdit = activeNote.lastEdit;
  }

  useInput((input, key) => {
    if (viewMode === 'quit-confirm') {
      if (input === 'y') {
        store.syncWithDb();
        store.saveState();
        exit();
      } else if (input === 'n' || key.escape) {
        store.setViewMode('application');
      }
      return;
    }

    if (viewMode === 'help') {
      store.setViewMode('application');
      return;
    }

    // ctrl+c: quit confirm
    if (key.ctrl && input === 'c') {
      store.setViewMode('quit-confirm');
      return;
    }

    // ctrl+q: sync
    if (key.ctrl && input === 'q') {
      store.setStatus('Syncing...');
      store.syncWithDb();
      store.setStatus('Synced with database');
      return;
    }

    // H: help
    if (input === 'H' && focus !== 'edit') {
      store.setViewMode('help');
      return;
    }

    // Tab: cycle focus or exit edit
    if (key.tab) {
      if (focus === 'edit') {
        store.setStatus('Saving edit...');
        if (store.previousFocus === 'threads') store.setThreadSummary(editingContent);
        else if (store.previousFocus === 'branches') store.setBranchSummary(editingContent);
        else if (store.previousFocus === 'notes') store.setNoteContent(editingContent);
        store.setFocus(store.previousFocus);
      } else {
        const cycle: FocusState[] = ['threads', 'branches', 'notes'];
        const idx = cycle.indexOf(focus as FocusState);
        if (idx !== -1) store.setFocus(cycle[(idx + 1) % 3] as FocusState);
      }
      return;
    }

    // Edit mode keys
    if (focus === 'edit') {
      if (key.ctrl && input === 's') {
        store.setStatus('Saving...');
        if (store.previousFocus === 'threads') store.setThreadSummary(editingContent);
        else if (store.previousFocus === 'branches') store.setBranchSummary(editingContent);
        else if (store.previousFocus === 'notes') store.setNoteContent(editingContent);
        store.setFocus(store.previousFocus);
        store.setStatus('Saved');
      } else if (key.ctrl && input === 'x') {
        store.setFocus(store.previousFocus);
        store.setStatus('Cancelled');
      }
      return;
    }

    // Changelog mode
    if (focus === 'changelog') {
      if (key.escape) { store.setFocus('notes'); }
      return;
    }

    // Recent mode
    if (focus === 'recent') {
      if (key.upArrow || input === 'k') setRecentCursor(c => Math.max(0, c - 1));
      else if (key.downArrow || input === 'j') setRecentCursor(c => Math.min(edits.noteEditStack.length - 1, c + 1));
      else if (key.return) {
        const reversed = [...edits.noteEditStack].reverse();
        const entry = reversed[recentCursor];
        if (entry) {
          data.switchThreadById(entry.link.threadId);
          data.switchBranchById(entry.link.branchId);
          data.switchNoteById(entry.link.noteId);
          store.setFocus('notes');
          forceUpdate(n => n + 1);
        }
      } else if (input === 'D') store.setFocus('diff');
      else if (input === 'R' || key.escape) store.setFocus('notes');
      return;
    }

    // Diff mode
    if (focus === 'diff') {
      if (input === 'R' || key.escape) store.setFocus('recent');
      return;
    }

    // Table mode: threads/branches/notes
    const isThreads = focus === 'threads';
    const isBranches = focus === 'branches';
    const isNotes = focus === 'notes';

    // Navigation up/down
    if (key.upArrow || input === 'k') {
      if (isThreads && data.activeThreadIdx > 0) { data.switchThread(data.activeThreadIdx - 1); forceUpdate(n => n + 1); }
      else if (isBranches && data.activeBranchIdx > 0) { data.switchBranch(data.activeBranchIdx - 1); forceUpdate(n => n + 1); }
      else if (isNotes && data.activeNoteIdx > 0) { data.switchNote(data.activeNoteIdx - 1); forceUpdate(n => n + 1); }
      return;
    }
    if (key.downArrow || input === 'j') {
      if (isThreads && data.activeThreadIdx < data.threads.length - 1) { data.switchThread(data.activeThreadIdx + 1); forceUpdate(n => n + 1); }
      else if (isBranches && data.activeBranchIdx < data.branches.length - 1) { data.switchBranch(data.activeBranchIdx + 1); forceUpdate(n => n + 1); }
      else if (isNotes && data.activeNoteIdx < data.notes.length - 1) { data.switchNote(data.activeNoteIdx + 1); forceUpdate(n => n + 1); }
      return;
    }

    // Enter: drill down / open edit for notes
    if (key.return) {
      if (isThreads) { store.setFocus('branches'); store.setStatus('Thread selected'); }
      else if (isBranches) { store.setFocus('notes'); store.setStatus('Branch selected'); }
      else if (isNotes) {
        store.previousFocus = 'notes';
        setEditingContent(activeNote?.content ?? '');
        store.setFocus('edit');
        store.setStatus('Editing note');
      }
      return;
    }

    // Escape: go back
    if (key.escape) {
      if (isNotes) store.setFocus('branches');
      else if (isBranches) store.setFocus('threads');
      return;
    }

    // e: enter edit
    if (input === 'e') {
      if (isThreads) {
        store.previousFocus = 'threads';
        setEditingContent(activeThread?.summary ?? '');
        store.setFocus('edit');
        store.setStatus('Editing thread');
      } else if (isBranches) {
        store.previousFocus = 'branches';
        setEditingContent(activeBranch?.summary ?? '');
        store.setFocus('edit');
        store.setStatus('Editing branch');
      } else if (isNotes) {
        store.previousFocus = 'notes';
        setEditingContent(activeNote?.content ?? '');
        store.setFocus('edit');
        store.setStatus('Editing note');
      }
      return;
    }

    // n: create new
    if (input === 'n') {
      if (isThreads) { store.createThread(); store.setStatus('Thread created'); }
      else if (isBranches) { store.createBranch(); store.setStatus('Branch created'); }
      else if (isNotes) { store.createNote(); store.setStatus('Note created'); }
      return;
    }

    // ctrl+d: delete
    if (key.ctrl && input === 'd') {
      if (isThreads) { store.deleteThread(); store.setStatus('Thread deleted'); }
      else if (isBranches) { store.deleteBranch(); store.setStatus('Branch deleted'); }
      else if (isNotes) { store.deleteNote(); store.setStatus('Note deleted'); }
      return;
    }

    // ctrl+h: highlight
    if (key.ctrl && input === 'h') {
      if (isThreads) { store.toggleThreadHighlight(); store.setStatus('Highlight toggled'); }
      else if (isBranches) { store.toggleBranchHighlight(); store.setStatus('Highlight toggled'); }
      else if (isNotes) { store.toggleNoteHighlight(); store.setStatus('Highlight toggled'); }
      return;
    }

    // ctrl+p: private
    if (key.ctrl && input === 'p') {
      if (isThreads) { store.toggleThreadPrivate(); store.setStatus('Private toggled'); }
      else if (isBranches) { store.toggleBranchPrivate(); store.setStatus('Private toggled'); }
      else if (isNotes) { store.toggleNotePrivate(); store.setStatus('Private toggled'); }
      return;
    }

    // R: recent
    if (input === 'R') {
      setRecentCursor(0);
      store.setFocus('recent');
      store.setStatus('Recent edits');
      return;
    }

    // ctrl+l: changelog
    if (key.ctrl && input === 'l' && isNotes) {
      store.setFocus('changelog');
      return;
    }

    // l/h: move between table levels
    if (input === 'l') {
      if (isThreads) store.setFocus('branches');
      else if (isBranches) store.setFocus('notes');
    }
    if (input === 'h') {
      if (isNotes) store.setFocus('branches');
      else if (isBranches) store.setFocus('threads');
    }
  });

  // Quit confirm overlay
  if (viewMode === 'quit-confirm') {
    return (
      <Box flexDirection="column" padding={2}>
        <Text bold color="yellow">Quit ntkpr?</Text>
        <Text>Unsaved changes will be synced to the database.</Text>
        <Text> </Text>
        <Text>Press <Text color="green" bold>y</Text> to quit, <Text color="red" bold>n</Text> or <Text color="red" bold>Esc</Text> to cancel.</Text>
      </Box>
    );
  }

  // Help overlay
  if (viewMode === 'help') {
    return (
      <Box flexDirection="column" padding={2}>
        <Text bold color="green">ntkpr — keyboard shortcuts</Text>
        <Text> </Text>
        <Text bold>Global</Text>
        <Text>  <Text color="yellow">Ctrl+C</Text>   Quit (with confirm)</Text>
        <Text>  <Text color="yellow">Ctrl+Q</Text>   Sync to database</Text>
        <Text>  <Text color="yellow">Tab</Text>      Cycle focus Threads→Branches→Notes</Text>
        <Text>  <Text color="yellow">H</Text>        This help</Text>
        <Text> </Text>
        <Text bold>Navigation</Text>
        <Text>  <Text color="yellow">j/k / ↑↓</Text>  Move up/down in table</Text>
        <Text>  <Text color="yellow">Enter</Text>    Drill down / Enter edit</Text>
        <Text>  <Text color="yellow">Esc</Text>      Go up one level</Text>
        <Text>  <Text color="yellow">l/h</Text>      Move focus to next/prev table</Text>
        <Text> </Text>
        <Text bold>Actions</Text>
        <Text>  <Text color="yellow">n</Text>        Create new item</Text>
        <Text>  <Text color="yellow">e</Text>        Edit current item</Text>
        <Text>  <Text color="yellow">Ctrl+D</Text>   Delete current item</Text>
        <Text>  <Text color="yellow">Ctrl+H</Text>   Toggle highlight</Text>
        <Text>  <Text color="yellow">Ctrl+P</Text>   Toggle private</Text>
        <Text>  <Text color="yellow">R</Text>        Recent edits panel</Text>
        <Text> </Text>
        <Text bold>Edit mode</Text>
        <Text>  <Text color="yellow">Ctrl+S</Text>   Save and return</Text>
        <Text>  <Text color="yellow">Ctrl+X</Text>   Cancel edit</Text>
        <Text> </Text>
        <Text color="grey">Press any key to dismiss</Text>
      </Box>
    );
  }

  // Main application view
  const showRecent = focus === 'recent' || focus === 'diff';
  const recentNote = (() => {
    const reversed = [...edits.noteEditStack].reverse();
    const entry = reversed[recentCursor];
    return entry ? data.findNoteById(entry.link.noteId) : undefined;
  })();

  const changelogRows = [...store.edits.editMap.entries()].map(([, e]) => ({
    type: e.editType.split('-')[0] ?? '',
    id: String(e.id),
    entity: e.entityType,
    desc: e.editType,
  }));

  return (
    <Box flexDirection="column" width={W} height={H}>
      <Box flexDirection="row" flexGrow={1}>
        {/* Left: three stacked tables */}
        <Box flexDirection="column" width={tableW}>
          <SelectTable
            title="Threads"
            columns={threadCols}
            rows={threadRows}
            cursor={data.activeThreadIdx}
            focused={focus === 'threads'}
            height={tableH}
          />
          <SelectTable
            title="Branches"
            columns={branchCols}
            rows={branchRows}
            cursor={data.activeBranchIdx}
            focused={focus === 'branches'}
            height={tableH}
          />
          <SelectTable
            title="Notes"
            columns={noteCols}
            rows={noteRows}
            cursor={data.activeNoteIdx}
            focused={focus === 'notes'}
            height={Math.max(3, H - tableH * 2 - 4)}
          />
        </Box>

        {/* Right: editor / recent / changelog / diff */}
        <Box flexDirection="column" flexGrow={1}>
          {!showRecent && focus !== 'changelog' && (
            <Editor
              value={focus === 'edit' ? editingContent : (
                focus === 'threads' ? (activeThread?.summary ?? '') :
                focus === 'branches' ? (activeBranch?.summary ?? '') :
                activeNote?.content ?? ''
              )}
              onChange={focus === 'edit' ? setEditingContent : () => undefined}
              active={focus === 'edit'}
              width={editorW}
              height={editorH}
              placeholder={
                focus === 'threads' ? 'Thread summary (first line = name)' :
                focus === 'branches' ? 'Branch summary (first line = name)' :
                'Note content'
              }
            />
          )}
          {focus === 'changelog' && (
            <Box flexDirection="column" borderStyle="round" borderColor="cyan" width={editorW} height={editorH}>
              <Text bold color="cyan"> Changelog </Text>
              <Box>
                <Text bold color="yellow">{'Type'.padEnd(10)}</Text>
                <Text bold color="yellow">{'Entity'.padEnd(10)}</Text>
                <Text bold color="yellow">{'ID'.padEnd(6)}</Text>
              </Box>
              {changelogRows.map((r, i) => (
                <Box key={i}>
                  <Text>{r.type.padEnd(10)}</Text>
                  <Text>{r.entity.padEnd(10)}</Text>
                  <Text>{r.id.padEnd(6)}</Text>
                </Box>
              ))}
              {changelogRows.length === 0 && <Text color="grey">  No pending changes</Text>}
              <Text color="grey">  Esc: close</Text>
            </Box>
          )}
          {showRecent && (
            <Box flexDirection="row">
              <RecentPanel
                noteEdits={edits.noteEditStack}
                cursor={recentCursor}
                focused={focus === 'recent'}
                dataStore={data}
                width={Math.floor(editorW / 2)}
                height={editorH}
              />
              <DiffView
                html={recentNote?.diff ?? ''}
                width={editorW - Math.floor(editorW / 2)}
                height={editorH}
                focused={focus === 'diff'}
              />
            </Box>
          )}
        </Box>
      </Box>

      {/* Status bar */}
      <StatusBar
        focus={focus}
        synced={synced}
        entityId={statusId}
        frequency={statusFreq}
        lastEdit={statusLastEdit}
        actionStatus={actionStatus}
        width={W}
      />

      {/* Help hint */}
      {focus !== 'edit' && (
        <Box>
          <Text color="grey">Tab:cycle  j/k:move  Enter:select  e:edit  n:new  ^D:del  ^H:hl  ^P:priv  R:recent  ^Q:sync  ^C:quit  H:help</Text>
        </Box>
      )}
    </Box>
  );
}
