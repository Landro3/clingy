import { TextAttributes } from '@opentui/core';
import { useState } from 'react';
import useScrollKeys from '../hooks/useScrollKeys';
import { useKeyboard } from '@opentui/react';
import { Pages } from '../context/navigation';
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
  const [id, setId] = useState('');
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
      setId('');
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
          setUpdatingId(contacts[index].id);
          setId(contacts[index].id);
          return;
        case 'r':
          setMode('delete');
          if (!contacts[index]) {
            return;
          }
          setUsername(contacts[index].username);
          setId(contacts[index].id);
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

    if (mode && ['create', 'update'].includes(mode) && (!username || !id) && name === 'return') {
      clearMode();
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
    <box margin={1}>
      {/* <text marginBottom={1}>Contacts</text> */}
      {!loading && <text attributes={TextAttributes.DIM}>username - uuid</text>}
      <box marginBottom={1}>
        {loading && <text>Loading...</text>}
        {!!contacts && !loading && contacts.map((c, i) => (
          <ArrowFocusText
            focused={i === index && !['create', 'delete'].includes(mode ?? '')}
            text={`${c.username} - ${c.id}`}
            key={c.id}
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
                value={id}
                onInput={setId}
                focused={inputFocus === 1}
              />
            </box>
          </box>
          <text attributes={TextAttributes.DIM}>enter to {mode === 'create' ? 'add' : 'edit'} user | esc-cancel</text>
        </box>
      )}
      {mode === 'delete' && (
        <text>Are you sure you want to remove {username} ({id})? (y / n)</text>
      )}
      {!mode && (
        <box>
          <text attributes={TextAttributes.DIM}>a-add | e-edit | r-remove | enter-open chat | esc-close</text>
        </box>
      )}
    </box>
  );
}
