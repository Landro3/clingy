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
          <text fg="#ff00ff">{title}</text>
          <box flexGrow={1} flexDirection="column">
            {children}
          </box>
        </box>
      </box>
    </box>
  );
};

export default Modal;

