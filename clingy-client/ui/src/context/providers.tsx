import type { PropsWithChildren } from 'react';
import { NavigationProvider } from './navigation';
import { ChatProvider } from './chat';

export default function Providers({ children }: PropsWithChildren) {

  return (
    <NavigationProvider>
      <ChatProvider>
        {children}
      </ChatProvider>
    </NavigationProvider>
  );
}
