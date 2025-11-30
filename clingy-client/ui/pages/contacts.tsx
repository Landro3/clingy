import { TextAttributes } from '@opentui/core';
import { useState } from 'react';
import useScrollKeys from '../src/hooks/useScrollKeys';
import { useKeyboard } from '@opentui/react';
import { Pages, useNavigation } from '../src/context/navigation';
import {
  type Contact,
  createContact as createContactApi,
  updateContact as updateContactApi,
  deleteContact as deleteContactApi,
  getContacts,
} from '../src/api/contacts';
import { useMutation, useQuery } from '../src/hooks/api';

export default function Contacts() {
  const { navigate } = useNavigation();
  const [mode, setMode] = useState<'create' | 'update' | 'delete' | null>(null);
  const [inputFocus, setInputFocus] = useState(0);
  const [username, setUsername] = useState('');
  const [id, setId] = useState('');
  const [updatingId, setUpdatingId] = useState('');

  const { data: contacts, loading: gettingContacts, refetch } = useQuery<Contact[]>(getContacts);
  const index = useScrollKeys(contacts?.length ?? 0);

  const { mutate: createContact, loading: creatingContact } = useMutation(createContactApi);
  const { mutate: updateContact, loading: updatingContact } = useMutation(updateContactApi);
  const { mutate: deleteContact, loading: deletingContact } = useMutation(deleteContactApi);

  useKeyboard(({ name }) => {
    if (!contacts) return;

    const clearMode = () => {
      setMode(null);
      setUsername('');
      setId('');
      setUpdatingId('');
      setInputFocus(0);
    };

    if (mode === null) {
      switch (name) {
        case 'a':
          setMode('create');
          break;
        case 'e':
          setMode('update');
          if (!contacts[index]) {
            return;
          }
          setUsername(contacts[index].username);
          setUpdatingId(contacts[index].id);
          setId(contacts[index].id);
          break;
        case 'r':
          setMode('delete');
          if (!contacts[index]) {
            return;
          }
          setUsername(contacts[index].username);
          setId(contacts[index].id);
          break;
        case 'escape':
          navigate(Pages.Chat);
          break;
      }
    } else {
      if (name === 'escape') {
        clearMode();
        return;
      }
    }

    if ((mode === 'create' || mode === 'update') && name === 'tab') {
      if (inputFocus === 1) {
        setInputFocus(0);
      } else {
        setInputFocus(inputFocus + 1);
      }

      return;
    }

    if (mode === 'create' && username && id && name === 'return') {
      createContact({ username, id })
        .then(clearMode)
        .then(refetch);
    }

    if (mode === 'update' && username && id && updatingId && name === 'return') {
      updateContact({ username, id, currentId: updatingId })
        .then(clearMode)
        .then(refetch);
    }

    if (mode === 'delete' && id && name === 'y') {
      deleteContact(id)
        .then(clearMode)
        .then(refetch);
    }

    if (mode === 'delete' && id && name === 'n') {
      clearMode();
    }
  });

  const loading = gettingContacts || creatingContact || updatingContact || deletingContact;

  return (
    <box>
      <text>Contacts</text>
      <box flexDirection="row" gap={3} marginBottom={1}>
        <text>a - add</text>
        <text>e - edit</text>
        <text>r - remove</text>
      </box>
      <box marginLeft={1}>
        {loading && <text>Loading...</text>}
        {!!contacts && !loading && contacts.map((c, i) => (
          <box
            border
            borderStyle="rounded"
            borderColor={(i === index && !mode) ? 'white' : 'transparent'}
            paddingLeft={2}
            key={c.id}
          >
            <text>{c.username}</text>
            <text attributes={TextAttributes.DIM}>{c.id}</text>
          </box>
        ))}
      </box>
      {(mode === 'create' || mode === 'update') && (
        <box>
          <box flexDirection="row">
            <box title="Username" border height={3} flexGrow={1}>
              <input
                value={username}
                onInput={setUsername}
                focused={inputFocus === 0}
              />
            </box>
            <box title="UUID" border height={3} flexGrow={1}>
              <input
                value={id}
                onInput={setId}
                focused={inputFocus === 1}
              />
            </box>
          </box>
          <box>
            <text attributes={TextAttributes.DIM}>Enter to {mode === 'create' ? 'add' : 'edit'} user</text>
            <text attributes={TextAttributes.DIM}>esc to exit</text>
          </box>
        </box>
      )}
      {mode === 'delete' && (
        <text>Are you sure you want to remove {username} ({id})? (y / n)</text>
      )}
      {!mode && <text attributes={TextAttributes.DIM}>esc to return to chat</text>}
    </box>
  );
}
