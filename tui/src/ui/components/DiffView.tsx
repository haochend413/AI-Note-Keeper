import React from 'react';
import { Box, Text } from 'ink';

interface Props {
  html: string; // diff HTML from diff-match-patch
  width: number;
  height: number;
  focused: boolean;
}

// Strip HTML tags and convert to plain text with some color hints
function parseDiffHtml(html: string): Array<{ text: string; type: 'equal' | 'insert' | 'delete' }> {
  const result: Array<{ text: string; type: 'equal' | 'insert' | 'delete' }> = [];
  const parts = html.split(/(<[^>]+>)/);
  let currentType: 'equal' | 'insert' | 'delete' = 'equal';
  for (const part of parts) {
    if (part.startsWith('<ins')) { currentType = 'insert'; }
    else if (part.startsWith('<del')) { currentType = 'delete'; }
    else if (part === '</ins>' || part === '</del>') { currentType = 'equal'; }
    else if (!part.startsWith('<')) {
      const text = part
        .replace(/&amp;/g, '&')
        .replace(/&lt;/g, '<')
        .replace(/&gt;/g, '>')
        .replace(/&para;/g, '¶')
        .replace(/&nbsp;/g, ' ');
      if (text) result.push({ text, type: currentType });
    }
  }
  return result;
}

export function DiffView({ html, width, height, focused }: Props) {
  const borderColor = focused ? 'cyan' : 'white';
  const parts = html ? parseDiffHtml(html) : [];

  return (
    <Box flexDirection="column" borderStyle="round" borderColor={borderColor} width={width} height={height}>
      <Text bold color={borderColor}> Diff </Text>
      {parts.length === 0 && <Text color="grey">  No diff available</Text>}
      <Box flexDirection="column">
        {parts.map((p, i) => {
          const color = p.type === 'insert' ? 'green' : p.type === 'delete' ? 'red' : undefined;
          const prefix = p.type === 'insert' ? '+' : p.type === 'delete' ? '-' : ' ';
          return (
            <Text key={i} color={color}>{prefix}{p.text}</Text>
          );
        })}
      </Box>
    </Box>
  );
}
