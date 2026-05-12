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
    const res = await api.post("/v1/users", data);
    return res.data.user;
  }

  async listUsers(): Promise<User[]> {
    const res = await api.get("/v1/users");
    return res.data.users;
  }

  async listManagers(): Promise<User[]> {
    const users = await this.listUsers();
    const managerRoles = ['FARM_MANAGER', 'FARM_ADMIN', 'ADMIN'];
    return users.filter(u => managerRoles.includes(u.role.toUpperCase()));
  }

  async createFarm(data: CreateFarmRequest): Promise<Farm> {
    const res = await api.post("/v1/farms", data);
    return res.data.farm;
  }
}

export const adminService = new AdminService();
