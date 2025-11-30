import { useKeyboard } from '@opentui/react';
import { useState } from 'react';
import { Pages, useNavigation } from '../context/navigation';

export default function Chat() {
  const { navigate } = useNavigation();
  const [focus, setFocus] = useState(0);

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
        if (focus === 1) navigate(Pages.Contacts);
        if (focus === 2) navigate(Pages.Config);
        break;
    }
  });

  return (
    <box>
      <box title=" Username:server1 " border borderStyle="rounded" height={20}>

      </box>
      <box>
        <box border height={3} borderStyle="rounded">
          <input
            placeholder="Type your message here..."
            focused={focus === 0}
          />
        </box>
        <box flexDirection="row">
          <box alignItems="center" border borderStyle="rounded" flexGrow={1} borderColor={focus === 1 ? 'red' : 'white'}>
            <text fg={focus === 1 ? 'red' : 'white'}>Contact</text>
          </box>
          <box alignItems="center" border borderStyle="rounded" flexGrow={1} borderColor={focus === 2 ? 'red' : 'white'}>
            <text fg={focus === 2 ? 'red' : 'white'}>Config</text>
          </box>
        </box>
      </box>
    </box>
  );
}
