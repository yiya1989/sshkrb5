# sshkrb5

# Usage
```
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
```

# Related issue/repo
https://github.com/golang/go/issues/25899
https://github.com/golang/crypto
