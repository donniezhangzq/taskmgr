package FastSort

import (
    "errors"
)

type T_item struct {
    Id int
    Item_name string
    Status string
    Pri int
    Type_id int
}

func FastSortFun(t *[]int, low int, high int){
    //快速排序
    var key int
    key = low
    if low >= high {
        return
    } else {
        i := low
        j := high
        value := (*t)[low]
        for i<=j {
            for j >= key && (*t)[j] >= value {
                j--
            }
            if j >= key {
                (*t)[key] = (*t)[j]
                key = j
                (*t)[j] = value
            }

            for i <= key && (*t)[i] <= value {
                i++
            }
            if i <= key {
                (*t)[key] = (*t)[i]
                key = i
                (*t)[i] = value
            }
        }
        FastSortFun(t, low, key)
        FastSortFun(t, key+1, high)
    }

}

func GetMaxPri(s []T_item) (result T_item,err error) {
    //快速排序实现返回最大优先级的事项
    var t []int
    err = errors.New("not found most pri")
    for _,v := range(s) {
        t = append(t, v.Pri)
    }
    FastSortFun(&t, 0, len(t)-1)
    MaxPri := t[len(t)-1]
    for _,v := range(s) {
        if v.Pri == MaxPri {
            result = v
            err = nil
        }
    }
    return result,err
}
