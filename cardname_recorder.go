package pokergame

import (
	"strconv"
	"github.com/wqtapp/poker"
	"sync"
)


//记牌器，记录玩家之外带各牌待出的数量
type cardNameRecorder struct {
	sync.RWMutex
	dic map[string]int
}

//初始化玩家记牌器，发牌后调用
func (re *cardNameRecorder)InitRecorder(){
	re.dic = make(map[string]int)
	re.dic[poker.CARD_SYMBOL_THREE] = 0
	re.dic[poker.CARD_SYMBOL_FOUR] = 0
	re.dic[poker.CARD_SYMBOL_FIVE] = 0
	re.dic[poker.CARD_SYMBOL_SIX] = 0
	re.dic[poker.CARD_SYMBOL_SEVEN] = 0
	re.dic[poker.CARD_SYMBOL_EIGHT] = 0
	re.dic[poker.CARD_SYMBOL_NINE] = 0
	re.dic[poker.CARD_SYMBOL_TEN] = 0
	re.dic[poker.CARD_SYMBOL_JACK] = 0
	re.dic[poker.CARD_SYMBOL_QUEEN] = 0
	re.dic[poker.CARD_SYMBOL_KING] = 0
	re.dic[poker.CARD_SYMBOL_ACE] = 0
	re.dic[poker.CARD_SYMBOL_TWO] = 0
	re.dic[poker.CARD_SYMBOL_BLACK_JOKER] = 0
	re.dic[poker.CARD_SYMBOL_RED_JOKER] = 0
}
//增加记录器
func (re *cardNameRecorder)AddPokerSet(playerPokers ...poker.PokerSet){
	re.Lock()
	defer re.Unlock()
	for _,pokerSet := range playerPokers{
		pokerSet.DoOnEachPokerCard(func(index int, card *poker.PokerCard){
			re.dic[card.GetCardName()]++
		})
	}

}
//更新玩家记牌器,玩家出牌后调用
func (re *cardNameRecorder)RemovePokerSet(cards poker.PokerSet){
	re.Lock()
	defer re.Unlock()
	cards.DoOnEachPokerCard(func(index int,card *poker.PokerCard){
		re.dic[card.GetCardName()]--
	})
}
//根据牌的顺序从大到小排序的记牌器json对象
func (re *cardNameRecorder) SequenceJsonEncode() string{
	re.Lock()
	defer re.Unlock()
	jsonString := ""
	jsonString += "{"
	jsonString += "'"+poker.CARD_SYMBOL_RED_JOKER+"'"+":"+strconv.Itoa(re.getPokerNum(poker.CARD_SYMBOL_RED_JOKER))
	jsonString += ","
	jsonString += "'"+poker.CARD_SYMBOL_BLACK_JOKER+"'"+":"+strconv.Itoa(re.getPokerNum(poker.CARD_SYMBOL_BLACK_JOKER))
	jsonString += ","
	jsonString += "'"+poker.CARD_SYMBOL_TWO+"'"+":"+strconv.Itoa(re.getPokerNum(poker.CARD_SYMBOL_TWO))
	jsonString += ","
	jsonString += "'"+poker.CARD_SYMBOL_ACE+"'"+":"+strconv.Itoa(re.getPokerNum(poker.CARD_SYMBOL_ACE))
	jsonString += ","
	jsonString += "'"+poker.CARD_SYMBOL_KING+"'"+":"+strconv.Itoa(re.getPokerNum(poker.CARD_SYMBOL_KING))
	jsonString += ","
	jsonString += "'"+poker.CARD_SYMBOL_QUEEN+"'"+":"+strconv.Itoa(re.getPokerNum(poker.CARD_SYMBOL_QUEEN))
	jsonString += ","
	jsonString += "'"+poker.CARD_SYMBOL_JACK+"'"+":"+strconv.Itoa(re.getPokerNum(poker.CARD_SYMBOL_JACK))
	jsonString += ","
	jsonString += "'"+poker.CARD_SYMBOL_TEN+"'"+":"+strconv.Itoa(re.getPokerNum(poker.CARD_SYMBOL_TEN))
	jsonString += ","
	jsonString += "'"+poker.CARD_SYMBOL_NINE+"'"+":"+strconv.Itoa(re.getPokerNum(poker.CARD_SYMBOL_NINE))
	jsonString += ","
	jsonString += "'"+poker.CARD_SYMBOL_EIGHT+"'"+":"+strconv.Itoa(re.getPokerNum(poker.CARD_SYMBOL_EIGHT))
	jsonString += ","
	jsonString += "'"+poker.CARD_SYMBOL_SEVEN+"'"+":"+strconv.Itoa(re.getPokerNum(poker.CARD_SYMBOL_SEVEN))
	jsonString += ","
	jsonString += "'"+poker.CARD_SYMBOL_SIX+"'"+":"+strconv.Itoa(re.getPokerNum(poker.CARD_SYMBOL_SIX))
	jsonString += ","
	jsonString += "'"+poker.CARD_SYMBOL_FIVE+"'"+":"+strconv.Itoa(re.getPokerNum(poker.CARD_SYMBOL_FIVE))
	jsonString += ","
	jsonString += "'"+poker.CARD_SYMBOL_FOUR+"'"+":"+strconv.Itoa(re.getPokerNum(poker.CARD_SYMBOL_FOUR))
	jsonString += ","
	jsonString += "'"+poker.CARD_SYMBOL_THREE+"'"+":"+strconv.Itoa(re.getPokerNum(poker.CARD_SYMBOL_THREE))
	jsonString += "}"
	return jsonString
}

func (re *cardNameRecorder) getPokerNum(key string) int{
	num,ok := re.dic[key]
	if ok {
		return num
	}else{
		return 0
	}
}


