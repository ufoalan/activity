package control_ev3

import (
	"errors"
	"fmt"
	//"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/ev3go/ev3dev"
)

// log is the default package logger
var log = logger.GetLogger("activity-tibco-rest")

const (
	method         = "method"
	pinNumber      = "pinNumber"
	directionState = "direction"
	state          = "state"
	direction      = "Direction"
	setState       = "Set State"
	readState      = "Read State"
	pull           = "Pull"
	start          = "start"
	stop           = "stop"
	auto           = "auto"

	input = "Input"
	//output = "Output"

	high = "High"
	//low = "Low"

	up   = "Up"
	down = "Down"
	//off = "off"

	//ouput

	result = "result"
)

type GPIOActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new GPIOActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &GPIOActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *GPIOActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *GPIOActivity) Eval(context activity.Context) (done bool, err error) {
	//getmethod
	log.Debug("Running control_ev3 activity.")
	methodInput := context.GetInput(method)

	ivmethod, ok := methodInput.(string)
	if !ok {
		return true, errors.New("Method field not set.")
	}

	//get pinNumber
	ivPinNumber, ok := context.GetInput(pinNumber).(int)

	if !ok {
		return true, errors.New("Pin number must exist")
	}

	log.Debugf("Method '%s' and pin number '%d'", methodInput, ivPinNumber)
	//Open pin
	//openErr := rpio.Open()
	//if openErr != nil {
	//	log.Errorf("Open RPIO error: %+v", openErr.Error())
	//	return true, errors.New("Open RPIO error: " + openErr.Error())
	//}

	//pin := rpio.Pin(ivPinNumber)

        outA, err := ev3dev.TachoMotorFor("outA", "lego-ev3-l-motor")
        if err != nil {
                log.Debugf("failed to find large motor on outA: %v", err)
        }
        err = outA.SetStopAction("brake").Err()
        if err != nil {
                log.Debugf("failed to set brake stop for large motor on outA: %v", err)
        }
        maxMedium := outA.MaxSpeed()


	switch ivmethod {
	case start:
                outA.SetSpeedSetpoint(50 * maxMedium / 100).Command("run-forever")
                checkErrors(outA)
	case stop:
                outA.Command("stop")
                checkErrors(outA)
	case auto:
                for i := 0; i < 2; i++ {

                        // Run medium motor on outA at speed 50, wait for 0.5 second and then brake.
                        outA.SetSpeedSetpoint(50 * maxMedium / 100).Command("run-forever")
                        time.Sleep(time.Second / 2)
                        outA.Command("stop")
                        checkErrors(outA)

                        // Run medium motor on outA at speed -75, wait for 0.5 second and then brake.
                        outA.SetSpeedSetpoint(-75 * maxMedium / 100).Command("run-forever")
                        time.Sleep(time.Second / 2)
                        outA.Command("stop")
                        checkErrors(outA)
                }
	default:
		log.Errorf("Cannot found method %s ", ivmethod)
		return true, errors.New("Cannot found method %s " + ivmethod)
	}

	context.SetOutput(result, 0)
	return true, nil
}

func checkErrors(devs ...ev3dev.Device) {
        for _, d := range devs {
                err := d.(*ev3dev.TachoMotor).Err()
                if err != nil {
                        drv, dErr := ev3dev.DriverFor(d)
                        if dErr != nil {
                                drv = fmt.Sprintf("(missing driver name: %v)", dErr)
                        }
                        addr, aErr := ev3dev.AddressOf(d)
                        if aErr != nil {
                                drv = fmt.Sprintf("(missing port address: %v)", aErr)
                        }
                        log.Debugf("motor error for %s:%s on port %s: %v", d, drv, addr, err)
                }
        }
}
