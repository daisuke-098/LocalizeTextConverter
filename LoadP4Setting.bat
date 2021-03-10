:: ------------------------------------------------------------
:: このバッチを別のバッチからコールすると、
:: TresEditorで設定したソースコントロール(Perforce)の設定を変数に保存します。
:: ------------------------------------------------------------
:: 変数: 説明
:: ------------------------------------------------------------
:: TEMP_P4AVAILABLE: Perforceが有効かどうか
:: TEMP_P4PORT: P4PORT
:: TEMP_P4USER: P4USER
:: TEMP_P4CLIENT: P4CLIENT
:: ------------------------------------------------------------
@echo off
pushd %~dp0

set FIND_CMD=C:\Windows\System32\findstr.exe /i /l

set TEMP_P4AVAILABLE=0
set TEMP_P4PORT=dummy
set TEMP_P4USER=dummy
set TEMP_P4CLIENT=dummy

set SETTING_INI_PATH=..\..\Saved\Config\Windows\SourceControlSettings.ini

if exist %SETTING_INI_PATH% (
	for /f "delims=" %%I in (%SETTING_INI_PATH%) do (
		set line=%%I
		call :CHECK
	)
) else (
	echo SourceControlSettings.iniが見つかりませんでした。
)
:: 変数削除
set SETTING_INI_PATH=

if %TEMP_P4PORT% == dummy (
	echo Perforceの設定の取得に失敗しました。
) else (
	set TEMP_P4AVAILABLE=1
	echo ----------------------------
	echo TEMP_P4PORT=%TEMP_P4PORT%
	echo TEMP_P4USER=%TEMP_P4USER%
	echo TEMP_P4CLIENT=%TEMP_P4CLIENT%
	echo ----------------------------
)

popd
exit /b 0

:CHECK
echo %line% | %FIND_CMD% "Port=" >NUL
if not ERRORLEVEL 1 (
	set TEMP_P4PORT=%line:~5%
)
echo %line% | %FIND_CMD% "UserName=" >NUL
if not ERRORLEVEL 1 (
	set TEMP_P4USER=%line:~9%
)
echo %line% | %FIND_CMD% "Workspace=" >NUL
if not ERRORLEVEL 1 (
	set TEMP_P4CLIENT=%line:~10%
)
goto :EOF
