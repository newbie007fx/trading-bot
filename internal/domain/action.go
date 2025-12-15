package domain

type Action string

const (
	ActionBuy        Action = "BUY"
	ActionSell       Action = "SELL"
	ActionCheck      Action = "CHECK"
	ActionAdjustRule Action = "ADJUST_RULE"
)
