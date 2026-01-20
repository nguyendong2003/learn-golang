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
	sliceEx1()
}
