import React from 'react';
import { Box, Text } from 'ink';

export interface Column {
  title: string;
  width: number;
  key: string;
}

export interface Row {
  [key: string]: string;
}

interface Props {
  columns: Column[];
  rows: Row[];
  cursor: number;
  focused: boolean;
  height: number;
  title?: string;
}

function truncate(s: string, n: number): string {
  if (s.length <= n) return s.padEnd(n);
  return s.slice(0, n - 1) + '…';
}

export function SelectTable({ columns, rows, cursor, focused, height, title }: Props) {
  // Compute scroll offset to keep cursor visible
  const visibleRows = Math.max(1, height - 3); // header + border
  const scrollOffset = Math.max(0, Math.min(cursor - Math.floor(visibleRows / 2), rows.length - visibleRows));
  const visible = rows.slice(scrollOffset, scrollOffset + visibleRows);

  const borderColor = focused ? 'green' : 'white';

  return (
    <Box flexDirection="column" borderStyle="round" borderColor={borderColor} height={height}>
      {title && <Text color={borderColor} bold> {title} </Text>}
      {/* Header */}
      <Box>
        {columns.map((col, i) => (
          <Text key={i} bold color="yellow">{truncate(col.title, col.width)}{' '}</Text>
        ))}
      </Box>
      {/* Rows */}
      {visible.map((row, i) => {
        const rowIdx = scrollOffset + i;
        const isSelected = rowIdx === cursor;
        const bg = isSelected ? (focused ? 'blue' : 'blackBright') : undefined;
        return (
          <Box key={rowIdx}>
            {columns.map((col, j) => (
              <Text key={j} backgroundColor={bg} color={isSelected && focused ? 'white' : undefined}>
                {truncate(row[col.key] ?? '', col.width)}{' '}
              </Text>
            ))}
          </Box>
        );
      })}
      {visible.length === 0 && <Text color="grey">  (empty)</Text>}
    </Box>
  );
}
