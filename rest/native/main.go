package main

// #include <stdbool.h>
import "C"
import (
	"io/ioutil"
	"net/http"
	"strings"
	"unsafe"

	"github.com/spec-tacles/go/rest"
)

var clients = make(map[*rest.Client]struct{})

func clientFromPtr(clientPtr uintptr) *rest.Client {
	return (*rest.Client)(unsafe.Pointer(clientPtr))
}

//export rest_create_client
func rest_create_client(token *C.char) uintptr {
	client := rest.NewClient(C.GoString(token))
	clients[client] = struct{}{}
	return uintptr(unsafe.Pointer(client))
}

//export rest_destroy_client
func rest_destroy_client(clientPtr uintptr) {
	client := clientFromPtr(clientPtr)
	delete(clients, client)
}

//export rest_globally_limited
func rest_globally_limited(clientPtr uintptr) C.bool {
	client := clientFromPtr(clientPtr)
	return C.bool(client.GloballyLimited())
}

//export rest_do
func rest_do(clientPtr uintptr, method *C.char, url *C.char, reqChars *C.char, resChars *C.char) C.bool {
	client := clientFromPtr(clientPtr)

	body := strings.NewReader(C.GoString(reqChars))
	req, err := http.NewRequest(C.GoString(method), C.GoString(url), body)
	if err != nil {
		resChars = C.CString(err.Error())
		return C.bool(false)
	}

	res, err := client.Do(req)
	if err != nil {
		resChars = C.CString(err.Error())
		return C.bool(false)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		resChars = C.CString(err.Error())
		return C.bool(false)
	}
	res.Body.Close()

	resChars = C.CString(string(resBody))
	return C.bool(true)
}

func main() {}
