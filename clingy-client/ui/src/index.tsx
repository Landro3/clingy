import { createCliRenderer } from '@opentui/core';
import { createRoot } from '@opentui/react';
import Pages from './pages';
import Providers from './context/providers';

function App() {
  return (
    <Providers>
      <Pages />
    </Providers>
  );
}

const renderer = await createCliRenderer();
renderer.console.show();
createRoot(renderer).render(<App />);
