import React from 'react';
import { Box, Text } from 'ink';
import type { FocusState } from '../../models/types.js';

interface Props {
  focus: FocusState;
  synced: boolean;
  entityId: number;
  frequency: number;
  lastEdit: Date | null;
  actionStatus: string;
  width: number;
}

function timeAgo(d: Date | null): string {
  if (!d) return 'Never';
  const diff = Date.now() - d.getTime();
  if (diff < 1000) return 'just now';
  if (diff < 60000) return `${Math.floor(diff / 1000)}s ago`;
  if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`;
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`;
  return d.toLocaleDateString();
}

export function StatusBar({ focus, synced, entityId, frequency, lastEdit, actionStatus }: Props) {
  const focusLabel = ({
    threads: 'Threads', branches: 'Branches', notes: 'Notes',
    edit: 'Edit', changelog: 'Changelog', recent: 'Recent', diff: 'Diff',
  } as Record<FocusState, string>)[focus] ?? focus;
  const syncLabel = synced ? ' Synced ' : ' Unsaved ';
  const syncColor = synced ? 'green' : 'yellow';
  const time = new Date().toLocaleTimeString();

  return (
    <Box flexDirection="row">
      <Text backgroundColor="grey" color="white" bold> {focusLabel} </Text>
      <Text backgroundColor="blackBright" color="white"> #{entityId || '-'} </Text>
      <Text backgroundColor="blackBright" color="white"> {timeAgo(lastEdit)} </Text>
      <Text backgroundColor="blackBright" color="white"> {frequency} edits </Text>
      <Text backgroundColor="blackBright" color="yellow"> {actionStatus || ' '} </Text>
      <Text backgroundColor={syncColor} color="black">{syncLabel}</Text>
      <Text backgroundColor="blackBright" color="white"> {time} </Text>
    </Box>
  );
}
