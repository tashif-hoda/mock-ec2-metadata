package metadata

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServe(t *testing.T) {
	metadataService := NewMetaDataService()
	go metadataService.Serve()

	b, err := json.Marshal(metadataService.config.MetadataValues)
	if err != nil {
		panic(err)
	}
	var validation map[string]interface{}
	err = json.Unmarshal(b, &validation)
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)

	for _, metadataPrefix := range metadataService.config.MetadataPrefixes {
		for _, endpoint := range listOfEndpoints["MetadataPrefix"] {
			endpoint = strings.ReplaceAll(endpoint, "{username}", metadataService.config.MetadataValues.User)
			resp, err := http.Get("http://localhost" + metadataPrefix + endpoint)
			if !assert.Equal(t, 200, resp.StatusCode, fmt.Sprintf("StatusCode != 200 status code: %v", resp.StatusCode)) {
				continue
			}
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				assert.Fail(t, fmt.Sprintf("Failed to read response. Error %v", err))
				continue
			}
			bodyString := string(bodyBytes)
			log.Println("response", bodyString)
			// assert.Equal(t, validation[endpoint], bodyString)
		}
	}
}
