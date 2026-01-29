import { TextAttributes } from '@opentui/core';
import { useKeyboard } from '@opentui/react';
import { useEffect, useState } from 'react';
import { type ServerConfig, getServerConfig, setServerConfig as setServerConfigApi } from '../api/config';
import { useMutation, useQuery } from '../hooks/api';
import ArrowFocusText from '#/components/arrow-focus';

enum Focus {
  Server,
  Username,
  Register,
}

interface ConfigProps {
  handleClose: () => void;
}

export default function Config({ handleClose }: ConfigProps) {
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
    if (name === 'escape') {
      handleClose();
      return;
    }

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
        .then(refetch)
        .then(handleClose);

      return;
    }
  });

  const loading = loadingServerConfig || settingServerConfig;

  return (
    <box margin={1}>
      {loading && <text>Loading...</text>}
      {!loading && (
        <box>
          <text>TODO: Connection explanation and also feedback on register attempt before closing modal</text>
          <box title="Clingy Server" borderStyle="rounded" border width={40} height={3}>
            <input
              placeholder="Enter server address..."
              value={serverAddr}
              onInput={setServerAddr}
              focused={focus === Focus.Server}
            />
          </box>
          <box title="Username" borderStyle="rounded" border width={40} height={3}>
            <input
              placeholder="Enter username..."
              value={username}
              onInput={setUsername}
              focused={focus === Focus.Username}
            />
          </box>
          <box flexDirection="row" alignItems="center">
            <ArrowFocusText text="Register" focused={focus === Focus.Register} />
            <text attributes={TextAttributes.DIM}>
              Current ID: {uniqueId}
            </text>
          </box>
        </box>
      )}
    </box>
  );
}
