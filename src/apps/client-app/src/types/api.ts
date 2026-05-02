export interface ApiResponse<T> {
  data: T;
  status: number;
}

export interface DemoPingResponse {
  id: string;
  message: string;
}

export interface SystemError {
  code: number;
  msg: string;
}
