// Copyright 2020 yiya1989. All rights reserved.
// Apache License Version 2.0
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/yiya1989/sshkrb5/krb5forssh"

	"golang.org/x/crypto/ssh"
)

func mustReadFile(filename string) []byte {
	fd, err := os.Open(filename)
	if err != nil {
		log.Fatalf("fail to open %s: %s", filename, err)
	}
	defer fd.Close()

	b, err := io.ReadAll(fd)
	if err != nil {
		log.Fatalf("fail to read from %s: %s", filename, err)
	}

	return b
}

func mustGetKrb5CCacheFilename() string {
	v, err := krb5forssh.GetKrb5CCacheFilename()
	if err != nil {
		log.Fatalf("fail to GetKrb5CCacheFilename(): %v", err)
	}
	return v
}

func testSShClientWithCCache(sshUser, sshHost, sshPort string) error {
	krb5Conf := string(mustReadFile("/etc/krb5.conf"))

	sshGSSAPIClient, err := krb5forssh.NewKrb5InitiatorClientWithCCache(krb5Conf, mustGetKrb5CCacheFilename())
	if err != nil {
		return err
	}
	sshCfg := &ssh.ClientConfig{
		User:            sshUser,
		Auth:            []ssh.AuthMethod{ssh.GSSAPIWithMICAuthMethod(&sshGSSAPIClient, sshHost)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", sshHost+":"+sshPort, sshCfg)
	if err != nil {
		log.Fatalln(err)
	}
	defer sshClient.Close()
	fmt.Printf("Dial sucess\n")

	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatalln(err)
		return err
	}
	defer session.Close()
	fmt.Printf("NewSession success\n")

	r, err := session.Output("pwd")
	if err != nil {
		log.Fatalln(err)
		return err
	}
	fmt.Printf("cmd: %s, result: %s", "pwd", r)

	return nil
}

func testSShClient(sshUser, sshHost, sshPort string) error {
	krb5Conf := "[libdefaults]\n\tdefault_realm = TESTREALM.COM\n\tdns_canonicalize_hostname = false\n\tdns_lookup_realm = false\n\tdns_lookup_kdc = false\n\n\tkdc_timesync = 1\n\tccache_type = 4\n\tforwardable = true\n\tproxiable = true\n\n\trdns = false\n\tignore_acceptor_hostname = true\n\n[realms]\n\tTESTREALM.COM = {\n\t\tkdc = krb5auth1.test.org\n\t\tkdc = krb5auth2.test.org\n\t\tkdc = krb5auth3.test.org\n\t\tkdc = krb5auth.test.org\n\t\tmaster_kdc = krb5auth.test.org\n\t\tadmin_server = krb5auth.test.org\n\t\tdefault_domain = test.org\n\t}\n\n[domain_realm]\n\t.test.org = TESTREALM.COM\n\ttest.org = TESTREALM.COM\n\n[login]\n\tkrb4_convert = true\n\tkrb4_get_tickets = false"
	user := "testuser"
	realm := "TESTREALM.COM"
	keytabPath := "/Users/admin/workspace/conf/testuser.keytab"
	keytabConf, err := ioutil.ReadFile(keytabPath)
	if err != nil {
		return err
	}

	sshGSSAPIClient, err := krb5forssh.NewKrb5InitiatorClientWithKeytab(krb5Conf, user, realm, keytabConf)
	if err != nil {
		return err
	}
	sshCfg := &ssh.ClientConfig{
		User:            sshUser,
		Auth:            []ssh.AuthMethod{ssh.GSSAPIWithMICAuthMethod(&sshGSSAPIClient, sshHost)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", sshHost+":"+sshPort, sshCfg)
	if err != nil {
		log.Fatalln(err)
	}
	defer sshClient.Close()
	fmt.Printf("Dial sucess\n")

	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatalln(err)
		return err
	}
	defer session.Close()
	fmt.Printf("NewSession success\n")

	//r, err := session.Output("pwd")
	//if err != nil {
	//	log.Fatalln(err)
	//	return err
	//}
	//fmt.Printf("cmd: %s, result: %s", "pwd", r)

	if err = session.RequestPty("xterm-256color", 25, 80, ssh.TerminalModes{}); err != nil {
		log.Fatalln(err)
		return err
	}

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	if err = session.Shell(); err != nil {
		log.Fatalln(err)
		return err
	}
	if err = session.Wait(); err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

func main() {
	var clientType, sshUser, sshHost, sshPort string
	flag.StringVar(&clientType, "client-type", "keytab", "client type to test (keytab or ccache)")
	flag.StringVar(&sshUser, "ssh-user", "root", "ssh remote user")
	flag.StringVar(&sshHost, "ssh-host", "localhost", "ssh remote host")
	flag.StringVar(&sshPort, "ssh-port", "22", "ssh remote port")
	flag.Parse()

	var err error
	switch clientType {
	case "keytab":
		err = testSShClient(sshUser, sshHost, sshPort)

	case "ccache":
		err = testSShClientWithCCache(sshUser, sshHost, sshPort)

	default:
		log.Fatalf("unexpected client-type %s", clientType)
	}

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
}
