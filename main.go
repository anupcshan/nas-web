package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

var header = `<html>
<head>
	<title>NAS Status</title>
<style>
h1 {
	font-family: monospace;
}

pre {
	background-color: #ddd;
	border-radius: 1em;
	padding: 1em;
}
</style>
</head>

<body>
`

var footer = `</body>

</html>
`

type command struct {
	header string
	cmd    string
	args   []string
}

func main() {
	commands := []command{
		{
			"Uptime",
			"/usr/bin/uptime",
			nil,
		},
		{
			"Zpool status",
			"/usr/sbin/zpool",
			[]string{"status"},
		},
		{
			"ZFS status",
			"/usr/sbin/zfs",
			[]string{"list"},
		},
		{
			"Sanoid status",
			"/usr/sbin/sanoid",
			[]string{"--monitor-snapshots"},
		},
		{
			"Samba status",
			"/usr/bin/smbstatus",
			[]string{"-v"},
		},
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(header))
		for _, command := range commands {
			w.Write([]byte(fmt.Sprintf("<h1>%s</h1>\n", command.header)))
			w.Write([]byte("<pre>\n"))
			cmd := exec.Command(command.cmd, command.args...)
			cmd.Stdout = w
			if err := cmd.Run(); err != nil {
				log.Println(err)
				return
			}
			w.Write([]byte("</pre>\n"))
		}
	})
	http.ListenAndServe(":9999", nil)
}
