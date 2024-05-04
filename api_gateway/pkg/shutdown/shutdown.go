package shutdown

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
)

// gracful shutdown
// вспомогательная функция для корректного освобождения ресурсов
// при получении сигналов от ОС
func Graceful(signals []os.Signal, log *slog.Logger, closeItems ...io.Closer) {

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, signals...)
	sig := <-sigc
	log.Warn(fmt.Sprintf("Caught signal %s. Shutting down...", sig))

	// Here we can do graceful shutdown (close connections and etc)
	for _, closer := range closeItems {
		if err := closer.Close(); err != nil {
			log.Error("failed to close %v: %v", closer, err)
		}
	}
}
