import { useKeyboard } from '@opentui/react';
import { useEffect, useState } from 'react';
import { Pages, useNavigation } from '../context/navigation';
import { useMutation, useQuery } from '../hooks/api';
import FocusTextBox from '../components/FocusTextBox';
import { TextAttributes } from '@opentui/core';
import { useChat } from '../context/chat';
import { getServerConfig, type ServerConfig } from '../api/config';

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

  const messages = chatMap[chatUser ?? ''] ?? [];

  const { data: serverConfig, /* loading: loadingServerConfig, refetch */ } = useQuery<ServerConfig>(getServerConfig);

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
    <box margin={1}>
      <box height={20} backgroundColor="#393939cc">
        <text>{chatUser ? `${chatUser}:server` : ''}</text>
        {!chatUser && <text>Select a user to chat with</text>}
        {!messages.length && <text>No messages yet</text>}
        <box padding={1}>
          {messages.map((m, i) => {
            const fromSelf = m.from === serverConfig?.username;
            return (
              <box flexDirection="row" alignItems="center" justifyContent={fromSelf ? "flex-end" : "flex-start"} key={i}>
                {!fromSelf && <text fg="red">-</text>}
                <box paddingLeft={1} paddingRight={1}>
                  <text>
                    {m.message}
                  </text>
                </box>
                {fromSelf && <text fg="blue">-</text>}
              </box>
            )
          })}
        </box>
      </box>
      <box>
        <box height={1}>
          <input
            placeholder="Type your message here..."
            focused={focus === Focus.ChatBox}
            value={message}
            onInput={setMessage}
          />
        </box>
        <box flexDirection="row">
          <box flexGrow={1}>
            <FocusTextBox text="Contact" focused={focus === Focus.Contact} />
          </box>
          <box flexGrow={1}>
            <FocusTextBox text="Config" focused={focus === Focus.Config} />
          </box>
        </box>
        <box>
          <text attributes={TextAttributes.DIM}>{process.env.API_URL}</text>
        </box>
      </box>
    </box>
  );
}
