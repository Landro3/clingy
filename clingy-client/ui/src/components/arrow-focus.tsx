import type { PropsWithChildren } from 'react';

interface ArrowFocusTextProps extends PropsWithChildren {
  text: string;
  focused: boolean;
}

const highlightColor = '#ff00ff';

export default function ArrowFocusText({ text, focused }: ArrowFocusTextProps) {

  return (
    <box flexDirection="row" alignItems="center" marginRight={1}>
      <text fg={highlightColor}>
        {focused ? <b>&gt;</b> : <>&nbsp;</>}
      </text>
      <text fg={focused ? highlightColor : '#ffffff'}>
        {text}
      </text>
    </box>
  );
}
