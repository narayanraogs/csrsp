package das

import (
	"csrsp/server/db"
	"csrsp/server/global"
	"csrsp/server/utils/binary"
	"csrsp/server/utils/netutil"
	"sync"
	"time"
)

type DASMonitor struct {
	DasName   string
	DpuNumber int
	Status    string
	Alarm     bool
}

func GetStatus(acqMode string) ([]DASMonitor, error) {
	compCode := global.App.DasLockStatus.ComponentCode
	lockOffset := global.App.DasLockStatus.LockOffset*4 + 4
	lockValue := global.App.DasLockStatus.LockValue
	alarmOffset := global.App.DasLockStatus.AlarmOffset*4 + 4
	alarmValue := global.App.DasLockStatus.AlarmValue
	monitorPort := global.App.DasPorts.Monitor

	dasPaths, err := db.GetAcquisitionChainDetails(acqMode)
	if err != nil {
		return nil, err
	}

	var dasMonitor = make([]DASMonitor, len(dasPaths))
	var wg sync.WaitGroup

	for i, path := range dasPaths {
		wg.Add(1)
		idx := i
		p := path

		go func() {
			defer wg.Done()

			var d DASMonitor
			d.DpuNumber = int(p.Dpunumber)

			// Get DAS Name
			dasInfo, err := db.GetDASDetails(int(p.Dasid))
			if err != nil {
				d.Status = "Error"
				dasMonitor[idx] = d
				return
			}
			d.DasName = dasInfo.Dasname
			ipAddress := dasInfo.Ipaddress
			dasType := dasInfo.Dastype

			conn, err := netutil.DialTCP(ipAddress, monitorPort, time.Second*5)
			if err != nil {
				d.Status = "NotConnected"
				d.Alarm = false
			} else {
				defer conn.Close()
				tmConn := netutil.NewTimeoutConn(conn, 5*time.Second)
				readPacket := getMonitorValue(tmConn, compCode, string(dasType))

				if readPacket == nil || len(readPacket) == 20 {
					d.Status = "NotConnected"
					d.Alarm = false
				} else {
					receivedAlarmValue, alarmErr := binary.BytesToUint32BE(readPacket[alarmOffset : alarmOffset+4])
					receivedLockValue, lockErr := binary.BytesToUint32BE(readPacket[lockOffset : lockOffset+4])

					if lockErr != nil || int(receivedLockValue) != lockValue {
						d.Status = "NotLocked"
					} else {
						d.Status = "Locked"
					}

					if alarmErr != nil || int(receivedAlarmValue) != alarmValue {
						d.Alarm = true
					} else {
						d.Alarm = false
					}
				}
			}
			dasMonitor[idx] = d
		}()
	}
	wg.Wait()
	return dasMonitor, nil
}
