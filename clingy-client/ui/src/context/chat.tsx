import { createContext, useCallback, useContext, useEffect, useState, type Dispatch, type PropsWithChildren, type SetStateAction } from 'react';
import { useMutation } from '../hooks/api';
import { sendChatMessage as sendChatMessageApi } from '../api/chat';
import { useServerConfig } from '../context/server-config';

export type Message = { from: string; text: string; fromSelf: boolean };
export type ChatMap = Record<string, Message[]>

interface ChatContextType {
  chatUser: string | null;
  setChatUser: Dispatch<SetStateAction<string | null>>;
  chatMap: ChatMap;
  sendChatMessage: (to: string, message: string) => Promise<void>;
}

const ChatContext = createContext<ChatContextType | undefined>(undefined);

export function ChatProvider({ children }: PropsWithChildren) {
  const [chatUser, setChatUser] = useState<string | null>(null);
  const [chatMap, setChatMap] = useState<ChatMap>({});

  const { mutate: sendChatMessageMutation } = useMutation(sendChatMessageApi);
  const { serverConfig } = useServerConfig();

  const addMessageToMap = useCallback((otherUser: string, text: string, fromSelf: boolean) => {
    if (!serverConfig) {
      return;
    }

    setChatMap((prev) => {
      if (!prev[otherUser]) {
        prev[otherUser] = [];
      }

      prev[otherUser].push({
        text,
        from: fromSelf ? serverConfig.username : otherUser,
        fromSelf,
      });

      return { ...prev };
    });
  }, [serverConfig]);

  const sendChatMessage = async (to: string, message: string) => {
    sendChatMessageMutation({ to, message })
      .then(() => addMessageToMap(to, message, true));
  };

  useEffect(() => {
    const controller = new AbortController();
    const connect = async () => {
      try {
        const response = await fetch(`${process.env.API_URL}/chat/stream`, {
          headers: { 'Accept': 'text/event-stream' },
          signal: controller.signal,
        });

        if (!response.body) return;

        let buffer = '';
        for await (const chunk of response.body) {
          buffer += new TextDecoder().decode(chunk);

          const lines = buffer.split('\n\n');
          buffer = lines.pop() || '';

          for (const line of lines) {
            if (line.startsWith('data: ')) {
              const rawJson = line.replace('data: ', '').trim();
              const { from, message } = JSON.parse(rawJson);
              addMessageToMap(from, message, false);
            }
          }
        }
      } catch (err) {
        if (err instanceof Error && err.name === 'AbortError') {
          return;
        }

        if (err instanceof Error) {
          console.error(err.name);
          connect();
          return;
        }

        console.error(err);
      }
    };

    connect();

    return () => controller.abort();
  }, [addMessageToMap]);

  const value: ChatContextType = {
    chatUser,
    setChatUser,
    chatMap,
    sendChatMessage,
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

