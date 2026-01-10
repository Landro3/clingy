import api from '.';

export const sendChatMessage = (body: { to: string; message: string }) => api.post('/chat', body);
