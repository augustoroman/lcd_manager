package main

import "github.com/augustoroman/serial_lcd"

func birthday_chars(lcd LCD) {
	// balloon
	lcd.CreateCustomChar(0, serial_lcd.MakeChar([8]string{
		".##..",
		"#..#.",
		"#..#.",
		".##..",
		"..#..",
		".#...",
		".#...",
		"..#..",
	}))

	// balloon 2
	lcd.CreateCustomChar(1, serial_lcd.MakeChar([8]string{
		"..#..",
		".....",
		".##..",
		"#..#.",
		"#..#.",
		".##..",
		"..#..",
		".#...",
	}))

	// smiley
	lcd.CreateCustomChar(2, serial_lcd.MakeChar([8]string{
		".....",
		".#.#.",
		".#.#.",
		".....",
		"..#..",
		"#...#",
		"#####",
		".###.",
	}))

	// confetti
	lcd.CreateCustomChar(3, serial_lcd.MakeChar([8]string{
		"....#",
		".#...",
		".....",
		"...#.",
		"#....",
		".....",
		"...#.",
		".#...",
	}))

	// confetti
	lcd.CreateCustomChar(4, serial_lcd.MakeChar([8]string{
		"#....",
		"...#.",
		".....",
		".....",
		".#...",
		"....#",
		".....",
		"#....",
	}))

	// confetti
	lcd.CreateCustomChar(5, serial_lcd.MakeChar([8]string{
		"..#..",
		".....",
		"#....",
		".....",
		"....#",
		".....",
		".#...",
		"...#.",
	}))

}

func heart_and_snowman(lcd LCD) {
	lcd.CreateCustomChar(0, serial_lcd.MakeChar([8]string{
		".....",
		".*.*.",
		"*.*.*",
		"*...*",
		"*...*",
		".*.*.",
		"..*..",
		".....",
	}))
	lcd.CreateCustomChar(1, serial_lcd.MakeChar([8]string{
		"..o..",
		".o.o.",
		".ooo.",
		"o...o",
		".ooo.",
		"o...o",
		"o...o",
		".ooo.",
	}))
}

func invaderChars2(lcd LCD) {
	lcd.CreateCustomChar(0, serial_lcd.MakeChar([8]string{
		"..#..",
		"#..#.",
		"#.###",
		"###.#",
		"#####",
		".####",
		"..#..",
		".#...",
	}))

	lcd.CreateCustomChar(1, serial_lcd.MakeChar([8]string{
		"..#..",
		".#..#",
		"###.#",
		"#.###",
		"#####",
		"####.",
		"..#..",
		"...#.",
	}))

	lcd.CreateCustomChar(3, serial_lcd.MakeChar([8]string{
		"..#..",
		"...#.",
		"..###",
		"###.#",
		"#.###",
		"#.###",
		"..#..",
		"...#.",
	}))

	lcd.CreateCustomChar(4, serial_lcd.MakeChar([8]string{
		"..#..",
		".#...",
		"###..",
		"#.###",
		"###.#",
		"###.#",
		"..#..",
		".#...",
	}))
}

func watcherChars(lcd LCD) {
	a1 := serial_lcd.MakeChar([8]string{
		".....",
		".....",
		".....",
		".....",
		".....",
		".....",
		".*.*.",
		"*.*.*",
	})

	a2 := serial_lcd.MakeChar([8]string{
		"**...",
		"..*..",
		"...*.",
		"**.*.",
		"*..**",
		"...**",
		"**.*.",
		"...*.",
	})

	a3 := serial_lcd.MakeChar([8]string{
		"...**",
		"..*..",
		".*...",
		".*.**",
		"**..*",
		"**...",
		".*.**",
		".*...",
	})

	lcd.CreateCustomChar(0, a1)
	lcd.CreateCustomChar(1, a2)
	lcd.CreateCustomChar(2, a3)
}

func pacmanChars(lcd LCD) {
	a1 := serial_lcd.MakeChar([8]string{
		".....",
		"**...",
		"***..",
		"****.",
		"****.",
		"*****",
		"*....",
		".....",
	})
	a2 := serial_lcd.MakeChar([8]string{
		".....",
		"...**",
		"..***",
		".****",
		"*****",
		"*****",
		"*****",
		"*****",
	})
	a3 := serial_lcd.MakeChar([8]string{
		".....",
		"*....",
		"*****",
		"****.",
		"***..",
		"**...",
		".....",
		".....",
	})
	a4 := serial_lcd.MakeChar([8]string{
		"*****",
		"*****",
		"*****",
		".****",
		"..***",
		"...**",
		".....",
		".....",
	})
	b1 := serial_lcd.MakeChar([8]string{
		".....",
		".....",
		".....",
		".....",
		".....",
		".....",
		"....*",
		"...**",
	})
	b2 := serial_lcd.MakeChar([8]string{
		".....",
		".....",
		".....",
		".....",
		".....",
		".....",
		"*....",
		"**...",
	})
	b3 := serial_lcd.MakeChar([8]string{
		"**...",
		"*....",
		".....",
		".....",
		".....",
		".....",
		".....",
		".....",
	})
	b4 := serial_lcd.MakeChar([8]string{
		"...**",
		"....*",
		".....",
		".....",
		".....",
		".....",
		".....",
		".....",
	})

	lcd.CreateCustomChar(0, a1)
	lcd.CreateCustomChar(1, a2)
	lcd.CreateCustomChar(2, a3)
	lcd.CreateCustomChar(3, a4)
	lcd.CreateCustomChar(4, b1)
	lcd.CreateCustomChar(5, b2)
	lcd.CreateCustomChar(6, b3)
	lcd.CreateCustomChar(7, b4)
}
