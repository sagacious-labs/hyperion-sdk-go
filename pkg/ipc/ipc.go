package ipc

import (
	"fmt"
	"os"

	"github.com/sagacious-labs/hyperion-sdk-go/pkg/mail"
)

// Logf provides similar interface to printf - it takes in
// format and arguments parses it into hyperion mail object and
// prints it into stdout
//
// NOTE: This method is supposed to act like an "IPC" and is
// not meant for pretty printing. Hence this method isn't
// designed for development, it is advisable to use another
// logger in development and switch to this logger when running
// with hyperion
func Logf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)

	mail.New(mail.LOG, []byte(msg)).EncodeTo(os.Stdout)
}

// SendData takes in the data as byte then encodes
// it as a string and prints it to the stdout, which has
// the effect of sending the data to hyperion
func SendData(data []byte) {
	mail.New(mail.DATA, data).EncodeTo(os.Stdout)
}
