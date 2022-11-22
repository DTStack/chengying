#管控server ip
$SERVER_IP_ADDRESS='{{.SERVER_IP_ADDRESS}}'
#管控server port
$SERVER_PORT=9890
#安装包下载host#ip:port
$EASYAGENT_DOWNLOAD_HOST='{{.callback}}'

$UUID='{{.UUID}}'
$CALLBACK_TO_PROD='{{.CALLBACK_TO_PROD}}'
$SIDECAR_INSTALL_PATH='{{.SIDECAR_INSTALL_PATH}}'
$DEBUG_MODE='{{.DEBUG_MODE}}'
$TARGET_OS='{{.TARGET_OS}}'

$SIDECAR_DIR="c:\dtstack\easyagent"

If(![String]::IsNullOrEmpty($SIDECAR_INSTALL_PATH))
{
	$SIDECAR_DIR=$SIDECAR_INSTALL_PATH
}

$SIDECAR_DOWNLOAD_URL="http://$EASYAGENT_DOWNLOAD_HOST/easyagent/easyagent_for_win_x86/easyagent.zip"
$INSTALL_ZIP="easyagent.zip"
$SIDECAR_APP_DIR="easyagent"

# Is this a 64 bit process
function Test-Win64() {
    return [IntPtr]::size -eq 8
}

# Is this a 32 bit process
function Test-Win32() {
    return [IntPtr]::size -eq 4
}

if(Test-Win64)
{
    $SIDECAR_DOWNLOAD_URL="http://$EASYAGENT_DOWNLOAD_HOST/easyagent/easyagent_for_win_x64/easyagent.zip"
}

$INSTALL_TYPE='sidecar'
$CHECK_FLG='{{.create_time}}'
$CALLBACK_URL="http://$EASYAGENT_DOWNLOAD_HOST/api/v1/deploy/sidecar/install/callback?check_flg=$CHECK_FLG&install_type=$INSTALL_TYPE&client_id=$UUID"

# success or failed
$INSTALL_RES=''

# install comount msg
$MSG=''


Function DownloadFile([string]$url, [string]$dest)
{
    Write-Host "-->Download File: $url, save to $dest"
	$client = new-object System.Net.WebClient
	$client.DownloadFile($url, $dest)
}

Function Clean_Install_Pkg
{
	If(Test-Path "$SIDECAR_DIR\$INSTALL_ZIP")
	{
		Write-Host "-->clean $SIDECAR_DIR\$INSTALL_ZIP "
		del "$SIDECAR_DIR\$INSTALL_ZIP"
	}
}

# Convert string to base64 string
function ConvertTo-Base64String([string]$str)
{
    $byteArray = [System.Text.UnicodeEncoding]::Unicode.GetBytes($str)
    [Convert]::ToBase64String( $byteArray )
}

# Convert base64 string to string
function ConvertFrom-Base64String([string]$string)
{
    $byteArray = [Convert]::FromBase64String($string)
    [System.Text.UnicodeEncoding]::Unicode.GetString($byteArray)
}


Function Unzip-File
{
    param([string]$ZipFile,[string]$TargetFolder)
    if(!(Test-Path $TargetFolder))
    {
		mkdir $TargetFolder
    }
    $shellApp = New-Object -ComObject Shell.Application
    $files = $shellApp.NameSpace($ZipFile).Items()
	$files|%{Remove-Item ("$TargetFolder/{0}" -f  $_.name ) -Force -Recurse -ErrorAction SilentlyContinue}
    $shellApp.NameSpace($TargetFolder).CopyHere($files)
}

Function Install_Sidecar
{
	Write-Host '-->Install easyagent sidecar...'
	if(!(Test-Path $SIDECAR_DIR))
	{
		md $SIDECAR_DIR
	}
    cd "$SIDECAR_DIR"

	If(!(Test-Path "$SIDECAR_DIR\$INSTALL_ZIP"))
	{
		Write-Host "[+]download $INSTALL_ZIP from $SIDECAR_DOWNLOAD_URL"
		DownloadFile -url "$SIDECAR_DOWNLOAD_URL" -dest "$SIDECAR_DIR\$INSTALL_ZIP"
	}
	If(!(Test-Path "$SIDECAR_DIR\$SIDECAR_APP_DIR\easyagent-sidecar.exe"))
	{
	    Write-Host "-->unzip file $SIDECAR_DIR\$INSTALL_ZIP"
	    Unzip-File -ZipFile "$SIDECAR_DIR\$INSTALL_ZIP" -TargetFolder "$SIDECAR_DIR"
	}
	If(!(Test-Path "$SIDECAR_DIR\$SIDECAR_APP_DIR\easyagent-sidecar.exe"))
	{
		Write-Host "[+]unzip file $SIDECAR_DIR\$INSTALL_ZIP failed"
		Clean_Install_Pkg
		$MSG = ConvertTo-Base64String -str "sidecar agent unzip file $SIDECAR_DIR\$INSTALL_ZIP failed"
        $INSTALL_RES='failed'
		DownloadFile -url "$CALLBACK_URL&msg=$MSG&install_res=$INSTALL_RES" -dest "$SIDECAR_DIR\callback.info"
		exit 1
	}
}

Function Config_Sidecar
{
    Write-Host "-->config EasyManage agent (easyagent sidecar)..."

	if((Test-Path $SIDECAR_DIR) -and (Test-Path $SIDECAR_DIR\$SIDECAR_APP_DIR\easyagent.yml))
	{
		$config = "$SIDECAR_DIR\$SIDECAR_APP_DIR\easyagent.yml"
		$is_in = select-string $config -pattern 'Callback='
		if(!$is_in)
		{
			(type $config) -replace ('\A\s{1,10}uuid:.*$',"  uuid: $UUID")|out-file $config
		}
		Else
		{
			Write-Host '[+] easyagent-sidecar is already installed...'
		}
		if($CALLBACK_TO_PROD)
		{
			Write-Host '[+]Config CALLBACK_TO_PROD'
			(type $config) -replace ('\A.{1,100}CallBack=.*$',"")|out-file $config
			"  - $CALLBACK_TO_PROD"|Out-File -Append $config
		}
		if($SIDECAR_DIR)
		{
			Write-Host '[+]Config SIDECAR_INSTALL_PATH'
			(type $config) -replace ('\A.{1,10}dir:.*$',"  dir: $SIDECAR_DIR\$SIDECAR_APP_DIR")|out-file $config
		}
		if($SERVER_IP_ADDRESS)
		{
			Write-Host '[+]Config SERVER_IP_ADDRESS'
			(type $config) -replace ('\A.{1,10}server:.*$',"  server: $SERVER_IP_ADDRESS")|out-file $config
		}
		if($SERVER_PORT)
		{
			Write-Host '[+]Config SERVER_PORT'
			(type $config) -replace ('\A.{1,10}port:.*$',"  port: $SERVER_PORT")|out-file $config
		}
	}
}

#----------------------------begin start servcie------------------------
function StartTheService()
{
        param([string]$service_name)
		$Service = Get-Service $service_name
		$status = $Service.Status
		Write-Host "service $service_name status $status"
        if($Service.Status -eq "Running")
        {
                Write-Host $Service "The service the running now,will restart the service"
                Write-Host $Service "Starting service..."
                Restart-Service -Name $service_name
        }
        elseif($Service.Status -eq "Disabled")
        {
                #Set the status of service, cause the cmd Start-Service can not start the service if the status is disabled
                Set-Service -Name $service_name -StartupType automatic Manual
                Write-Host $status "Set the status to manual successful"
                Write-Host $Service "Starting service..."
                Start-Service -Name $service_name
        }
        else
        {
                Write-Host $Service.Name "The service the stopped now,will start the service"
                Write-Host $Service.Name "Starting service..."
                Start-Service -Name $service_name
                $Service = Get-Service $service_name
                $step = 1
                do
                {
                    $step = $step + 1
                    Write-Host $service_name "Starting service..."
                    Start-Service -Name $service_name
                    $Service = Get-Service $service_name
                }while($Service.Status -ne "Running" -and $step -le 3)
        }
}
#----------------------------end start servcie------------------------
#---------------------------------------------------------------------
#----------------------------begin stop servcie-----------------------
function StopTheService()
{
		param([string]$service_name)
		$Service = Get-Service $service_name
        if($Service.Status -eq "Running")
        {
                do
                {
                        Write-Host $Service "Stopping service..."
                        Stop-Service -Name $Service.Name
                }while($_.Status -eq "Running") #need to wait until it stopped
                if($_.Status -eq "Running")
                {
                        # stop the servcie failed,send the error mail
                        #sendMail here

                }
                else
                {
                        Write-Host $Service "Stop the servcie successful"
                }
        }
        else
        {
                Write-Host $Service.Status "The service is already stopped now"
        }
}

Function Install_Service
{
	$service_name = "dtstackeasyagent"
	$workdir = "$SIDECAR_DIR\$SIDECAR_APP_DIR"
	# create new service
	if(Check_Service)
	{
	    Write-Host $service_name "service already exist!"
	}
	else
	{
	    New-Service -name $service_name `
        		-displayName $service_name `
        		-Description "dtstack easyagent service" `
        		-StartupType Automatic `
        		-binaryPathName "`"$workdir\easyagent-sidecar.exe`" -c `"$workdir\easyagent.yml`" --agents `"$workdir\agents-file.yml`""
	}

	$Service = Get-Service $service_name
	if([String]::IsNullOrEmpty($Service))
	{
        Write-Host "New service for easyagent failed, service_name: $service_name"
		$MSG = ConvertTo-Base64String -str "New service for easyagent failed, service_name: $service_name"
        $INSTALL_RES='failed'
		DownloadFile -url "$CALLBACK_URL&msg=$MSG&install_res=$INSTALL_RES" -dest "$SIDECAR_DIR\callback.info"
		exit 1
	}
	Write-Host "The dtstackeasyagent starting!"
	Start-Sleep -Seconds 10
	StartTheService($service_name)
	$service_check = Get-Service $service_name
	if($service_check.Status -eq "Running")
    {
        Write-Host $service_check "The easyagent install and start success!"
		$MSG = ConvertTo-Base64String -str "The easyagent install and start success!"
        $INSTALL_RES='success'
		DownloadFile -url "$CALLBACK_URL&msg=$MSG&install_res=$INSTALL_RES" -dest "$SIDECAR_DIR\callback.info"
    }
	Else
	{
		Write-Host $service_check "The easyagent servcie start error!"
		$MSG = ConvertTo-Base64String -str "The easyagent servcie start error!"
        $INSTALL_RES='failed'
		DownloadFile -url "$CALLBACK_URL&msg=$MSG&install_res=$INSTALL_RES" -dest "$SIDECAR_DIR\callback.info"
	}
}

Function CheckAdmin
{
	$currentWi = [Security.Principal.WindowsIdentity]::GetCurrent()
    $currentWp = [Security.Principal.WindowsPrincipal]$currentWi
	if( -not $currentWp.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator))
	{
		return $false
	}
	return $true
}

Function Check_Service
{
    $service_name = "dtstackeasyagent"
    # delete service if it already exists = Get-Service $service_name
    $Service = Get-Service $service_name -ErrorAction SilentlyContinue
    if(![String]::IsNullOrEmpty($Service))
    {
        return $true
    }
    return $false
}

#if(!CheckAdmin)
#{
#    Write-Host "Need administrator privileges to run this install script!"
#    exit 1;
#}

Check_Service
Install_Sidecar
Config_Sidecar
Install_Service