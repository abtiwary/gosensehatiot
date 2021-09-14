package SenseHatIoT

/*
#include <stdint.h>
#include <unistd.h>
#include <linux/i2c-dev.h>
#include <stdlib.h>
#include <fcntl.h>

#define DEV_ID 0x5c
#define DEV_PATH "/dev/i2c-1"
#define WHO_AM_I 0x0F
#define TEMP_OUT_L 0x2B
#define TEMP_OUT_H 0x2C
#define CTRL_REG1 0x20
#define CTRL_REG2 0x21

void delay(int t) {
    usleep(t * 1000);
}

double GetTemperature() {
    int fd = 0;
    uint8_t temp_out_l = 0;
    uint8_t temp_out_h = 0;
    int16_t temp_out = 0;
    double t_c = 0.0;
    uint8_t status = 0;

    if((fd = open(DEV_PATH, O_RDWR)) < 0) {
        return t_c;
    }

    if (ioctl(fd, I2C_SLAVE, DEV_ID) < 0) {
        close(fd);
        return t_c;
    }

    if (i2c_smbus_read_byte_data(fd, WHO_AM_I) != 0xBD) {
         close(fd);
         return t_c;
    }

    i2c_smbus_write_byte_data(fd, CTRL_REG2, 0x01);

    do {
        delay(25);
        status = i2c_smbus_read_byte_data(fd, CTRL_REG2);
    }
    while (status != 0);

    temp_out_l = i2c_smbus_read_byte_data(fd, TEMP_OUT_L);
    temp_out_h = i2c_smbus_read_byte_data(fd, TEMP_OUT_H);

    temp_out = temp_out_h << 8 | temp_out_l;
    t_c = 42.5 + (temp_out / 480.0);

    close(fd);
    return t_c;
}
*/
import "C"

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GetCpuTemperature() float64 {
	out, err := exec.Command("/usr/bin/vcgencmd", "measure_temp").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	outstr := string(out)
	outstr = strings.Replace(outstr, "temp=", "", 1)
	outstr = strings.Replace(outstr, "'C\n", "", 1)

	f64out, _ := strconv.ParseFloat(outstr, 64)

	return f64out
}

func GetSenseHatTemperature() float64 {
	cputemp := GetCpuTemperature()
	temp := float64(C.GetTemperature())

	finaltemp := temp - (cputemp - temp)

	return finaltemp
}
