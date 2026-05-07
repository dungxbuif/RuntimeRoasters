package identity

import "context"

/*
Cơ chế tránh xung đột Key (contextKey):
Trong Go, context.Value sử dụng một interface{} làm key.
Nếu dùng một chuỗi (string) như "user_id", các package khác cũng có thể dùng chung chuỗi đó, dẫn đến việc ghi đè dữ liệu của nhau (collision).
Bằng cách định nghĩa một kiểu dữ liệu mới (contextKey) và không export => đảm bảo rằng chỉ có package identity mới có thể truy cập vào vùng dữ liệu đó trong Context.
struct{} chiếm 0 byte trong bộ nhớ => cực kỳ tối ưu về mặt hiệu năng.
*/
type contextKey struct {}

/* Tự handle set identity context từ middleware */
func InjectContext(ctx context.Context, c Claims) context.Context {
	return context.WithValue(ctx, contextKey{}, c)
}

/* Lấy identity context */
func FromContext(ctx context.Context) (Claims, bool) {
	c, ok := ctx.Value(contextKey{}).(Claims)
	return c, ok
}

/* MustFromContext lấy identity context hoặc panic. Cho các case required Auth */
func MustFromContext(ctx context.Context) Claims {
	c, ok := ctx.Value(contextKey{}).(Claims)
	if !ok {
		panic("identity: no identity in context — ensure the Auth Middleware is applied")
	}
	return c
}