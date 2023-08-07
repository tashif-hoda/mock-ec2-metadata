package metadata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type (
	SecurityCredentials struct {
		AccessKeyId     string `json:"AccessKeyId"`
		SecretAccessKey string `json:"SecretAccessKey"`
		Token           string `json:"Token"`
		Expiration      string `json:"Expiration"`
		Code            string `json:"Code"`
	}

	Network struct {
		Interfaces map[string][]string `json:"interfaces"`
	}

	MetadataValues struct {
		AmiId               string              `json:"ami-id"`
		AmiLaunchIndex      string              `json:"ami-launch-index"`
		AmiManifestPath     string              `json:"ami-manifest-path"`
		AvailabilityZone    string              `json:"availability-zone"`
		Hostname            string              `json:"hostname"`
		InstanceAction      string              `json:"instance-action"`
		InstanceId          string              `json:"instance-id"`
		InstanceType        string              `json:"instance-type"`
		LocalHostName       string              `json:"local-hostname"`
		LocalIpv4           string              `json:"local-ipv4"`
		Mac                 string              `json:"mac"`
		Profile             string              `json:"profile"`
		ReservationId       string              `json:"reservation-id"`
		User                string              `json:"user"`
		SecurityGroups      []string            `json:"security-groups"`
		SecurityCredentials SecurityCredentials `json:"security-credentials"`
		Network             Network             `json:"network"`
	}

	Config struct {
		router           *mux.Router
		Server           struct{ HTTPPort int }
		MetadataValues   *MetadataValues
		MetadataPrefixes []string
		UserdataValues   map[string]string
		UserdataPrefixes []string
		NetworkPrefixes  []string
	}

	MetadataService struct {
		config *Config
	}
)

var listOfEndpoints = map[string][]string{
	"MetadataPrefix": {
		"/",
		"/ami-id",
		"/ami-launch-index",
		"/ami-manifest-path",
		"placement/availability-zone",
		"/hostname",
		"/instance-action",
		"/instance-id",
		"/instance-type",
		"/iam/",
		"/iam/security-credentials",
		"/iam/security-credentials/",
		"/iam/security-credentials/{username}",
		"/local-hostname",
		"/local-ipv4",
		"/mac",
		"/profile",
		"/reservation-id",
		"/security-groups",
		"/network/interfaces/macs/00:00:00:00:00:00/security-group-ids",
	},
	"UserdataPrefix": {
		"/",
	},
}

// middleware for adding plaintext content type
func plainText(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		h(w, r)
	}
}

func movedPermanently(redirectPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectPath, http.StatusMovedPermanently)
	}
}

func (s *MetadataService) GetAmiId(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.AmiId)
}

func (s *MetadataService) GetAmiLaunchIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.AmiLaunchIndex)
}

func (s *MetadataService) GetAmiManifestPath(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.AmiManifestPath)
}

func (s *MetadataService) GetAvailabilityZone(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.AvailabilityZone)
}

func (s *MetadataService) GetHostName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.Hostname)
}

func (s *MetadataService) GetInstanceAction(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.InstanceAction)
}

func (s *MetadataService) GetInstanceId(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.InstanceId)
}

func (s *MetadataService) GetInstanceType(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.InstanceType)
}

func (s *MetadataService) GetLocalHostName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.LocalHostName)
}

func (s *MetadataService) GetLocalIpv4(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.LocalIpv4)
}

func (s *MetadataService) GetIAM(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "security-credentials/")
}

func (s *MetadataService) GetMac(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.Mac)
}

func (s *MetadataService) GetProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.Profile)
}

func (s *MetadataService) GetReservationId(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.ReservationId)
}

func (s *MetadataService) GetSecurityCredentials(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.MetadataValues.User)
}

func (s *MetadataService) GetSecurityGroups(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, strings.Join(s.config.MetadataValues.SecurityGroups, "\n"))
}

func (s *MetadataService) GetSecurityGroupIds(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, strings.Join(s.config.MetadataValues.Network.Interfaces["00:00:00:00:00:00"], "\n"))
}

func (s *MetadataService) GetSecurityCredentialDetails(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	if username != s.config.MetadataValues.User {
		log.Println("error, IAM user not found")
		http.Error(w, "", http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(w).Encode(s.config.MetadataValues.SecurityCredentials)
	if err != nil {
		log.Println("error converting security credentails to json: ", err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	log.Printf("GetSecurityCredentialDetails returning: %#v",
		s.config.MetadataValues.SecurityCredentials)
}

func (s *MetadataService) GetMetadataIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `hostname
instance-id
instance-type
iam`)
}

func (s *MetadataService) GetUserData(w http.ResponseWriter, r *http.Request) {

	for index, value := range s.config.UserdataValues {
		fmt.Fprintf(w, fmt.Sprint(index+"="+value+"\n"))
	}
}

func (s *MetadataService) GetIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Mock EC2 Metadata Service")
}

func (service *MetadataService) Endpoints() (handlers map[string]map[string]http.HandlerFunc) {
	handlers = make(map[string]map[string]http.HandlerFunc)
	for index, value := range service.config.MetadataPrefixes {
		log.Println("adding Metadata prefix (", index, ") ", value)
		handlers[value+listOfEndpoints["MetadataPrefix"][0]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetMetadataIndex),
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][1]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetAmiId),
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][2]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetAmiLaunchIndex),
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][3]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetAmiManifestPath),
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][4]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetAvailabilityZone),
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][5]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetHostName),
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][6]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetInstanceAction),
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][7]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetInstanceId),
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][8]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetInstanceType),
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][9]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetIAM),
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][10]] = map[string]http.HandlerFunc{
			"GET": movedPermanently(value + listOfEndpoints["MetadataPrefix"][11]),
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][11]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetSecurityCredentials),
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][12]] = map[string]http.HandlerFunc{
			"GET": service.GetSecurityCredentialDetails,
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][13]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetLocalHostName),
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][14]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetLocalIpv4),
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][15]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetMac),
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][16]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetProfile),
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][17]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetReservationId),
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][18]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetSecurityGroups),
		}
		handlers[value+listOfEndpoints["MetadataPrefix"][19]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetSecurityGroupIds),
		}
	}

	for index, value := range service.config.UserdataPrefixes {
		log.Println("adding Userdata prefix (", index, ") ", value)

		handlers[value+listOfEndpoints["UserdataPrefix"][0]] = map[string]http.HandlerFunc{
			"GET": plainText(service.GetUserData),
		}
	}

	handlers["/"] = map[string]http.HandlerFunc{
		"GET": service.GetIndex,
	}
	return
}

func (service *MetadataService) RegisterEndpoints() {
	for endpoint, handlers := range service.Endpoints() {
		service.config.router.HandleFunc(endpoint, handlers["GET"])
	}
}

func (service *MetadataService) Serve() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", service.config.Server.HTTPPort), service.config.router)
}

func NewMetaDataService() (metadataService *MetadataService) {
	metadataService = &MetadataService{
		config: &Config{router: mux.NewRouter()},
	}
	if _, err := os.Stat("./mock-ec2-metadata-config.json"); err == nil {
		LoadJSONFile("./mock-ec2-metadata-config.json", metadataService.config)
	} else if _, err := os.Stat("/etc/mock-ec2-metadata-config.json"); err == nil {
		LoadJSONFile("/etc/mock-ec2-metadata-config.json", metadataService.config)
	} else {
		log.Fatal("unable to locate config file. Please place mock-ec2-metadata-config.json in current directory.")
	}
	metadataService.RegisterEndpoints()
	http.Handle("/", metadataService.config.router)
	return metadataService
}

func LoadJSONFile(fileName string, cfg interface{}) {
	cb, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Unable to read config file '%s': %s", fileName, err)
	}

	if err = json.Unmarshal(cb, &cfg); err != nil {
		log.Fatalf("Unable to parse JSON in config file '%s': %s", fileName, err)
	}
}
