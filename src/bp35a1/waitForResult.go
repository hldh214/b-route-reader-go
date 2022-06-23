package bp35a1

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

var STD_TIMEOUT_DURATION = 15 * time.Second
var LONG_TIMEOUT_DURATION = 60 * time.Second

// OK等を返すコマンドの応答を返す
func waitForResult() ([]string, error) {
	return waitForResultImpl(RET_STOP_WORDS, STD_TIMEOUT_DURATION)
}

// SKLL64の応答を返す
// このコマンドはいきなりIPv6アドレスだけを返してくる
func waitForResultSKLL64() ([]string, error) {
	return waitForResultImpl([]string{}, STD_TIMEOUT_DURATION)
}

func waitForResultSKSCAN() ([]string, error) {
	return waitForResultImpl([]string{RET_SCAN_COMPLETE}, LONG_TIMEOUT_DURATION)
}

func waitForResultSKJOIN() ([]string, error) {
	return waitForResultImpl([]string{RET_JOIN_COMPLETE}, LONG_TIMEOUT_DURATION)
}

func waitForResultERXUDP() ([]string, error) {
	return waitForResultImpl([]string{RET_ERXUDP}, LONG_TIMEOUT_DURATION)
}

func waitForResultImpl(stopWords []string, timeoutDuration time.Duration) ([]string, error) {

	log.Debug().Msgf("Response start. timeout=%s stop words=[%s]",
		timeoutDuration, strings.Join(stopWords, "|"))
	BYTE_CR := []byte("\r")
	BYTE_LF := []byte("\n")
	BYTE_CRLF := []byte("\r\n")
	LF := byte(0xa)

	port.SetReadTimeout(300 * time.Millisecond)

	timeoutTime := time.Now().Add(timeoutDuration)

	var result []string

	// 応答を全部溜め込むバッファ
	var byteBuf []byte

	stopFlag := false
	if len(stopWords) == 0 {
		stopFlag = true
	}

	for {
		// 一回のreadで読んだデータのバッファ
		readBuf := make([]byte, 200)

		n, err := port.Read(readBuf)
		if err != nil {
			log.Err(err).Msg("Read error")
			return nil, err
		}
		if n == 0 && len(byteBuf) == 0 && stopFlag {
			if len(result) > 0 {
				break
			}
		}

		// タイムアウトまわり
		if n == 0 {
			if time.Now().After(timeoutTime) {
				return result, fmt.Errorf("waitForResult timeout reached")
			}
		} else {
			timeoutTime = time.Now().Add(timeoutDuration) // タイムアウト延長
		}

		byteBuf = append(byteBuf, readBuf[:n]...)

		// byteBufを改行で分割する
		for {
			if bytes.Contains(byteBuf, BYTE_CR) {
				crPos := bytes.Index(byteBuf, BYTE_CR)
				lineBuf := byteBuf[:crPos]
				lineStr := trimResponse(string(lineBuf))

				byteBuf = byteBuf[crPos:] // 改行コードは削除

				// CRLFで区切られている場合、LFが残るので削除
				switch {
				case bytes.HasPrefix(byteBuf, BYTE_CRLF):
					byteBuf = byteBuf[len(BYTE_CRLF):]
					log.Debug().Msgf("<-- %s<CRLF>", lineStr)
				case bytes.HasPrefix(byteBuf, BYTE_LF):
					byteBuf = byteBuf[len(BYTE_LF):]
					log.Debug().Msgf("<-- %s<LF>", lineStr)
				case bytes.HasPrefix(byteBuf, BYTE_CR):
					byteBuf = byteBuf[len(BYTE_CR):]
					log.Debug().Msgf("<-- %s<CR>", lineStr)
				}

				result = append(result, lineStr)

				// ストップワード（通常、コマンド応答の末尾に来るワード）を見つけたら終了フラグを立てる
				// OK を返したあとにさらに応答を返すコマンドがあるため（しかし、そのコマンドは使わない）
				for _, stopWord := range stopWords {
					if strings.Contains(lineStr, stopWord) {
						stopFlag = true
						break
					}
				}
			} else if len(byteBuf) == 1 && byteBuf[0] == LF {
				// SKLL64の時、LFだけがバッファに残ってしまい無限ループすることへの対策
				log.Warn().Msgf("executing start with LF workaround")
				byteBuf = byteBuf[1:]
				break
			} else {
				// 改行コードが含まれない
				// 行の途中でバッファがいっぱいになったか、応答の末尾まで読み切った
				// どちらにしてももう一度readを呼ぶ
				break
			}

		} // 無限ループ
	}

	log.Debug().Msg("Response done")
	return result, nil
}

func trimResponse(str string) string {
	result := str
	result = strings.TrimPrefix(result, " ")
	result = strings.TrimPrefix(result, "\r")
	result = strings.TrimPrefix(result, "\n")

	return result

}
