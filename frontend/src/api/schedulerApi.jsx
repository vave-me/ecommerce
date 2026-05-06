import axiosInstance from './axiosInstance';

// Scheduler Management API Client
// Implements the Swagger API specification provided

export const schedulerApi = {
  // GET /api/scheduler - List all schedulers
  async getSchedulers(userId = undefined) {
    try {
      const params = userId ? { userId } : {};
      const response = await axiosInstance.get('/api/scheduler', { params });
      return response.data;
    } catch (error) {
      // Error: 'Error fetching schedulers:', error...
      throw error;
    }
  },

  // POST /api/scheduler - Create a new scheduler
  async createScheduler(data) {
    try {
      const response = await axiosInstance.post('/api/scheduler', data);
      return response.data;
    } catch (error) {
      // Error: 'Error creating scheduler:', error...
      throw error;
    }
  },

  // GET /api/scheduler/{userId} - Retrieve a specific scheduler
  async getScheduler(userId) {
    try {
      const response = await axiosInstance.get(`/api/scheduler/${userId}`);
      return response.data;
    } catch (error) {
      // Error: 'Error fetching scheduler:', error...
      throw error;
    }
  },

  // GET /api/scheduler/{schedulerId}/actions - List all actions for a scheduler
  async getActions(schedulerId) {
    try {
      const response = await axiosInstance.get(`/api/scheduler/${schedulerId}/actions`);
      return response.data;
    } catch (error) {
      // Error: 'Error fetching actions:', error...
      throw error;
    }
  },

  // POST /api/scheduler/actions - Add a new action
  async addAction(data) {
    try {
      const response = await axiosInstance.post('/api/scheduler/actions', data);
      return response.data;
    } catch (error) {
      // Error: 'Error adding action:', error...
      throw error;
    }
  },

  // GET /api/scheduler/actions/{id} - Retrieve a specific action
  async getAction(actionId) {
    try {
      const response = await axiosInstance.get(`/api/scheduler/actions/${actionId}`);
      return response.data;
    } catch (error) {
      // Error: 'Error fetching action:', error...
      throw error;
    }
  },

  // DELETE /api/scheduler/actions/{id} - Remove an action
  async removeAction(actionId) {
    try {
      const response = await axiosInstance.delete(`/api/scheduler/actions/${actionId}`);
      return response.data;
    } catch (error) {
      // Error: 'Error removing action:', error...
      throw error;
    }
  },
};

export default schedulerApi; 