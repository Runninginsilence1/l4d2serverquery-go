package service

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/spf13/cast"
)

// udp, 并返回数据
func serverInfo(address string, timeout time.Duration) ([]byte, error) {
	var err error
	defer func() {
		if err != nil {
			err = fmt.Errorf("serverInfo: %w", err)
		}
	}()

	conn, err := net.Dial("udp", address)
	if err != nil {
		err = fmt.Errorf("无法连接到 %v: %w", address, err)
		return nil, err
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(timeout))

	header := []byte{0xFF, 0xFF, 0xFF, 0xFF}
	packet := []byte("\x54Source Engine Query\x00")

	_, err = conn.Write(append(header, packet...))
	if err != nil {
		err = fmt.Errorf("无法发送数据: %w", err)
		return nil, err
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("无法读取信息: %w", err)
	}

	if !bytes.HasPrefix(buffer, header) {
		err = fmt.Errorf("返回的响应中不包含header: %v", header)
		return nil, err
	}

	if buffer[4] == 65 {
		challengeNumber := buffer[5:]
		packet = append(header, packet...)
		packet = append(packet, challengeNumber...)
		binary.LittleEndian.PutUint32(challengeNumber, 0)
		_, err = conn.Write(packet)
		if err != nil {
			return nil, err
		}

		n, err = conn.Read(buffer)
		if err != nil {
			return nil, err
		}
	}

	if buffer[4] == 73 {
		return buffer[4:n], nil
	} else if string(buffer[4]) == "A" {
		fmt.Println("Server is using Challenge Number")
	} else {
		fmt.Println("Invalid Response")
	}

	return buffer, nil
}

// 解析从 serverInfo 获取的数据
// []any 用来接收所有的解析后的数据
// 这里应该直接根据字节解析, 把处理字段信息的逻辑耦合进去了, 不是很好
func unpackInfo(data []byte) []string {

	checkStringRes := []string{}

	var msg []interface{}
	mainLog := 0
	_ = mainLog

	// Plan B
	sock := 0

	// index bound handle
	if len(data) == 0 {
		return checkStringRes
	}

	for n := 1; n < 17; n++ {
		switch n {
		// 推测为单独一个字节, 转换后添加的结果集中
		case 1, 2, 8, 9, 10, 11, 12, 13, 14:
			tempUnit := int8(data[sock])
			msg = append(msg, tempUnit)
			sock++
		case 7:
			// 推测长度为2, 作用未知
			tempUnit := int16(binary.LittleEndian.Uint16(data[sock : sock+2]))
			msg = append(msg, tempUnit)
			sock += 2
		case 3, 4, 5, 6, 15, 16:
			// checkString作用
			res, dataLen := checkString(data, sock)
			checkStringRes = append(checkStringRes, res)
			sock = sock + dataLen + 1
		}

	}
	// if n == 17 { // long long
	//     tempUnit := binary.LittleEndian.Uint64(data[s : s+8])
	//     msg = append(msg, tempUnit)
	//     s += 8
	// }

	curPlayerNum := cast.ToString(msg[3])
	maxPlayerNum := cast.ToString(msg[4])
	botPlayerNum := cast.ToString(msg[5])

	checkStringRes = append(checkStringRes, curPlayerNum, maxPlayerNum, botPlayerNum)
	return checkStringRes
}

// 将 sock 后的字节弄出来, 然后解析
// 原理我不管了
func checkString(data []byte, sock int) (output string, dataLen int) {
	// 将17后的数据解析出来
	dataThis := data[sock:]
	dataRight := bytes.Split(dataThis, []byte("\x00"))
	dataLen = len(dataRight[0])
	//tag := fmt.Sprintf("<%ds", dataLen)
	var newBytes = make([]byte, dataLen)
	err := binary.Read(bytes.NewReader(data[sock:sock+dataLen]), binary.LittleEndian, newBytes)
	if err != nil {
		return "", 0
	}
	//fmt.Printf("打印解析后的数据: %v\n", string(newBytes))
	return string(newBytes), dataLen
}

type L4d2SeverInfo struct {
	Name          string `json:"name"`
	Map           string `json:"map"`
	Version       string `json:"version"`
	OnlinePlayers int    `json:"online_players"`
	MaxPlayers    int    `json:"max_players"`
	BotPlayers    int    `json:"bot_players"`
}

func ParseL4d2SeverInfo(data []byte) (*L4d2SeverInfo, error) {
	return &L4d2SeverInfo{}, nil
}

func NewL4d2SeverInfo(name string, Map string, version string, onlinePlayers int, maxPlayers int, botPlayers int) *L4d2SeverInfo {
	return &L4d2SeverInfo{Name: name, Map: Map, Version: version, OnlinePlayers: onlinePlayers, MaxPlayers: maxPlayers, BotPlayers: botPlayers}
}

func Query(addr string) (L4d2SeverInfo, error) {
	var r L4d2SeverInfo

	data, err := serverInfo(addr, 3*time.Second)
	if err != nil {
		return L4d2SeverInfo{}, err
	}
	res := unpackInfo(data)
	if len(res) == 0 {
		return r, fmt.Errorf("no data in unpackInfo")
	}

	r.Name = res[0]
	r.Map = res[1]
	r.Version = res[4]
	r.OnlinePlayers = cast.ToInt(res[6])
	r.MaxPlayers = cast.ToInt(res[7])
	r.BotPlayers = cast.ToInt(res[8])

	return r, nil
}

type name struct {
}
