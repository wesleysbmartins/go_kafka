package usecases

import "fmt"

func SyscallUsecase(msg string) {
	fmt.Printf("Reveived message SIGNALL: %v\n", msg)
}
