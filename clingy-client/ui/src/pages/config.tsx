import { TextAttributes } from '@opentui/core';
import { useKeyboard } from '@opentui/react';
import { useEffect, useState } from 'react';
import { Pages, useNavigation } from '../context/navigation';
import { type ServerConfig, getServerConfig, setServerConfig as setServerConfigApi } from '../api/config';
import { useMutation, useQuery } from '../hooks/api';
import FocusTextBox from '../components/FocusTextBox';

enum Focus {
  Server,
  Username,
  Register,
}

export default function Config() {
  const { navigate } = useNavigation();

  const { data: serverConfig, loading: loadingServerConfig, refetch } = useQuery<ServerConfig>(getServerConfig);

  const { mutate: setServerConfig, loading: settingServerConfig } = useMutation(setServerConfigApi);

  const [focus, setFocus] = useState(0);
  const [serverAddr, setServerAddr] = useState('');
  const [username, setUsername] = useState('');
  const [uniqueId, setUniqueId] = useState('');

  useEffect(() => {
    if (serverConfig) {
      setServerAddr(serverConfig.serverAddr);
      setUsername(serverConfig.username);
      setUniqueId(serverConfig.uniqueId);
    }

  }, [serverConfig]);

  useKeyboard(({ name }) => {
    if (name === 'tab') {
      if (focus >= Focus.Register) {
        setFocus(0);
      } else {
        setFocus(focus + 1);
      }

      return;
    }

    if (name === 'return' && focus === Focus.Register) {
      setServerConfig({ serverAddr, username })
        .then(refetch);

      return;
    }

    if (name === 'escape') {
      navigate(Pages.Chat);

      return;
    }
  });

  const loading = loadingServerConfig || settingServerConfig;

  return (
    <box>
      <text marginBottom={1}>Config</text>
      {loading && <text>Loading...</text>}
      {!loading && (
        <box>
          <box title="Clingy Server" style={{ border: true, width: 40, height: 3 }}>
            <input
              placeholder="Enter server address..."
              value={serverAddr}
              onInput={setServerAddr}
              focused={focus === Focus.Server}
            />
          </box>
          <box title="Username" style={{ border: true, width: 40, height: 3 }}>
            <input
              placeholder="Enter username..."
              value={username}
              onInput={setUsername}
              focused={focus === Focus.Username}
            />
          </box>
          <box flexDirection="row" alignItems="center">
            <FocusTextBox text="Regsister" focused={focus === Focus.Register} />
            <text attributes={TextAttributes.DIM}>
              Current ID: {uniqueId}
            </text>
          </box>
          <box>
            <text attributes={TextAttributes.DIM}>esc to return to chat</text>
          </box>
        </box>
      )}
    </box>
  );
}
