import type { PropsWithChildren } from 'react';
import { NavigationProvider } from './navigation';

export default function Providers({ children }: PropsWithChildren) {

  return (
    <NavigationProvider>
      {children}
    </NavigationProvider>
  );
}
