package surah

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var (
	surahCache any
	cacheMutex sync.RWMutex
)

type SurahHandler struct{}

func NewSurahHandler() *SurahHandler {
	return &SurahHandler{}
}

func (h *SurahHandler) GetSurahs(c *fiber.Ctx) error {
	data := getSurahs()

	if data == nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data dari API pusat",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func getSurahs() any {

	// Read lock
	cacheMutex.RLock()
	if surahCache != nil {
		defer cacheMutex.RUnlock()
		return surahCache
	}
	cacheMutex.RUnlock()

	// Write lock
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Double check
	if surahCache != nil {
		return surahCache
	}

	agent := fiber.Get("https://equran.id/api/v2/surat")

	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    []struct {
			Nomor       int               `json:"nomor"`
			Nama        string            `json:"nama"`
			NamaLatin   string            `json:"namaLatin"`
			JumlahAyat  int               `json:"jumlahAyat"`
			TempatTurun string            `json:"tempatTurun"`
			Arti        string            `json:"arti"`
			Deskripsi   string            `json:"deskripsi"`
			AudioFull   map[string]string `json:"audioFull"`
		} `json:"data"`
	}

	agent.Timeout(10 * time.Second)

	statusCode, body, errs := agent.Bytes()

	if len(errs) > 0 {
		log.Error().Errs("errors", errs).Msg("Error Fetching Surahs")
		return nil
	}

	if statusCode != 200 {
		log.Error().Int("status_code", statusCode).Msg("Status Code Not 200 when fetching Surahs")
		return nil
	}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Error().Err(err).Msg("Unmarshal Error when fetching Surahs")
		return nil
	}

	surahCache = response.Data
	return surahCache
}
