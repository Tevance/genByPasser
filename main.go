package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/process"
	"golang.org/x/sys/windows/registry"
)

func main() {
	fmt.Println("==============================")
	fmt.Println("genByPasser :: 광주광역시교육청 보호프로그램 무력화")
	fmt.Println("==============================")

	duration := 5 * time.Second
	time.Sleep(duration)
	processNames := []string{"iPCGuard.exe", "iAgent.exe", "iAgent32.exe", "iWatcher.exe", "iService.exe", "systemama.exe", "systemams.exe"}

	processList, err := process.Processes()
	if err != nil {
		fmt.Printf("[WARNING] 프로세스 목록을 가져올 수 없습니다.(재부팅 권장): %v\n", err)
		return
	}

	for _, processName := range processNames {
		for _, p := range processList {
			name, err := p.Name()
			if err != nil {
				continue
			}

			if strings.Contains(strings.ToLower(name), strings.ToLower(processName)) {
				err := p.Suspend()

				if err == nil {
					fmt.Printf("[SUCCESS] 프로세스 %s (PID %d)를 정지했습니다.\n", name, p.Pid)
				} else {
					fmt.Printf("[FAIL] 프로세스 %s (PID %d)를 정지하지 못했습니다.(관리자 권한으로 실행): %v\n", name, p.Pid, err)
				}
			}
		}
	}

	keyPath := `Software\Microsoft\Windows\CurrentVersion\Internet Settings`

	key, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		fmt.Printf("[FAIL] 레지스트리 키를 여는데 실패했습니다.(관리자 권한으로 실행): %v\n", err)
		return
	}
	defer key.Close()

	err = key.SetDWordValue("ProxyEnable", 0)
	if err != nil {
		fmt.Printf("[FAIL] 프록시 서버 설정 변경에 실패했습니다.(관리자 권한으로 실행): %v\n", err)
		return
	}

	fmt.Println("[SUCCESS] 프록시 서버 설정이 비활성화 되었습니다.");
	fmt.Println("[INFO] 우회 작업이 완료되었습니다.");
	fmt.Println("[INFO] 이 프로그램 효력은 재부팅 전까지만 유효하며 재부팅 후에는 다시 실행해야 합니다.");
	fmt.Println("[INFO] 이 프로그램으로 12시 제한, 게임 사이트 제한, 계정 로그인 제한 등이 해제 되었습니다.");
	fmt.Scanln();
}
