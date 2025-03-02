package viewer

import (
	"fmt"
	"net/http"
	"net/url"

	"log"
	"os"

	yaml "gopkg.in/yaml.v3"

	"github.com/streletsa/savings-aggreagator/collector"
)

type Esp32SavingsViewer struct {
	Config Esp32SavingsViewerConfig
}

func (x Esp32SavingsViewer) View(info *collector.SavingsCollectionInfo) {
	u, err := url.Parse(x.Config.ESP32URL)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	fmt.Println(u)

	q := u.Query()
	q.Set("total", fmt.Sprintf("%.2f", info.Total))
	q.Set("source", string(info.SourceType))

	u.RawQuery = q.Encode()

	_, err = http.Post(u.String(), "application/json", nil)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
}

type Esp32SavingsViewerConfig struct {
	ESP32URL string `yaml:"ESP32URL"`
}

func LoadConfig(filename string) Esp32SavingsViewerConfig {
	var c Esp32SavingsViewerConfig
	input, err := os.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return Esp32SavingsViewerConfig{}
	}
	err = yaml.Unmarshal(input, &c)
	if err != nil {
		log.Println(err)
	}
	return c
}
