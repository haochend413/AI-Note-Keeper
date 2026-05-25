import React from 'react';
import { Box, Text } from 'ink';
import type { NoteEditEntry } from '../../models/types.js';
import type { DataStore } from '../../store/dataStore.js';

interface Props {
  noteEdits: NoteEditEntry[];
  cursor: number;
  focused: boolean;
  dataStore: DataStore;
  width: number;
  height: number;
}

export function RecentPanel({ noteEdits, cursor, focused, dataStore, width, height }: Props) {
  const reversed = [...noteEdits].reverse();
  const borderColor = focused ? 'green' : 'white';

  return (
    <Box flexDirection="column" borderStyle="round" borderColor={borderColor} width={width} height={height}>
      <Text bold color={borderColor}> Recent Edits </Text>
      <Box>
        <Text bold color="yellow">{'Thread'.padEnd(20)}</Text>
        <Text bold color="yellow">{'Branch'.padEnd(20)}</Text>
        <Text bold color="yellow">{'Note'.padEnd(30)}</Text>
      </Box>
      {reversed.slice(0, height - 4).map((entry, i) => {
        const { link } = entry;
        const thread = dataStore.findThreadById(link.threadId);
        const branch = dataStore.findBranchById(link.branchId);
        const note = dataStore.findNoteById(link.noteId);
        const tName = (thread?.name || `T#${link.threadId}`).slice(0, 19).padEnd(20);
        const bName = (branch?.name || `B#${link.branchId}`).slice(0, 19).padEnd(20);
        const nContent = (note?.content || `N#${link.noteId}`).split('\n')[0]?.slice(0, 29).padEnd(30) ?? '';
        const isSelected = i === cursor;
        const bg = isSelected ? (focused ? 'blue' : 'blackBright') : undefined;
        return (
          <Box key={i}>
            <Text backgroundColor={bg}>{tName}</Text>
            <Text backgroundColor={bg}>{bName}</Text>
            <Text backgroundColor={bg}>{nContent}</Text>
          </Box>
        );
      })}
      {reversed.length === 0 && <Text color="grey">  No recent edits</Text>}
    </Box>
  );
}
