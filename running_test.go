package running_in_docker

import "testing"

func Test_getContainerID(t *testing.T) {

	realFile := "test_file/cgroup_real.txt"
	inDockerFile1 := "test_file/cgroup1.txt"
	inDockerFile2 := "test_file/cgroup2.txt"

	if getContainerID(realFile) != "" {
		t.Fatal("realFile error")
	}
	containId := getContainerID(inDockerFile1)
	println("containId1:", containId)
	if containId == "" {
		t.Fatal("inDockerFile1 error, containId:", containId)
	}
	containId = getContainerID(inDockerFile2)
	println("containId2:", containId)
	if containId == "" {
		t.Fatal("inDockerFile2 error, containId:", containId)
	}
}
