// Copyright 2020 yiya1989. All rights reserved.
// Apache License Version 2.0
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/yiya1989/sshkrb5/krb5forssh"
	
	"golang.org/x/crypto/ssh"
)


func testSShClient() error {
	sshUser := "root"
	sshHost := "x.xx.x.x"
	sshPort := "22"

	krb5Conf := "[libdefaults]\n\tdefault_realm = TESTREALM.COM\n\tdns_canonicalize_hostname = false\n\tdns_lookup_realm = false\n\tdns_lookup_kdc = false\n\n\tkdc_timesync = 1\n\tccache_type = 4\n\tforwardable = true\n\tproxiable = true\n\n\trdns = false\n\tignore_acceptor_hostname = true\n\n[realms]\n\tTESTREALM.COM = {\n\t\tkdc = krb5auth1.test.org\n\t\tkdc = krb5auth2.test.org\n\t\tkdc = krb5auth3.test.org\n\t\tkdc = krb5auth.test.org\n\t\tmaster_kdc = krb5auth.test.org\n\t\tadmin_server = krb5auth.test.org\n\t\tdefault_domain = test.org\n\t}\n\n[domain_realm]\n\t.test.org = TESTREALM.COM\n\ttest.org = TESTREALM.COM\n\n[login]\n\tkrb4_convert = true\n\tkrb4_get_tickets = false"
	user := "testuser"
	realm := "TESTREALM.COM"
	keytabPath := "/Users/admin/workspace/conf/testuser.keytab"
	keytabConf, err := ioutil.ReadFile(keytabPath)
	if err != nil {
		return err
	}

	sshGSSAPIClient, err := krb5forssh.NewKrb5InitiatorClient(krb5Conf, user, realm, keytabConf)
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
	err := testSShClient()
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
}
