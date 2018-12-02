package pokergame

import (
	"testing"
	"reflect"
	"strconv"
	"github.com/wqtapp/poker"
	"fmt"
)

var dec poker.PokerDeck

func init(){
	dec = poker.CreateDeck()
}

type Check struct{
	setStr []string   //牌型字符串
	isTrue bool       //是否检测的类型
	setType int       //扑克集类型
}

func checkBool(t *testing.T,c []Check,funcName string){
	for i := range c{
		set := getPokerset(c[i].setStr)
		setInfo := reflect.ValueOf(&LandLordChecker{})
		method := setInfo.MethodByName(funcName)
		if method.String() == "<invalid Value>"{
			t.Error("no this func "+funcName)
		}else{
			params := make([]reflect.Value, 1)                 // 参数
			params[0] = reflect.ValueOf(set)                  // 参数设置为20
			res := method.Call(params)
			if res[0].Interface().(bool) != c[i].isTrue{
				t.Error(funcName+strconv.Itoa(i))
			}
		}
	}
}

func checkBoolWithType(t *testing.T,c []Check,funcName string){
	for i := range c{
		set := getPokerset(c[i].setStr)
		setInfo := reflect.ValueOf(LandLordChecker{})
		fmt.Println(setInfo.NumMethod())
		return
		method := setInfo.MethodByName(funcName)
		if method.String() == "<invalid Value>"{
			t.Error("no this func "+funcName)
		}else{
			params := make([]reflect.Value, 1)                 // 参数
			params[0] = reflect.ValueOf(set)
			res := method.Call(params)
			tempBool := false
			if res[1].Interface() == nil{
				tempBool = true
			}

			if tempBool != c[i].isTrue {
				t.Error(funcName+" err"+strconv.Itoa(i))
			}else{

				if !c[i].isTrue{
					continue
				}

				if res[0].Interface() != nil{
					setTypeInfo,ok := res[0].Interface().(*SetInfo)
					if ok{
						if setTypeInfo == nil || setTypeInfo.setType != c[i].setType{
							if setTypeInfo == nil{
								t.Error(funcName+" typeInfo nil"+strconv.Itoa(i))
							}else{
								t.Error(funcName+" typeInfo "+strconv.Itoa(setTypeInfo.setType)+" "+strconv.Itoa(i))
							}
						}
					}else{
						t.Error(funcName+" typeInfo "+strconv.Itoa(setTypeInfo.setType)+" "+strconv.Itoa(i))
					}
				}else{
					t.Error(funcName+" typeInfo result nil"+strconv.Itoa(i))
				}
			}
		}
	}
}

//根据传入的字符数组，生成相应的扑克集，便于写测试用例
//此处不区分花色
func getPokerset(setString []string) poker.PokerSet{
	set := poker.PokerSet{}
	for _,name := range setString{
		for i,card := range dec.GetAllCards(){
			if card.GetCardName() == name{
				set = append(set,dec.GetCard(i))
				break
			}
		}
	}
	return set
}
//先测试生成扑克集的正确性
func TestPokerSet_GetPokerset(t *testing.T){
	test := []string{"3","4","5"}
	set := getPokerset(test)
	if set[0].GetCardName() != test[0]{
		t.Error("pokerset creator err")
	}

	if set[1].GetCardName() != test[1]{
		t.Error("pokerset creator err")
	}

	if set[2].GetCardName() != test[2]{
		t.Error("pokerset creator err")
	}

}

func TestPokerSet_GetSetTypeInfo(t *testing.T) {
	checks := []Check{
		{[]string{"3"}, true, LANDLORD_SET_TYPE_SINGLE},
		{[]string{"3", "3"}, true, LANDLORD_SET_TYPE_PAIR},
		{[]string{"3", "3", "4", "4", "5", "5"}, true, LANDLORD_SET_TYPE_MULIT_PAIRS},
		{[]string{"3", "3", "3"}, true, LANDLORD_SET_TYPE_THREE},
		{[]string{"3", "3", "3", "4"}, true, LANDLORD_SET_TYPE_THREE_PLUS_ONE},
		{[]string{"3", "3", "3", "4", "4"}, true, LANDLORD_SET_TYPE_THREE_PLUS_TWO},
		{[]string{"3", "3", "3", "4", "4", "4"}, true, LANDLORD_SET_TYPE_MULITY_THREE},
		{[]string{"3", "3", "3", "4", "4", "4", "5", "6"}, true, LANDLORD_SET_TYPE_MULITY_THREE_PLUS_ONE},
		{[]string{"3", "3", "3", "4", "4", "4", "5", "5","5","7","7","7"}, true, LANDLORD_SET_TYPE_MULITY_THREE_PLUS_ONE},
		{[]string{"3", "3", "3", "4", "4", "4", "5", "6", "5", "6"}, true, LANDLORD_SET_TYPE_MULITY_THREE_PLUS_TWO},
		{[]string{"3", "3", "3", "4", "4", "4", "5", "5", "5", "5"}, true, LANDLORD_SET_TYPE_MULITY_THREE_PLUS_TWO},

		{[]string{"3", "3", "3", "3"}, true, LANDLORD_SET_TYPE_COMMON_BOMB},
		{[]string{"3", "3", "3", "3", "4", "5"}, true, LANDLORD_SET_TYPE_FOUR_PLUS_TWO},
		{[]string{"3", "3", "3", "3", "4", "4", "5", "5"}, true, LANDLORD_SET_TYPE_FOUR_PLUS_FOUR},
		{[]string{"3", "3", "3", "3", "5", "5", "5", "5"}, true, LANDLORD_SET_TYPE_FOUR_PLUS_FOUR},
		{[]string{"3", "3", "3", "3", "4", "4", "4", "4"}, true, LANDLORD_SET_TYPE_MULITY_FOUR},
		{[]string{"3", "3", "3", "3", "4", "4", "4", "4", "5", "6", "7", "8"}, true, LANDLORD_SET_TYPE_MULITY_FOUR_PLUS_TWO},
		{[]string{"3", "3", "3", "3", "4", "4", "4", "4", "5", "5", "7", "8"}, true, LANDLORD_SET_TYPE_MULITY_FOUR_PLUS_TWO},
		{[]string{"3", "3", "3", "3", "4", "4", "4", "4", "5", "5", "6", "6", "7", "7", "8", "8"}, true, LANDLORD_SET_TYPE_MULITY_FOUR_PLUS_FOUR},

		{[]string{"3", "4", "5", "6", "7", "8"}, true, LANDLORD_SET_TYPE_DRAGON},
		{[]string{"BlackJoker", "RedJoker"}, true, LANDLORD_SET_TYPE_JOKER_BOMB},
	}

	checkBoolWithType(t, checks, "GetSetInfo")
}