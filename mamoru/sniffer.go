package mamoru

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Mamoru-Foundation/mamoru-sniffer-go/mamoru_sniffer"
	"github.com/ethereum/go-ethereum/log"
)

var (
	connectMutex       sync.Mutex
	sniffer            *mamoru_sniffer.Sniffer
	SnifferConnectFunc = mamoru_sniffer.Connect
)

const DefaultDelta int64 = 50

type statusProgress interface {
	Process() bool
}

func mapToInterfaceSlice(m map[string]string) []interface{} {
	var result []interface{}
	for key, value := range m {
		result = append(result, key, value)
	}

	return result
}

func init() {
	mamoru_sniffer.InitLogger(func(entry mamoru_sniffer.LogEntry) {
		kvs := mapToInterfaceSlice(entry.Ctx)
		msg := "Mamoru core: " + entry.Message
		switch entry.Level {
		case mamoru_sniffer.LogLevelDebug:
			log.Debug(msg, kvs...)
		case mamoru_sniffer.LogLevelInfo:
			log.Info(msg, kvs...)
		case mamoru_sniffer.LogLevelWarning:
			log.Warn(msg, kvs...)
		case mamoru_sniffer.LogLevelError:
			log.Error(msg, kvs...)
		}
	})
}

type Sniffer struct {
	mu     sync.Mutex
	status statusProgress
	synced bool
}

func NewSniffer() *Sniffer {
	return &Sniffer{}
}

func (s *Sniffer) checkSynced() bool {
	//if status is nil, we assume that node is synced
	if s.status == nil {
		return true
	}

	s.synced = s.status.Process()

	log.Info("Mamoru Sniffer sync", "syncing", s.synced)

	return s.synced
}

func (s *Sniffer) SetDownloader(downloader statusProgress) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status = downloader
}

func (s *Sniffer) CheckRequirements() bool {
	return isSnifferEnable() && connect() && s.checkSynced()
}

func isSnifferEnable() bool {
	val, ok := os.LookupEnv("MAMORU_SNIFFER_ENABLE")
	if !ok {
		return false
	}
	isEnable, err := strconv.ParseBool(val)
	if err != nil {
		log.Error("Mamoru Sniffer env parse error", "err", err)
		return false
	}

	return ok && isEnable
}

func connect() bool {
	connectMutex.Lock()
	defer connectMutex.Unlock()
	if sniffer != nil {
		return true
	}
	var err error
	if sniffer == nil {
		sniffer, err = SnifferConnectFunc()
		if err != nil {
			erst := strings.Replace(err.Error(), "\t", "", -1)
			erst = strings.Replace(erst, "\n", "", -1)
			log.Error("Mamoru Sniffer connect", "err", erst)
			return false
		}
	}
	return true
}

func getDelta() int64 {
	val, ok := os.LookupEnv("MAMORU_SNIFFER_DELTA")
	if !ok {
		return DefaultDelta
	}

	delta, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		log.Error("Mamoru Sniffer env parse error", "err", err)
		return DefaultDelta
	}

	return delta
}
