package running_in_docker

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func IsRunningInDocker(pid ...int32) bool {

	nowPid := 1
	if len(pid) > 0 {
		nowPid = int(pid[0])
	}
	bok := pathExist(DockerEnvFPath)
	if bok == false {
		cGroupPath := fmt.Sprintf("/proc/%s/cgroup", strconv.Itoa(nowPid))
		containID := getContainerID(cGroupPath)
		if len(containID) > 0 {
			return true
		} else {
			return false
		}
	} else {
		return true
	}
}

func pathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func getContainerID(cGroupFPath string) string {

	containerID := ""
	content, err := ioutil.ReadFile(cGroupFPath)
	if err != nil {
		return containerID
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		field := strings.Split(line, ":")
		if len(field) < 3 {
			continue
		}
		CGroupPath := field[2]
		CGroupPath = strings.TrimSpace(CGroupPath)
		if len(CGroupPath) < 64 {
			continue
		}
		// Non-systemd Docker
		//5:net_prio,net_cls:/docker/de630f22746b9c06c412858f26ca286c6cdfed086d3b302998aa403d9dcedc42
		//3:net_cls:/kubepods/burstable/pod5f399c1a-f9fc-11e8-bf65-246e9659ebfc/9170559b8aadd07d99978d9460cf8d1c71552f3c64fefc7e9906ab3fb7e18f69
		pos := strings.LastIndex(CGroupPath, "/")
		if pos > 0 {

			idLen := len(CGroupPath) - pos - 1
			if idLen == 64 {
				//p.InDocker = true
				// docker id
				containerID = CGroupPath[pos+1 : pos+1+64]
				// logs.Debug("pid:%v in docker id:%v", pid, id)
				return containerID
			}
		}
		// systemd Docker
		//5:net_cls:/system.slice/docker-afd862d2ed48ef5dc0ce8f1863e4475894e331098c9a512789233ca9ca06fc62.scope
		dockerStr := "docker-"
		pos = strings.Index(CGroupPath, dockerStr)
		if pos > 0 {
			posScope := strings.Index(CGroupPath, ".scope")
			idLen := posScope - pos - len(dockerStr)
			if posScope > 0 && idLen == 64 {
				containerID = CGroupPath[pos+len(dockerStr) : pos+len(dockerStr)+64]
				return containerID
			}
		}
	}
	return containerID
}

const DockerEnvFPath string = "/.dockerenv"
