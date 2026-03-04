# 🐛 Error Log

## Lỗi 1: mockery not found in $PATH

### Triệu chứng
```
exec: "mockery": executable file not found in $PATH
```

### Nguyên nhân
`go install` cài binary vào `~/go/bin/`, nhưng thư mục này **chưa có trong PATH**.

### Cách fix

**Fix cho terminal** — thêm dòng này vào `~/.zshrc`:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```
Sau đó chạy `source ~/.zshrc` hoặc mở terminal mới.

**Fix cho IDE** — IDE không đọc `~/.zshrc`, cần tạo symlink:
```bash
sudo ln -s /home/sen/go/bin/mockery /usr/local/bin/mockery
```

---

## Lỗi 2: mockery v3 không dùng flag `--name`

### Triệu chứng
```
Error: unknown flag: --name
```

### Nguyên nhân
Mockery v3 bỏ flag `--name`, dùng file config `.mockery.yml` thay thế.

### Cách fix
Hạ về v2:
```bash
go install github.com/vektra/mockery/v2@latest
```

Cú pháp v2 trong `go:generate`:
```go
//go:generate mockery --name Password --filename pass_service.go
```

---

## Lỗi 3: File .env không hoạt động

### Triệu chứng
Server luôn dùng port default thay vì port trong `.env`.

### Nguyên nhân
Thư viện `envconfig` (kelseyhightower) chỉ đọc **biến môi trường** của OS, **không đọc file `.env`**.

### Cách fix
Dùng `godotenv` để load file `.env` trước:
```go
import "github.com/joho/godotenv"

func main() {
    godotenv.Load()  // load .env vào OS environment
    cfg, err := api.NewConfig("")  // envconfig giờ đọc được
}
```

### Lưu ý thêm
- `export APP_PORT=8081` ưu tiên hơn `.env` — vì `godotenv` không ghi đè biến đã tồn tại
- `export APP_PORT=` (rỗng) ≠ `unset APP_PORT` (xoá) — rỗng khiến server chạy port `:`
- Xoá biến đã export: `unset APP_PORT`

---

## Lỗi 4: swag.Spec — unknown field LeftDelim/RightDelim

### Triệu chứng
```
docs/docs.go:60:2: unknown field LeftDelim in struct literal of type "github.com/swaggo/swag".Spec
docs/docs.go:61:2: unknown field RightDelim in struct literal of type "github.com/swaggo/swag".Spec
```

### Nguyên nhân
`swag` CLI (trên máy) và `swag` library (trong go.mod) **khác version**. CLI mới tạo field `LeftDelim`/`RightDelim`, nhưng library cũ không có.

### Cách fix
```bash
go get -u github.com/swaggo/swag
go mod tidy
swag init -g cmd/api/main.go
```

---

## Lỗi 5: Swagger import sai

### Triệu chứng
```
undefined: swaggerfiles
```

### Nguyên nhân
Import sai tên package và thiếu alias.

### Cách fix
Import đúng:
```go
import (
    _ "github.com/senn404/bookmark-managent/docs"       // swagger docs (chạy init)
    swaggerFiles "github.com/swaggo/files"               // alias: swaggerFiles
    ginSwagger "github.com/swaggo/gin-swagger"            // alias: ginSwagger
)
```

Sử dụng:
```go
a.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
//                                                 ↑ chữ F hoa
```

Cài package:
```bash
go get github.com/swaggo/files
go get github.com/swaggo/gin-swagger
```

---

## Lỗi 6: Gin không hiển thị request log

### Triệu chứng
Server chạy nhưng không hiển thị log khi có request đến.

### Nguyên nhân
`gin.New()` tạo server **trống**, không có middleware nào (kể cả Logger).

### Cách fix
Đổi sang `gin.Default()`:
```go
app: gin.Default()   // có Logger + Recovery middleware
```

| | `gin.New()` | `gin.Default()` |
|---|---|---|
| Logger | ❌ | ✅ |
| Recovery | ❌ | ✅ |

---

## Lỗi 7: Response không trả JSON

### Triệu chứng
Response trả về text thuần thay vì JSON.

### Nguyên nhân
Dùng `c.String()` thay vì `c.JSON()`.

### Cách fix
```go
// Trước — text
c.String(http.StatusOK, pass)
// → Qf8@CZd&xK2m

// Sau — JSON
c.JSON(http.StatusOK, gin.H{"pass": pass})
// → {"pass": "Qf8@CZd&xK2m"}
```

---

## Lỗi 8: Handler thiếu return sau error response

### Triệu chứng
Khi service trả lỗi, handler gửi response lỗi nhưng vẫn tiếp tục chạy → gửi thêm response thành công → lỗi double response.

### Cách fix
Thêm `return` sau error response:
```go
if err != nil {
    c.String(http.StatusInternalServerError, "err")
    return    // ← PHẢI có return ở đây
}
```
