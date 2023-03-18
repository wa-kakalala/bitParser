package utils

import (
	"bufio"
	"errors"
	"fmt"
)

type bitStructure struct {
	bitBuf         []byte // save data from file
	bitFifo        []byte // save bit info
	bitBufCapcity  uint8  // bitbuf's capcity
	bitFifoCapcity uint8  // bitfifo's capcity
	frontIndex     uint8  // the fifo's front index
	tailIndex      uint8  // the fifo's tail index
	bitSize        uint8  // the size of element in fifo
}

/**
* @func: create bitStructure
* @param bitBufSize : bitBuf's Size maybe    1
* @param bitFifoSize : bitFifo's Size  maybe 16
* @return *bitStructure
 */
func CreateBitStructure(bitBufSize uint8, bitFifoSize uint8) *bitStructure {
	return &bitStructure{
		bitBuf:         make([]byte, bitBufSize),
		bitFifo:        make([]byte, bitFifoSize),
		bitBufCapcity:  bitBufSize,
		bitFifoCapcity: bitFifoSize,
		frontIndex:     uint8(0),
		tailIndex:      uint8(0),
		bitSize:        uint8(0),
	}
}

/**
* @fund: enqueue bit data
* @param bitData
 */
func (bs *bitStructure) BitFifoEnqueue(bitData uint8) error {
	if bs.bitSize >= bs.bitFifoCapcity {
		fmt.Println("bit fifo is full")
		return errors.New("enqueue bit fifo error")
	}
	bs.bitFifo[bs.tailIndex] = bitData
	bs.tailIndex = (bs.tailIndex + 1) % bs.bitFifoCapcity
	bs.bitSize++
	return nil
}

func (bs *bitStructure) BitFifoDequeue() (uint8, error) {
	if bs.bitSize <= 0 {
		fmt.Print("bit fifo is empty")
		return 0, errors.New("dequeue bit fifo error")
	}

	bitData := bs.bitFifo[bs.frontIndex]
	bs.frontIndex = (bs.frontIndex + 1) % bs.bitFifoCapcity
	bs.bitSize--

	return bitData, nil
}

func (bs *bitStructure) ReadtoBitBuf(bufr *bufio.Reader) error {
	_, err := bufr.Read(bs.bitBuf)
	if err != nil {
		fmt.Println("read from file error", err)
	}
	return err
}

func (bs *bitStructure) BitBuf2BitFifo() error {
	for j := uint8(0); j < bs.bitBufCapcity; j++ {
		bitData := bs.bitBuf[j]
		for i := 0; i < 8; i++ {
			if err := bs.BitFifoEnqueue((bitData >> (7 - i)) & 0x01); err != nil {
				return err
			}
		}
	}
	return nil
}

func (bs *bitStructure) DebufPrintbitBuf() {
	fmt.Print("bitbuf:")
	for i := uint8(0); i < bs.bitBufCapcity; i++ {
		fmt.Printf("%08b\t", bs.bitBuf[i])
	}
	fmt.Println("")
}

func (bs *bitStructure) DebufPrintbitFifo() {
	fmt.Print("bitfifo:")
	for i := uint8(0); i < bs.bitFifoCapcity; i++ {
		fmt.Printf("%v\t", bs.bitFifo[i])
	}
	fmt.Println("")
}

func (bs *bitStructure) ReadtoBitFifo(bytes int, bufr *bufio.Reader) error {
	for i := 0; i < bytes/int(bs.bitBufCapcity); i++ {
		if err := bs.ReadtoBitBuf(bufr); err != nil {
			return errors.New("read to bitbuf error")
		}
		if err := bs.BitBuf2BitFifo(); err != nil {
			return errors.New("bitbuf to bitfifo error")
		}
	}
	return nil
}

/**
* @func: parse n bits to a number of uint64 type
* @param nbits : the amount of bit to parse
* @param bytes : control the amount of byte to read at once
* @param bufr  : bufio.Reader type
* @param debug : debug mode or not , true or false
* @return a uint64 number and error
 */
func (bs *bitStructure) BitFifo2NBitsNum(nbits uint8, bytes int, bufr *bufio.Reader, debug bool) (uint64, error) {
	bit_count := uint8(0)
	nbitsNum := uint64(0)

	for bit_count < nbits {
		if bs.bitSize+bit_count >= nbits {
			times := int(nbits - bit_count)
			for i := 0; i < times; i++ {
				bitData, err := bs.BitFifoDequeue()
				if err != nil {
					fmt.Println("bitfifo to nbits num error", err)
					return 0, err
				}

				nbitsNum = nbitsNum << 1
				nbitsNum += uint64(bitData)
				bit_count++
			}
			break
		} else {
			for bs.bitSize != uint8(0) {
				bitData, err := bs.BitFifoDequeue()
				if err != nil {
					fmt.Println("bitfifo to nbits num error", err)
					return 0, err
				}

				nbitsNum = nbitsNum << 1
				nbitsNum += uint64(bitData)
				bit_count++
			}
			if err := bs.ReadtoBitFifo(bytes, bufr); err != nil {
				fmt.Println("bitfifo to nbits num error", err)
				return 0, err
			}
			if debug {
				bs.DebufPrintbitBuf()
				bs.DebufPrintbitFifo()
			}

		}
	}
	return nbitsNum, nil
}
