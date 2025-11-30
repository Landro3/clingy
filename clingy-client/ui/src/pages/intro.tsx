import { TextAttributes } from '@opentui/core';
import { useKeyboard } from '@opentui/react';
import { Pages, useNavigation } from '../context/navigation';

export default function Intro() {
  const { navigate } = useNavigation();

  useKeyboard((key) => {
    if (key.name === 'return') {
      navigate(Pages.Config);
    }
  });

  return (
    <box alignItems="center" justifyContent="center" flexGrow={1}>
      <box justifyContent="center" alignItems="flex-start">
        <ascii-font font="tiny" text="clingy" />
        <text attributes={TextAttributes.DIM}>E2E Encrypted & AI-Powered Terminal Messaging</text>
      </box>
      <box marginTop={2} border={true} borderStyle="rounded">
        <text>Get Started</text>
      </box>
    </box>
  );
}
