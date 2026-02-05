/*
Ví dụ Decorator trong Go: Notifier (Email → SMS → Slack → Facebook)
Bài toán
# Ban đầu chỉ gửi thông báo qua Email
# Sau này muốn gửi thông báo thêm qua SMS, Slack, Facebook
Có thể bật/tắt, kết hợp linh hoạt
# Không sửa code EmailNotifier gốc
👉 Decorator là lựa chọn chuẩn bài.
*/
package main

// 1. Component
type Notifier interface {
	Send(message string)
}

// 2. ConcreteComponent – Email
type EmailNotifier struct{}

func (e *EmailNotifier) Send(message string) {
	println("Gửi Email:", message)
}

// 3. Base Decorator
type NotifierDecorator struct {
	notifier Notifier
}

// 4. ConcreteDecorator – SMS
type SMSDecorator struct {
	NotifierDecorator
}

func NewSMSDecorator(n Notifier) Notifier {
	return &SMSDecorator{
		NotifierDecorator{notifier: n},
	}
}

func (s *SMSDecorator) Send(message string) {
	s.notifier.Send(message)
	println("Gửi SMS:", message)
}

// 5. Concrete Decorator - Facebook
type FacebookDecorator struct {
	NotifierDecorator
}

func NewFacebookDecorator(n Notifier) Notifier {
	return &FacebookDecorator{
		NotifierDecorator{notifier: n},
	}
}

func (s *FacebookDecorator) Send(message string) {
	s.notifier.Send(message)
	println("Gửi Facebook:", message)
}

// 6. ConcreteDecorator – Slack
type SlackDecorator struct {
	NotifierDecorator
}

func NewSlackDecorator(n Notifier) Notifier {
	return &SlackDecorator{
		NotifierDecorator{notifier: n},
	}
}

func (s *SlackDecorator) Send(message string) {
	s.notifier.Send(message)
	println("Gửi Slack:", message)
}

// 7. Client sử dụng
func main() {
	// Vào trước thì gửi trước: SMS -> Slack -> Facebook
	var notifier Notifier = &EmailNotifier{}

	notifier = NewSMSDecorator(notifier)
	notifier = NewSlackDecorator(notifier)
	notifier = NewFacebookDecorator(notifier)

	notifier.Send("Server down!")
}

/*
OUTPUT:
Gửi Email: Server down!
Gửi SMS: Server down!
Gửi Slack: Server down!
Gửi Facebook: Server down!
*/
