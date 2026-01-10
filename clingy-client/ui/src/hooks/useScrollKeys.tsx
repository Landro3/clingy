import { useKeyboard } from '@opentui/react';
import { useState } from 'react';

export default function useScrollKeys(numItems: number) {
  const [index, setIndex] = useState(0);

  const handleUp = () => {
    if (index <= 0) {
      setIndex(numItems - 1);
      return;
    }
    setIndex(index - 1);
  };

  const handleDown = () => {
    if (index >= numItems - 1) {
      setIndex(0);
      return;
    }
    setIndex(index + 1);
  };

  useKeyboard(({ name }) => {
    switch (name) {
      case 'up':
      case 'k':
        handleUp();
        break;
      case 'down':
      case 'j':
        handleDown();
        break;
    }
  });

  return index;
};
