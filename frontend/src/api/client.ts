import axios from 'axios';

// Create axios instance with default config
const api = axios.create({
  baseURL: '/api',
  withCredentials: true, // Important for session cookies
  headers: {
    'Content-Type': 'application/json',
  },
});

// Types
export interface User {
  id: number;
  name: string;
  email: string;
  created_at: string;
}

export interface Board {
  ID: number;
  Name: string;
  Description: string;
  CreatedBy: number;
  CreatedAt: string;
  UpdatedAt: string;
}

export interface List {
  ID: number;
  Name: string;
  BoardID: number;
  Position: number;
  CreatedAt: string;
  UpdatedAt: string;
}

export interface Card {
  ID: number;
  Title: string;
  Description: string;
  ListID: number;
  Position: number;
  CreatedAt: string;
  UpdatedAt: string;
}

// Auth API
export const authAPI = {
  register: async (name: string, email: string, password: string): Promise<User> => {
    const response = await api.post('/register', { name, email, password });
    return response.data;
  },

  login: async (email: string, password: string): Promise<User> => {
    const response = await api.post('/login', { email, password });
    return response.data;
  },

  me: async (): Promise<User> => {
    const response = await api.get('/me');
    return response.data;
  },
};

// Boards API
export const boardsAPI = {
  getAll: async (): Promise<Board[]> => {
    const response = await api.get('/boards');
    return response.data;
  },

  getById: async (id: number): Promise<Board> => {
    const response = await api.get(`/boards/${id}`);
    return response.data;
  },

  create: async (name: string, description: string): Promise<Board> => {
    const response = await api.post('/boards', { name, description });
    return response.data;
  },

  update: async (id: number, name: string, description: string): Promise<Board> => {
    const response = await api.put(`/boards/${id}`, { name, description });
    return response.data;
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/boards/${id}`);
  },
};

// Lists API
export const listsAPI = {
  getByBoard: async (boardId: number): Promise<List[]> => {
    const response = await api.get(`/boards/${boardId}/lists`);
    return response.data;
  },

  create: async (boardId: number, name: string, position: number): Promise<List> => {
    const response = await api.post(`/boards/${boardId}/lists`, { name, position });
    return response.data;
  },

  update: async (id: number, name: string, position: number): Promise<List> => {
    const response = await api.put(`/lists/${id}`, { name, position });
    return response.data;
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/lists/${id}`);
  },
};

// Cards API
export const cardsAPI = {
  getByList: async (listId: number): Promise<Card[]> => {
    const response = await api.get(`/lists/${listId}/cards`);
    return response.data;
  },

  create: async (listId: number, title: string, description: string, position: number): Promise<Card> => {
    const response = await api.post(`/lists/${listId}/cards`, { title, description, position });
    return response.data;
  },

  update: async (id: number, title: string, description: string): Promise<Card> => {
    const response = await api.put(`/cards/${id}`, { title, description });
    return response.data;
  },

  move: async (id: number, listId: number, position: number): Promise<Card> => {
    const response = await api.put(`/cards/${id}/move`, { list_id: listId, position });
    return response.data;
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/cards/${id}`);
  },
};

export default api;
