import api from "@/lib/axios";

export interface CreateUserRequest {
  email: string;
  password: string;
  name: string;
  role: string;
}

export interface User {
  id: string;
  email: string;
  name: string;
  role: string;
}

export interface CreateFarmRequest {
  name: string;
  location: string;
  area: number;
  farm_type: string;
  owner_id: string;
}

export interface Farm {
  id: string;
  name: string;
  location: string;
  area: number;
  farm_type: string;
  owner_id: string;
}

class AdminService {
  async createUser(data: CreateUserRequest): Promise<User> {
    const res = await api.post("/v1/admin/users", data);
    return res.data.user;
  }

  async listUsers(): Promise<User[]> {
    const res = await api.get("/v1/admin/users");
    return res.data.users;
  }

  async listManagers(): Promise<User[]> {
    const users = await this.listUsers();
    return users.filter(u => u.role === 'farm_manager');
  }

  async createFarm(data: CreateFarmRequest): Promise<Farm> {
    const res = await api.post("/v1/farms", data);
    return res.data.farm;
  }
}

export const adminService = new AdminService();
