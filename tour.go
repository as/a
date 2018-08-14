package main

var saydelay = time.Second / 32

func say(text string) {
	for _, v := range text {
		time.Sleep(saydelay)
		D.Key <- key.Event{
			Rune: v,
		}
	}
}
