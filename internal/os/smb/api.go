package smb

import (
	"fmt"
	"os"
	"os/exec"

	"k8s.io/klog"
)

type APIImplementor struct{}

func New() APIImplementor {
	return APIImplementor{}
}

func (APIImplementor) NewSmbGlobalMapping(remotePath, localPath, username, password string) error {
	klog.V(4).Infof("NewSmbGlobalMapping: remotePath:%q, localPath:%q", remotePath, localPath)

	env := append(os.Environ(),
		fmt.Sprintf("smbuser=%s", username),
		fmt.Sprintf("smbpassword=%s", password),
		fmt.Sprintf("smbremotepath=%s", remotePath))
	// use PowerShell Environment Variables to store user input string to prevent command line injection
	// https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_environment_variables?view=powershell-5.1
	cmdLine := fmt.Sprintf(`$PWord = ConvertTo-SecureString -String $Env:smbpassword -AsPlainText -Force` +
		`;$Credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $Env:smbuser, $PWord;`)

	if len(localPath) == 0 {
		cmdLine = fmt.Sprintf("%s%s", cmdLine, `New-SmbGlobalMapping -RemotePath $Env:smbremotepath -Credential $Credential`)
	} else {
		env = append(env, fmt.Sprintf("smblocalpath=%s", localPath))
		cmdLine = fmt.Sprintf("%s%s", cmdLine, `New-SmbGlobalMapping -RemotePath $Env:smbremotepath -LocalPath $Env:smblocalpath -Credential $Credential`)
	}

	cmd := exec.Command("powershell", "/c", cmdLine)
	cmd.Env = env

	// TODO: Remove this debug log.
	klog.Infof("cmdline: %s. env: %v", cmdLine, cmd.Env)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("NewSmbGlobalMapping failed: %v, output: %q", err, string(output))
	}
	return nil
}

func (APIImplementor) RemoveSmbGlobalMapping(remotePath string) error {
	klog.V(4).Infof("RemoveSmbGlobalMapping remotePath (%q)", remotePath)
	cmd := exec.Command("powershell", "/c", `Remove-SmbGlobalMapping -RemotePath $Env:smbremotepath -Force`)
	cmd.Env = append(os.Environ(), fmt.Sprintf("smbremotepath=%s", remotePath))
	if output, err := cmd.CombinedOutput(); err != nil {
		klog.Errorf("Remove-SmbGlobalMapping failed: %v, output: %q", err, output)
		return fmt.Errorf("UnmountSmbShare failed: %v, output: %q", err, output)
	}
	return nil
}
