package sample

import (
	"math/rand"
	"proto_demo/pb"
	"time"

	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func randomKeyboardLayout() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY
	}
}

func randomLaptopBrand() string {
	return randomstringFromSet("Apple", "Dell", "Lenovo")
}
func randomCPUBrand() string {
	return randomstringFromSet("Intel", "AMD")
}
func randomGPUBrand() string {
	return randomstringFromSet("NVIDIA", "AMD")
}
func randomstringFromSet(a ...string) string {
	if len(a) == 0 {
		return ""
	}
	return a[rand.Intn(len(a))]
}

func randomLaptopName(brand string) string {
	switch brand {
	case "Apple":
		return randomstringFromSet("Macbook Air", "Macbook pro")
	case "Deli":
		return randomstringFromSet("Latitude", "Vostro", "XPS", "Alienware")
	default:
		return randomstringFromSet("Thinkpad X1", "Thinkpad P1", "Thinkpad P53")
	}
}
func rangdomGPUName(brand string) string {
	if brand == "NVIDIA" {
		return randomstringFromSet(
			"RTX 2060",
			"RTX 2070",
			"GTX 1660-Ti",
			"GTX 1070",
		)
	}
	return randomstringFromSet(
		"RX 590",
		"RX 580",
		"RX 5700-XT",
		"RX Vega-56",
	)
}
func randomCPUName(brand string) string {
	if brand == "Intel" {
		return randomstringFromSet(
			"Xeon E-2286M",
			"Core i9-9980HK",
			"core i7-9750H",
			"core i5-9400F",
			"core i3-1005G1",
		)
	}
	return randomstringFromSet(
		"Ryzen 7 PRO 2700U",
		"Ryzen 5 PRO 3500U",
		"Ryzen 3 PRO 3200GE",
	)
}
func randomBool() bool {
	return rand.Intn(2) == 1
}
func randomInt(min int, max int) int {
	return min + rand.Intn(max-min+1)
}
func randomFloag64(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
func randomFloag32(min float32, max float32) float32 {
	return min + rand.Float32()*(max-min)
}
func randomScreenPanel() pb.Screen_Panel {
	if rand.Intn(2) == 1 {
		return pb.Screen_IPS
	}
	return pb.Screen_OLED
}
func randomScreenResolution() *pb.Screen_Resolution {
	height := randomInt(1080, 4320)
	width := height * 16 / 9
	resolution := &pb.Screen_Resolution{
		Height: uint32(height),
		Width:  uint32(width),
	}
	return resolution
}
func randomID() string {
	return uuid.New().String()
}
