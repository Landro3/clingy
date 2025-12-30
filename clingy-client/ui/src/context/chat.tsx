import { createContext, useContext, useState, type Dispatch, type PropsWithChildren, type SetStateAction } from 'react';

interface ChatContextType {
  chatUser: string | null;
  setChatUser: Dispatch<SetStateAction<string | null>>;
}

const ChatContext = createContext<ChatContextType | undefined>(undefined);

export function ChatProvider({ children }: PropsWithChildren) {
  const [chatUser, setChatUser] = useState<string | null>(null);

  const value: ChatContextType = {
    chatUser,
    setChatUser,
  };

  return (
    <ChatContext.Provider value={value}>
      {children}
    </ChatContext.Provider>
  );
}

export function useChat() {
  const context = useContext(ChatContext);
  if (context === undefined) {
    throw new Error('useChat must be used within a ChatProvider');
  }
  return context;
}

