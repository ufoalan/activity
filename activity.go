package ev3

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/ev3go/ev3dev"
)

// THIS IS ADDED
// log is the default package logger which we'll use to log
var log = logger.GetLogger("activity-ev3")

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, evalError *activity.Error) {

	// do eval
	// Get the activity data from the context
	action := context.GetInput("action").(string)
	speed := context.GetInput("speed").(int)

	// Use the log object to log the greeting
	log.Debugf("The Flogo engine says motor [%s] at [%s] speed", action, speed)

	        // Get the handle for the medium motor on outA.
        outA, err := ev3dev.TachoMotorFor("outA", "lego-ev3-l-motor")
        if err != nil {
                log.Fatalf("failed to find large motor on outA: %v", err)
        }
        err = outA.SetStopAction("brake").Err()
        if err != nil {
                log.Fatalf("failed to set brake stop for large motor on outA: %v", err)
        }
        maxMedium := outA.MaxSpeed()
       	for i := 0; i < 2; i++ {
                // Render the gopher to the screen.
                //draw.Draw(ev3.LCD, ev3.LCD.Bounds(), gopher, gopher.Bounds().Min, draw.Src)

                // Run medium motor on outA at speed 50, wait for 0.5 second and then brake.
                outA.SetSpeedSetpoint(50 * maxMedium / 100).Command("run-forever")
                time.Sleep(time.Second / 2)
                outA.Command("stop")
                checkErrors(outA)

                // Run large motors on B+C at speed 70, wait for 2 second and then brake.
                //outB.SetSpeedSetpoint(70 * maxLarge / 100).Command("run-forever")
                //outC.SetSpeedSetpoint(70 * maxLarge / 100).Command("run-forever")
                //checkErrors(outB, outC)
                //time.Sleep(2 * time.Second)
                //outB.Command("stop")
                //outC.Command("stop")
                //checkErrors(outB, outC)

                // Run medium motor on outA at speed -75, wait for 0.5 second and then brake.
                outA.SetSpeedSetpoint(-75 * maxMedium / 100).Command("run-forever")
                time.Sleep(time.Second / 2)
                outA.Command("stop")
                checkErrors(outA)
	}
	
	// Set the result as part of the context
	context.SetOutput("result", "The Flogo engine says motor "+action+" at "+speed+" speed")

	// Signal to the Flogo engine that the activity is completed

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
                        log.Fatalf("motor error for %s:%s on port %s: %v", d, drv, addr, err)
                }                     
        }                                
}
