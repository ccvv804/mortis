package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

const (
	BYTEMASK = 0x000000ff
	N        = 314
	T        = (2*N - 1)
)

var HuffCode = []uint16{
	0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000,
	0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000,
	0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000,
	0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000,
	0x040, 0x040, 0x040, 0x040, 0x040, 0x040, 0x040, 0x040,
	0x040, 0x040, 0x040, 0x040, 0x040, 0x040, 0x040, 0x040,
	0x080, 0x080, 0x080, 0x080, 0x080, 0x080, 0x080, 0x080,
	0x080, 0x080, 0x080, 0x080, 0x080, 0x080, 0x080, 0x080,
	0x0c0, 0x0c0, 0x0c0, 0x0c0, 0x0c0, 0x0c0, 0x0c0, 0x0c0,
	0x0c0, 0x0c0, 0x0c0, 0x0c0, 0x0c0, 0x0c0, 0x0c0, 0x0c0,
	0x100, 0x100, 0x100, 0x100, 0x100, 0x100, 0x100, 0x100,
	0x140, 0x140, 0x140, 0x140, 0x140, 0x140, 0x140, 0x140,
	0x180, 0x180, 0x180, 0x180, 0x180, 0x180, 0x180, 0x180,
	0x1c0, 0x1c0, 0x1c0, 0x1c0, 0x1c0, 0x1c0, 0x1c0, 0x1c0,
	0x200, 0x200, 0x200, 0x200, 0x200, 0x200, 0x200, 0x200,
	0x240, 0x240, 0x240, 0x240, 0x240, 0x240, 0x240, 0x240,
	0x280, 0x280, 0x280, 0x280, 0x280, 0x280, 0x280, 0x280,
	0x2c0, 0x2c0, 0x2c0, 0x2c0, 0x2c0, 0x2c0, 0x2c0, 0x2c0,
	0x300, 0x300, 0x300, 0x300, 0x340, 0x340, 0x340, 0x340,
	0x380, 0x380, 0x380, 0x380, 0x3c0, 0x3c0, 0x3c0, 0x3c0,
	0x400, 0x400, 0x400, 0x400, 0x440, 0x440, 0x440, 0x440,
	0x480, 0x480, 0x480, 0x480, 0x4c0, 0x4c0, 0x4c0, 0x4c0,
	0x500, 0x500, 0x500, 0x500, 0x540, 0x540, 0x540, 0x540,
	0x580, 0x580, 0x580, 0x580, 0x5c0, 0x5c0, 0x5c0, 0x5c0,
	0x600, 0x600, 0x640, 0x640, 0x680, 0x680, 0x6c0, 0x6c0,
	0x700, 0x700, 0x740, 0x740, 0x780, 0x780, 0x7c0, 0x7c0,
	0x800, 0x800, 0x840, 0x840, 0x880, 0x880, 0x8c0, 0x8c0,
	0x900, 0x900, 0x940, 0x940, 0x980, 0x980, 0x9c0, 0x9c0,
	0xa00, 0xa00, 0xa40, 0xa40, 0xa80, 0xa80, 0xac0, 0xac0,
	0xb00, 0xb00, 0xb40, 0xb40, 0xb80, 0xb80, 0xbc0, 0xbc0,
	0xc00, 0xc40, 0xc80, 0xcc0, 0xd00, 0xd40, 0xd80, 0xdc0,
	0xe00, 0xe40, 0xe80, 0xec0, 0xf00, 0xf40, 0xf80, 0xfc0,
}

var HuffLength = []int16{
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
	4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
	4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
}

var (
	lzahBuf       [4096]byte
	lzahBufPtr    int
	lzahBitsAvail int
	lzahBits      int
	Frequ         [1000]int
	ForwTree      [1000]int
	BackTree      [1000]int

	outPtr []byte
	inPtr  []byte
	inIdx  int
)

func lzahGetByte() int {
	if inIdx >= len(inPtr) {
		return -1
	}
	b := int(inPtr[inIdx])
	inIdx++
	return b & 0xFF
}

func lzahInitHuf() {
	var i, j int
	for i = 0; i < N; i++ {
		Frequ[i] = 1
		ForwTree[i] = i + T
		BackTree[i+T] = i
	}
	for i, j = 0, N; j < T; i += 2 {
		Frequ[j] = Frequ[i] + Frequ[i+1]
		ForwTree[j] = i
		BackTree[i] = j
		BackTree[i+1] = j
		j++
	}
	Frequ[T] = 0xffff
	BackTree[T-1] = 0
}

func lzahReorder() {
	var i, j, k, l int
	j = 0
	for i = 0; i < T; i++ {
		if ForwTree[i] >= T {
			Frequ[j] = (Frequ[i] + 1) >> 1
			ForwTree[j] = ForwTree[i]
			j++
		}
	}
	for i, j = 0, N; i < T; i += 2 {
		k = i + 1
		l = Frequ[i] + Frequ[k]
		Frequ[j] = l
		k = j - 1
		for l < Frequ[k] {
			k--
		}
		k++
		copy(Frequ[k+1:], Frequ[k:j])
		Frequ[k] = l
		copy(ForwTree[k+1:], ForwTree[k:j])
		ForwTree[k] = i
		j++
	}
	for i = 0; i < T; i++ {
		k = ForwTree[i]
		if k >= T {
			BackTree[k] = i
		} else {
			BackTree[k] = i
			BackTree[k+1] = i
		}
	}
}

func lzahGetBit() {
	if lzahBitsAvail != 0 {
		lzahBits = lzahBits << 1
		lzahBitsAvail--
	} else {
		lzahBits = lzahGetByte() & BYTEMASK
		lzahBitsAvail = 7
	}
}

func lzahOutChar(ch byte) {
	outPtr = append(outPtr, ch)
	lzahBuf[lzahBufPtr] = ch
	lzahBufPtr++
	lzahBufPtr &= 0xfff
}

func deLZAH(in []byte, outlen int) ([]byte, error) {
	var i, i1, j, ch, bbyte, offs, skip int

	inPtr = in
	inIdx = 0

	outPtr = make([]byte, 0, outlen)

	lzahInitHuf()
	lzahBitsAvail = 0

	for i = 0; i < 4036; i++ {
		lzahBuf[i] = ' '
	}
	lzahBufPtr = 4036

	obytes := outlen

	for obytes != 0 {
		ch = ForwTree[T-1]
		for ch < T {
			lzahGetBit()
			if (lzahBits & 0x80) != 0 {
				ch = ch + 1
			}
			ch = ForwTree[ch]
		}
		ch -= T
		if Frequ[T-1] >= 0x8000 {
			lzahReorder()
		}
		i = BackTree[ch+T]
		for {
			Frequ[i]++
			j = Frequ[i]
			i1 = i + 1
			if Frequ[i1] < j {
				for {
					i1++
					if Frequ[i1] >= j {
						break
					}
				}
				i1--
				Frequ[i] = Frequ[i1]
				Frequ[i1] = j

				j = ForwTree[i]
				BackTree[j] = i1
				if j < T {
					BackTree[j+1] = i1
				}
				ForwTree[i] = ForwTree[i1]
				ForwTree[i1] = j
				j = ForwTree[i]
				BackTree[j] = i
				if j < T {
					BackTree[j+1] = i
				}
				i = i1
			}
			i = BackTree[i]
			if i == 0 {
				break
			}
		}
		if ch < 256 {
			lzahOutChar(byte(ch))
			obytes--
		} else {
			if lzahBitsAvail != 0 {
				bbyte = ((lzahBits << 1) & BYTEMASK)
				lzahBits = lzahGetByte() & BYTEMASK
				bbyte |= (lzahBits >> lzahBitsAvail)
				lzahBits = lzahBits << (7 - lzahBitsAvail)
			} else {
				bbyte = lzahGetByte() & BYTEMASK
			}
			offs = int(HuffCode[bbyte])
			skip = int(HuffLength[bbyte]) - 2
			for skip != 0 {
				skip--
				bbyte = bbyte << 1
				lzahGetBit()
				if (lzahBits & 0x80) != 0 {
					bbyte++
				}
			}
			offs |= (bbyte & 0x3f)
			offs = ((lzahBufPtr - offs - 1) & 0xfff)
			ch = ch - 253
			for ch > 0 {
				lzahOutChar(lzahBuf[offs&0xfff])
				offs++
				ch--
				obytes--
				if obytes == 0 {
					break
				}
			}
		}
	}

	return outPtr, nil
}

func mashiro(inputdata []byte, filename string) {
	outsize := int(binary.LittleEndian.Uint32(inputdata[740:744]))
	soksize := int(binary.LittleEndian.Uint32(inputdata[208:212]))
	dataread := inputdata[748:]
	outdata := nanami(dataread, outsize)
	mididata := outdata[soksize:]
	name1 := strings.Replace(filename, ".KYC", ".mid", 1)
	name2 := strings.Replace(name1, ".kyc", ".mid", 1)
	err := ioutil.WriteFile(name2, mididata, 0755)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func nanami(indata []byte, outdataintorg int) (newdata []byte) {
	newdata, err := deLZAH(indata, outdataintorg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return newdata
}
func main() {
	file := flag.String("file", "03000.KYC", "Input file")
	flag.Parse()
	dat, err := ioutil.ReadFile(*file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Mortis V1 (public)")
	kanefile := path.Base(*file)
	mashiro(dat, kanefile)
}
