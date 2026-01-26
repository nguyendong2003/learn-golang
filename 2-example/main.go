/*
  - fmt.Printf():
    %T :  type
    %v :  value		(%v chỉ in ra giá trị của các giá trị của field trong struct)
    %+v: dùng để in giá trị chi tiết hơn so với %v (%+v in cả tên field + giá trị của field trong struct)
    %q :  In ra string hoặc rune kèm dấu " hoặc '
    %d : so nguyen
    %f : so thuc
    %.2d : float 2 chữ số thập phân
    %c : In ra 1 ký tự (character)
    %s: string
    %p : địa chỉ con trỏ
    %g : format verb dùng cho số thực (float)

- fmt.Sprint: hàm trong Go dùng để ghép (convert) giá trị thành chuỗi (string), KHÔNG in ra màn hình.  (trả về string)
*/
package main

import (
	"fmt"
	"learngo/0-package/product"
	structembedding11 "learngo/11-struct-embedding"
	structembedding12 "learngo/12-struct-embedding"
	profitrevenue "learngo/13-profit-revenue"
	"math"
	"math/cmplx"
	"runtime"
)

func add(x int, y int) int {
	return x + y
}

func sub(x, y int) (int, int) {
	return x, y
}

func swap(x, y string) (string, string) {
	return y, x
}

func split(sum int) (x, y int) {
	x = sum + 1
	y = sum - x
	return
}

func variables() {
	var x, y bool
	var i, j int = 1, 2
	var i1, j1 string = "hello", "world"
	var a, b, c = true, false, "no"
	//
	k1 := 3
	k2, k3 := "hello", "world"
	//
	var (
		tobe   bool       = false
		maxInt uint64     = 1<<64 - 1
		z      complex128 = cmplx.Sqrt(-5 + 12i)
		s      string     = "messi"
		t      string
	)

	fmt.Println(x, y, i, j, i1, j1, a, b, c, k1, k2, k3)
	fmt.Printf("Type: %T, Value: %v\n", tobe, tobe)
	fmt.Printf("Type: %T, Value: %v\n", maxInt, maxInt)
	fmt.Printf("Type: %T, Value: %v\n", z, z)
	fmt.Printf("Type: %T, Value: %q\n", s, s)
	fmt.Printf("Type: %T, Value: %q\n", t, t)
}

func typeConversion() {
	var i int = 42
	var f float64 = float64(i)
	var u uint = uint(f)

	i2 := 42
	f2 := float64(i2)
	u2 := uint(f2)

	t := i2

	fmt.Println(i, f, u, i2, f2, u2, t)
}

func constValue() {
	const a = "hello"
	const (
		b = 1
		c = "xin"
		d = true
		f = 3.142
	)

	// Big không bị tràn số vì Big là hằng số (constant) chưa có kiểu – untyped constant
	const (
		Big   = 1 << 100  // 2^100
		Small = Big >> 99 // 2
	)

	m := float64(Big)

	fmt.Println(a, b, c, d, f, m, Small)
}

func calSum() int {
	n := 10
	sum := 0

	for i := 1; i <= n; i++ {
		sum += i
	}

	return sum
}

func switchEx1() {
	fmt.Print("Go runs on ")
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("macOS.")
	case "linux":
		fmt.Println("Linux.")
	default:
		// freebsd, openbsd,
		// plan9, windows...
		fmt.Printf("%s.\n", os)
	}
}

func deferEx1() {
	fmt.Println("counting")

	for i := 0; i < 10; i++ {
		defer fmt.Print(i, " ")
	}

	fmt.Println("done")
}

func pointerEx1() {
	i, j := 30, 40

	p := &i
	fmt.Println(i, &i)
	fmt.Println(p, &p, *p)
	*p = 35
	fmt.Println(p, &p, *p)

	p = &j
	*p = *p + 1
	fmt.Println(p, &p, *p)
}

type Vertex struct {
	X int
	Y int
}

func vertexEx1() {
	v := Vertex{1, 2}
	fmt.Println(v) // {1,2}

	v.X = 4
	fmt.Println(v)        // {4,2}
	fmt.Println(v.X, v.Y) // 4 2
}

func vertexEx2() {
	v := Vertex{1, 2}
	p := &v            // p: Pointer to struct v
	var t *Vertex = &v // t: Pointer to struct v (like p)

	p.X = 1e9  // cach 1
	(*p).Y = 5 // cach 2

	fmt.Println(v)  // {1000000000 5}
	fmt.Println(*p) // {1000000000 5}
	fmt.Println(*t) // {1000000000 5}

	(*t).X = 6
	t.Y = 5

	fmt.Println(v)  // {6 5}
	fmt.Println(*p) // {6 5}
	fmt.Println(*t) // {6 5}
}

func vertexEx3() {
	// Struct Literals
	var (
		v1 = Vertex{1, 2}  // has type Vertex
		v2 = Vertex{X: 1}  // Y:0 is implicit
		v3 = Vertex{}      // X:0, Y:0
		p  = &Vertex{1, 2} // has type *Vertex
	)

	fmt.Println(v1, v2, v3, p)
}

func arrayEx1() {
	var a [2]string
	a[0] = "Hello"
	a[1] = "World"
	fmt.Println(a[0], a[1])
	fmt.Println(a)
	fmt.Printf("%T\n", a)

	primes := [6]int{1, 2, 3, 4, 5}
	fmt.Println(primes)
}

func sliceEx1() {
	primes := [6]int{1, 2, 3, 4, 5, 6}

	var s []int = primes[1:4]
	var t []int = primes[2:6]

	fmt.Println(primes) // [1 2 3 4 5 6]
	fmt.Println(s)      // [2 3 4]
	fmt.Println(t)      // [3 4 5 6]
	s[2] = 9            // thay đổi phần tử giá trị 4 thành 9
	fmt.Println(primes) // [1 2 3 9 5 6]
	fmt.Println(s)      // [2 3 9]
	fmt.Println(t)      // [3 9 5 6]

	primes[3] = 99
	fmt.Println(primes) // [1 2 3 99 5 6]
	fmt.Println(s)      // [2 3 99]
	fmt.Println(t)      // [3 99 5 6]
}

func sliceEx2() {
	// Bước 1: Khởi tạo Slice
	data := [5]int{10, 20, 30, 40, 50}
	s1 := data[1:4] // Lấy từ index 1 đến 3

	fmt.Println("Buoc 1: ", s1, len(s1), cap(s1))
	fmt.Println("Buoc 1: ", data, len(data), cap(data))

	// Bước 2: Cắt tiếp từ slice (Reslicing)
	s2 := s1[1:3] // Lấy từ index 1 đến 2 của s1
	fmt.Println("Buoc 2: ", s2, len(s2), cap(s2))

	// Bước 3: Trường hợp 1 -> Append khi còn Capacity
	s2 = append(s2, 100)
	fmt.Println("Buoc 3: ", s2, len(s2), cap(s2))

	// Bước 4: Trường hợp 2 -> Append khi hết Capacity
	s2 = append(s2, 200)

	fmt.Println("Buoc 4: ", s2, len(s2), cap(s2))

	fmt.Println(data) // [10 20 30 40 100]
	s2[0] = 9999
	fmt.Println("Buoc 4:", s2, len(s2), cap(s2)) // [9999 40 100 200]
	fmt.Println(data)                            // [10 20 30 40 100]
}

func sliceEx3() {
	// Khởi tạo slice có len=3 nhưng cap=10
	a := make([]int, 3, 10)
	a[0], a[1], a[2] = 1, 2, 3

	fmt.Println(a) // [1 2 3]

	// a còn dư chỗ, append không tạo mảng mới mà dùng chung mảng của a
	b := append(a, 4) // b ghi số 4 vào vị trí index 3 của mảng ẩn
	c := append(a, 5) // c ghi đè số 5 vào đúng vị trí index 3 đó!

	fmt.Println(a) // [1 2 3]
	fmt.Println(b) // [1 2 3 5] -> Số 4 đã bị biến thành số 5!
	fmt.Println(c) // [1 2 3 5]
}

func makeCopyFullSliceEx1() {
	s := make([]int, 3, 5)
	fmt.Printf("%T\n", s) // []int
	fmt.Println(s, len(s), cap(s))

	p := append(s, 3)
	fmt.Println(p, len(p), cap(p))
	fmt.Printf("%T\n", p) // []int

	s = append(s, 4)
	s = append(s, 5)
	fmt.Println(s, len(s), cap(s))

	s = append(s, 6)
	fmt.Println(s, len(s), cap(s))
}

func makeCopyFullSliceEx2() {
	a := make([]int, 4, 5)
	fmt.Println(a, len(a), cap(a)) // [0 0 0 0] 4 5

	ptr := append(a, 4)
	fmt.Println(ptr, len(ptr), cap(ptr)) // [0 0 0 0 4] 5 5
	fmt.Println(a, len(a), cap(a))       // [0 0 0 0] 4 5

	ptr = append(a, 6)
	fmt.Println(ptr, len(ptr), cap(ptr)) // [0 0 0 0 6] 5 5
	fmt.Println(a, len(a), cap(a))       // [0 0 0 0] 4 5
}

func makeCopyFullSliceEx3() {
	a := [5]int{1, 2, 3, 4, 5}
	p := a[1:3]
	p = append(p, 6)
	fmt.Println(a, len(a), cap(a)) // [1 2 3 6 5] 5 5
	fmt.Println(p, len(p), cap(p)) // [2 3 6] 3 4

	t := append(p, 7)
	fmt.Println(a, len(a), cap(a)) // [1 2 3 6 7] 5 5
	fmt.Println(p, len(p), cap(p)) // [2 3 6] 3 4
	fmt.Println(t, len(t), cap(t)) // [2 3 6 7] 4 4

	u := append(t, 8)
	fmt.Println(a, len(a), cap(a)) // [1 2 3 6 8] 5 5
	fmt.Println(t, len(t), cap(t)) // [2 3 6] 3 4
	fmt.Println(u, len(u), cap(u)) // [2 3 6 7 8] 5 8
}

func makeCopyFullSliceEx4() {
	src := []int{1, 2, 3}
	dest := make([]int, len(src)) // Phải tạo dest có cùng len
	copy(dest, src)

	fmt.Println(src, len(src), cap(src))    // [1 2 3] 3 3
	fmt.Println(dest, len(dest), cap(dest)) // [1 2 3] 3 3

	src[0] = 99
	dest[0] = 55

	fmt.Println(src, len(src), cap(src))    // [99 2 3] 3 3
	fmt.Println(dest, len(dest), cap(dest)) // [55 2 3] 3 3
}

// append luôn trả về slice header mới, còn array có thể cũ hoặc mới.

// Ví dụ 1: Header mới – array CŨ (không realloc)
func appendEx1() {
	a := make([]int, 2, 4)
	a[0], a[1] = 1, 2

	b := append(a, 3)

	fmt.Println("a:", a, ", len: ", len(a), ", cap: ", cap(a))
	fmt.Println("b:", b, ", len: ", len(b), ", cap: ", cap(b))

	fmt.Printf("&a = %p\n", &a)
	fmt.Printf("&b = %p\n", &b)

	fmt.Printf("&a[0] = %p\n", &a[0])
	fmt.Printf("&b[0] = %p\n", &b[0])
}

// Ví dụ 2: Header mới – array MỚI (realloc)
func appendEx2() {
	a := make([]int, 2, 2)
	a[0], a[1] = 1, 2

	b := append(a, 3)

	fmt.Println("a:", a, ", len: ", len(a), ", cap: ", cap(a))
	fmt.Println("b:", b, ", len: ", len(b), ", cap: ", cap(b))

	fmt.Printf("&a = %p\n", &a)
	fmt.Printf("&b = %p\n", &b)

	fmt.Printf("&a[0] = %p\n", &a[0])
	fmt.Printf("&b[0] = %p\n", &b[0])
}

// Ví dụ 3: Header mới – array CŨ nhưng làm “đổi” dữ liệu gốc
func appendEx3() {
	a := make([]int, 2, 3)
	a[0], a[1] = 1, 2

	b := append(a, 99)

	fmt.Printf("&a[0] = %p\n", &a[0])
	fmt.Printf("&b[0] = %p\n", &b[0])

	fmt.Println("a:", a)
	fmt.Println("b:", b)
}

// Ví dụ 4: append mà KHÔNG gán lại → slice cũ không đổi
// func appendEx4(s []int) {
// 	append(s, 100)
// }

func rangeEx1() {
	arr := []int{1, 2, 3, 4, 5}
	for index, value := range arr {
		fmt.Println(index, value)
	}

	// omit index
	for _, value := range arr {
		fmt.Println(value)
	}

	// omit value
	for index, _ := range arr {
		fmt.Println(index)
	}

	// omit value
	for index := range arr {
		fmt.Println(index)
	}
}

// Khai bao map
func mapEx1() {
	// Cach 1
	var mp map[string]Vertex = make(map[string]Vertex)
	mp["Messi"] = Vertex{1, 2}
	mp["Ronaldo"] = Vertex{2, 3}
	fmt.Println(mp["Messi"])

	for key, value := range mp {
		fmt.Println(key, value)
	}

	// Cach 2
	mp2 := make(map[int]string)
	mp2[1] = "Neymar"
	mp2[5] = "LM10"
	mp2[5] = "Cong Phuong"

	fmt.Println(mp2[9]) // key = 9 khong co trong mp2 nen mp2[9] = zerod value
	for key, value := range mp2 {
		fmt.Println(key, value)
	}

	// Cach 3: Map literals
	var mp3 = map[Vertex]bool{
		{1, 2}: true,
		{2, 3}: false,
		{3, 4}: true,
	}

	fmt.Println(mp3[Vertex{2, 3}])

	for key, value := range mp3 {
		fmt.Println(key, value)
	}
}

func mapEx2() {
	mp := make(map[string]int)
	mp["Messi"] = 100
	fmt.Println(mp["Messi"])

	mp["Messi"] = 200
	fmt.Println(mp["Messi"])

	delete(mp, "Messi")
	fmt.Println(mp["Messi"])

	value, ok := mp["Messi"]
	fmt.Println(value, ok)

	a := 'h'
	fmt.Printf("%T %c %v\n", a, a, a)

	t := "messi"
	for _, value := range t {
		fmt.Printf("%T %c %v\n", value, value, value)
	}
}

// Function Example
func compute(fn func(float64, float64) float64) float64 {
	return fn(3, 4)
}

func calculator(a, b float64, fn func(a, b float64) float64) float64 {
	return fn(a, b)
}

func add1(a, b float64) float64 {
	return a + b
}

func sub1(a, b float64) float64 {
	return math.Abs(a - b)
}

//

func funcEx1() {
	// Function cũng là giá trị (value)
	func1 := func(x, y float64) float64 {
		return math.Sqrt(x*x + y*y)
	}
	fmt.Println(compute(func1))
	fmt.Println(compute(math.Pow))

	product := func(x, y float64) float64 {
		return x * y
	}

	divide := func(x, y float64) float64 {
		return x / y
	}

	fmt.Println(calculator(3, 4, add1))
	fmt.Println(calculator(3, 4, sub1))
	fmt.Println(calculator(3, 4, product))
	fmt.Println(calculator(3, 4, divide))
}

// Function Closure
func counter() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

func funcClosureEx1() {
	fmt.Println("Vi du 1:")
	c1 := counter()
	fmt.Println(c1()) // 1
	fmt.Println(c1()) // 2
	fmt.Println(c1()) // 3

	c2 := counter()
	fmt.Println(c2()) // 1
	fmt.Println(c2()) // 2

	////
	fmt.Println("Vi du 2:")
	counter := func() func() int {
		i := 0
		return func() int {
			i++
			return i
		}
	}

	cnt1 := counter()
	fmt.Println(cnt1()) // 1
	fmt.Println(cnt1()) // 2
	fmt.Println(cnt1()) // 3

	cnt2 := counter()
	fmt.Println(cnt2()) // 1
	fmt.Println(cnt2()) // 2

	////
	fmt.Println("Vi du 3:")
	multiplyBy := func(n int) func(int) int {
		// n = 5
		return func(x int) int {
			return x * n
		}
	}

	// double, triple la mot closure function, nó nhớ biến n được truyền vào hàm multiplyBy
	double := multiplyBy(2)
	triple := multiplyBy(3)

	fmt.Println(double(5), triple(20)) // 10 60
}

func funcClosureEx2() {
	adder := func() func(int) int {
		sum := 0
		return func(x int) int {
			sum += x
			return sum
		}
	}

	pos, neg := adder(), adder()
	fmt.Println(pos(1), neg(-2)) // 1 -2
	fmt.Println(pos(2), neg(-3)) // 3 -5
	fmt.Println(pos(4), neg(-5)) // 7 -10
}

// Fibonacci closure (cach 1)
func funcClosureEx3() {
	fibonacci := func() func() int {
		a, b := 0, 1
		return func() int {
			res := a
			a, b = b, a+b
			return res
		}
	}

	fibo := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Print(fibo(), " ")
	}
	fmt.Println()
}

// Fibonacci closure (cach 2)
func funcClosureEx4() {
	fibonacci := func() func() int {
		fibo := []int{0, 1}
		i := 0

		return func() int {
			if i < len(fibo) {
				value := fibo[i]
				i++
				return value
			}

			n := len(fibo)
			next := fibo[n-1] + fibo[n-2]
			fibo = append(fibo, next)
			i++
			return next
		}
	}

	fibo := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Print(fibo(), " ")
	}
	fmt.Println()
}

// Fibonacci closure (cach 1) => sinh so fibonacci sau do dua vao mang (slice)
func funcClosureEx5() {
	fibonacci := func() func() int {
		a, b := 0, 1
		return func() int {
			value := a
			a, b = b, a+b
			return value
		}
	}

	f := fibonacci()
	fibo := make([]int, 0, 50)
	for i := 0; i < 50; i++ {
		fibo = append(fibo, f())
	}

	for _, value := range fibo {
		fmt.Print(value, " ")
	}
	fmt.Println()
}

// Fibonacci khong dung function closure (cach 1)
func fibonacciEx6() {
	fibonacci := func() []int {
		n := 50
		fibo := make([]int, n)
		// fmt.Println(len(fibo), cap(fibo))		// 50 50
		fibo[0] = 0
		fibo[1] = 1
		for i := 2; i < n; i++ {
			fibo[i] = fibo[i-1] + fibo[i-2]
		}
		return fibo
	}

	fibo := fibonacci()
	for _, value := range fibo {
		fmt.Print(value, " ")
	}
	fmt.Println()
}

// Fibonacci khong dung function closure (cach 2)
func fibonacciEx7() {
	fibonacci := func() []int {
		n := 50
		fibo := make([]int, 0, n)
		// fmt.Println(len(fibo), cap(fibo)) // 0 50
		fibo = append(fibo, 0)
		fibo = append(fibo, 1)

		for i := 2; i < n; i++ {
			fibo = append(fibo, fibo[i-1]+fibo[i-2])
		}

		return fibo
	}

	fibo := fibonacci()
	for _, value := range fibo {
		fmt.Print(value, " ")
	}
	fmt.Println()
}

// Ví dụ 1: method với receiver type la struct
type Rect struct {
	width, height int
}

// Day la method
func (r Rect) Area() int {
	return r.width * r.height
}

// Day la function
func Area2(r Rect) int {
	return r.width * r.height
}

func methodEx1() {
	r := Rect{3, 4}
	fmt.Println("Dien tich: ", r.Area()) // Gọi method giống như trong OOP. Day la method
	fmt.Println("Dien tich: ", Area2(r)) // Day la function
}

// Ví dụ 2: method với receiver type khong phai struct
type MyFloat float64

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

func methodEx2() {
	f := MyFloat(-2.5)
	fmt.Println(f.Abs())
}

// Ví dụ 3: method với Pointer receiver
type CustomVertex struct {
	x, y int
}

func (v CustomVertex) notChange() {
	v.x += 10
	v.y += 10
}

func (v *CustomVertex) increase() {
	v.x += 20
	v.y += 20
}

func methodEx3() {
	v := CustomVertex{3, 4}
	fmt.Println(v) // {3 4}

	v.notChange()
	fmt.Println(v) // {3 4}

	v.increase()
	fmt.Println(v) // {23 24}
}

// Ví dụ 4: Methods and pointer indirection
type CustomVertex2 struct {
	x, y int
}

func (v CustomVertex2) notChange2(f int) {
	v.x *= f
	v.y *= f
}

func (v *CustomVertex2) scale2(f int) {
	v.x *= f
	v.y *= f
}

func ScaleFunc2(v *CustomVertex2, f int) {
	v.x *= f
	v.y *= f
}

func methodEx4() {
	v := CustomVertex2{3, 4}
	fmt.Println(v) // {3 4}

	v.notChange2(2)
	fmt.Println(v) // {3 4}

	v.scale2(2)
	fmt.Println(v) // {6 8}

	ScaleFunc2(&v, 2)
	fmt.Println(v) // {12 16}

	//
	p := &v

	p.notChange2(2)
	fmt.Println(v, *p) // {12 16} {12 16}

	p.scale2(2)        // được thông dịch (interpreted) thành: (*p).scale2(2)
	fmt.Println(v, *p) // {24 32} {24 32}

	(*p).scale2(2)
	fmt.Println(v, *p) // {48 64} {48 64}

	ScaleFunc2(p, 2)
	fmt.Println(v, *p) // {96 128} {96 128}
}

func importPackageEx1() {
	product.RunProduct()
}

// Ví dụ: đọc dữ liệu từ bàn phím và chuỗi sử dụng fmt.Scanln, fmt.Scanf, fmt.Scan, fmt.Sscanf
/*
	fmt.Scan: đọc dữ liệu từ bàn phím, kết thúc khi gặp ký tự trắng (space, tab, newline)
	fmt.Scanln: đọc dữ liệu từ bàn phím, kết thúc khi gặp ký tự xuống dòng (Enter)
	fmt.Scanf: đọc dữ liệu từ bàn phím với định dạng cụ thể
	fmt.Sscanf: đọc dữ liệu từ chuỗi với định dạng cụ thể
*/
func fmtScanEx1() {
	// Su dung ffmt.Scan, fmt.Scanln, fmt.Scanf
	var a int
	var b float64
	var s string

	fmt.Print("Nhap so nguyen: ")
	fmt.Scan(&a)

	fmt.Print("Nhap so thuc: ")
	fmt.Scanln(&b)

	fmt.Print("Nhap chuoi: ")
	fmt.Scanf("%s", &s)

	fmt.Printf("Ban vua nhap: a = %d, b = %.2f, s = %q\n", a, b, s)

	// Su dung fmt.Sscanf
	var x int
	var y float64
	var str string

	input := "100 3.14 HelloGo"
	n, err := fmt.Sscanf(input, "%d %f %s", &x, &y, &str)
	if err != nil {
		fmt.Println("Loi:", err)
	} else {
		fmt.Printf("Da doc duoc %d gia tri: x = %d, y = %.2f, str = %q\n", n, x, y, str)
	}
}

// Nhap du lieu tu ban phim vao array
func fmtScanEx2() {
	var n int
	fmt.Scan(&n)

	var arr [100]int // array cố định

	for i := 0; i < n; i++ {
		fmt.Scan(&arr[i])
	}

	for i := 0; i < n; i++ {
		fmt.Print(arr[i], " ")
	}
}

// Nhap du lieu tu ban phim vao slice
func fmtScanEx3() {
	var n int
	fmt.Scan(&n)

	arr := make([]int, n) // tạo slice

	for i := 0; i < n; i++ {
		fmt.Scan(&arr[i])
	}

	fmt.Println(arr)
}

// struct embedding
func structEmbeddingEx1() {
	structembedding11.ExampleEmbedding()
}

func structEmbeddingEx2() {
	structembedding12.Main()
}

func profitRevenueEx1() {
	profitrevenue.Main()
}

func main() {
	// variables()
	// typeConversion()
	// constValue()
	// fmt.Println(calSum())
	// switchEx1()
	// deferEx1()
	// pointerEx1()
	// vertexEx1()
	// vertexEx2()
	// vertexEx3()
	// arrayEx1()
	// sliceEx1()
	// sliceEx2()
	// sliceEx3()
	// makeCopyFullSliceEx3()
	// appendEx1()
	// appendEx2()
	// appendEx3()
	// rangeEx1()
	// mapEx1()
	// mapEx2()
	// funcEx1()
	// funcClosureEx1()
	// funcClosureEx2()
	// funcClosureEx3()
	// funcClosureEx4()
	// funcClosureEx5()
	// fibonacciEx6()
	// fibonacciEx7()
	// methodEx1()
	// methodEx2()
	// methodEx3()
	// methodEx4()
	// importPackageEx1()
	// fmtScanEx1()
	// fmtScanEx2()
	// fmtScanEx3()
	// structEmbeddingEx1()
	structEmbeddingEx2()
	// profitRevenueEx1()
}
