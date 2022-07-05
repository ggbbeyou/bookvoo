package models

type Period string

const (
	PERIOD_M1  Period = "m1"
	PERIOD_M3  Period = "m3"
	PERIOD_M5  Period = "m5"
	PERIOD_M15 Period = "m15"
	PERIOD_M30 Period = "m30"
	PERIOD_H1  Period = "h1"
	PERIOD_H2  Period = "h2"
	PERIOD_H4  Period = "h4"
	PERIOD_H6  Period = "h6"
	PERIOD_H8  Period = "h8"
	PERIOD_H12 Period = "h12"
	PERIOD_D1  Period = "d1"
	PERIOD_D3  Period = "d3"
	PERIOD_W1  Period = "w1"
	PERIOD_MN  Period = "mn"
)

func Periods() []Period {
	return []Period{
		PERIOD_M1, PERIOD_M3, PERIOD_M5, PERIOD_M15, PERIOD_M30,
		PERIOD_H1, PERIOD_H2, PERIOD_H4, PERIOD_H6, PERIOD_H8, PERIOD_H12,
		PERIOD_D1, PERIOD_D3,
		PERIOD_W1,
		PERIOD_MN,
	}
}
