package main

import (
	"encoding/binary"
	"log"
	"math"
	"os"
	"sync/atomic"
	"time"

	"github.com/goburrow/modbus"

	"github.com/tbrandon/mbserver"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %s", err)
		os.Exit(1)
	}
}

func run() error {
	serv := mbserver.NewServer()
	startTime := uint32(time.Now().Unix() & 0xffffffff)
	staMsb := uint16((startTime >> 16) & 0xffff)
	staLsb := uint16(startTime & 0xffff)
	log.Printf("msb:%d lsb:%d", staMsb, staLsb)

	serv.HoldingRegisters[100] = 0xff00 // 65280
	serv.HoldingRegisters[101] = 0xffff // 65535 or -1
	serv.HoldingRegisters[102] = 0x0000 // 0
	serv.HoldingRegisters[200] = 0x0000 // "artificially generates error" 201-208,210-211
	serv.HoldingRegisters[300] = 0x0000 // uptime msb
	serv.HoldingRegisters[301] = 0x0000 // uptime lsb
	serv.HoldingRegisters[302] = staMsb // application start time msb
	serv.HoldingRegisters[303] = staLsb // application start time lsb
	serv.HoldingRegisters[400] = 0x0000 // unixtime msb
	serv.HoldingRegisters[401] = 0x0000 // unixtime lsb
	serv.HoldingRegisters[500] = 0x0000 // math.pi msb
	serv.HoldingRegisters[501] = 0x0000 // math.pi lsb

	serv.InputRegisters[100] = 0xff00 // 65280
	serv.InputRegisters[101] = 0xffff // 65535 or -1
	serv.InputRegisters[102] = 0x0000 // 0
	serv.InputRegisters[200] = 0x0000 // "artificially generates error" 201-208,210-211
	serv.InputRegisters[300] = 0x0000 // uptime msb
	serv.InputRegisters[301] = 0x0000 // uptime lsb
	serv.InputRegisters[302] = staMsb // application start time msb
	serv.InputRegisters[303] = staLsb // application start time lsb
	serv.InputRegisters[400] = 0x0000 // unixtime msb
	serv.InputRegisters[401] = 0x0000 // unixtime lsb
	serv.InputRegisters[500] = 0x0000 // math.pi msb
	serv.InputRegisters[501] = 0x0000 // math.pi lsb

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		var uptime uint32 = 0
		for {
			<-ticker.C
			atomic.AddUint32(&uptime, 1)
		}
	}()

	serv.RegisterFunctionHandler(modbus.FuncCodeReadInputRegisters, ReadHoldingRegisters) //note this is a hack
	serv.RegisterFunctionHandler(modbus.FuncCodeReadHoldingRegisters, ReadHoldingRegisters)

	listenAddr := "0.0.0.0:1502"

	log.Printf("Modbus Server listening on %s", listenAddr)
	err := serv.ListenTCP(listenAddr)
	if err != nil {
		return err
	}
	defer serv.Close()

	// Wait forever
	for {
		time.Sleep(1 * time.Second)
	}
}

//func ReadInputRegisters(s *mbserver.Server, frame mbserver.Framer) ([]byte, *mbserver.Exception) {
//	register, numRegs, endRegister := registerAddressAndNumber(frame)
//	if endRegister > 65536 {
//		return []byte{}, &mbserver.IllegalDataAddress
//	}
//	return append([]byte{byte(numRegs * 2)}, mbserver.Uint16ToBytes(s.InputRegisters[register:endRegister])...), &mbserver.Success
//}

func ReadHoldingRegisters(s *mbserver.Server, frame mbserver.Framer) ([]byte, *mbserver.Exception) {
	register, numRegs, endRegister := registerAddressAndNumber(frame)
	if endRegister > 65536 {
		return []byte{}, &mbserver.IllegalDataAddress
	}
	log.Printf("Read r:%d, count:%d", register, numRegs)

	// get the current unix timestamp, converted as a 32-bit unsigned integer for simplicity
	unixTs := uint32(time.Now().Unix() & 0xffffffff)

	switch register {
	//200 "artificially generates error"
	case 201:
		return []byte{}, &mbserver.IllegalFunction
	case 202:
		return []byte{}, &mbserver.IllegalDataAddress
	case 203:
		return []byte{}, &mbserver.IllegalDataValue
	case 204:
		return []byte{}, &mbserver.SlaveDeviceFailure
	case 205:
		return []byte{}, &mbserver.AcknowledgeSlave
	case 206:
		return []byte{}, &mbserver.SlaveDeviceBusy
	case 207:
		return []byte{}, &mbserver.NegativeAcknowledge
	case 208:
		return []byte{}, &mbserver.MemoryParityError
	case 210:
		return []byte{}, &mbserver.GatewayPathUnavailable
	case 211:
		return []byte{}, &mbserver.GatewayTargetDeviceFailedtoRespond

	//300 uptime msb
	//301 uptime lsb
	case 300:
		fallthrough
	case 301:
		if numRegs > 2 {
			return []byte{}, &mbserver.IllegalDataAddress
		}
		time32 := uint32(time.Now().Unix() & 0xffffffff)

		staMsb := s.HoldingRegisters[302] // application start time msb
		staLsb := s.HoldingRegisters[303] // application start time lsb
		startTime := uint32(staMsb)<<16 + uint32(staLsb)

		uptime := time32 - startTime

		msb := uint16((uptime >> 16) & 0xffff)
		lsb := uint16(uptime & 0xffff)
		bits := []uint16{msb, lsb}
		return append([]byte{byte(numRegs * 2)}, mbserver.Uint16ToBytes(bits[register-300:endRegister-300])...), &mbserver.Success

	//400 unixtime msb
	//401 unixtime lsb
	case 400:
		fallthrough
	case 401:
		startingAddress := 400
		if numRegs > 2 {
			return []byte{}, &mbserver.IllegalDataAddress
		}
		msb := uint16((unixTs >> 16) & 0xffff)
		lsb := uint16(unixTs & 0xffff)
		bits := []uint16{msb, lsb}
		return append([]byte{byte(numRegs * 2)}, mbserver.Uint16ToBytes(bits[register-startingAddress:endRegister-startingAddress])...), &mbserver.Success

	//500 math.pi msb
	case 500:
		fallthrough
	//501 math.pi lsb
	case 501:
		startingAddress := 500
		if numRegs > 2 {
			return []byte{}, &mbserver.IllegalDataAddress
		}
		pi32 := math.Float32bits(math.Pi)
		msb := uint16((pi32 >> 16) & 0xffff)
		lsb := uint16(pi32 & 0xffff)
		bits := []uint16{msb, lsb}
		return append([]byte{byte(numRegs * 2)}, mbserver.Uint16ToBytes(bits[register-startingAddress:endRegister-startingAddress])...), &mbserver.Success
	}
	return append([]byte{byte(numRegs * 2)}, mbserver.Uint16ToBytes(s.HoldingRegisters[register:endRegister])...), &mbserver.Success
}

func registerAddressAndNumber(frame mbserver.Framer) (register int, numRegs int, endRegister int) {
	data := frame.GetData()
	register = int(binary.BigEndian.Uint16(data[0:2]))
	numRegs = int(binary.BigEndian.Uint16(data[2:4]))
	endRegister = register + numRegs
	return register, numRegs, endRegister
}
