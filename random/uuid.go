package random


import (
    "os/exec"
    "net/http"
)



func rootHandler( w http  .ResponseWriter, r *http.Request) {
     uuid, err := GenerateUUID()
    if err != nil {
        http.Error(w, "error generating UUID", 500)
        return
    }
    w.Write([]byte(uuid))
}



func GenerateUUID() (string, error) {
    out, err := exec.Command("uuidgen").Output()
    if err != nil  {
    return "", err
    }
    return string(out), nil
}
