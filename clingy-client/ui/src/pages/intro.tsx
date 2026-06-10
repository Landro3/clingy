import { TextAttributes } from '@opentui/core';
import { useKeyboard } from '@opentui/react';
import { Pages, useNavigation } from '../context/navigation';
import { registerWithServer as registerWithServerApi } from '../api/config';
import { useMutation } from '../hooks/api';
import { useServerConfig } from '../context/server-config';
import { useEffect } from 'react';

export default function Intro() {
  const { navigate } = useNavigation();

  const { serverConfig, loading: loadingServerConfig } = useServerConfig();
  const { mutate: registerWithServer } = useMutation(registerWithServerApi);

  useKeyboard((key) => {
    if (key.name === 'return' && !loadingServerConfig) {
      navigate(Pages.Chat);
    }
  });

  // TODO: Modal on error here
  useEffect(() => {
    if (loadingServerConfig) return;
    if (serverConfig) {
      registerWithServer({})
        .then(() => navigate(Pages.Chat))
        .catch(() => navigate(Pages.Chat));
      return;
    }

    if (!serverConfig && !loadingServerConfig) {
      navigate(Pages.Chat);
    }
  }, [serverConfig, loadingServerConfig, navigate, registerWithServer]);

  return (
    <box alignItems="center" justifyContent="center" flexGrow={1} height={30}>
      <box justifyContent="center" alignItems="flex-start">
        <ascii-font font="tiny" text="clingy" />
        <text attributes={TextAttributes.DIM}>E2E Encrypted (soon) & AI-Powered (soon again) Terminal Messaging (sorry)</text>
        <text attributes={TextAttributes.DIM}>connecting to server...</text>
      </box>
    </box>
  );
}
