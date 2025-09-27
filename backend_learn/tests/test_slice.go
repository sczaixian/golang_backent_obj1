package tests

import "fmt"

func printArr[T any](arr []T, name string) {
	fmt.Println("---------- ", name, " ----------len:", len(arr), "  cap:", cap(arr))
	for _, value := range arr {
		fmt.Print(value, ", ")
	}
	fmt.Println()

}

/*
 * 切片有 长度限制
 */
func TestSlice() {
	myNum := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	printArr(myNum, "newNum")

	newNum := myNum[1:3]

	printArr(newNum, "newNum")
	printArr(myNum, "myNum")

	newNum = append(newNum, 1)
	printArr(newNum, "newNum")
	printArr(myNum, "myNum")

	fmt.Println("--------------------------------------------------")

	myNum = []int{10, 20, 30, 40}

	newNum = append(myNum, 50)
	printArr(newNum, "newNum")
	printArr(myNum, "myNum")

	fmt.Println("--------------------------------------------------")

	fruit := []string{"apple", "banana", "pear", "orange", "cherry"}
	myFruit := fruit[2:3:4]
	// 起始索引 low = 2（包含 "pear"）
	// 结束索引 high = 3（不包含 "orange"）
	// 容量上限 max = 4（不包含 "cherry"）
	// 长度（len） = high - low= 3 - 2= 1  切片元素：["pear"]
	// 容量（cap） = max - low= 4 - 2= 2   底层数组可访问的范围：["pear", "orange"]
	printArr(fruit, "fruit")
	printArr(myFruit, "myFruit")
	myFruit = fruit[2:3:3]
	printArr(myFruit, "myFruit")
	myFruit = append(myFruit, "Kiwi")
	printArr(myFruit, "myFruit")

	fmt.Println("--------------------------------------------------")

	num1 := []int{1, 2}
	num2 := []int{3, 4}
	num3 := append(num1, num2...)
	printArr(num3, "num3")

	fmt.Println("--------------------------------------------------")
	num1 = []int{10, 20, 30}
	num2 = make([]int, 5)
	count := copy(num2, num1)
	fmt.Println(count)
	printArr(num1, "newNum")
	printArr(num2, "newNum")
	count = copy(num1, num2)
	fmt.Println(count)
	printArr(num1, "newNum")
	printArr(num2, "newNum")
	// ------------    go  中 int  在32位机器上占 4 字节，  64位机器上占 8 字节-----------------------
}
