package pokergame

import (
	"github.com/wqtapp/poker"
	"sync"
)

//定义玩家的扑克牌分析器map的索引为poker的value,value为改值得扑克牌在玩家牌中的索引
type LandLordAnalyzer struct{
	sync.RWMutex
	dic map[int]poker.PokerSet
}

//根据给定的扑克集初始化分析器
func (ana *LandLordAnalyzer) InitAnalyzer(){
	ana.dic[poker.CARD_VALUE_THREE] = poker.PokerSet{}
	ana.dic[poker.CARD_VALUE_FOUR] = poker.PokerSet{}
	ana.dic[poker.CARD_VALUE_FIVE] = poker.PokerSet{}
	ana.dic[poker.CARD_VALUE_SIX] = poker.PokerSet{}
	ana.dic[poker.CARD_VALUE_SEVEN] = poker.PokerSet{}
	ana.dic[poker.CARD_VALUE_EIGHT] = poker.PokerSet{}
	ana.dic[poker.CARD_VALUE_NINE] = poker.PokerSet{}
	ana.dic[poker.CARD_VALUE_TEN] = poker.PokerSet{}
	ana.dic[poker.CARD_VALUE_JACK] = poker.PokerSet{}
	ana.dic[poker.CARD_VALUE_QUEEN] = poker.PokerSet{}
	ana.dic[poker.CARD_VALUE_KING] =  poker.PokerSet{}
	ana.dic[poker.CARD_VALUE_ACE] = poker.PokerSet{}
	ana.dic[poker.CARD_VALUE_TWO] = poker.PokerSet{}
	ana.dic[poker.CARD_VALUE_BLACK_JOKER] = poker.PokerSet{}
	ana.dic[poker.CARD_VALUE_RED_JOKER] = poker.PokerSet{}
}
//根据给定的扑克集更新记牌器,出牌时调用
func (ana *LandLordAnalyzer) RemovePokerSet(pokers poker.PokerSet){
	ana.Lock()
	defer ana.Unlock()
	pokers.DoOnEachPokerCard(func(index int,card *poker.PokerCard){
		ana.dic[card.GetValue()],_ = ana.dic[card.GetValue()].DelPokers(poker.PokerSet{card})
	})
}

func (ana *LandLordAnalyzer) AddPokerSet(pokers poker.PokerSet){
	ana.Lock()
	defer ana.Unlock()
	pokers.DoOnEachPokerCard(func(index int,card *poker.PokerCard){
		ana.dic[card.GetValue()] = ana.dic[card.GetValue()].AddPokers(poker.PokerSet{card})
	})
}

func (ana *LandLordAnalyzer) GetMinPlayableCards() poker.PokerSet{
	ana.Lock()
	defer ana.Unlock()
	for i:= poker.CARD_VALUE_THREE;i<=poker.CARD_VALUE_RED_JOKER;i++{
		set,_ := ana.dic[i]
		if set.CountCards() > 0{
			return set
		}
	}
	return poker.PokerSet{}
}
//根据最后一次出牌的牌型信息，返回可出的扑克集
func (ana *LandLordAnalyzer) GetUseableCards(setType *SetInfo) []poker.PokerSet{
	ana.Lock()
	defer ana.Unlock()

	var useableSets []poker.PokerSet

	switch setType.setType {
	case SET_TYPE_SINGLE:
		useableSets = ana.getSingleValueSet(1,setType.GetMinValue())
	case SET_TYPE_DRAGON:
		useableSets = ana.getMultiValueSet(1,setType.GetMinValue(),setType.GetMaxValue())
	case SET_TYPE_PAIR:
		useableSets = ana.getSingleValueSet(2,setType.GetMinValue())
	case SET_TYPE_MULIT_PAIRS:
		useableSets = ana.getMultiValueSet(2,setType.GetMinValue(),setType.GetMaxValue())
	case SET_TYPE_THREE:
		useableSets = ana.getSingleValueSet(3,setType.GetMinValue())
	case SET_TYPE_THREE_PLUS_ONE:
		useableSets = ana.getSingleValueSet(3,setType.GetMinValue())
		for i,tempset := range useableSets{
			tempsetPlus := ana.getPlusSet(1,1,tempset)
			if tempsetPlus.CountCards() >0 {
				useableSets[i] = tempset.AddPokers(tempsetPlus)
			}else{//没有牌可以带，将之前的主牌移除可出牌集合
				useableSets[i] = nil
			}
		}
	case SET_TYPE_THREE_PLUS_TWO:
		useableSets = ana.getSingleValueSet(3,setType.GetMinValue())
		for i,tempset := range useableSets{
			tempsetPlus := ana.getPlusSet(2,1,tempset)
			if tempsetPlus.CountCards() >0 {
				useableSets[i] = tempset.AddPokers(tempsetPlus)
			}else{//没有牌可以带，将之前的主牌移除可出牌集合
				useableSets[i] = nil
			}
		}
	case SET_TYPE_MULITY_THREE:
		useableSets = ana.getMultiValueSet(3,setType.GetMinValue(),setType.GetMaxValue())
	case SET_TYPE_MULITY_THREE_PLUS_ONE:
		useableSets = ana.getMultiValueSet(3,setType.GetMinValue(),setType.GetMaxValue())
		for i,tempset := range useableSets{
			tempsetPlus := ana.getPlusSet(1,setType.GetRangeWidth(),tempset)
			if tempsetPlus.CountCards() >0 {
				useableSets[i] = tempset.AddPokers(tempsetPlus)
			}else{//没有牌可以带，将之前的主牌移除可出牌集合
				useableSets[i] = nil
			}
		}
	case SET_TYPE_MULITY_THREE_PLUS_TWO:
		useableSets = ana.getMultiValueSet(3,setType.GetMinValue(),setType.GetMaxValue())
		for i,tempset := range useableSets{
			tempsetPlus := ana.getPlusSet(2,setType.GetRangeWidth(),tempset)
			if tempsetPlus.CountCards() >0 {
				useableSets[i] = tempset.AddPokers(tempsetPlus)
			}else{//没有牌可以带，将之前的主牌移除可出牌集合
				useableSets[i] = nil
			}
		}
	case SET_TYPE_FOUR_PLUS_TWO:
		useableSets = ana.getSingleValueSet(4,setType.GetMinValue())
		for i,tempset := range useableSets{
			//带两个单牌
			tempsetPlus := ana.getPlusSet(1,2,tempset)
			if tempsetPlus.CountCards() >0 {
				useableSets[i] = tempset.AddPokers(tempsetPlus)
			}else{
				//带一对牌，看做两个单牌
				tempsetPlus := ana.getPlusSet(2,1,tempset)
				if tempsetPlus.CountCards() >0 {
					useableSets[i] = tempset.AddPokers(tempsetPlus)
				}else{//没有牌可以带，将之前的主牌移除可出牌集合
					useableSets[i] = nil
				}
			}
		}
	case SET_TYPE_FOUR_PLUS_FOUR:
		useableSets = ana.getSingleValueSet(4,setType.GetMinValue())
		for i,tempset := range useableSets{
			tempsetPlus := ana.getPlusSet(2,2,tempset)
			if tempsetPlus.CountCards() >0 {
				useableSets[i] = tempset.AddPokers(tempsetPlus)
			}else{//没有牌可以带，将之前的主牌移除可出牌集合
				useableSets[i] = nil
			}
		}
	case SET_TYPE_MULITY_FOUR:
		useableSets = ana.getMultiValueSet(4,setType.GetMinValue(),setType.GetMaxValue())
	case SET_TYPE_MULITY_FOUR_PLUS_TWO:
		useableSets = ana.getMultiValueSet(4,setType.GetMinValue(),setType.GetMaxValue())
		for i,tempset := range useableSets{
			//带两个单牌
			tempsetPlus := ana.getPlusSet(1,2*setType.GetRangeWidth(),tempset)
			if tempsetPlus.CountCards() >0 {
				useableSets[i] = tempset.AddPokers(tempsetPlus)
			}else{
				//带一对牌，看做两个单牌
				tempsetPlus := ana.getPlusSet(2,setType.GetRangeWidth(),tempset)
				if tempsetPlus.CountCards() >0 {
					useableSets[i] = tempset.AddPokers(tempsetPlus)
				}else{//没有牌可以带，将之前的主牌移除可出牌集合
					useableSets[i] = nil
				}
			}
		}
	case SET_TYPE_MULITY_FOUR_PLUS_FOUR:
		useableSets = ana.getMultiValueSet(4,setType.GetMinValue(),setType.GetMaxValue())
		for i,tempset := range useableSets{
			tempsetPlus := ana.getPlusSet(2,2*setType.GetRangeWidth(),tempset)
			if tempsetPlus.CountCards() >0 {
				useableSets[i] = tempset.AddPokers(tempsetPlus)
			}else{//没有牌可以带，将之前的主牌移除可出牌集合
				useableSets[i] = nil
			}
		}
	case SET_TYPE_COMMON_BOMB:
		useableSets = ana.getSingleValueSet(4,setType.GetMinValue())
	case SET_TYPE_JOKER_BOMB:
		useableSets = []poker.PokerSet{}
	default:
		useableSets = []poker.PokerSet{}
	}
	//去掉nil元素
	newUseableSets := []poker.PokerSet{}
	for _,sets := range useableSets{
		if sets != nil{
			newUseableSets = append(newUseableSets,sets)
		}
	}
	//上一次出牌不是炸弹，则直接将炸弹加入可出的排中
	if setType.setType != SET_TYPE_COMMON_BOMB && setType.setType != SET_TYPE_JOKER_BOMB{
		//王炸
		jokerBombSet := ana.getJokerBomb()
		if jokerBombSet.CountCards() > 0{
			newUseableSets = append(newUseableSets,jokerBombSet)
		}
		//普通炸弹
		for _,tempSet := range ana.getSingleValueSet(4,-1){
			if tempSet.CountCards() > 0{
				newUseableSets = append(newUseableSets,tempSet)
			}
		}
	}
	return newUseableSets
}
//获取单值牌组成的扑克集的切片，单排对牌三牌四排等等
//count表示单值牌的张数
//minValue表示上家出牌的最小的牌的大小
func (ana *LandLordAnalyzer) getSingleValueSet(count int,minValue int) []poker.PokerSet{
	sets := []poker.PokerSet{}
	se := poker.NewPokerSet()
	//先不拆牌的情况下查找
	for i:=minValue+1;i<= poker.CARD_VALUE_RED_JOKER;i++{
		if ana.dic[i].CountCards() == count {
			se = se.AddPokers(ana.dic[i])
			sets = append(sets,se)
			se =  poker.NewPokerSet()
		}
	}
	//不拆牌的情况下找不到可出的牌，再考虑拆牌的情况
	if len(sets) == 0{
		for i:=minValue+1;i<= poker.CARD_VALUE_RED_JOKER;i++{
			if ana.dic[i].CountCards() > count {
				se = se.AddPokers(ana.dic[i][:count])
				sets = append(sets,se)
				se =  poker.NewPokerSet()
			}
		}
	}
	return sets
}
//获取多种不同值组成的扑克集的切片,2连3连4连5连等
func (ana *LandLordAnalyzer) getMultiValueSet(count int,minValue int,maxValue int) []poker.PokerSet{
	sets := []poker.PokerSet{}
	se := poker.NewPokerSet()
	valueRange := maxValue-minValue+1
	//先考虑不拆拍的情况
	for i:=minValue+1;i<= poker.CARD_VALUE_TWO-valueRange;i++{
		for j:=i;j<i+valueRange;j++{
			if ana.dic[j].CountCards() == count {
				se = se.AddPokers(ana.dic[j])
			}
		}
		//该范围内连续的牌的张数符合要求
		if se.CountCards() == valueRange*count{
			sets = append(sets,se)
			se =  poker.NewPokerSet()
		}else{
			se = poker.NewPokerSet()
		}
	}
	//如果不拆拍找不到可出的牌，则考虑拆牌
	if len(sets) == 0{
		for i:=minValue+1;i<= poker.CARD_VALUE_TWO-valueRange;i++{
			for j:=i;j<i+valueRange;j++{
				if ana.dic[j].CountCards() > count {
					se = se.AddPokers(ana.dic[j][:count])
				}
			}
			//该范围内连续的牌的张数符合要求
			if se.CountCards() == valueRange*count{
				sets = append(sets,se)
				se =  poker.NewPokerSet()
			}else{
				se =  poker.NewPokerSet()
			}
		}
	}

	return sets
}
//获取附牌，比如三带一中的一，四带二中二，只获取一种可能即可
//不拆牌为第一原则，可能会带出去大牌
//num张数count系列数exceptset不能包含在内的扑克集
func (ana *LandLordAnalyzer) getPlusSet(num int,count int,exceptSet poker.PokerSet) poker.PokerSet{
	resSet := poker.NewPokerSet()
	//第一原则不拆牌原则
	for i:= poker.CARD_VALUE_THREE;i<= poker.CARD_VALUE_RED_JOKER;i++{
		if ana.dic[i].CountCards() == num{
			if !ana.dic[i][:num].HasSameValueCard(exceptSet) {
				resSet = resSet.AddPokers(ana.dic[i])
			}
		}
		if resSet.CountCards() == num*count{
			return resSet
		}
	}
	//不拆牌找不到则，考虑拆牌
	if resSet.CountCards() == 0{
		for i:= poker.CARD_VALUE_THREE;i<= poker.CARD_VALUE_RED_JOKER;i++{
			if ana.dic[i].CountCards() > num{
				if !ana.dic[i][:num].HasSameValueCard(exceptSet) {
					resSet = resSet.AddPokers(ana.dic[i][:num])
				}
			}
			if resSet.CountCards() == num*count{
				return resSet
			}
		}
	}

	return poker.PokerSet{}
}
func (ana *LandLordAnalyzer) getJokerBomb() poker.PokerSet{
	resSet := poker.NewPokerSet()
	for i:= poker.CARD_VALUE_BLACK_JOKER;i<= poker.CARD_VALUE_RED_JOKER;i++ {
		if ana.dic[i].CountCards() > 0 {
			resSet = resSet.AddPokers(ana.dic[i])
		}
	}
	if resSet.CountCards() > 1{
		return resSet
	}else{
		return poker.NewPokerSet()
	}
}


