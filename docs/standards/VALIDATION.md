# Validation Standard

This file defines validation policy. Runtime proof state belongs in `docs/work/VALIDATION_MATRIX.md`.

## Principle

Every accepted behavior needs proof. The proof can be automated, manual, or explicitly not required with a reason, but it must be recorded.

## Proof Types

| Proof Type | Use When | Evidence |
| --- | --- | --- |
| Unit | Pure domain logic, helpers, parsing, validation, isolated business rules | Test command and result |
| Integration | Backend enforcement, persistence, provider behavior, jobs, service contracts | Test command, fixture/setup notes, result |
| E2E | User-visible browser/app flow | Browser/app test command or manual steps |
| UAT | User-facing acceptance criteria or workflow sign-off | Expected behavior, verified behavior, sign-off or not-required reason |
| Platform/manual | Runtime, shell, deployment, mobile, desktop, environment-specific behavior | Manual steps, screenshots/log summary, or command output |
| Docs review | Behavior, contract, architecture, data, or decision changes | Docs review checklist result |

## Required Proof By Work Type

| Work Type | Required Proof |
| --- | --- |
| Tiny docs-only | Docs review or not-needed reason |
| Feature/ticket | Unit or integration where applicable, UAT for user-facing behavior, docs review |
| Bug | Reproduction evidence, regression proof, docs review |
| API/contract | Integration or contract proof, docs review, UAT if user-visible |
| Data/schema | Migration/schema-sensitive proof, docs review, ADR when ownership/tradeoff changes |
| Security/auth | Integration or E2E proof, negative cases, docs review, ADR when model changes |
| Deployment/runtime | Platform/manual proof, rollback or release check, docs review |
| Standards/framework | Docs review and trace links |

## Evidence Rules

- Record the exact command when a command was run.
- Record `pass`, `fail`, or `skipped`.
- If skipped, record why and what residual risk remains.
- Do not mark behavior `implemented` in `docs/work/VALIDATION_MATRIX.md` without evidence.
- If validation requirements change, update the matrix and create an ADR when the change is durable or risky.

## Project-Specific Unit Test Rules

# 🧪 Unit Test Agent Guardrails: Robust Backend & Frontend Test Cases

Tài liệu này đóng vai trò là **nguyên tắc và luật bất di bất dịch (Guardrails)** dành cho AI Agent khi thiết kế, sinh mã (generate) hoặc cập nhật Unit Tests trong repository **Runtime Roasters** (bao gồm Go backend và TypeScript/Next.js frontend). 

Mục tiêu tối thượng của Unit Test không phải là đạt 100% coverage một cách vô nghĩa bằng các happy path đơn giản, mà là **bắt được bug trước khi user hoặc hệ thống phát hiện ra**.

---

## 🎯 Checklist nhanh cho Agent khi viết Unit Test
Trước khi kết thúc lượt code hoặc tự đánh giá bài test đã đạt yêu cầu chưa, Agent **BẮT BUỘC** phải rà soát danh sách sau:
- [ ] **Happy path** (Case lý tưởng hoạt động hoàn hảo).
- [ ] **Edge cases** (Mảng rỗng, giá trị `0`, rỗng, cực đại/cực tiểu).
- [ ] **Boundary cases** (Giá trị ranh giới `min-1`, `min`, `min+1`, `max-1`, `max`, `max+1`).
- [ ] **Invalid input** (Thiếu trường bắt buộc, sai định dạng, sai kiểu dữ liệu).
- [ ] **Error cases & Dependency Failures** (Mocking repository trả lỗi, timeout, context cancel, gRPC/HTTP failure).
- [ ] **Nil / Empty cases** (Nil pointer, nil slice, empty map, optional field bị khuyết).
- [ ] **Concurrency-related cases** (Data race, shared state, goroutine leaks - *nếu có sử dụng concurrency*).
- [ ] **Idempotency cases** (Gọi lặp lại nhiều lần cùng một tham số).
- [ ] **Regression cases** (Kiểm thử các lỗi lịch sử đã được fix).
- [ ] **Security-related logic cases** (Rào chắn phân quyền ở mức nghiệp vụ, token sai, truy cập tài nguyên không thuộc sở hữu).

---

## 🧱 Chi tiết 10 Loại Test Cases Bắt Buộc

### 1️⃣ Happy path – Case lý tưởng
✔ **Input hợp lệ:** Dữ liệu hoàn hảo, đầy đủ các trường bắt buộc và tùy chọn.  
✔ **Điều kiện đầy đủ:** Hệ thống ngoại vi (DB, Cache) trả kết quả thành công.  
✔ **Kết quả đúng như mong đợi:** Trả về mã thành công (`200 OK`, `nil error`) kèm payload chính xác.

> [!NOTE]
> **Mục đích:** Xác nhận luồng logic chính hoạt động mượt mà và làm baseline (điểm so sánh chuẩn) cho các trường hợp đặc biệt khác.

#### 💡 Ví dụ minh họa (Go):
```go
func TestCalculateRoastLoss_HappyPath(t *testing.T) {
    // Input hợp lệ: Hạt Arabica, độ ẩm chuẩn, hao hụt thông thường
    inputWeight := 100.0
    outputWeight := 85.0
    
    loss, err := RoastLossCalculator(inputWeight, outputWeight)
    
    assert.NoError(t, err)
    assert.Equal(t, 15.0, loss)
}
```

---

### 2️⃣ Edge cases – Các trường hợp biên & dị biệt
Edge case là nơi các bug logic tiềm ẩn xuất hiện nhiều nhất do lập trình viên thường chỉ tập trung vào dữ liệu đẹp.

*   **List/Slice:** Rỗng (`len == 0`), chỉ có đúng 1 phần tử (`len == 1`).
*   **Giá trị số:** Bằng `0`, số âm, giá trị cực đại (`math.MaxInt`), cực tiểu.
*   **String:** String rỗng (`""`), string chỉ có khoảng trắng (`"   "`), chuỗi ký tự unicode/UTF-8 phức tạp.

> [!IMPORTANT]
> Agent phải luôn giả định rằng người dùng hoặc hệ thống khác sẽ truyền vào những giá trị "kỳ quặc" nhất có thể.

#### 💡 Ví dụ minh họa (Go):
```go
func TestProcessBatch_EmptySlice(t *testing.T) {
    var emptyBatch []Harvest
    result, err := ProcessHarvestBatch(emptyBatch)
    
    assert.NoError(t, err) // Không crash panic
    assert.Len(t, result, 0) // Trả về slice rỗng an toàn
}
```

---

### 3️⃣ Boundary cases – Ranh giới giữa Đúng và Sai
Boundary case là điểm giao thoa nhạy cảm của các biểu thức so sánh (`<`, `<=`, `>`, `>=`).

*   **Độ dài ký tự (e.g. Min length = 8):** Test độ dài `7` (sai), `8` (đúng), `9` (đúng).
*   **Hạn mức số lượng (e.g. Max items = 100):** Test `99`, `100`, `101`.
*   **Phân trang (Pagination):** Limit bằng `0`, bằng `1`, vượt quá kích thước tối đa của trang.

> [!TIP]
> Hãy luôn kiểm tra sát sườn các điểm giới hạn thay vì chỉ test ngẫu nhiên một con số nằm giữa khoảng cho phép.

#### 💡 Ví dụ minh họa (Go):
```go
func TestValidatePasswordLength_Boundaries(t *testing.T) {
    tests := []struct {
        name     string
        password string
        isValid  bool
    }{
        {"Below boundary (7 chars)", "1234567", false},
        {"At boundary (8 chars)",    "12345678", true},
        {"Above boundary (9 chars)", "123456789", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidatePassword(tt.password)
            if tt.isValid {
                assert.NoError(t, err)
            } else {
                assert.Error(t, err)
            }
        })
    }
}
```

---

### 4️⃣ Invalid input – Dữ liệu không hợp lệ
Backend xuất sắc không chỉ xử lý đúng mà phải từ chối cái sai một cách nhanh chóng, minh bạch và an toàn (Fail-fast).

*   **Thiếu trường bắt buộc:** Payload thiếu ID, thiếu email, v.v.
*   **Sai định dạng (Format):** Định dạng email không hợp lệ, UUID sai chuẩn, định dạng ngày tháng không đúng.
*   **Kiểu dữ liệu (Type):** Truyền string vào trường số, truyền float vào trường integer.
*   **Logic nghiệp vụ:** Giá trị số tiền thanh toán hoặc số lượng hàng hóa là số âm.

#### 💡 Ví dụ minh họa (Go):
```go
func TestCreateOrder_InvalidInputs(t *testing.T) {
    tests := []struct {
        name  string
        req   CreateOrderRequest
        errID string
    }{
        {"Negative quantity", CreateOrderRequest{Quantity: -5, Price: 100}, "ERR_INVALID_QUANTITY"},
        {"Zero price",        CreateOrderRequest{Quantity: 2, Price: 0},     "ERR_INVALID_PRICE"},
        {"Missing ProductID", CreateOrderRequest{Quantity: 2, Price: 100},   "ERR_MISSING_PRODUCT_ID"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := CreateOrder(context.Background(), tt.req)
            assert.Error(t, err)
            assert.Contains(t, err.Error(), tt.errID)
        })
    }
}
```

---

### 5️⃣ Error cases – Khi Dependency thất bại
Hệ thống microservices của Runtime Roasters phụ thuộc vào nhiều thành phần (DB, Kafka, gRPC Client, Redis). Unit test của một service bắt buộc phải mô phỏng (Mock) các kịch bản lỗi từ các dependency này để đảm bảo service xử lý lỗi trơn tru.

*   **Database error:** Lỗi kết nối, lỗi khóa ngoại (Foreign key violation), khóa trùng lặp.
*   **Cache miss & Cache error:** Redis bị sập đột ngột.
*   **Context cancellation:** Lần theo vòng đời request khi client ngắt kết nối giữa chừng hoặc quá thời gian chờ (Timeout).
*   **API/gRPC Call Failures:** Các service hạ tầng phía sau bị lỗi `500` hoặc offline.

> [!WARNING]
> Cần đảm bảo hệ thống có cơ chế rollback giao dịch (Transaction Rollback) hoặc dọn dẹp (Cleanup) tài nguyên khi có lỗi xảy ra giữa chừng để tránh rò rỉ dữ liệu.

#### 💡 Ví dụ minh họa (Go - Mocking):
```go
func TestGetFarm_RepositoryError(t *testing.T) {
    mockRepo := new(MockFarmRepository)
    // Giả lập Repository quăng ra lỗi kết nối Database
    mockRepo.On("FindByID", "farm-123").Return(nil, errors.New("db connection timed out"))
    
    service := NewFarmService(mockRepo)
    farm, err := service.GetFarm(context.Background(), "farm-123")
    
    assert.Nil(t, farm)
    assert.Error(t, err)
    assert.Equal(t, "db connection timed out", err.Error())
}
```

---

### 6️⃣ Nil / Empty cases – Giá trị rỗng và con trỏ rỗng
Trong ngôn ngữ Go, `nil pointer dereference` là nguyên nhân hàng đầu gây ra lỗi hệ thống nghiêm trọng (Panic sập server).

*   **Con trỏ Nil:** Truyền `nil` thay vì địa chỉ vùng nhớ của Struct.
*   **Slices/Maps rỗng hoặc Nil:** Gọi các hàm ghi/đọc trên Slice hoặc Map chưa được khởi tạo (`make`).
*   **Optional fields:** Các tham số tùy chọn không được truyền (phải kiểm tra tính an toàn trước khi sử dụng).

#### 💡 Ví dụ minh họa (Go):
```go
func TestProcessProfile_NilPointerSafety(t *testing.T) {
    // Tuyệt đối không để xảy ra panic khi truyền nil pointer
    var nilProfile *UserProfile = nil
    
    err := UpdateUserProfile(nilProfile)
    
    assert.Error(t, err)
    assert.Equal(t, ErrProfileCannotBeNil, err)
}
```

---

### 7️⃣ Concurrency-related cases – Xử lý bất đồng bộ & Song song
*Nếu logic của bạn có sử dụng goroutines, channels, mutexes, sync.Map, hoặc shared state:*

*   **Race Conditions:** Nhiều thread/goroutine đọc ghi cùng lúc trên một tài nguyên.
*   **Deadlocks:** Các goroutine bị chặn lẫn nhau mãi mãi do cơ chế khóa tài nguyên không hợp lý.
*   **Goroutine Leaks:** Goroutine chạy ngầm bị treo và không bao giờ thoát giải phóng bộ nhớ.

> [!TIP]
> Hãy tận dụng cờ `-race` khi chạy Go test để phát hiện sớm các lỗi xung đột dữ liệu bất đồng bộ.

---

### 8️⃣ Idempotency cases – Tính an toàn khi gọi lặp
Một HTTP REST API hoặc sự kiện Kafka có thể bị kích hoạt/gọi lại nhiều lần do cơ chế retry tự động của hệ thống mạng hoặc lỗi trùng lặp sự kiện (Network duplicate/retry).

*   **Không nhân bản dữ liệu:** Gọi API tạo đơn hàng 2 lần với cùng một `Idempotency-Key` thì đơn hàng thứ hai không được tạo mới, mà chỉ trả về kết quả của đơn hàng đầu tiên.
*   **Không tạo Side Effects trùng lặp:** Không gửi 2 email xác nhận, không trừ tiền 2 lần.

#### 💡 Ví dụ minh họa (Go):
```go
func TestSubmitPayment_Idempotent(t *testing.T) {
    mockPaymentService := NewPaymentService(mockRepo)
    req := PaymentRequest{
        IdempotencyKey: "unique-key-123",
        Amount:         100.0,
    }
    
    // Gọi lần 1: Tạo mới thành công
    res1, err1 := mockPaymentService.Process(req)
    assert.NoError(t, err1)
    assert.Equal(t, "SUCCESS", res1.Status)
    
    // Gọi lần 2: Trả về kết quả cũ, không trừ tiền thêm
    res2, err2 := mockPaymentService.Process(req)
    assert.NoError(t, err2)
    assert.Equal(t, res1.TransactionID, res2.TransactionID)
}
```

---

### 9️⃣ Regression cases – Chặn đứng bug lịch sử tái diễn
Có một nguyên tắc vàng trong phát triển phần mềm chuyên nghiệp: **Bug nào đã từng được sửa thì BẮT BUỘC phải có Unit Test đi kèm để nó không bao giờ quay lại.**

*   Khi fix một bug thực tế (từ QC báo hay production log), hãy viết ngay một Unit Test tái hiện đúng tình trạng bug đó trước khi tiến hành code sửa đổi.
*   Test case này sẽ bảo vệ code khỏi bị "vỡ" khi có các đợt refactor lớn trong tương lai.

---

### 🔟 Security-related cases – Bảo mật mức Logic Nghiệp Vụ
Không cần phải chạy test thâm nhập (Penetration Test) ở cấp Unit Test, nhưng logic phân quyền và sở hữu dữ liệu phải được bọc chặt chẽ:

*   **Unauthorized Access:** Gửi token sai, token hết hạn, hoặc không gửi token thì hệ thống phải từ chối.
*   **Resource Ownership:** User A có quyền `read` nhưng không được phép đọc/sửa dữ liệu thuộc sở hữu của User B.
*   **Role Scoping:** Quyền hạn phải khớp chặt chẽ với cấu hình (Ví dụ: `FARM_MANAGER` chỉ quản lý nông trại được gán, không được truy cập toàn hệ thống như `ADMIN`).

#### 💡 Ví dụ minh họa (Go với Casbin/JWT scoping):
```go
func TestGetFarmDetails_AccessControl(t *testing.T) {
    // Cấu hình Subject "user-A" (Manager của Farm A) muốn truy cập "Farm B"
    ctx := context.WithValue(context.Background(), "user_id", "user-A")
    ctx = context.WithValue(ctx, "user_roles", []string{"FARM_MANAGER"})
    
    service := NewFarmService(mockRepo, mockCasbinEnforcer)
    farm, err := service.GetFarmDetails(ctx, "farm-B") // Cố tình truy cập Farm B
    
    assert.Nil(t, farm)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "permission denied")
}
```

---

## 🛠️ Quy định Thực thi cho Agents

1.  **Đọc tệp tin này trước khi viết code:** Agent luôn kiểm tra sự tồn tại của [UNIT_TEST_RULES.md](./UNIT_TEST_RULES.md) và coi đây là tiêu chuẩn cao nhất cho chất lượng Test.
2.  **Rà soát file đích:** Khi tạo hoặc cập nhật file test `*_test.go` hoặc `*.spec.ts`, hãy cố gắng cấu trúc bảng dữ liệu đầu vào (Table-driven Tests) để bao phủ đầy đủ ít nhất 5 lớp đầu tiên (Happy Path, Edge, Boundary, Invalid Input, Error cases).
3.  **Chạy lệnh verify cục bộ để kiểm tra lỗi:**
    ```bash
    go test -v -race ./...
    ```
