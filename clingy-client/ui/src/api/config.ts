import api from '.';

export interface ServerConfig {
  serverAddr: string;
  username: string;
  uniqueId: string;
}

export const getServerConfig = () => api.get('/config/server');

export const setServerConfig = (body: { serverAddr: string; username: string }) => api.post('/config/server', body);

export const registerWithServer = () => api.post('/register');

