import type { Message } from '#/context/chat';
import { TextAttributes } from '@opentui/core';

interface MessageBoxProps {
  message: Message;
}

const MessageBox = ({ message }: MessageBoxProps) => {
  const { text, from, fromSelf } = message;

  return (
    <box
      flexDirection="row"
      backgroundColor="#111111"
      minWidth={20}
      maxWidth="60%"
      paddingLeft={fromSelf ? 1 : 0}
      paddingRight={fromSelf ? 0 : 1}
    >
      {!fromSelf && (
        <box
          backgroundColor="red"
          marginRight={1}
          width={1}
        />
      )}
      <box flexGrow={1}>
        <box alignSelf="flex-start">
          <text attributes={TextAttributes.DIM}>
            {from}
          </text>
        </box>
        <box>
          <text>
            {text}
          </text>
        </box>
      </box>
      {fromSelf && (
        <box
          backgroundColor="blue"
          marginLeft={1}
          width={1}
        />
      )}
    </box>
  );
};

export default MessageBox;
