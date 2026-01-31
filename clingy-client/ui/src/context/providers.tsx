import type { PropsWithChildren } from 'react';
import { NavigationProvider } from './navigation';
import { ChatProvider } from './chat';
import { ServerConfigProvider } from './server-config';

export default function Providers({ children }: PropsWithChildren) {

  return (
    <NavigationProvider>
      <ServerConfigProvider>
        <ChatProvider>
          {children}
        </ChatProvider>
      </ServerConfigProvider>
    </NavigationProvider>
  );
}
