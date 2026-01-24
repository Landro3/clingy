import { createCliRenderer } from '@opentui/core';
import { createRoot } from '@opentui/react';
import Pages from './pages';
import Providers from './context/providers';

function App() {
  return (
    <Providers>
      <box backgroundColor="#222222" height="100%">
        <Pages />
      </box>
    </Providers>
  );
}

const renderer = await createCliRenderer();
// renderer.console.show();
createRoot(renderer).render(<App />);
