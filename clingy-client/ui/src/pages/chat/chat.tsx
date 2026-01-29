import { useKeyboard } from '@opentui/react';
import { useState } from 'react';
import { TextAttributes } from '@opentui/core';
import { useChat } from '#/context/chat';
import MessageBox from './message';
import ArrowFocusText from '#/components/arrow-focus';
import Modal from '#/components/modal';
import Contacts from '../contacts';
import Config from '../config';

enum Focus {
  ChatBox,
  Contact,
  Config,
}

export default function Chat() {
  const { chatUser, chatMap, sendChatMessage } = useChat();

  const [message, setMessage] = useState('');
  const [focus, setFocus] = useState(0);
  const [modal, setModal] = useState<'contacts' | 'config' | null>(null);

  const messages = chatMap[chatUser ?? ''] ?? [];

  const handleClose = () => setModal(null);

  useKeyboard((key) => {
    if (modal) {
      return;
    }

    switch (key.name) {
      case 'tab':
        if (focus >= 2) {
          setFocus(0);
        } else {
          setFocus(focus + 1);
        }
        break;
      case 'return':
        if (focus === Focus.Contact) setModal('contacts');
        if (focus === Focus.Config) setModal('config');
        if (focus === Focus.ChatBox && !!message && !!chatUser) {
          sendChatMessage(chatUser, message).then(() => setMessage(''));
        }
        break;
    }
  });

  return (
    <box>
      <box flexGrow={1} margin={1}>
        <box marginBottom={1} flexDirection="row">
          {!!chatUser && (
            <box border={['bottom']} borderColor="#ff00ff">
              <text>{chatUser}:uuid</text>
            </box>
          )}
          {!chatUser && <text>Select a user to chat with</text>}
        </box>
        <scrollbox>
          {!messages.length && !!chatUser && <text>No messages yet</text>}
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
            cursorColor="#ff00ff"
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
      {modal === 'contacts' && (
        <Modal title="Contacts">
          <Contacts handleClose={handleClose} />
        </Modal>
      )}
      {modal === 'config' && (
        <Modal title="Config">
          <Config handleClose={handleClose} />
        </Modal>
      )}
    </box>
  );
}
