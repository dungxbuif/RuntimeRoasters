# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: sprint3.spec.ts >> Sprint 3: High-Fidelity Flow & Role Constraints >> Persona: Farm Manager - Rights & Restrictions
- Location: e2e/sprint3.spec.ts:93:7

# Error details

```
Error: expect(page).toHaveURL(expected) failed

Expected pattern: /\/dashboard\/users/
Received string:  "http://localhost:3000/login"
Timeout: 15000ms

Call log:
  - Expect "toHaveURL" with timeout 15000ms
    34 × unexpected value "http://localhost:3000/login"

```

```yaml
- banner:
  - text: coffee
  - heading "System Access" [level=1]
  - paragraph: Please authenticate to manage cluster nodes
- main:
  - text: Email
  - textbox "Enter your email...": mgr.1778651280600@runtimeroasters.com
  - text: Password
  - textbox "Enter your password...": Hello@123
  - button "Authenticate"
- alert
```

# Test source

```ts
  1   | import { test, expect } from '@playwright/test';
  2   | 
  3   | /**
  4   |  * Sprint 3 E2E Test: Full Vertical Slice (Refined)
  5   |  * Rules:
  6   |  * - Admin creates Manager
  7   |  * - Admin creates Farm & assigns to Manager
  8   |  * - Manager logs in -> Can view/edit own farm
  9   |  * - Manager CANNOT delete any farm (Role restriction)
  10  |  * - Data isolation: Manager cannot see other's farms
  11  |  */
  12  | 
  13  | const ADMIN_EMAIL = 'admin@runtimeroasters.com';
  14  | const ADMIN_PASS = 'Hello@123';
  15  | 
  16  | const TEST_MANAGER_EMAIL = `mgr.${Date.now()}@runtimeroasters.com`;
  17  | const TEST_MANAGER_PASS = 'Hello@123';
  18  | const TEST_FARM_NAME = `Hectare Node ${Date.now()}`;
  19  | 
  20  | // Use the same strings as e2eSelectors to ensure consistency
  21  | const SELECTORS = {
  22  |   LOGIN_SUBMIT: '[data-e2e="login-submit"]',
  23  |   USER_LIST_TABLE: '[data-e2e="user-list-table"]',
  24  |   CREATE_USER_BTN: '[data-e2e="create-user-btn"]',
  25  |   USER_EMAIL_INPUT: '[data-e2e="user-email-input"]',
  26  |   USER_PASS_INPUT: '[data-e2e="user-password-input"]',
  27  |   USER_ROLE_SELECT: '[data-e2e="user-role-select"]',
  28  |   USER_SUBMIT_BTN: '[data-e2e="user-submit-btn"]',
  29  |   
  30  |   CREATE_FARM_BTN: '[data-e2e="create-farm-btn"]',
  31  |   FARM_NAME_INPUT: '[data-e2e="farm-name-input"]',
  32  |   FARM_AREA_INPUT: '[data-e2e="farm-area-input"]',
  33  |   FARM_OWNER_SELECT: '[data-e2e="farm-owner-select"]',
  34  |   FARM_SUBMIT_BTN: '[data-e2e="farm-submit-btn"]',
  35  |   FARM_DELETE_BTN: '[data-e2e="farm-delete-btn"]',
  36  |   FARM_LIST_TABLE: '[data-e2e="farm-list-table"]',
  37  |   
  38  |   AUTH_LOADER: '[data-e2e="auth-loader"]',
  39  |   DASHBOARD_HEADER: '[data-e2e="dashboard-header"]',
  40  | };
  41  | 
  42  | test.describe('Sprint 3: High-Fidelity Flow & Role Constraints', () => {
  43  | 
  44  |   test.beforeEach(async ({ page }) => {
  45  |     // Clear cookies to ensure clean session for each test
  46  |     await page.context().clearCookies();
  47  |   });
  48  | 
  49  |   test('Persona: System Admin - Full Provisioning Flow', async ({ page }) => {
  50  |     // 1. Login as Admin
  51  |     await page.goto('/login');
  52  |     await page.waitForSelector('input[name="identifier"]');
  53  |     await page.fill('input[name="identifier"]', ADMIN_EMAIL);
  54  |     await page.fill('input[name="password"]', ADMIN_PASS);
  55  |     await page.click(SELECTORS.LOGIN_SUBMIT);
  56  |     
  57  |     // Redirect check
  58  |     await expect(page).toHaveURL(/\/dashboard\/users/, { timeout: 15000 });
  59  |     
  60  |     // Wait for AuthGuard loader to disappear
  61  |     await page.waitForSelector(SELECTORS.AUTH_LOADER, { state: 'hidden', timeout: 10000 });
  62  |     await expect(page.locator(SELECTORS.DASHBOARD_HEADER)).toBeVisible({ timeout: 10000 });
  63  | 
  64  |     // 2. Create a Manager
  65  |     await page.click(SELECTORS.CREATE_USER_BTN);
  66  |     await page.fill(SELECTORS.USER_EMAIL_INPUT, TEST_MANAGER_EMAIL);
  67  |     await page.fill(SELECTORS.USER_PASS_INPUT, TEST_MANAGER_PASS);
  68  |     await page.selectOption(SELECTORS.USER_ROLE_SELECT, 'FARM_MANAGER');
  69  |     await page.click(SELECTORS.USER_SUBMIT_BTN);
  70  | 
  71  |     // Verify
  72  |     await expect(page.locator(`[data-e2e="user-row-${TEST_MANAGER_EMAIL}"]`)).toBeVisible({ timeout: 10000 });
  73  | 
  74  |     // 3. Create a Farm for this Manager
  75  |     await page.goto('/dashboard/farm-ops/registry');
  76  |     await page.waitForSelector(SELECTORS.AUTH_LOADER, { state: 'hidden' });
  77  |     await page.click(SELECTORS.CREATE_FARM_BTN);
  78  |     await page.fill(SELECTORS.FARM_NAME_INPUT, TEST_FARM_NAME);
  79  |     await page.fill(SELECTORS.FARM_AREA_INPUT, '88.5');
  80  |     
  81  |     // Select the manager
  82  |     await page.selectOption(SELECTORS.FARM_OWNER_SELECT, { label: TEST_MANAGER_EMAIL });
  83  |     await page.click(SELECTORS.FARM_SUBMIT_BTN);
  84  | 
  85  |     // Verify farm appears in table
  86  |     await expect(page.locator(`[data-e2e="farm-row-${TEST_FARM_NAME}"]`)).toBeVisible({ timeout: 10000 });
  87  |     
  88  |     // Logout
  89  |     await page.click('[title="Logout"]');
  90  |     await expect(page).toHaveURL(/\/login/);
  91  |   });
  92  | 
  93  |   test('Persona: Farm Manager - Rights & Restrictions', async ({ page }) => {
  94  |     // 1. Login as the newly created Manager
  95  |     await page.goto('/login');
  96  |     await page.fill('input[name="identifier"]', TEST_MANAGER_EMAIL);
  97  |     await page.fill('input[name="password"]', TEST_MANAGER_PASS);
  98  |     await page.click(SELECTORS.LOGIN_SUBMIT);
  99  | 
> 100 |     await expect(page).toHaveURL(/\/dashboard\/users/, { timeout: 15000 });
      |                        ^ Error: expect(page).toHaveURL(expected) failed
  101 |     await page.waitForSelector(SELECTORS.AUTH_LOADER, { state: 'hidden', timeout: 10000 });
  102 | 
  103 |     // 2. Verify own farm visibility
  104 |     await page.goto('/dashboard/farm-ops/registry');
  105 |     await page.waitForSelector(SELECTORS.AUTH_LOADER, { state: 'hidden' });
  106 |     const farmRow = page.locator(`[data-e2e="farm-row-${TEST_FARM_NAME}"]`);
  107 |     await expect(farmRow).toBeVisible({ timeout: 10000 });
  108 |     
  109 |     // 3. Verify Edit is possible
  110 |     await expect(farmRow.locator('button >> internal:has-text="edit"')).toBeEnabled();
  111 | 
  112 |     // 4. Test Restriction: DELETE must fail
  113 |     page.on('dialog', async dialog => {
  114 |       console.log('Dialog detected:', dialog.message());
  115 |       await dialog.accept();
  116 |     });
  117 | 
  118 |     // Click delete
  119 |     await page.click(`[data-e2e="farm-row-${TEST_FARM_NAME}"] ${SELECTORS.FARM_DELETE_BTN}`);
  120 |     
  121 |     // Check if farm is STILL there (backend forbidden enforcement)
  122 |     // After the failure alert, the farm should still be visible in the list
  123 |     await expect(page.locator(`[data-e2e="farm-row-${TEST_FARM_NAME}"]`)).toBeVisible({ timeout: 5000 });
  124 |   });
  125 | 
  126 |   test('Persona: Anonymous - Showcase Only', async ({ page }) => {
  127 |     await page.goto('/');
  128 |     await expect(page.locator('h1:has-text("System")')).toBeVisible();
  129 |     // Sidebar should be HIDDEN for anonymous at root
  130 |     await expect(page.locator('aside')).not.toBeVisible();
  131 |   });
  132 | });
  133 | 
```