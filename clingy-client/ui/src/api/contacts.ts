import api from '.';

export interface Contact {
  id: string;
  username: string;
}

export const getContacts = () => api.get('contacts');

export const createContact = (body: Contact) => api.post('/contacts', body);

export const updateContact = (body: Contact & { currentId: string }) => api.put('/contacts', body);

export const deleteContact = (id: string) => api.delete('/contacts', { params: { id } });
