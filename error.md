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

> **Ghi nhớ**: Lỗi này xảy ra với **mọi tool** cài bằng `go install`, không chỉ mockery. Fix PATH 1 lần là xong cho terminal. Riêng IDE có thể cần symlink hoặc restart.

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
