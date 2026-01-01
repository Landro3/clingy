import { TextAttributes } from '@opentui/core';
import { useKeyboard } from '@opentui/react';
import { Pages, useNavigation } from '../context/navigation';
import { type ServerConfig, getServerConfig, setServerConfig as setServerConfigApi } from '../api/config';
import { useMutation, useQuery } from '../hooks/api';
import { useEffect } from 'react';

export default function Intro() {
  const { navigate } = useNavigation();

  const { data: serverConfig, loading: loadingServerConfig, refetch } = useQuery<ServerConfig>(getServerConfig);
  const { mutate: setServerConfig } = useMutation(setServerConfigApi);

  // TODO: Modal on error here
  useEffect(() => {
    if (loadingServerConfig) return;
    if (serverConfig) {
      setServerConfig({
        username: serverConfig.username,
        serverAddr: serverConfig. serverAddr,
      }).then(() => navigate(Pages.Chat));
      return;
    }
  
    if (!serverConfig && !loadingServerConfig) {
      navigate(Pages.Config);
    }
  }, [serverConfig]);

  return (
    <box alignItems="center" justifyContent="center" flexGrow={1}>
      <box justifyContent="center" alignItems="flex-start">
        <ascii-font font="tiny" text="clingy" />
        <text attributes={TextAttributes.DIM}>E2E Encrypted & AI-Powered Terminal Messaging</text>
        <text attributes={TextAttributes.DIM}>connecting to server...</text>
      </box>
    </box>
  );
}
