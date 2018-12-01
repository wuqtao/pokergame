package pokergame
//冒泡排序法对整数排序
func BubbleSortIntMin2Max(ints []int){
	length := len(ints)
	for i :=0;i<length;i++{
		for j := i;j < length;j++{
			if ints[i] > ints[j]{
				ints[i],ints[j] = ints[j],ints[i]
			}
		}
	}
}
