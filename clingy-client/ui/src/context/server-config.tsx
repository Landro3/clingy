import { createContext, useCallback, useContext, type PropsWithChildren } from 'react';
import { useMutation, useQuery } from '../hooks/api';
import { getServerConfig, setServerConfig, type ServerConfig } from '../api/config';

interface ServerConfigContextType {
  serverConfig: ServerConfig | null;
  loading: boolean;
  updating: boolean;
  error: string | null;
  updateConfig: (config: { serverAddr: string; username: string }) => Promise<void>;
}

const ServerConfigContext = createContext<ServerConfigContextType | undefined>(undefined);

export function ServerConfigProvider({ children }: PropsWithChildren) {
  const { data: serverConfig, loading, error, refetch } = useQuery<ServerConfig>(getServerConfig);
  const { mutate: updateConfigMutation, loading: updating } = useMutation(setServerConfig);

  const updateConfig = useCallback(async (config: { serverAddr: string; username: string }) => {
    await updateConfigMutation(config);
    refetch();
  }, [updateConfigMutation, refetch]);

  const value: ServerConfigContextType = {
    serverConfig,
    loading,
    updating,
    error,
    updateConfig,
  };

  return (
    <ServerConfigContext.Provider value={value}>
      {children}
    </ServerConfigContext.Provider>
  );
}

export function useServerConfig() {
  const context = useContext(ServerConfigContext);
  if (context === undefined) {
    throw new Error('useServerConfig must be used within a ServerConfigProvider');
  }
  return context;
}
