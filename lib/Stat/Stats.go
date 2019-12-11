package Stats

import (
	"fmt"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
)

// Method to get current disk usgae
func GetDiskUsage() (diskUsed float64) {
	parts, _ := disk.Partitions(false)

	var usage []*disk.UsageStat
	var diskU float64

	for _, part := range parts {
		u, _ := disk.Usage(part.Mountpoint)
		usage = append(usage, u)
		fmt.Printf(fmt.Sprintf("%f", u.UsedPercent))
		fmt.Printf("\n")
		fmt.Printf(u.Path)
		fmt.Printf("\n")

		diskU = u.UsedPercent

		break
	}
	return diskU

}

func GetHostName() string {

	devInfo, _ := host.Info()

	fmt.Printf(devInfo.Hostname)
	return devInfo.Hostname
}
