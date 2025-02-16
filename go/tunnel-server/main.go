package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"time"

	chclient "github.com/jpillora/chisel/client"
	chshare "github.com/jpillora/chisel/share"
	"github.com/jpillora/chisel/share/cos"

	"github.com/keyurboss/Generic-Servers/go/validator"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
)

type (
	tunnelService struct{}
	Config        struct {
		ServerIp           string         `json:"serverIp" validate:"required"`
		ServerPort         int            `json:"serverPort" validate:"required,port"`
		RemoteTunnelConfig []TunnelConfig `json:"remoteTunnelConfig" validate:"required,min=1"`
	}
	TunnelConfig struct {
		LocalPort  int `json:"localPort" validate:"required,port"`
		RemotePort int `json:"remotePort" validate:"required,port"`
	}
)

var (
	Client          *chclient.Client
	ChiselContext   context.Context
	ShouldBeRunning = false
)

// Execute(args []string, r <-chan ChangeRequest, s chan<- Status) (svcSpecificEC bool, exitCode uint32)

func (m *tunnelService) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	tick := time.Tick(30 * time.Second)

	status <- svc.Status{State: svc.StartPending}

	status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	StartChisel()

	// if IsDebug() {
	// 	log.Print("In Debug Mode")
	// 	go func() {
	// 		time.Sleep(15 * time.Second)
	// 		r <- svc.ChangeRequest{Cmd: svc.Pause, CurrentStatus: svc.Status{State: svc.Paused}}
	// 		time.Sleep(15 * time.Second)
	// 		r <- svc.ChangeRequest{Cmd: svc.Continue, CurrentStatus: svc.Status{State: svc.Running}}
	// 		time.Sleep(15 * time.Second)
	// 		r <- svc.ChangeRequest{Cmd: svc.Stop, CurrentStatus: svc.Status{State: svc.Stopped}}
	// 	}()
	// }
loop:
	for {
		select {
		case <-tick:
			log.Print("Tick Handled...!")
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				log.Print("Service Interrogate.....!", c.CurrentStatus)
				status <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				log.Print("Shutting service...!")
				StopChisel()
				break loop
			case svc.Pause:
				log.Print("Service Paused.....!")
				StopChisel()
				status <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				log.Print("Service Continue.....!")
				StartChisel()
				status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			default:
				log.Printf("Unexpected service control request #%d", c)
			}
		}
	}

	status <- svc.Status{State: svc.StopPending}
	return false, 1
}

func runService(name string, isDebug bool) {
	structS := &tunnelService{}
	if isDebug {
		err := debug.Run(name, structS)
		if err != nil {
			log.Fatalln("Error running service in debug mode.")
		}
	} else {
		err := svc.Run(name, structS)
		if err != nil {
			log.Fatalln("Error running service in Service Control mode.")
		}
	}
}

func main() {
	f, err := os.OpenFile(filepath.Join(FindAndReturnCurrentDir(), "debug.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(fmt.Errorf("error opening file: %v", err))
	}
	defer f.Close()

	log.SetOutput(f)
	start()
	runService("tunnelService", IsDebug())
	log.Print("Finished Function")
}

func start() {
	fmt.Println(len(os.Args), os.Args)
	currentDir := FindAndReturnCurrentDir()
	configFilePAth := filepath.Join(currentDir, "configs.json")
	ThisConfig := &Config{}

	fmt.Printf("Config path %s\n", configFilePAth)
	if _, err := os.Stat(configFilePAth); errors.Is(err, os.ErrNotExist) {
		panic(fmt.Errorf("Config Not Exist on Path %s", configFilePAth))
	}
	dat, err := os.ReadFile(configFilePAth)
	Check(err)
	json.Unmarshal(dat, ThisConfig)

	if errs := validator.Validator.Validate(ThisConfig); len(errs) > 0 {
		panic(fmt.Errorf("Config Error %#v", errs))
	}

	remote := make([]string, len(ThisConfig.RemoteTunnelConfig))
	for i, r := range ThisConfig.RemoteTunnelConfig {
		remote[i] = fmt.Sprintf("R:%d:%d", r.RemotePort, r.LocalPort)
	}
	chshare.BuildVersion = "1.9.1"
	client, err := chclient.NewClient(&chclient.Config{
		Remotes:          remote,
		MaxRetryCount:    100,
		MaxRetryInterval: 60 * time.Minute,
		Server:           fmt.Sprintf("%s:%d", ThisConfig.ServerIp, ThisConfig.ServerPort),
	})
	if err != nil {
		log.Print(err.Error())
	}
	Client = client
	ChiselContext = cos.InterruptContext()
}

func StartChisel() {
	ShouldBeRunning = true
	go runningChisel()
}

func runningChisel() {
	time.Sleep(3 * time.Second)
	if ShouldBeRunning {
		log.Print("Chisel Started...!")
		Client.Start(ChiselContext)
		Client.Wait()
		log.Print("Chisel Stopped...!")
		runningChisel()
	}
}

func StopChisel() {
	ShouldBeRunning = false
	Client.Close()
	log.Print("Chisel Stopped...!")
}

func IsDebug() bool {
	return slices.Contains(os.Args, "--dev")
}

func FindAndReturnCurrentDir() string {
	currentDir := ""
	if IsDebug() {
		current, err := os.Getwd()
		Check(err)
		currentDir = current
	} else {
		exePath, err := os.Executable()
		currentDir = filepath.Dir(exePath)
		Check(err)
	}
	return currentDir
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
