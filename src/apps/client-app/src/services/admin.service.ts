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
  id: number;
  name: string;
  location: string;
  area: number;
  farm_type: string;
  owner_id: string;
}

export interface CreateWarehouseRequest {
  name: string;
  code: string;
  location: string;
  manager_id: string;
  manager_email: string;
  capacity: string;
}

export interface Warehouse {
  id: string;
  name: string;
  code: string;
  location: string;
  manager_id: string;
  manager_email: string;
  capacity: string;
  status: string;
}

export interface CreateStoreRequest {
  name: string;
  city: string;
  address: string;
  manager_id: string;
  manager_email: string;
}

export interface Store {
  id: string;
  name: string;
  city: string;
  address: string;
  manager_id: string;
  manager_email: string;
  status: string;
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
    const managerRoles = ['FARM_MANAGER', 'FARM_ADMIN', 'ADMIN', 'WAREHOUSE_MGR', 'STORE_MGR'];
    return users.filter(u => managerRoles.includes(u.role.toUpperCase()));
  }

  async createFarm(data: CreateFarmRequest): Promise<Farm> {
    const res = await api.post("/v1/farms", data);
    return res.data.farm;
  }

  async listWarehouses(): Promise<Warehouse[]> {
    const res = await api.get("/v1/warehouse/warehouses");
    return res.data.warehouses || [];
  }

  async createWarehouse(data: CreateWarehouseRequest): Promise<Warehouse> {
    const res = await api.post("/v1/warehouse/warehouses", data);
    return res.data.warehouse;
  }

  async listStores(): Promise<Store[]> {
    const res = await api.get("/v1/retail/stores");
    return res.data.stores || [];
  }

  async createStore(data: CreateStoreRequest): Promise<Store> {
    const res = await api.post("/v1/retail/stores", data);
    return res.data.store;
  }
}

export const adminService = new AdminService();
