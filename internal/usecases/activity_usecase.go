package usecases

import (
	"fmt"
	"go_kafka/internal/entities"
)

func ActivityUsecase(activity entities.Activity) {
	fmt.Println("Activity received: ", activity)
}
