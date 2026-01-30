import { TextAttributes } from '@opentui/core';
import { useState } from 'react';
import useScrollKeys from '../hooks/useScrollKeys';
import { useKeyboard } from '@opentui/react';
import {
  type Contact,
  createContact as createContactApi,
  updateContact as updateContactApi,
  deleteContact as deleteContactApi,
  getContacts,
} from '../api/contacts';
import { useMutation, useQuery } from '../hooks/api';
import { useChat } from '../context/chat';
import ArrowFocusText from '#/components/arrow-focus';

interface ContactsProps {
  handleClose: () => void;
}

export default function Contacts({ handleClose }: ContactsProps) {
  const { setChatUser } = useChat();

  const [mode, setMode] = useState<'create' | 'update' | 'delete' | null>(null);
  const [inputFocus, setInputFocus] = useState(0);
  const [username, setUsername] = useState('');
  const [uniqueId, setUniqueId] = useState('');
  const [updatingId, setUpdatingId] = useState('');

  const { data: contacts, loading: gettingContacts, refetch } = useQuery<Contact[]>(getContacts);
  const index = useScrollKeys(contacts?.length ?? 0, !!mode);

  const { mutate: createContact, loading: creatingContact } = useMutation(createContactApi);
  const { mutate: updateContact, loading: updatingContact } = useMutation(updateContactApi);
  const { mutate: deleteContact, loading: deletingContact } = useMutation(deleteContactApi);

  useKeyboard(({ name }) => {
    if (!contacts) return;

    const clearMode = () => {
      setMode(null);
      setUsername('');
      setUniqueId('');
      setUpdatingId('');
      setInputFocus(0);
    };

    if (mode === null) {
      switch (name) {
        case 'a':
          setMode('create');
          return;
        case 'e':
          setMode('update');
          if (!contacts[index]) {
            return;
          }
          setUsername(contacts[index].username);
          setUpdatingId(contacts[index].uniqueId);
          setUniqueId(contacts[index].uniqueId);
          return;
        case 'r':
          setMode('delete');
          if (!contacts[index]) {
            return;
          }
          setUsername(contacts[index].username);
          setUniqueId(contacts[index].uniqueId);
          return;
        case 'escape':
          handleClose();
          return;
      }
    }

    if (!!mode && name === 'escape') {
      clearMode();
      return;
    }

    if ((mode === 'create' || mode === 'update') && name === 'tab') {
      if (inputFocus === 1) {
        setInputFocus(0);
      } else {
        setInputFocus(inputFocus + 1);
      }

      return;
    }

    if (!mode && contacts[index] && name === 'return') {
      setChatUser(contacts[index].username);
      handleClose();
      return;
    }

    if (mode && ['create', 'update'].includes(mode) && (!username || !uniqueId) && name === 'return') {
      clearMode();
      return;
    }

    if (mode === 'create' && username && uniqueId && name === 'return') {
      createContact({ username, uniqueId })
        .then(clearMode)
        .then(refetch);
    }

    if (mode === 'update' && username && uniqueId && updatingId && name === 'return') {
      updateContact({ username, uniqueId, currentId: updatingId })
        .then(clearMode)
        .then(refetch);
    }

    if (mode === 'delete' && uniqueId && name === 'y') {
      deleteContact(uniqueId)
        .then(clearMode)
        .then(refetch);
    }

    if (mode === 'delete' && uniqueId && name === 'n') {
      clearMode();
    }
  });

  const loading = gettingContacts || creatingContact || updatingContact || deletingContact;

  return (
    <box margin={1}>
      {!loading && <text attributes={TextAttributes.DIM}>username - uuid</text>}
      <box marginBottom={1}>
        {loading && <text>Loading...</text>}
        {!!contacts && !loading && contacts.map((c, i) => (
          <ArrowFocusText
            focused={i === index && mode !== 'create'}
            text={`${c.username} - ${c.uniqueId}`}
            key={c.uniqueId}
          />
        ))}
      </box>
      {(mode === 'create' || mode === 'update') && (
        <box>
          <box flexDirection="row">
            <box title="Username" border borderStyle="rounded" height={3} flexGrow={1}>
              <input
                value={username}
                onInput={setUsername}
                focused={inputFocus === 0}
              />
            </box>
            <box title="UUID" border borderStyle="rounded" height={3} flexGrow={1}>
              <input
                value={uniqueId}
                onInput={setUniqueId}
                focused={inputFocus === 1}
              />
            </box>
          </box>
          <text attributes={TextAttributes.DIM}>enter to {mode === 'create' ? 'add' : 'edit'} user | esc-cancel</text>
        </box>
      )}
      {mode === 'delete' && (
        <text>Are you sure you want to remove {username} ({uniqueId})? (y / n)</text>
      )}
      {!mode && (
        <box>
          <text attributes={TextAttributes.DIM}>a-add | e-edit | r-remove | enter-open chat | esc-close</text>
        </box>
      )}
    </box>
  );
}
