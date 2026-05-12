/**
 * Utility for consistent data-e2e attributes
 */
export const testId = (id: string) => {
  return { 'data-e2e': id };
};

export const e2eSelectors = {
  // Common
  LOGIN_SUBMIT: 'login-submit',
  LOGOUT_BTN: 'logout-btn',
  DASHBOARD_LINK: 'dashboard-link',
  
  // User Management
  USER_LIST_TABLE: 'user-list-table',
  CREATE_USER_BTN: 'create-user-btn',
  USER_MODAL: 'user-modal',
  USER_EMAIL_INPUT: 'user-email-input',
  USER_PASS_INPUT: 'user-password-input',
  USER_ROLE_SELECT: 'user-role-select',
  USER_SUBMIT_BTN: 'user-submit-btn',
  
  // Farm Management
  FARM_LIST_TABLE: 'farm-list-table',
  CREATE_FARM_BTN: 'create-farm-btn',
  FARM_MODAL: 'farm-modal',
  FARM_NAME_INPUT: 'farm-name-input',
  FARM_AREA_INPUT: 'farm-area-input',
  FARM_OWNER_SELECT: 'farm-owner-select',
  FARM_SUBMIT_BTN: 'farm-submit-btn',
  FARM_DELETE_BTN: 'farm-delete-btn',
  FARM_EDIT_BTN: 'farm-edit-btn',
};
