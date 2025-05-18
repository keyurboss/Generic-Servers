package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"syscall"
	"time"

	"github.com/keyurboss/Generic-Servers/go/validator"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
)

type (
	tunnelService struct{}
	Config        struct {
		ServicePath string   `json:"servicePath" validate:"required"`
		Args        []string `json:"args"`
	}
)

var (
	ShouldBeRunning = false
	BgCommand       *exec.Cmd
	ThisConfig      *Config
)

func (m *tunnelService) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	tick := time.Tick(30 * time.Second)

	status <- svc.Status{State: svc.StartPending}

	status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
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
	StartService()
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
				StopService()
				break loop
			case svc.Pause:
				log.Print("Service Paused.....!")
				StopService()
				status <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				log.Print("Service Continue.....!")
				StartService()
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
	if isDebug {
		err := debug.Run(name, &tunnelService{})
		if err != nil {
			log.Fatalln("Error running service in debug mode.")
			log.Fatalf("Error running service in debug mode. %s", err.Error())
		}
	} else {
		err := svc.Run(name, &tunnelService{})
		if err != nil {
			log.Fatalln("Error running service in Service Control mode.")
		}
	}
}

func main() {
	if !IsDebug() {
		f, err := os.OpenFile(filepath.Join(FindAndReturnCurrentDir(), "debug.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln(fmt.Errorf("error opening file: %v", err))
		}
		defer f.Close()
		log.SetOutput(f)
	}
	start()
	runService("tunnelService", IsDebug())
	log.Print("Finished Function")
}

func start() {
	fmt.Println(len(os.Args), os.Args)
	currentDir := FindAndReturnCurrentDir()
	configFilePAth := filepath.Join(currentDir, "configs.json")
	ThisConfig = &Config{}

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
}

func StartService() {
	ShouldBeRunning = true
	go runningService()
}

func CreateBgService() {
	BgCommand = exec.Command(ThisConfig.ServicePath, ThisConfig.Args...)
	BgCommand.Stdout = log.Default().Writer()
	BgCommand.Stderr = log.Default().Writer()
	log.Print("Service Created...!")
}

func runningService() {
	time.Sleep(3 * time.Second)
	if ShouldBeRunning {
		CreateBgService()
		BgCommand.Start()
		log.Print("Service Started...!")
		BgCommand.Wait()
		log.Print("Service Stopped Automatically...!")
		BgCommand.Process.Kill()
		runningService()
	}
}
func StopService() {
	ShouldBeRunning = false
	BgCommand.Process.Signal(syscall.SIGINT)
	BgCommand.Process.Kill()
	log.Print("Service Stopped...!")
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
