// Copyright Jetstack Ltd. See LICENSE for details.
package ssh

import (
	//"bufio"
	"bytes"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	//	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/jetstack/tarmak/pkg/tarmak/interfaces"
	"github.com/jetstack/tarmak/pkg/tarmak/utils"
)

var _ interfaces.SSH = &SSH{}

type SSH struct {
	tarmak interfaces.Tarmak
	log    *logrus.Entry

	controlPaths []string
}

func New(tarmak interfaces.Tarmak) *SSH {
	s := &SSH{
		tarmak: tarmak,
		log:    tarmak.Log(),
	}

	return s
}

func (s *SSH) WriteConfig(c interfaces.Cluster) error {

	hosts, err := c.ListHosts()
	if err != nil {
		return err
	}

	var sshConfig bytes.Buffer
	sshConfig.WriteString(fmt.Sprintf("# ssh config for tarmak cluster %s\n", c.ClusterName()))

	for _, host := range hosts {
		_, err = sshConfig.WriteString(host.SSHConfig())
		if err != nil {
			return err
		}

		s.controlPaths = append(s.controlPaths, host.SSHControlPath())
	}

	err = utils.EnsureDirectory(filepath.Dir(c.SSHConfigPath()), 0700)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.SSHConfigPath(), sshConfig.Bytes(), 0600)
	if err != nil {
		return err
	}

	return nil
}

func (s *SSH) args() []string {
	return []string{
		"ssh",
		"-F",
		s.tarmak.Cluster().SSHConfigPath(),
	}
}

// Pass through a local CLI session
func (s *SSH) PassThrough(argsAdditional []string) error {
	hosts, err := s.tarmak.Cluster().ListHosts()
	if err != nil {
		return err
	}

	var host interfaces.Host
	var bastion interfaces.Host
	for _, h := range hosts {
		if h.Aliases()[0] == argsAdditional[0] {
			host = h
			continue
		}
		if h.Aliases()[0] == "bastion" {
			bastion = h
		}
	}

	b, err := ioutil.ReadFile(s.tarmak.Environment().SSHPrivateKeyPath())
	if err != nil {
		return err
	}
	signer, err := ssh.ParsePrivateKey(b)
	if err != nil {
		return err
	}

	conf := &ssh.ClientConfig{
		Timeout:         time.Minute * 10,
		User:            host.User(),
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	confProxy := &ssh.ClientConfig{
		Timeout:         time.Minute * 10,
		User:            host.User(),
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	proxyClient, err := ssh.Dial("tcp", net.JoinHostPort(bastion.Hostname(), "22"), confProxy)
	if err != nil {
		return fmt.Errorf("failed to set up connection to bastion: %s", err)
	}

	conn, err := proxyClient.Dial("tcp", net.JoinHostPort(host.Hostname(), "22"))
	if err != nil {
		return fmt.Errorf("failed to set up connection to %s from basiton: %s", host.Hostname(), err)
	}

	ncc, chans, reqs, err := ssh.NewClientConn(conn, net.JoinHostPort(host.Hostname(), "22"), conf)
	if err != nil {
		return fmt.Errorf("failed to set up ssh client: %s", err)
	}

	client := ssh.NewClient(ncc, chans, reqs)
	sess, err := client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	sess.Stderr = os.Stderr
	sess.Stdout = os.Stdout
	sess.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	fileDescriptor := int(os.Stdin.Fd())
	if terminal.IsTerminal(fileDescriptor) {
		originalState, err := terminal.MakeRaw(fileDescriptor)
		if err != nil {
			return err
		}
		defer terminal.Restore(fileDescriptor, originalState)

		termWidth, termHeight, err := terminal.GetSize(fileDescriptor)
		if err != nil {
			return err
		}

		err = sess.RequestPty("xterm-256color", termHeight, termWidth, modes)
		if err != nil {
			return err
		}
	}

	if err := sess.Shell(); err != nil {
		return err
	}

	if err := sess.Wait(); err != nil {
		return err
	}

	return nil
}

func (s *SSH) Execute(host string, command string, argsAdditional []string) (returnCode int, err error) {
	args := append(s.args(), host, "--", command)
	args = append(args, argsAdditional...)

	cmd := exec.Command(args[0], args[1:len(args)]...)

	err = cmd.Start()
	if err != nil {
		return -1, err
	}

	err = cmd.Wait()
	if err != nil {
		perr, ok := err.(*exec.ExitError)
		if ok {
			if status, ok := perr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus(), nil
			}
		}
		return -1, err
	}

	return 0, nil

}

func (s *SSH) Validate() error {
	// no environment in tarmak so we have no SSH to validate
	if s.tarmak.Environment() == nil {
		return nil
	}

	keyPath := s.tarmak.Environment().SSHPrivateKeyPath()
	f, err := os.Stat(keyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("failed to read ssh file status: %v", err)
	}

	if f.IsDir() {
		return fmt.Errorf("expected ssh file location '%s' is directory", keyPath)
	}

	if f.Mode() != os.FileMode(0600) && f.Mode() != os.FileMode(0400) {
		s.log.Warnf("ssh file '%s' holds incorrect permissions (%v), setting to 0600", keyPath, f.Mode())
		if err := os.Chmod(keyPath, os.FileMode(0600)); err != nil {
			return fmt.Errorf("failed to set ssh private key file permissions: %v", err)
		}
	}

	bytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("unable to read ssh private key: %s", err)
	}

	block, _ := pem.Decode(bytes)
	if block == nil {
		return errors.New("failed to parse PEM block containing the ssh private key")
	}

	return nil
}

func (s *SSH) Cleanup() error {
	var result *multierror.Error

	for _, c := range utils.RemoveDuplicateStrings(s.controlPaths) {
		if err := os.RemoveAll(c); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}
