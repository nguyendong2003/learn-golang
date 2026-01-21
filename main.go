/*
  - fmt.Printf():
    %T :  type
    %v :  value
    %q :  In ra string hoặc rune kèm dấu " hoặc '
    %d : so nguyen
    %f : so thuc
    %.2d : float 2 chữ số thập phân
    %s: string
    %p : địa chỉ con trỏ
    %g : format verb dùng cho số thực (float)

- fmt.Sprint: hàm trong Go dùng để ghép (convert) giá trị thành chuỗi (string), KHÔNG in ra màn hình.  (trả về string)
*/
package main

import (
	"fmt"
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
	appendEx3()
}
