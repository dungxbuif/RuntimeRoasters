# [RR-18] Farm Control Plane (UI)

- **Summary:** Xây dựng giao diện quản trị Nông trại trên ứng dụng Web (Client App).
- **Priority:** `MEDIUM`
- **Type:** Feature

---

## 📖 User Story
> As a farm manager, I want a user-friendly dashboard to view and manage my farms visually so that I can easily keep track of my assets without using technical tools.

## 💰 Business Value
Cải thiện trải nghiệm người dùng (UX) và giảm rào cản kỹ thuật cho nhân viên vận hành. Tăng khả năng quan sát (visibility) toàn bộ hệ thống nông trại.

## 🔍 Acceptance Criteria

### Scenario 1: Farm Listing View
- **Given:** A logged-in farm manager.
- **When:** Navigating to the "Farms" menu.
- **Then:** The system displays a table/list of all farms they own.
- **And:** Each row shows Name, Location, and Area.

### Scenario 2: Farm Creation Form
- **Given:** The Farm Listing page.
- **When:** Clicking "Create New Farm".
- **Then:** A form appears with fields for Name, Location, Area, and Type.
- **And:** Submitting the form with valid data adds the farm to the list.

### Scenario 3: Real-time Feedback
- **Given:** A form submission.
- **When:** The API returns success or failure.
- **Then:** The UI shows a clear notification (Toast) to the user.

### Scenario 4: Access Control in UI
- **Given:** A user with `guest` role.
- **When:** Attempting to access the Farm management page.
- **Then:** The UI redirects them to a "Permission Denied" page or hides the menu.
