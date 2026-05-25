import React, { useState, useEffect } from 'react';
import { Box, Text, useInput } from 'ink';

interface Props {
  value: string;
  onChange: (val: string) => void;
  active: boolean;
  width: number;
  height: number;
  placeholder?: string;
}

export function Editor({ value, onChange, active, width, height, placeholder }: Props) {
  const [lines, setLines] = useState<string[]>(() => value.split('\n'));
  const [cursorLine, setCursorLine] = useState(0);
  const [cursorCol, setCursorCol] = useState(0);
  const [scrollTop, setScrollTop] = useState(0);

  // Sync external value → lines (when focus switches in)
  useEffect(() => {
    const newLines = value.split('\n');
    setLines(newLines);
    setCursorLine(prev => Math.min(prev, newLines.length - 1));
  }, [value]);

  const commitChange = (newLines: string[]) => {
    setLines(newLines);
    onChange(newLines.join('\n'));
  };

  useInput((input, key) => {
    if (!active) return;

    if (key.return) {
      const newLines = [...lines];
      const before = newLines[cursorLine]?.slice(0, cursorCol) ?? '';
      const after = newLines[cursorLine]?.slice(cursorCol) ?? '';
      newLines.splice(cursorLine, 1, before, after);
      setCursorLine(cl => cl + 1);
      setCursorCol(0);
      commitChange(newLines);
      return;
    }

    if (key.backspace || key.delete) {
      const newLines = [...lines];
      if (cursorCol > 0) {
        const line = newLines[cursorLine] ?? '';
        newLines[cursorLine] = line.slice(0, cursorCol - 1) + line.slice(cursorCol);
        setCursorCol(cc => cc - 1);
      } else if (cursorLine > 0) {
        const prevLen = newLines[cursorLine - 1]?.length ?? 0;
        newLines[cursorLine - 1] = (newLines[cursorLine - 1] ?? '') + (newLines[cursorLine] ?? '');
        newLines.splice(cursorLine, 1);
        setCursorLine(cl => cl - 1);
        setCursorCol(prevLen);
      }
      commitChange(newLines);
      return;
    }

    if (key.upArrow) {
      setCursorLine(cl => {
        if (cl > 0) {
          const newCl = cl - 1;
          setCursorCol(cc => Math.min(cc, lines[newCl]?.length ?? 0));
          return newCl;
        }
        return cl;
      });
      return;
    }

    if (key.downArrow) {
      setCursorLine(cl => {
        if (cl < lines.length - 1) {
          const newCl = cl + 1;
          setCursorCol(cc => Math.min(cc, lines[newCl]?.length ?? 0));
          return newCl;
        }
        return cl;
      });
      return;
    }

    if (key.leftArrow) {
      if (cursorCol > 0) {
        setCursorCol(cc => cc - 1);
      } else if (cursorLine > 0) {
        const prevLineLen = lines[cursorLine - 1]?.length ?? 0;
        setCursorLine(cl => cl - 1);
        setCursorCol(prevLineLen);
      }
      return;
    }

    if (key.rightArrow) {
      const lineLen = lines[cursorLine]?.length ?? 0;
      if (cursorCol < lineLen) {
        setCursorCol(cc => cc + 1);
      } else if (cursorLine < lines.length - 1) {
        setCursorLine(cl => cl + 1);
        setCursorCol(0);
      }
      return;
    }

    // ctrl+a: line start, ctrl+e: line end
    if (key.ctrl && input === 'a') { setCursorCol(0); return; }
    if (key.ctrl && input === 'e') { setCursorCol(lines[cursorLine]?.length ?? 0); return; }
    if (key.ctrl) return; // ignore other ctrl combos

    if (input && !key.meta) {
      const newLines = [...lines];
      const line = newLines[cursorLine] ?? '';
      newLines[cursorLine] = line.slice(0, cursorCol) + input + line.slice(cursorCol);
      setCursorCol(cc => cc + input.length);
      commitChange(newLines);
    }
  }, { isActive: active });

  // Scroll to keep cursor visible
  const visibleLines = Math.max(1, height - 2);
  let newScrollTop = scrollTop;
  if (cursorLine < scrollTop) {
    newScrollTop = cursorLine;
  } else if (cursorLine >= scrollTop + visibleLines) {
    newScrollTop = cursorLine - visibleLines + 1;
  }

  const displayLines = lines.slice(newScrollTop, newScrollTop + visibleLines);

  const borderColor = active ? 'cyan' : 'white';

  // Update scrollTop after render if needed
  useEffect(() => {
    if (newScrollTop !== scrollTop) setScrollTop(newScrollTop);
  }, [newScrollTop, scrollTop]);

  return (
    <Box flexDirection="column" borderStyle="round" borderColor={borderColor} width={width} height={height}>
      {displayLines.length === 0 && !value && (
        <Text color="grey">{placeholder ?? 'Type here...'}</Text>
      )}
      {displayLines.map((line, i) => {
        const absLine = newScrollTop + i;
        const isCursorLine = absLine === cursorLine && active;
        if (!isCursorLine) {
          return <Text key={i}>{line || ' '}</Text>;
        }
        const before = line.slice(0, cursorCol);
        const cursorChar = line[cursorCol] ?? ' ';
        const after = line.slice(cursorCol + 1);
        return (
          <Text key={i}>
            {before}
            <Text backgroundColor="white" color="black">{cursorChar}</Text>
            {after}
          </Text>
        );
      })}
    </Box>
  );
}
