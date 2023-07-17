package sample

import (
	"proto_demo/pb"

	"github.com/golang/protobuf/ptypes"
)

func NewKeyboard() *pb.Keyboard {
	keyboard := &pb.Keyboard{
		Layout:  randomKeyboardLayout(),
		Backlit: randomBool(),
	}
	return keyboard
}
func NewCPU() *pb.CPU {
	brand := randomCPUBrand()
	name := randomCPUName(brand)
	numberCores := randomInt(2, 8)
	numberThreads := randomInt(numberCores, 12)
	minghz := randomFloag64(2.0, 3.5)
	maxghz := randomFloag64(minghz, 5.0)
	cpu := &pb.CPU{
		Brand:         brand,
		Name:          name,
		NumberCores:   uint32(numberCores),
		NumberThreads: uint32(numberThreads),
		MinGhz:        minghz,
		MaxGhz:        maxghz,
	}
	return cpu
}
func NEWGPU() *pb.GPU {
	brand := randomGPUBrand()
	name := rangdomGPUName(brand)
	minGhz := randomFloag64(1.0, 1.5)
	maxGhz := randomFloag64(minGhz, 2.0)

	memory := &pb.Memory{
		Value: uint64(randomInt(2, 6)),
		Unit:  pb.Memory_GIGABYTE,
	}
	gpu := &pb.GPU{
		Brand:  brand,
		Name:   name,
		MinGhz: minGhz,
		MaxGhz: maxGhz,
		Memory: memory,
	}
	return gpu
}
func NewRAM() *pb.Memory {

	ram := &pb.Memory{
		Value: uint64(randomInt(4, 64)),
		Unit:  pb.Memory_GIGABYTE,
	}
	return ram
}
func NewSSD() *pb.Storage {
	ssd := &pb.Storage{
		Driver: pb.Storage_SSD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(128, 1024)),
			Unit:  pb.Memory_GIGABYTE,
		},
	}
	return ssd
}
func NewHDD() *pb.Storage {
	hdd := &pb.Storage{
		Driver: pb.Storage_HDD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(1, 6)),
			Unit:  pb.Memory_TERABYTE,
		},
	}
	return hdd
}
func NewScreen() *pb.Screen {
	screen := &pb.Screen{
		SizeInch:   randomFloag32(13, 17),
		Resolution: randomScreenResolution(),
		Panel:      randomScreenPanel(),
		Multitouch: randomBool(),
	}
	return screen
}

func NewLaptop() *pb.Laptop {
	brand := randomLaptopBrand()
	name := randomLaptopName(brand)
	laptop := &pb.Laptop{
		Id:       randomID(),
		Brand:    brand,
		Name:     name,
		Cpu:      NewCPU(),
		Ram:      NewRAM(),
		Gpus:     []*pb.GPU{NEWGPU()},
		Storages: []*pb.Storage{NewHDD(), NewSSD()},
		Screen:   NewScreen(),
		Keyboard: NewKeyboard(),
		Weight: &pb.Laptop_WeightKg{
			WeightKg: randomFloag64(1.0, 3.0),
		},
		PriceUsd:    randomFloag64(1500.0, 3000.0),
		ReleaseYear: uint32(randomInt(2015, 2019)),
		UpdateAt:    ptypes.TimestampNow(),
	}
	return laptop
}

func RandomLaptopScore() float64 {
	return float64(randomInt(1, 10))
}
