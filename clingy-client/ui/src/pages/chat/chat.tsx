import { useKeyboard } from '@opentui/react';
import { useState } from 'react';
import { Pages, useNavigation } from '#/context/navigation';
import { TextAttributes } from '@opentui/core';
import { useChat } from '#/context/chat';
import MessageBox from './message';
import ArrowFocusText from '#/components/arrow-focus';

enum Focus {
  ChatBox,
  Contact,
  Config,
}

export default function Chat() {
  const { navigate } = useNavigation();
  const { chatUser, chatMap, sendChatMessage } = useChat();

  const [message, setMessage] = useState('');
  const [focus, setFocus] = useState(0);

  const messages = chatMap[chatUser ?? ''] ?? [
    {
      from: 'user2',
      text: 'hello',
      fromSelf: false,
    },
    {
      from: 'user1',
      text: 'hey what is up',
      fromSelf: true,
    },
    {
      from: 'user2',
      text: 'oh not much',
      fromSelf: false,
    },
    {
      from: 'user2',
      text: 'just checking styling',
      fromSelf: false,
    },
    {
      from: 'user1',
      text: 'oh i know',
      fromSelf: true,
    },
    {
      from: 'user1',
      text: 'he is never going to figure this out he has not artistic ability, just ask anyone',
      fromSelf: true,
    },
  ];

  useKeyboard((key) => {
    switch (key.name) {
      case 'tab':
        if (focus >= 2) {
          setFocus(0);
        } else {
          setFocus(focus + 1);
        }
        break;
      case 'return':
        if (focus === Focus.Contact) navigate(Pages.Contacts);
        if (focus === Focus.Config) navigate(Pages.Config);
        if (focus === Focus.ChatBox && !!message && !!chatUser) {
          sendChatMessage(chatUser, message).then(() => setMessage(''));
        }
        break;
    }
  });

  return (
    <box>
      <box flexGrow={1} margin={1}>
        <box>
          <text>{chatUser ? `${chatUser}:server` : ''}</text>
          {!chatUser && <text>Select a user to chat with</text>}
          {!messages.length && <text>No messages yet</text>}
        </box>
        <scrollbox>
          {messages.map((m, i) => {
            return (
              <box
                flexDirection="row"
                alignItems="center"
                justifyContent={m.fromSelf ? 'flex-end' : 'flex-start'}
                marginBottom={1}
                marginLeft={1}
                marginRight={1}
                key={i}
              >
                <MessageBox message={m} />
              </box>
            );
          })}
        </scrollbox>
      </box>
      <box>
        <box height={1} flexDirection="row" backgroundColor="black" paddingLeft={1} paddingRight={2}>
          <input
            placeholder="clingy"
            focused={focus === Focus.ChatBox}
            focusedBackgroundColor="black"
            value={message}
            onInput={setMessage}
            flexGrow={1}
          />
          <box flexDirection="row">
            <ArrowFocusText focused={focus === Focus.Contact} text="Contact" />
            <ArrowFocusText focused={focus === Focus.Config} text="Config"/>
          </box>
        </box>
        <box flexDirection="row" justifyContent="space-between" marginLeft={1} marginRight={1}>
          <box>
            <text attributes={TextAttributes.DIM}>API URL: {process.env.API_URL}</text>
          </box>
          <box marginRight={2}>
            <text attributes={TextAttributes.DIM}>tab → | enter - open page</text>
          </box>
        </box>
      </box>
    </box>
  );
}
