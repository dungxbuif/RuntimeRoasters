/**
 * Domain constants synchronized with Backend
 */

export const FARM_LOCATIONS = [
  { code: 'CAU_DAT', name: 'Cầu Đất, Đà Lạt' },
  { code: 'BUON_MA_THUOT', name: 'Buôn Ma Thuột, Đắk Lắk' },
  { code: 'PLEIKU', name: 'Pleiku, Gia Lai' },
  { code: 'GIA_NGHIA', name: 'Gia Nghĩa, Đắk Nông' },
  { code: 'KON_TUM', name: 'Kon Tum' },
] as const;

export const COFFEE_TYPES = [
  { code: 'ARABICA', name: 'Arabica' },
  { code: 'ROBUSTA', name: 'Robusta' },
  { code: 'CHERRY', name: 'Cherry' },
  { code: 'CULI', name: 'Culi' },
] as const;

export const HARVEST_STATUSES = {
  NEW: { label: 'New', color: 'bg-primary' },
  PROCESSING: { label: 'Processing', color: 'bg-tertiary' },
  COMPLETED: { label: 'Completed', color: 'bg-emerald-500' },
} as const;

export const USER_ROLE_LABELS: Record<string, string> = {
  'ADMIN': 'System Admin',
  'FARM_ADMIN': 'Agri Admin',
  'FARM_MANAGER': 'Farm Manager',
  'FARMER': 'Field Farmer',
  'PROCESSOR': 'Roast Master',
  'DRIVER': 'Logistics Driver',
  'GUEST': 'Guest',
};
