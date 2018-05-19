; PLM (Process Load Monitor) installer creator NSIS script
;
; Prerequisities:
;  - GOPATH environment variable needs to be set correctly
;  - plm needs to be installed (go install github.com/midstar/plm)
;  - plmc needs to be installed (go install github.com/midstar/plm/plmc)
;
; Usage:
;  - makensis -DVERSION=<version> plm_windows_installer.nsi
;     (<version> should be in the format 1.1.1.1)
;    
;
;-------------------------------------------------

;--------------------------------
;External dependencies / libraries

; Use the NSIS Modern UI 2
!include MUI2.nsh
!include x64.nsh

;--------------------------------
;Common definitions

!define APPLICATION_NAME "Process Load Monitor"
!define APPLICATION_FOLDER "plm"
!define APPLICATION_SOURCE "$%GOPATH%\src\github.com\midstar\plm"
!define APPLICATION_BINARY "$%GOPATH%\bin"

; The application version. Override with /DVERSION flag
!ifndef VERSION
!define VERSION "0.0.0.0"
!endif

; The name of the installer
Name "${APPLICATION_NAME} ${VERSION}"

; The file to write
OutFile "${APPLICATION_BINARY}\${APPLICATION_FOLDER}Setup-${VERSION}.exe"

; The default installation directory
InstallDir $PROGRAMFILES64\${APPLICATION_FOLDER}

; Registry key to check for directory (so if you install again, it will 
; overwrite the old one automatically)
InstallDirRegKey HKLM "Software\${APPLICATION_FOLDER}" "Install_Dir"

; Request application privileges
RequestExecutionLevel admin

;--------------------------------
;Interface Settings

!define MUI_ABORTWARNING
!define MUI_ICON "images\logo.ico"

;--------------------------------
;Pages

!insertmacro MUI_PAGE_LICENSE "LICENSE.txt"
!insertmacro MUI_PAGE_COMPONENTS
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
;!define MUI_FINISHPAGE_RUN "TBD <open link here>"
;!define MUI_FINISHPAGE_RUN_TEXT "Launch PLM User Interface"
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

;--------------------------------
;Languages
 
!insertmacro MUI_LANGUAGE "English"

;--------------------------------
;Version Information

VIProductVersion "${VERSION}"
VIAddVersionKey /LANG=${LANG_ENGLISH} "ProductName" "${APPLICATION_NAME}"
VIAddVersionKey /LANG=${LANG_ENGLISH} "Comments" "Monitor processes"
VIAddVersionKey /LANG=${LANG_ENGLISH} "CompanyName" "Joel Midstjarna"
VIAddVersionKey /LANG=${LANG_ENGLISH} "LegalTrademarks" "-"
VIAddVersionKey /LANG=${LANG_ENGLISH} "LegalCopyright" "Copyright Joel Midstjarna"
VIAddVersionKey /LANG=${LANG_ENGLISH} "FileDescription" "${APPLICATION_NAME} Setup"
VIAddVersionKey /LANG=${LANG_ENGLISH} "FileVersion" "${VERSION}"
VIAddVersionKey /LANG=${LANG_ENGLISH} "ProductVersion" "${VERSION}"

;-----------------------------------------------------------------------------
; Init function - executed before the installation starts
Function .onInit
  
  ;---------------------------------------------------------------------------
  ; Check if already installed 
 
  ReadRegStr $R0 HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPLICATION_FOLDER}"  "UninstallString"
  StrCmp $R0 "" noPreviousInstaller
 
  MessageBox MB_OKCANCEL|MB_ICONEXCLAMATION "${APPLICATION_NAME} is already installed. $\n$\n\
    Click `OK` to remove the previous version or `Cancel` to cancel this upgrade."  IDOK uninst
  Abort
 
  ;Run the uninstaller
  uninst:
     ClearErrors
      Exec $R0
   
  noPreviousInstaller:

FunctionEnd


;======================================================================================================================
; Application install section
Section "${APPLICATION_NAME}" SectionMain

  SectionIn RO
  
  ; Set output path to the installation directory.
  SetOutPath $INSTDIR
  
  ; Copy plm binary
  File "${APPLICATION_BINARY}\plm.exe"
	
	; Copy plm default configuration
  File "${APPLICATION_SOURCE}\plm.config"
	
	; Copy plm URL
  File "${APPLICATION_SOURCE}\Process Load Monitor.url"
	
	; Copy the templates
	SetOutPath "$INSTDIR\templates"
	File "${APPLICATION_SOURCE}\templates\*.*"
	
	; Copy plmc
	SetOutPath "$INSTDIR\client"
	File "${APPLICATION_BINARY}\plmc.exe"
	
	
  
  ; Write the installation path into the registry
  WriteRegStr HKLM SOFTWARE\${APPLICATION_FOLDER} "Install_Dir" "$INSTDIR"
  
  ; Write the version into the registry
  WriteRegStr HKLM SOFTWARE\${APPLICATION_FOLDER} "Version" "${VERSION}"
  
  ; Write the uninstall keys for Windows
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPLICATION_FOLDER}" "DisplayName" "${APPLICATION_NAME} ${VERSION}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPLICATION_FOLDER}" "Publisher" "Joel Midstjarna"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPLICATION_FOLDER}" "UninstallString" '"$INSTDIR\uninstall.exe"'
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPLICATION_FOLDER}" "NoModify" 1
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPLICATION_FOLDER}" "NoRepair" 1
  WriteUninstaller "uninstall.exe"
  
  ; Install and start Vanadis View service
  ; TBD ExecWait "install_service.bat"
  
SectionEnd


;======================================================================================================================
; Start menu shortcuts install section (can be disabled by the user)
Section "Start Menu Shortcuts" SectionStartMenu

  CreateDirectory "$SMPROGRAMS\${APPLICATION_FOLDER}"
  CreateShortcut "$SMPROGRAMS\${APPLICATION_FOLDER}\Uninstall.lnk" "$INSTDIR\uninstall.exe" "" "$INSTDIR\uninstall.exe" 0
  ; TBD CreateShortcut "$SMPROGRAMS\${APPLICATION_FOLDER}\VanadisViewManager.lnk" "$INSTDIR\view-gui\Vanadis.View.Manager.exe" "" "$INSTDIR\view-gui\Vanadis.View.Manager.exe" 0
  
SectionEnd


;======================================================================================================================
; Description of the sections
!insertmacro MUI_FUNCTION_DESCRIPTION_BEGIN
	!insertmacro MUI_DESCRIPTION_TEXT ${SectionMain} "Install and start ${APPLICATION_NAME}."
	!insertmacro MUI_DESCRIPTION_TEXT ${SectionStartMenu} "Create Shortcuts on Start Menu."
!insertmacro MUI_FUNCTION_DESCRIPTION_END


;======================================================================================================================
; Uninstaller section
Section "Uninstall"

  ; --------------------------------------------------------------------------  
  ; Uninstall and stop  Vanadis View service
  ; TBD execWait "$INSTDIR\view-main\uninstall_service.bat"
 
  ; Remove registry keys
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPLICATION_FOLDER}"
  DeleteRegKey HKLM SOFTWARE\${APPLICATION_FOLDER}

  ; Remove shortcuts, if any
  Delete "$SMPROGRAMS\${APPLICATION_FOLDER}\*.*"
  RMDir "$SMPROGRAMS\${APPLICATION_FOLDER}"
	
	; Remove installation directory
	RMDir /r $INSTDIR\*

SectionEnd


