package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

// Set by go build flags
var version string

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

pre.error {
	background-color: red;
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
	versionOnly := flag.Bool("version", false, "Print version number and exit")
	flag.Parse()

	if *versionOnly {
		fmt.Printf(version)
		os.Exit(0)
	}

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
		{
			"iSCSI status",
			"/usr/sbin/tgtadm",
			[]string{"--mode", "target", "--op", "show"},
		},
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(header))
		for _, command := range commands {
			w.Write([]byte(fmt.Sprintf("<h1>%s</h1>\n", command.header)))
			cmd := exec.Command(command.cmd, command.args...)
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			cmd.Stdout = stdout
			cmd.Stderr = stderr
			if err := cmd.Run(); err != nil {
				w.Write([]byte("<pre class='error'>\n"))
				w.Write([]byte(fmt.Sprintf("err: %s\n", err)))
				w.Write([]byte("</pre>\n"))

				w.Write([]byte("<pre class='error'>\n"))
				w.Write([]byte(fmt.Sprintf("stderr: %s\n", stderr.String())))
				w.Write([]byte("</pre>\n"))

				w.Write([]byte("<pre>\n"))
				w.Write([]byte(fmt.Sprintf("stdout: %s\n", stdout.String())))
				w.Write([]byte("</pre>\n"))
			} else {
				w.Write([]byte("<pre>\n"))
				io.Copy(w, stdout)
				w.Write([]byte("</pre>\n"))
			}
		}
	})
	http.ListenAndServe(":9999", nil)
}
