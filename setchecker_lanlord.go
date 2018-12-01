package pokergame

import (
	"errors"
	"github.com/wqtapp/poker"
)

type LandLordChecker struct {

}

func (self *LandLordChecker) GetSetInfo(set poker.PokerSet) (*SetInfo,error) {
	switch set.CountCards() {
	case 0:
		return nil,errors.New("玩家出牌为空")
		//单张
	case 1:
		return NewSetInfo(SET_TYPE_SINGLE,set[0].GetValue(),set[0].GetValue()),nil
		//对子或者王炸
	case 2:
		if self.isPair(set){
			return NewSetInfo(SET_TYPE_PAIR,set[0].GetValue(),set[0].GetValue()),nil
		}

		if self.isJokerBomb(set){
			return NewSetInfo(SET_TYPE_JOKER_BOMB,set[0].GetValue(),set[1].GetValue()),nil
		}

		return nil,errors.New("牌型不符合规则")
		//三张
	case 3:
		return self.checkThreePlus(set)
		//炸弹或三带一
	case 4:
		if self.isCommonBomb(set){
			return NewSetInfo(SET_TYPE_COMMON_BOMB,set[0].GetValue(),set[0].GetValue()),nil
		}
		return self.checkThreePlus(set)
		//三带二或者一条龙
	case 5:
		if self.isDragon(set){
			return NewSetInfo(SET_TYPE_DRAGON,set[0].GetValue(),set[set.CountCards()-1].GetValue()),nil
		}
		return self.checkThreePlus(set)
		//一条龙，或者四带二，或者四带二对
	default:
		if self.isDragon(set){
			return NewSetInfo(SET_TYPE_DRAGON,set[0].GetValue(),set[set.CountCards()-1].GetValue()),nil
		}

		if self.isMultiPair(set){
			return NewSetInfo(SET_TYPE_MULIT_PAIRS,set[0].GetValue(),set[set.CountCards()-1].GetValue()),nil
		}

		if cardsType,err := self.checkFourPlus(set);err == nil{
			return cardsType,err
		}else if cardsType,err := self.checkMultiFourPlus(set);err == nil{
			return cardsType,err
		}else{
			return self.checkMultiThreePlus(set)
		}
	}
}

func (self *LandLordChecker) isPair(set poker.PokerSet) bool{
	if set.CountCards() != 2 {
		return false
	}

	if(set[0].GetValue() == set[1].GetValue()){
		return true
	}
	return false
}

func (self *LandLordChecker) isMultiPair(set poker.PokerSet) bool{
	if set.CountCards()%2 != 0 || len(set) < 6 {
		return false
	}

	set.SortAsc()
	//2和王不能作为连对出牌
	if set[set.CountCards()-1].GetValue() >= poker.CARD_VALUE_TWO{
		return false
	}

	currValue := -1

	for i,card := range set{
		if i == 0 {
			currValue = card.GetValue()
		}else{
			if i % 2 == 1{
				if card.GetValue() != currValue{
					return false
				}
			}else{
				if card.GetValue() == currValue+1{
					currValue=card.GetValue()
				}else{
					return false
				}
			}
		}
	}
	return true
}

func (self *LandLordChecker) isJokerBomb(set poker.PokerSet) bool{
	if set.CountCards() != 2{
		return false
	}
	set.SortAsc()
	if(set[0].GetValue() == poker.CARD_VALUE_BLACK_JOKER && set[1].GetValue() == poker.CARD_VALUE_RED_JOKER){
		return true
	}else{

	}
	return false
}

func (self *LandLordChecker) isCommonBomb(set poker.PokerSet) bool{
	if set.CountCards() != 4{
		return false
	}
	if set[0].GetValue() == set[1].GetValue()  && set[2].GetValue() == set[3].GetValue() &&
		set[0].GetValue() == set[2].GetValue(){
		return true
	}else{
		return false
	}
}

func (self *LandLordChecker) isDragon(set poker.PokerSet) bool{
	if len(set) < 5 {
		return false
	}

	set.SortAsc()
	//2和王不能参与顺子出牌
	if set[set.CountCards()-1].GetValue() >= poker.CARD_VALUE_TWO{
		return false
	}

	tempValue := -1
	for i,card := range set{
		if i == 0 {
			tempValue = card.GetValue()
		}else{
			if card.GetValue() == tempValue+1{
				tempValue = card.GetValue()
			}else{
				return false
			}
		}
	}
	return true
}

func (self *LandLordChecker) checkThreePlus(set poker.PokerSet) (*SetInfo,error){
	pokersNum := set.CountCards()
	if pokersNum < 3 || pokersNum >5{
		return nil,errors.New("不是三带牌")
	}

	set.SortAsc()
	cardNum := set.AnalyzeEachCardValueNum()
	cardNumCount := len(cardNum)
	if pokersNum == 3{
		if cardNumCount == 1{
			return NewSetInfo(SET_TYPE_THREE,set[0].GetValue(),set[0].GetValue()),nil
		}else{
			return nil,errors.New("不是三带牌")
		}
	}else{
		if cardNumCount == 2 {

			for k,v := range cardNum{
				if v == 3{
					if(pokersNum == 4){
						return NewSetInfo(SET_TYPE_THREE_PLUS_ONE,k,k),nil
					}else{
						return NewSetInfo(SET_TYPE_THREE_PLUS_TWO,k,k),nil
					}
				}
			}
			return nil,errors.New("不是三带牌")
		}else{
			return nil,errors.New("不是三带牌")
		}
	}
}

//是否是四代一或者四代二
func (self *LandLordChecker) checkFourPlus(set poker.PokerSet) (*SetInfo,error){

	pokersNum := set.CountCards()
	if pokersNum != 6 && pokersNum != 8{
		return nil,errors.New("不是四带牌")
	}

	set.SortAsc()

	cardNum := set.AnalyzeEachCardValueNum()
	cardNumCount := len(cardNum)
	if cardNumCount == 2{
		k1 := -1
		v1 := -1
		k2 := -1
		v2 := -1
		i := 1
		for k,v := range cardNum{
			if i == 1{
				k1 = k
				v1 = v
			}else{
				k2 = k
				v2 = v
			}
			i++
		}

		if pokersNum == 6{
			//支持444455这种四带二
			if v1 == 4{
				return NewSetInfo(SET_TYPE_FOUR_PLUS_TWO,k1,k1),nil
			}else if v2 == 4{
				return NewSetInfo(SET_TYPE_FOUR_PLUS_TWO,k2,k2),nil
			}else{
				return nil,errors.New("不是四带牌")
			}
		}else{
			//不支持44445555这种连续四带四的牌,因为没法判断是四带四还是多连四
			//支持44446666这种牌型，以大的作为主牌，即四个6带两对4
			if v1 == 4 && k1 != k2+1 && k2 != k1+1{
				k := -1
				if k1 > k2 {
					k = k1
				}else{
					k = k2
				}
				return NewSetInfo(SET_TYPE_FOUR_PLUS_FOUR,k,k),nil
			}else{
				return nil,errors.New("不是四带牌")
			}
		}
	}else if cardNumCount == 3{
		mainValue := -1
		for k,v := range cardNum{
			if v == 4 {
				mainValue = k
			}
			if pokersNum == 6{
				if v == 2 || v == 3 || v == 5 || v == 6{
					return nil,errors.New("不是四带牌")
				}
			}else{
				if v == 1 || v == 3 || v == 5 || v == 6 {
					return nil,errors.New("不是四带牌")
				}
			}
		}
		if mainValue == -1{
			return nil,errors.New("不是四带牌")
		}else{
			if pokersNum == 6 {
				//四带二444456这种
				return NewSetInfo(SET_TYPE_FOUR_PLUS_TWO,mainValue,mainValue),nil
			}else{
				//四带四44445566这种
				return NewSetInfo(SET_TYPE_FOUR_PLUS_FOUR,mainValue,mainValue),nil
			}
		}
	} else{
		return nil,errors.New("不是四带牌")
	}
}

//是否多个三带一，或三代二，或不带
func (self *LandLordChecker) checkMultiThreePlus(set poker.PokerSet) (*SetInfo,error){
	pokerNum := set.CountCards()
	if pokerNum < 6 {
		return nil,errors.New("不是三顺")
	}

	set.SortAsc()
	cardNum := set.AnalyzeEachCardValueNum()

	//mainCardValue := -1      //暂存主牌的value，用于比较是否连续
	//mainCardNum := 0        //主牌的数量
	mainCardValues := []int{}  //存放主牌的值
	attachCardNum := 0    //附牌的数量
	attachCardNumMap := make(map[int]int)  //附牌的value和num的map

	for k,v := range cardNum{
		if v == 3{
			mainCardValues = append(mainCardValues,k)
		}else{
			attachCardNumMap[k] = v
			attachCardNum += v
		}
	}
	BubbleSortIntMin2Max(mainCardValues)
	//只包含连续的主牌的数量，不连续的同数量的当做附牌对待
	realMainCardValues := []int{}
	//主牌连续，且只有一个连续的，其他的间断连续作为附牌处理
	for i,value := range mainCardValues{
		if i < len(mainCardValues)-1 && mainCardValues[i] + 1 == mainCardValues[i+1]{
			if len(realMainCardValues) > 0 && value == realMainCardValues[len(realMainCardValues)-1]+1{
				realMainCardValues = append(realMainCardValues,value)
			}else if len(realMainCardValues) == 0{
				realMainCardValues = append(realMainCardValues,value)
			}else{
				attachCardNumMap[value] = 3
				attachCardNum += 3
			}
		}else if i > 0 && mainCardValues[i] == mainCardValues[i-1]+1{
			if len(realMainCardValues) > 0 && value == realMainCardValues[len(realMainCardValues)-1]+1{
				realMainCardValues = append(realMainCardValues,value)
			}else if len(realMainCardValues) == 0{
				realMainCardValues = append(realMainCardValues,value)
			}else{
				attachCardNumMap[value] = 3
				attachCardNum += 3
			}
		}else{//该值的牌作为附牌对待
			attachCardNumMap[value] = 3
			attachCardNum += 3
		}
	}

	mainCardNum := len(realMainCardValues)
	if mainCardNum < 2 {//未构成连续牌型
		return nil,errors.New("不是三顺")
	}

	//2和王不能参与连顺
	if mainCardNum > 1 && realMainCardValues[len(realMainCardValues)-1] > poker.CARD_VALUE_ACE{
		return nil,errors.New("不是三顺")
	}

	//没有附牌
	if attachCardNum == 0{
		return NewSetInfo(SET_TYPE_MULITY_THREE,realMainCardValues[0],realMainCardValues[len(realMainCardValues)-1]),nil
	}else if mainCardNum == attachCardNum{//三带一
		return NewSetInfo(SET_TYPE_MULITY_THREE_PLUS_ONE,realMainCardValues[0],realMainCardValues[len(realMainCardValues)-1]),nil
	}else if mainCardNum*2 == attachCardNum{//三带二
		for _,v := range attachCardNumMap{
			if v != 2 && v != 4{
				return nil,errors.New("不是三顺")
			}
		}
		return NewSetInfo(SET_TYPE_MULITY_THREE_PLUS_TWO,realMainCardValues[0],realMainCardValues[len(realMainCardValues)-1]),nil
	}else{
		return nil,errors.New("不是三顺")
	}
}

//是否多个四带一或四代二，或不带
func (self *LandLordChecker) checkMultiFourPlus(set poker.PokerSet) (*SetInfo,error){

	pokerNum := set.CountCards()
	if pokerNum < 8 || pokerNum%2 != 0 {
		return nil,errors.New("不是四顺")
	}

	set.SortAsc()
	cardNum := set.AnalyzeEachCardValueNum()

	//mainCardValue := -1      //暂存主牌的value，用于比较是否连续
	//mainCardNum := 0        //主牌的数量
	mainCardValues := []int{}  //存放主牌的值
	attachCardNum := 0      //附牌的数量
	attachCardNumMap := make(map[int]int)  //附牌的value和num的map

	for k,v := range cardNum{
		if v == 4{
			mainCardValues = append(mainCardValues,k)
		}else{
			attachCardNumMap[k] = v
			attachCardNum += v
		}
	}
	BubbleSortIntMin2Max(mainCardValues)

	//只包含连续的主牌的数量，不连续的同数量的当做附牌对待
	realMainCardValues := []int{}
	//主牌连续，且只有一个连续的，其他的间断连续作为附牌处理
	for i,value := range mainCardValues{
		if i < len(mainCardValues)-1 && mainCardValues[i] + 1 == mainCardValues[i+1]{
			if len(realMainCardValues) > 0 && value == realMainCardValues[len(realMainCardValues)-1]+1{
				realMainCardValues = append(realMainCardValues,value)
			}else if len(realMainCardValues) == 0{
				realMainCardValues = append(realMainCardValues,value)
			}else{
				attachCardNumMap[value] = 4
				attachCardNum += 4
			}
		}else if i > 0 && mainCardValues[i] == mainCardValues[i-1]+1{
			if len(realMainCardValues) > 0 && value == realMainCardValues[len(realMainCardValues)-1]+1{
				realMainCardValues = append(realMainCardValues,value)
			}else if len(realMainCardValues) == 0{
				realMainCardValues = append(realMainCardValues,value)
			}else{
				attachCardNumMap[value] = 4
				attachCardNum += 4
			}
		}else{//该值的牌作为附牌对待
			attachCardNumMap[value] = 4
			attachCardNum += 4
		}
	}

	mainCardNum := len(realMainCardValues)
	for mainCardNum < 2{
		return nil,errors.New("不是四顺")
	}

	//2和王不能参与连顺
	if mainCardNum > 1 && realMainCardValues[len(realMainCardValues)-1] > poker.CARD_VALUE_ACE{
		return nil,errors.New("不是四顺")
	}

	//没有附牌
	if attachCardNum == 0{//四不带
		return NewSetInfo(SET_TYPE_MULITY_FOUR,realMainCardValues[0],realMainCardValues[len(realMainCardValues)-1]),nil
	}else if mainCardNum*2 == attachCardNum{//四带二
		return NewSetInfo(SET_TYPE_MULITY_FOUR_PLUS_TWO,realMainCardValues[0],realMainCardValues[len(realMainCardValues)-1]),nil
	}else if mainCardNum*4 == attachCardNum{//四带四
		for _,v := range attachCardNumMap{
			if v != 2 && v != 4{
				return nil,errors.New("不是四顺")
			}
		}
		return NewSetInfo(SET_TYPE_MULITY_FOUR_PLUS_FOUR,realMainCardValues[0],realMainCardValues[len(realMainCardValues)-1]),nil
	}else{
		return nil,errors.New("不是四顺")
	}
}
