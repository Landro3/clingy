import { useKeyboard } from '@opentui/react';
import { useEffect, useState } from 'react';
import { Pages, useNavigation } from '../context/navigation';
import { sendChatMessage as sendChatMessageApi } from '../api/chat';
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
  const { chatUser } = useChat();

  const [focus, setFocus] = useState(0);
  const [message, setMessage] = useState('');
  const [messages, setMessages] = useState<{ from: string; message: string }[]>([]);

  const { data: serverConfig, /* loading: loadingServerConfig, refetch */ } = useQuery<ServerConfig>(getServerConfig);
  const { mutate: sendChatMessage } = useMutation(sendChatMessageApi);

  useEffect(() => {
    const controller = new AbortController();
    (async () => {
      try {
        const response = await fetch(`${process.env.API_URL}/chat/stream`, {
          headers: { "Accept": "text/event-stream" },
          signal: controller.signal
        });
  
        if (!response.body) return;
  
        let buffer = "";
        for await (const chunk of response.body) {
          buffer += new TextDecoder().decode(chunk);
    
          const lines = buffer.split("\n\n");
          buffer = lines.pop() || "";
  
          for (const line of lines) {
            if (line.startsWith("data: ")) {
                const rawJson = line.replace("data: ", "").trim();
                const parsed = JSON.parse(rawJson);
                setMessages((prev) => [...prev, parsed]);
            }
          }
        }
      } catch (err) {
        if (err instanceof Error && err.name === 'AbortError') {
          return;
        }

        console.log('Probably timing out');
        console.error(err);
      }
    })();


    return () => controller.abort();
  }, []);

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
          sendChatMessage({ to: chatUser, message })
            .then(() => setMessage(''))
            .then(() => setMessages((prev) => [...prev, { from: serverConfig!.username, message }]));
        }
        break;
    }
  });

  return (
    <box>
      <box title={chatUser ? ` ${chatUser}:server ` : ''} border borderStyle="rounded" height={20}>
        {!chatUser && <text>Select a user to chat with</text>}
        {!messages.length && <text>No messages yet</text>}
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
      <box>
        <box border height={3} borderStyle="rounded">
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
