import { TextAttributes } from '@opentui/core';
import type { PropsWithChildren } from 'react';

interface ModalProps extends PropsWithChildren {
  title: string;
}

const Modal = ({ title, children }: ModalProps) => {
  return (
    <box
      position="absolute"
      height="100%"
      width="100%"
      backgroundColor="#111111bb"
      flexDirection="column"
      alignItems="center"
      justifyContent="center"
    >
      <box
        backgroundColor="#111111"
        width="50%"
      >
        <box
          flexGrow={1}
          flexDirection="column"
          padding={1}
          backgroundColor="#010101"
        >
          <box flexDirection="row" justifyContent="space-between">
            <text fg="#ff00ff">{title}</text>
            <text fg="#ff00ff" attributes={TextAttributes.DIM}>
              esc-close
            </text>
          </box>
          <box flexGrow={1} flexDirection="column">
            {children}
          </box>
        </box>
      </box>
    </box>
  );
};

export default Modal;

